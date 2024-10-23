// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package repository

import (
	"context"
	"database/sql"
	"fmt"
)

type DBTX interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

func New(db DBTX) *Queries {
	return &Queries{db: db}
}

func Prepare(ctx context.Context, db DBTX) (*Queries, error) {
	q := Queries{db: db}
	var err error
	if q.deleteUrlByShortUrlStmt, err = db.PrepareContext(ctx, deleteUrlByShortUrl); err != nil {
		return nil, fmt.Errorf("error preparing query DeleteUrlByShortUrl: %w", err)
	}
	if q.findUrlByShortUrlStmt, err = db.PrepareContext(ctx, findUrlByShortUrl); err != nil {
		return nil, fmt.Errorf("error preparing query FindUrlByShortUrl: %w", err)
	}
	if q.insertUrlStmt, err = db.PrepareContext(ctx, insertUrl); err != nil {
		return nil, fmt.Errorf("error preparing query InsertUrl: %w", err)
	}
	if q.updateUrlStmt, err = db.PrepareContext(ctx, updateUrl); err != nil {
		return nil, fmt.Errorf("error preparing query UpdateUrl: %w", err)
	}
	if q.updateVisitedCountUrlStmt, err = db.PrepareContext(ctx, updateVisitedCountUrl); err != nil {
		return nil, fmt.Errorf("error preparing query UpdateVisitedCountUrl: %w", err)
	}
	return &q, nil
}

func (q *Queries) Close() error {
	var err error
	if q.deleteUrlByShortUrlStmt != nil {
		if cerr := q.deleteUrlByShortUrlStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing deleteUrlByShortUrlStmt: %w", cerr)
		}
	}
	if q.findUrlByShortUrlStmt != nil {
		if cerr := q.findUrlByShortUrlStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing findUrlByShortUrlStmt: %w", cerr)
		}
	}
	if q.insertUrlStmt != nil {
		if cerr := q.insertUrlStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing insertUrlStmt: %w", cerr)
		}
	}
	if q.updateUrlStmt != nil {
		if cerr := q.updateUrlStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updateUrlStmt: %w", cerr)
		}
	}
	if q.updateVisitedCountUrlStmt != nil {
		if cerr := q.updateVisitedCountUrlStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updateVisitedCountUrlStmt: %w", cerr)
		}
	}
	return err
}

func (q *Queries) exec(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) (sql.Result, error) {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).ExecContext(ctx, args...)
	case stmt != nil:
		return stmt.ExecContext(ctx, args...)
	default:
		return q.db.ExecContext(ctx, query, args...)
	}
}

func (q *Queries) query(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) (*sql.Rows, error) {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).QueryContext(ctx, args...)
	case stmt != nil:
		return stmt.QueryContext(ctx, args...)
	default:
		return q.db.QueryContext(ctx, query, args...)
	}
}

func (q *Queries) queryRow(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) *sql.Row {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).QueryRowContext(ctx, args...)
	case stmt != nil:
		return stmt.QueryRowContext(ctx, args...)
	default:
		return q.db.QueryRowContext(ctx, query, args...)
	}
}

type Queries struct {
	db                        DBTX
	tx                        *sql.Tx
	deleteUrlByShortUrlStmt   *sql.Stmt
	findUrlByShortUrlStmt     *sql.Stmt
	insertUrlStmt             *sql.Stmt
	updateUrlStmt             *sql.Stmt
	updateVisitedCountUrlStmt *sql.Stmt
}

func (q *Queries) WithTx(tx *sql.Tx) *Queries {
	return &Queries{
		db:                        tx,
		tx:                        tx,
		deleteUrlByShortUrlStmt:   q.deleteUrlByShortUrlStmt,
		findUrlByShortUrlStmt:     q.findUrlByShortUrlStmt,
		insertUrlStmt:             q.insertUrlStmt,
		updateUrlStmt:             q.updateUrlStmt,
		updateVisitedCountUrlStmt: q.updateVisitedCountUrlStmt,
	}
}
