package postgres

import (
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// isPgUniqueViolation checks if error is a PostgreSQL unique violation
func isPgUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		// 23505 is the PostgreSQL error code for unique_violation
		return pgErr.Code == "23505"
	}
	return false
}

// isPgForeignKeyViolation checks if error is a PostgreSQL foreign key violation
func isPgForeignKeyViolation(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		// 23503 is the PostgreSQL error code for foreign_key_violation
		return pgErr.Code == "23503"
	}
	return false
}

// isPgNoRows checks if error is pgx.ErrNoRows (no rows in result set)
func isPgNoRows(err error) bool {
	return errors.Is(err, pgx.ErrNoRows)
}

// isPgCheckViolation checks if error is a PostgreSQL check constraint violation
func isPgCheckViolation(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		// 23514 is the PostgreSQL error code for check_violation
		return pgErr.Code == "23514"
	}
	return false
}
