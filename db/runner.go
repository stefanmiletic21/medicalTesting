package db

import (
	"context"
	"database/sql"
	"database/sql/driver"
)

type runnerKey struct{}

// RunnerKey is a key for context.Context to extract Context from it.
var RunnerKey = runnerKey{}

// Runner is an interface to access database
type Runner interface {

	// Transact executes the function passed as a parameter in a database transaction.
	// It does not support nested transactions. It is possible to have
	// a Transact wrapped in within a Transact, but in this case only
	// the first Transact will create a database transaction. The second
	// Transact will run in the first Transact's transaction (with its txOptions).
	Transact(ctx context.Context, txOptions *sql.TxOptions, txFunc func() error) error

	// Conn executes the function passed as a parameter. All calls
	// to the Runner within it will use the same DB connection.
	Conn(ctx context.Context, connFunc func() error) error

	// Query is a wrapper around sql.QueryContext function
	Query(ctx context.Context, query string, args ...interface{}) (rows *sql.Rows, err error)

	// QueryRow is a wrapper around sql.QueryRowContext function
	QueryRow(ctx context.Context, query string, args ...interface{}) (row *sql.Row)

	// Exec is a wrapper around sql.ExecContext function
	Exec(ctx context.Context, query string, args ...interface{}) (res sql.Result, err error)

	// Prepare is a wrapper around sql.PrepareContext function
	Prepare(ctx context.Context, query string) (stmt *sql.Stmt, err error)
}

type dbRunner struct {
	db      *sql.DB
	tx      *sql.Tx
	conn    *sql.Conn
	txCount int
}

// CreateRunner returns database context that is used to access database
func CreateRunner() Runner {
	run := new(dbRunner)
	run.db = handle
	return run
}

func (run *dbRunner) Transact(ctx context.Context, txOptions *sql.TxOptions, txFunc func() error) (err error) {
	if run.tx == nil {
		var tx *sql.Tx

		if run.conn == nil {
			tx, err = run.db.BeginTx(ctx, txOptions)
		} else {
			tx, err = run.conn.BeginTx(ctx, txOptions)
		}

		if err != nil {
			return err
		}

		run.tx = tx
		run.txCount = 1
	} else {
		run.txCount++
	}

	defer func() {
		// Recover from panic
		p := recover()

		// Rollback the transaction in case of error or panic
		if (err != nil || p != nil) && run.tx != nil {
			run.tx.Rollback()
			run.tx = nil
			run.txCount = 0
		}

		// Re-panic if was panicking
		if p != nil {
			panic(p)
		}

		if err != nil {
			return
		}

		if run.tx == nil {
			panic("Transaction is already rolled back or commited")
		}

		// Decrement tx counter and commit tx if it's the last one in Transact chain
		run.txCount--
		if run.txCount == 0 {
			err = run.tx.Commit()
			run.tx = nil
		}
	}()

	err = txFunc()

	return
}

func (run *dbRunner) Conn(ctx context.Context, connFunc func() error) (err error) {
	// If in transaction or already using single connection just call the function
	if run.tx != nil || run.conn != nil {
		return connFunc()
	}

	run.conn, err = run.db.Conn(ctx)
	if err == driver.ErrBadConn {
		run.conn, err = run.db.Conn(ctx)
	}

	if err != nil {
		return
	}

	defer func() {
		errClose := run.conn.Close()
		run.conn = nil
		if err == nil {
			err = errClose
		}
	}()

	err = connFunc()

	return
}

func (run *dbRunner) Query(ctx context.Context, query string, args ...interface{}) (rows *sql.Rows, err error) {
	if run.tx != nil {
		rows, err = run.tx.QueryContext(ctx, query, args...)
	} else if run.conn != nil {
		rows, err = run.conn.QueryContext(ctx, query, args...)
	} else {
		rows, err = run.db.QueryContext(ctx, query, args...)
	}

	return
}

func (run *dbRunner) QueryRow(ctx context.Context, query string, args ...interface{}) (row *sql.Row) {
	if run.tx != nil {
		row = run.tx.QueryRowContext(ctx, query, args...)
	} else if run.conn != nil {
		row = run.conn.QueryRowContext(ctx, query, args...)
	} else {
		row = run.db.QueryRowContext(ctx, query, args...)
	}

	return
}

func (run *dbRunner) Exec(ctx context.Context, query string, args ...interface{}) (res sql.Result, err error) {
	if run.tx != nil {
		res, err = run.tx.ExecContext(ctx, query, args...)
	} else if run.conn != nil {
		res, err = run.conn.ExecContext(ctx, query, args...)
	} else {
		res, err = run.db.ExecContext(ctx, query, args...)
	}

	return
}

func (run *dbRunner) Prepare(ctx context.Context, query string) (stmt *sql.Stmt, err error) {
	if run.tx != nil {
		stmt, err = run.tx.PrepareContext(ctx, query)
	} else if run.conn != nil {
		stmt, err = run.conn.PrepareContext(ctx, query)
	} else {
		stmt, err = run.db.PrepareContext(ctx, query)
	}

	return
}

type dbMockRunner struct{}

// CreateMockRunner returns mocked database runner that panics if used to access database
func CreateMockRunner() Runner {
	return dbMockRunner{}
}

func (dbMockRunner) Transact(ctx context.Context, txOptions *sql.TxOptions, txFunc func() error) error {
	return txFunc()
}

func (dbMockRunner) Conn(ctx context.Context, connFunc func() error) error {
	return connFunc()
}

func (dbMockRunner) Query(context.Context, string, ...interface{}) (*sql.Rows, error) {
	panic("Can't Query in mocked runner")
}

func (dbMockRunner) QueryRow(context.Context, string, ...interface{}) *sql.Row {
	panic("Can't QueryRow in mocked runner")
}

func (dbMockRunner) Exec(context.Context, string, ...interface{}) (sql.Result, error) {
	panic("Can't Exec in mocked runner")
}

func (dbMockRunner) Prepare(context.Context, string) (*sql.Stmt, error) {
	panic("Can't Prepare in mocked runner")
}
