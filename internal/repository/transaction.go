package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// TransactionManager manages database transactions
type TransactionManager interface {
	WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}

// PgxTransactionManager implements TransactionManager for pgx
type PgxTransactionManager struct {
	pool *pgxpool.Pool
}

// NewPgxTransactionManager creates a new transaction manager
func NewPgxTransactionManager(pool *pgxpool.Pool) *PgxTransactionManager {
	return &PgxTransactionManager{pool: pool}
}

// WithTransaction executes a function within a database transaction
// If the function returns an error, the transaction is rolled back
// Otherwise, the transaction is committed
func (tm *PgxTransactionManager) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := tm.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Create a context with the transaction
	txCtx := context.WithValue(ctx, txKey{}, tx)

	// Execute the function
	err = fn(txCtx)
	if err != nil {
		// Rollback on error
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("failed to rollback transaction: %v (original error: %w)", rbErr, err)
		}
		return err
	}

	// Commit on success
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// txKey is a context key for storing transaction
type txKey struct{}

// Executor is a common interface for pool and transaction
type Executor interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
}

// GetTx retrieves a transaction from context, or returns the pool if no transaction exists
func GetTx(ctx context.Context, pool *pgxpool.Pool) Executor {
	if tx, ok := ctx.Value(txKey{}).(pgx.Tx); ok {
		return tx
	}
	return pool
}
