// Package v1inspect performs read-only migration preflight checks against a
// V1 SQLite database. It is intentionally separate from the V2 store: V1 is
// an external source and must never be opened through the V2 migration runner.
package v1inspect

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"

	_ "modernc.org/sqlite"
)

const driverName = "sqlite"

type Severity string

const (
	SeverityBlocker Severity = "blocker"
	SeverityWarning Severity = "warning"
)

type Finding struct {
	Severity Severity `json:"severity"`
	Code     string   `json:"code"`
	Table    string   `json:"table,omitempty"`
	Count    int64    `json:"count,omitempty"`
	Message  string   `json:"message"`
}

type Table struct {
	Name     string   `json:"name"`
	Columns  []string `json:"columns"`
	RowCount int64    `json:"rowCount"`
}

type MigrationLedger struct {
	Present    bool `json:"present"`
	Version    int  `json:"version,omitempty"`
	Contiguous bool `json:"contiguous"`
}

type Report struct {
	SourcePath string          `json:"sourcePath"`
	Tables     []Table         `json:"tables"`
	Migrations MigrationLedger `json:"migrations"`
	Findings   []Finding       `json:"findings"`
}

var expectedTables = map[string][]string{
	"activation":                       {"gym_id", "gym_name", "region"},
	"customers":                        {"id", "email", "identification_number", "external_id", "has_access", "status", "gym_id"},
	"customer_memberships":             {"id", "customer_id", "membership_plan_id", "start_date", "membership_price_usd", "visits_remaining"},
	"customer_progress_tracking":       {"id", "customer_id"},
	"customer_visits":                  {"id", "customer_id", "membership_id", "visit_date", "membership_status"},
	"device_config":                    {"gym_id", "ip", "username", "password"},
	"membership_plans":                 {"id", "type", "price", "duration", "visits", "gym_id"},
	"payments":                         {"id", "customer_id", "customer_membership_id", "amount_usd", "transaction_type"},
	"products":                         {"id", "price", "stock", "gym_id"},
	"sale_items":                       {"id", "sale_id", "product_id", "unit_price", "subtotal"},
	"sales":                            {"id", "total_amount", "customer_id", "gym_id"},
	"sessions":                         {"id", "user_id", "token", "expires_at"},
	"system_users":                     {"id", "username", "email", "role", "status", "gym_id"},
	"survey_questions":                 {"id", "question", "rating_type"},
	"audit_log":                        {"id", "action", "entity_type"},
	"alert_log":                        {"id", "alert_type", "gym_id"},
	"marketing_fidelity_notifications": {"id", "customer_id", "gym_id"},
	"marketing_survey_invitations":     {"id", "customer_id", "customer_membership_id", "gym_id", "token"},
	"survey_responses":                 {"id", "invitation_id", "question_text", "answer_value"},
}

// Inspect opens path in SQLite read-only mode and returns source-schema and
// data-quality findings needed before a V1-to-V2 dry run. It never alters the
// source database, creates journal files, or returns secret values.
func Inspect(ctx context.Context, path string) (Report, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return Report{}, fmt.Errorf("resolve source path: %w", err)
	}
	info, err := os.Stat(absPath)
	if err != nil {
		return Report{}, fmt.Errorf("stat source database: %w", err)
	}
	if info.IsDir() {
		return Report{}, errors.New("source database path is a directory")
	}

	db, err := sql.Open(driverName, readOnlyDatabaseURL(absPath))
	if err != nil {
		return Report{}, fmt.Errorf("open source database read-only: %w", err)
	}
	defer db.Close()
	db.SetMaxOpenConns(1)
	if err := db.PingContext(ctx); err != nil {
		return Report{}, fmt.Errorf("connect source database read-only: %w", err)
	}

	report := Report{SourcePath: absPath}
	tables, err := loadTables(ctx, db)
	if err != nil {
		return Report{}, err
	}
	report.Tables = tables
	tableByName := make(map[string]Table, len(tables))
	for _, table := range tables {
		tableByName[table.Name] = table
	}

	report.Migrations, report.Findings, err = inspectMigrationLedger(ctx, db, tableByName, report.Findings)
	if err != nil {
		return Report{}, err
	}
	report.Findings = append(report.Findings, schemaFindings(tableByName)...)
	report.Findings, err = appendDataFindings(ctx, db, tableByName, report.Findings)
	if err != nil {
		return Report{}, err
	}

	sort.Slice(report.Findings, func(i, j int) bool {
		if report.Findings[i].Severity != report.Findings[j].Severity {
			return report.Findings[i].Severity == SeverityBlocker
		}
		return report.Findings[i].Code < report.Findings[j].Code
	})
	return report, nil
}

func readOnlyDatabaseURL(path string) string {
	u := url.URL{Scheme: "file", Path: path}
	query := u.Query()
	query.Set("mode", "ro")
	query.Add("_pragma", "query_only(1)")
	u.RawQuery = query.Encode()
	return u.String()
}

func loadTables(ctx context.Context, db *sql.DB) ([]Table, error) {
	rows, err := db.QueryContext(ctx, `
SELECT name
FROM sqlite_master
WHERE type = 'table' AND name NOT LIKE 'sqlite_%'
ORDER BY name
`)
	if err != nil {
		return nil, fmt.Errorf("list source tables: %w", err)
	}

	var names []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, fmt.Errorf("scan source table: %w", err)
		}
		names = append(names, name)
	}
	if err := rows.Err(); err != nil {
		_ = rows.Close()
		return nil, fmt.Errorf("iterate source tables: %w", err)
	}
	if err := rows.Close(); err != nil {
		return nil, fmt.Errorf("close source table list: %w", err)
	}

	tables := make([]Table, 0, len(names))
	for _, name := range names {
		columns, err := loadColumns(ctx, db, name)
		if err != nil {
			return nil, err
		}
		var rowCount int64
		if err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM "+quoteIdentifier(name)).Scan(&rowCount); err != nil {
			return nil, fmt.Errorf("count rows in %s: %w", name, err)
		}
		tables = append(tables, Table{Name: name, Columns: columns, RowCount: rowCount})
	}
	return tables, nil
}

func loadColumns(ctx context.Context, db *sql.DB, table string) ([]string, error) {
	rows, err := db.QueryContext(ctx, "PRAGMA table_info("+quoteIdentifier(table)+")")
	if err != nil {
		return nil, fmt.Errorf("inspect columns for %s: %w", table, err)
	}
	defer rows.Close()

	var columns []string
	for rows.Next() {
		var cid int
		var name, columnType string
		var notNull, primaryKey int
		var defaultValue any
		if err := rows.Scan(&cid, &name, &columnType, &notNull, &defaultValue, &primaryKey); err != nil {
			return nil, fmt.Errorf("scan column for %s: %w", table, err)
		}
		columns = append(columns, name)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate columns for %s: %w", table, err)
	}
	return columns, nil
}

func inspectMigrationLedger(ctx context.Context, db *sql.DB, tables map[string]Table, findings []Finding) (MigrationLedger, []Finding, error) {
	ledger, ok := tables["migrations"]
	if !ok || !hasColumns(ledger, "version") {
		return MigrationLedger{}, append(findings, Finding{
			Severity: SeverityWarning,
			Code:     "source_migration_ledger_missing",
			Table:    "migrations",
			Message:  "V1 migration history is unavailable; source schema drift cannot be measured from the ledger.",
		}), nil
	}

	rows, err := db.QueryContext(ctx, "SELECT version FROM migrations ORDER BY version")
	if err != nil {
		return MigrationLedger{}, findings, fmt.Errorf("read V1 migration ledger: %w", err)
	}
	defer rows.Close()

	ledgerState := MigrationLedger{Present: true, Contiguous: true}
	expected := 1
	for rows.Next() {
		var version int
		if err := rows.Scan(&version); err != nil {
			return MigrationLedger{}, findings, fmt.Errorf("scan V1 migration version: %w", err)
		}
		if version != expected {
			ledgerState.Contiguous = false
		}
		expected = version + 1
		ledgerState.Version = version
	}
	if err := rows.Err(); err != nil {
		return MigrationLedger{}, findings, fmt.Errorf("iterate V1 migration versions: %w", err)
	}
	if !ledgerState.Contiguous || ledgerState.Version != 28 {
		findings = append(findings, Finding{
			Severity: SeverityWarning,
			Code:     "source_migration_ledger_incomplete",
			Table:    "migrations",
			Count:    int64(ledgerState.Version),
			Message:  "V1 migration ledger is non-contiguous or does not end at the expected version 28; inspect the actual source schema before migration.",
		})
	}
	return ledgerState, findings, nil
}

func schemaFindings(tables map[string]Table) []Finding {
	var findings []Finding
	for tableName, requiredColumns := range expectedTables {
		table, ok := tables[tableName]
		if !ok {
			findings = append(findings, Finding{
				Severity: SeverityWarning,
				Code:     "expected_table_missing",
				Table:    tableName,
				Message:  "Expected V1 table is missing; confirm whether this database is an incomplete or customized V1 installation.",
			})
			continue
		}
		for _, column := range requiredColumns {
			if hasColumns(table, column) {
				continue
			}
			findings = append(findings, Finding{
				Severity: SeverityWarning,
				Code:     "expected_column_missing",
				Table:    tableName,
				Message:  fmt.Sprintf("Expected V1 column %q is missing; automatic conversion for this table may be unavailable.", column),
			})
		}
	}
	return findings
}

func appendDataFindings(ctx context.Context, db *sql.DB, tables map[string]Table, findings []Finding) ([]Finding, error) {
	checks := []struct {
		table    string
		columns  []string
		code     string
		message  string
		query    string
		severity Severity
	}{
		{"customers", []string{"identification_number"}, "duplicate_member_identification", "Customers share an identification number; migration needs an operator resolution.", `SELECT COUNT(*) FROM (SELECT identification_number FROM customers WHERE identification_number IS NOT NULL AND TRIM(identification_number) <> '' GROUP BY identification_number HAVING COUNT(*) > 1)`, SeverityBlocker},
		{"customers", []string{"external_id"}, "duplicate_device_external_id", "Customers share a device external ID; device identities cannot be imported safely.", `SELECT COUNT(*) FROM (SELECT external_id FROM customers WHERE external_id IS NOT NULL AND TRIM(external_id) <> '' GROUP BY external_id HAVING COUNT(*) > 1)`, SeverityBlocker},
		{"customers", []string{"has_access"}, "invalid_legacy_access_flag", "Customers have a legacy has_access value outside 0 or 1; never infer a V2 override from it.", `SELECT COUNT(*) FROM customers WHERE has_access NOT IN (0, 1) OR has_access IS NULL`, SeverityBlocker},
		{"membership_plans", []string{"type", "duration", "visits"}, "invalid_membership_plan_validity", "Membership plans have an invalid V1 type/duration/visit combination.", `SELECT COUNT(*) FROM membership_plans WHERE (type = 'count-based' AND (visits IS NULL OR visits <= 0 OR duration IS NOT NULL)) OR (type <> 'count-based' AND (duration IS NULL OR duration <= 0 OR visits IS NOT NULL))`, SeverityBlocker},
		{"payments", []string{"amount_usd"}, "invalid_payment_amount", "Payments contain a negative or missing amount.", `SELECT COUNT(*) FROM payments WHERE amount_usd IS NULL OR amount_usd < 0`, SeverityBlocker},
		{"payments", []string{"amount_usd"}, "payment_amount_requires_rounding", "Payment amounts have more than two decimal places and need an approved cents-conversion rule.", `SELECT COUNT(*) FROM payments WHERE amount_usd IS NOT NULL AND ABS(amount_usd * 100 - ROUND(amount_usd * 100)) > 0.000001`, SeverityWarning},
		{"payments", []string{"transaction_type"}, "unknown_payment_transaction_type", "Payments contain a transaction type other than income or expense.", `SELECT COUNT(*) FROM payments WHERE transaction_type NOT IN ('income', 'expense') OR transaction_type IS NULL`, SeverityBlocker},
	}

	for _, check := range checks {
		table, ok := tables[check.table]
		if !ok || !hasColumns(table, check.columns...) {
			continue
		}
		count, err := queryCount(ctx, db, check.query)
		if err != nil {
			return nil, fmt.Errorf("run %s check: %w", check.code, err)
		}
		if count > 0 {
			findings = append(findings, Finding{Severity: check.severity, Code: check.code, Table: check.table, Count: count, Message: check.message})
		}
	}

	if table, ok := tables["payments_backup"]; ok {
		findings = append(findings, Finding{
			Severity: SeverityWarning,
			Code:     "payments_backup_present",
			Table:    table.Name,
			Count:    table.RowCount,
			Message:  "V1 migration 023 may have left a payments backup table; it is diagnostic only and must not be imported as a second payment source.",
		})
	}
	if _, ok := tables["device_config"]; ok {
		findings = append(findings, Finding{
			Severity: SeverityWarning,
			Code:     "device_credentials_require_secure_reentry",
			Table:    "device_config",
			Message:  "V1 device credentials are not exported by this inspector and must be re-entered or migrated through an approved secure mechanism.",
		})
	}

	foreignKeyViolations, err := queryCount(ctx, db, "SELECT COUNT(*) FROM pragma_foreign_key_check")
	if err != nil {
		return nil, fmt.Errorf("check V1 foreign keys: %w", err)
	}
	if foreignKeyViolations > 0 {
		findings = append(findings, Finding{
			Severity: SeverityBlocker,
			Code:     "foreign_key_violation",
			Count:    foreignKeyViolations,
			Message:  "The source database contains foreign-key violations; affected rows need resolution before automatic migration.",
		})
	}
	return findings, nil
}

func queryCount(ctx context.Context, db *sql.DB, query string) (int64, error) {
	var count int64
	if err := db.QueryRowContext(ctx, query).Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

func hasColumns(table Table, required ...string) bool {
	available := make(map[string]struct{}, len(table.Columns))
	for _, column := range table.Columns {
		available[column] = struct{}{}
	}
	for _, column := range required {
		if _, ok := available[column]; !ok {
			return false
		}
	}
	return true
}

func quoteIdentifier(identifier string) string {
	return `"` + strings.ReplaceAll(identifier, `"`, `""`) + `"`
}
