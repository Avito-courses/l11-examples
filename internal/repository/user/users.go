package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	model "github.com/Avito-courses/l11-examples/internal/model/user"
)

const tableUsers = "users"

type Repository struct {
	pool         *pgxpool.Pool
	queryBuilder squirrel.StatementBuilderType
}

func NewUserRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{
		pool:         pool,
		queryBuilder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

// GetByID возвращает пользователя по ID
func (r *Repository) GetByID(ctx context.Context, id int) (*model.User, error) {
	query := r.queryBuilder.
		Select("id", "name", "phone", "rating", "created_at", "updated_at").
		From(tableUsers).
		Where(squirrel.Eq{"id": id})

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	var user model.User
	err = r.pool.QueryRow(ctx, sql, args...).Scan(
		&user.ID,
		&user.Name,
		&user.Phone,
		&user.Rating,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, model.ErrUserNotFound
		}
		return nil, fmt.Errorf("database error: %w", err)
	}

	return &user, nil
}

// GetByRating возвращает всех пользователей c указанным рейтингом
func (r *Repository) GetByRating(ctx context.Context, ratingFrom, ratingTo int) ([]model.User, error) {
	query, args, err := r.queryBuilder.
		Select("id", "name", "phone", "rating", "created_at", "updated_at").
		From(tableUsers).
		Where(squirrel.And{squirrel.Gt{"rating": ratingFrom}, squirrel.Lt{"rating": ratingTo}}).
		OrderBy("id").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	defer rows.Close()

	users := make([]model.User, 0)
	for rows.Next() {
		var user model.User
		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Phone,
			&user.Rating,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error reading data: %w", err)
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}

	return users, nil
}

// Create создает нового пользователя
func (r *Repository) Create(ctx context.Context, user model.User) (id int, err error) {
	query, args, err := r.queryBuilder.
		Insert(tableUsers).
		Columns("name", "phone", "rating").
		Values(user.Name, user.Phone, user.Rating).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return id, fmt.Errorf("failed to build query: %w", err)
	}

	err = r.pool.QueryRow(ctx, query, args...).Scan(&id)
	if err != nil {
		return id, handleUniqueViolationError(err)
	}

	return id, nil
}

func handleUniqueViolationError(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
		return model.ErrPhoneExists
	}
	return fmt.Errorf("database error: %w", err)
}
