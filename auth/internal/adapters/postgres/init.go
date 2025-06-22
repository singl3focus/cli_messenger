package postgres

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/singl3focus/cli_messenger/auth/internal/core/interfaces"
	"github.com/singl3focus/cli_messenger/auth/internal/core/models"
)

var _ interfaces.Repostory = &Database{}

type Database struct {
	pool *pgxpool.Pool
}

var sqBuilder = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

func NewPostgres(dsn string) (*Database, error) {
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}

	return &Database{
		pool: pool,
	}, nil
}

func (d *Database) Close() {
	d.pool.Close()
}

const (
	tblUser = "auth.tbl_user"
)

func (d *Database) CreateUser(ctx context.Context, user models.User) (int64, error) {
	const op = "postgres.CreateUser"
	
	query, args, err := sqBuilder.Insert(tblUser).
			Columns("name", "email", "password_hash", "role", "created_at", "updated_at").
			Values(user.Name, user.Email, user.PasswordHash, user.Role, user.CreatedAt, user.UpdatedAt).
			Suffix("RETURNING id").
			ToSql()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	var id int64
	err = d.pool.QueryRow(ctx, query, args...).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (d *Database) GetUser(ctx context.Context, id int64) (models.User, error) {
	const op = "postgres.GetUser"

	var user models.User

	query, args, err := sqBuilder.
		Select("id", "name", "email", "password_hash", "role", "created_at", "updated_at").
		From(tblUser).
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return user, fmt.Errorf("%s: %w", op, err)
	}

	err = d.pool.QueryRow(ctx, query, args...).
		Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return user, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (d *Database) UpdateUser(ctx context.Context, id int64, name, email *string) error {
	const op = "postgres.UpdateUser"

	builder := sqBuilder.Update(tblUser).
		Where(sq.Eq{"id": id})

	if name != nil {
		builder = builder.Set("name", name)
	}
	if email != nil {
		builder = builder.Set("email", email)
	}

	query, args, err := builder.ToSql()	
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	result, err := d.pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if rowsAffected := result.RowsAffected(); rowsAffected == 0 {
    	return fmt.Errorf("%s: no rows affected (user not found?)", op)
	}

	return nil
}

func (d *Database) DeleteUser(ctx context.Context, id int64) error {
	const op = "postgres.DeleteUser"

	query, args, err := sqBuilder.Delete(tblUser).
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	result, err := d.pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if rowsAffected := result.RowsAffected(); rowsAffected == 0 {
    	return fmt.Errorf("%s: no rows affected (user not found?)", op)
	}

	return nil
}