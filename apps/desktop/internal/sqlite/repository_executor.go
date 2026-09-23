package sqlite

import (
	"context"
	"database/sql"
)

// sqlExecutor is the shared subset of *sql.DB and *sql.Tx used by repository
// inserts. It keeps transaction mechanics inside infrastructure.
type sqlExecutor interface {
	ExecContext(context.Context, string, ...any) (sql.Result, error)
}
