// internal/repository/postgres/user_postgres.go
package postgres

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"my-application/internal/domain"
	"my-application/internal/repository"
)

// Compile-time interface check.
var _ repository.UserRepository = (*UserPostgres)(nil)

// UserPostgres implements repository.UserRepository with PostgreSQL.
type UserPostgres struct {
	pool   *pgxpool.Pool
	logger *slog.Logger
}

// NewUserPostgres creates a new UserPostgres repository.
func NewUserPostgres(pool *pgxpool.Pool, logger *slog.Logger) *UserPostgres {
	return &UserPostgres{pool: pool, logger: logger}
}

// columns shared across single-row queries.
const userColumns = `id, username, email, password_hash, full_name, role, enterprise_id, firebase_uid, is_active, created_at, updated_at`

// scanUser scans a row into a domain.User.
func scanUser(row pgx.Row) (*domain.User, error) {
	var u domain.User
	err := row.Scan(
		&u.ID, &u.Username, &u.Email, &u.PasswordHash, &u.FullName,
		&u.Role, &u.EnterpriseID, &u.FirebaseUID, &u.IsActive, &u.CreatedAt, &u.UpdatedAt,
	)
	return &u, err
}

func (r *UserPostgres) GetByID(ctx context.Context, id int64) (*domain.User, error) {
	query := `SELECT ` + userColumns + ` FROM users WHERE id = $1`

	u, err := scanUser(r.pool.QueryRow(ctx, query, id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.NewAppError(domain.ErrNotFound, fmt.Sprintf("user with id %d not found", id))
		}
		return nil, domain.NewAppError(domain.ErrDatabaseOperation, err.Error())
	}
	return u, nil
}

func (r *UserPostgres) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `SELECT ` + userColumns + ` FROM users WHERE email = $1`

	u, err := scanUser(r.pool.QueryRow(ctx, query, email))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.NewAppError(domain.ErrNotFound, "user not found")
		}
		return nil, domain.NewAppError(domain.ErrDatabaseOperation, err.Error())
	}
	return u, nil
}

func (r *UserPostgres) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	query := `SELECT ` + userColumns + ` FROM users WHERE username = $1`

	u, err := scanUser(r.pool.QueryRow(ctx, query, username))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.NewAppError(domain.ErrNotFound, "user not found")
		}
		return nil, domain.NewAppError(domain.ErrDatabaseOperation, err.Error())
	}
	return u, nil
}

func (r *UserPostgres) GetByFirebaseUID(ctx context.Context, firebaseUID string) (*domain.User, error) {
	query := `SELECT ` + userColumns + ` FROM users WHERE firebase_uid = $1`

	u, err := scanUser(r.pool.QueryRow(ctx, query, firebaseUID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.NewAppError(domain.ErrNotFound, "user not found")
		}
		return nil, domain.NewAppError(domain.ErrDatabaseOperation, err.Error())
	}
	return u, nil
}

func (r *UserPostgres) List(ctx context.Context, filter domain.UserFilter) ([]domain.User, int64, error) {
	baseQuery := `SELECT ` + userColumns + ` FROM users WHERE 1=1`
	countQuery := `SELECT COUNT(*) FROM users WHERE 1=1`
	args := []interface{}{}
	argIdx := 1

	if filter.Role != "" {
		condition := fmt.Sprintf(" AND role = $%d", argIdx)
		baseQuery += condition
		countQuery += condition
		args = append(args, filter.Role)
		argIdx++
	}
	if filter.IsActive != nil {
		condition := fmt.Sprintf(" AND is_active = $%d", argIdx)
		baseQuery += condition
		countQuery += condition
		args = append(args, *filter.IsActive)
		argIdx++
	}

	// Total count.
	var total int64
	if err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, domain.NewAppError(domain.ErrDatabaseOperation, err.Error())
	}

	// Ensure pagination bounds (defensive — service layer normalizes first).
	filter.Normalize()

	baseQuery += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
	args = append(args, filter.Limit, filter.Offset)

	rows, err := r.pool.Query(ctx, baseQuery, args...)
	if err != nil {
		return nil, 0, domain.NewAppError(domain.ErrDatabaseOperation, err.Error())
	}
	defer rows.Close()

	users := make([]domain.User, 0)
	for rows.Next() {
		var u domain.User
		if err := rows.Scan(
			&u.ID, &u.Username, &u.Email, &u.PasswordHash, &u.FullName,
			&u.Role, &u.EnterpriseID, &u.FirebaseUID, &u.IsActive, &u.CreatedAt, &u.UpdatedAt,
		); err != nil {
			return nil, 0, domain.NewAppError(domain.ErrDatabaseOperation, err.Error())
		}
		users = append(users, u)
	}

	return users, total, nil
}

func (r *UserPostgres) Create(ctx context.Context, user *domain.User) error {
	query := `INSERT INTO users (username, email, password_hash, full_name, role, enterprise_id, firebase_uid, is_active)
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
			  RETURNING id, created_at, updated_at`

	err := r.pool.QueryRow(ctx, query,
		user.Username, user.Email, user.PasswordHash, user.FullName,
		user.Role, user.EnterpriseID, user.FirebaseUID, user.IsActive,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return domain.NewAppError(domain.ErrAlreadyExists, "user with this username or email already exists")
		}
		return domain.NewAppError(domain.ErrDatabaseOperation, err.Error())
	}
	return nil
}

func (r *UserPostgres) Update(ctx context.Context, user *domain.User) error {
	query := `UPDATE users SET username=$1, email=$2, full_name=$3, role=$4, enterprise_id=$5, is_active=$6
			  WHERE id=$7 RETURNING updated_at`

	err := r.pool.QueryRow(ctx, query,
		user.Username, user.Email, user.FullName, user.Role, user.EnterpriseID, user.IsActive, user.ID,
	).Scan(&user.UpdatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.NewAppError(domain.ErrNotFound, fmt.Sprintf("user with id %d not found", user.ID))
		}
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return domain.NewAppError(domain.ErrAlreadyExists, "user with this username or email already exists")
		}
		return domain.NewAppError(domain.ErrDatabaseOperation, err.Error())
	}
	return nil
}

func (r *UserPostgres) Delete(ctx context.Context, id int64) error {
	result, err := r.pool.Exec(ctx, `DELETE FROM users WHERE id = $1`, id)
	if err != nil {
		return domain.NewAppError(domain.ErrDatabaseOperation, err.Error())
	}
	if result.RowsAffected() == 0 {
		return domain.NewAppError(domain.ErrNotFound, fmt.Sprintf("user with id %d not found", id))
	}
	return nil
}
