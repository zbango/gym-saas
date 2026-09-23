package v1inspect

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"os"
	"path/filepath"
	"slices"
	"testing"
)

func TestInspectSchemaOnlyFixture(t *testing.T) {
	path := createFixtureDatabase(t, "schema-only.sql")

	report, err := Inspect(context.Background(), path)
	if err != nil {
		t.Fatalf("Inspect returned error: %v", err)
	}
	if !report.Migrations.Present || report.Migrations.Version != 28 || !report.Migrations.Contiguous {
		t.Fatalf("migration ledger = %#v, want contiguous V1 version 28", report.Migrations)
	}
	if len(report.Tables) != 20 {
		t.Fatalf("table count = %d, want 20", len(report.Tables))
	}
	if hasFinding(report, "expected_table_missing") || hasFinding(report, "expected_column_missing") {
		t.Fatalf("schema-only fixture has unexpected schema findings: %#v", report.Findings)
	}
}

func TestInspectRepresentativeFixtureReportsNoBlockers(t *testing.T) {
	path := createFixtureDatabase(t, "schema-only.sql", "representative.sql")

	report, err := Inspect(context.Background(), path)
	if err != nil {
		t.Fatalf("Inspect returned error: %v", err)
	}
	if tableRowCount(report, "customers") != 1 || tableRowCount(report, "payments") != 2 {
		t.Fatalf("unexpected representative counts: %#v", report.Tables)
	}
	if hasSeverity(report, SeverityBlocker) {
		t.Fatalf("representative fixture has blocker findings: %#v", report.Findings)
	}
	if !hasFinding(report, "device_credentials_require_secure_reentry") {
		t.Fatalf("representative fixture did not flag secure device credential re-entry: %#v", report.Findings)
	}
}

func TestInspectEdgeCaseFixtureReportsMigrationBlockers(t *testing.T) {
	path := createFixtureDatabase(t, "schema-only.sql", "edge-case.sql")

	report, err := Inspect(context.Background(), path)
	if err != nil {
		t.Fatalf("Inspect returned error: %v", err)
	}
	wantCodes := []string{
		"duplicate_device_external_id",
		"duplicate_member_identification",
		"foreign_key_violation",
		"invalid_legacy_access_flag",
		"invalid_membership_plan_validity",
		"invalid_payment_amount",
		"payment_amount_requires_rounding",
		"payments_backup_present",
		"unknown_payment_transaction_type",
	}
	for _, code := range wantCodes {
		if !hasFinding(report, code) {
			t.Errorf("missing finding %q in %#v", code, report.Findings)
		}
	}
}

func TestInspectDoesNotModifySourceDatabase(t *testing.T) {
	path := createFixtureDatabase(t, "schema-only.sql", "representative.sql")
	before := fileChecksum(t, path)

	if _, err := Inspect(context.Background(), path); err != nil {
		t.Fatalf("Inspect returned error: %v", err)
	}
	if after := fileChecksum(t, path); before != after {
		t.Fatal("Inspect modified the source database")
	}
}

func createFixtureDatabase(t *testing.T, files ...string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "v1-fixture.sqlite")
	db, err := sql.Open(driverName, "file:"+path)
	if err != nil {
		t.Fatalf("open fixture database: %v", err)
	}
	defer db.Close()
	for _, name := range files {
		sqlSource, err := os.ReadFile(filepath.Join("testdata", "v1", name))
		if err != nil {
			t.Fatalf("read fixture %s: %v", name, err)
		}
		if _, err := db.Exec(string(sqlSource)); err != nil {
			t.Fatalf("apply fixture %s: %v", name, err)
		}
	}
	return path
}

func tableRowCount(report Report, name string) int64 {
	for _, table := range report.Tables {
		if table.Name == name {
			return table.RowCount
		}
	}
	return -1
}

func hasFinding(report Report, code string) bool {
	return slices.ContainsFunc(report.Findings, func(finding Finding) bool { return finding.Code == code })
}

func hasSeverity(report Report, severity Severity) bool {
	return slices.ContainsFunc(report.Findings, func(finding Finding) bool { return finding.Severity == severity })
}

func fileChecksum(t *testing.T, path string) [sha256.Size]byte {
	t.Helper()
	contents, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read source database: %v", err)
	}
	return sha256.Sum256(contents)
}
