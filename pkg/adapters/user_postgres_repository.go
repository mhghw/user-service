package adapters

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	userPkg "github.com/mhghw/user-service/pkg/domain/user"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const (
	inserUserQuery     = `INSERT INTO users (id,username,firstname,lastname) VALUES (:id,:username,:firstname,:lastname)`
	deleteUserQuery    = `DELETE FROM users WHERE id=?`
	getUserQuery       = `SELECT id,username,firstname,lastname FROM users WHERE id=?`
	upsertQuery        = inserUserQuery + ` ON CONFLICT (id) DO UPDATE SET username=:username`
	checkUsernameQuery = `SELECT username FROM users WHERE username=?`
)

type UserPostgresRepository struct {
	db *sqlx.DB
}

func NewUserPostgresRepository(db *sqlx.DB) UserPostgresRepository {
	return UserPostgresRepository{
		db: db,
	}
}

type UserModel struct {
	ID        string `db:"id"`
	Username  string `db:"username"`
	Firstname string `db:"firstname"`
	Lastname  string `db:"lastname"`
}

func (r UserPostgresRepository) bindQueryVariables(q string) string {
	return r.db.Rebind(q)
}

func (r UserPostgresRepository) CreateUser(ctx context.Context, u userPkg.User) error {
	model := UserModel(u)

	_, err := r.db.NamedExecContext(ctx, inserUserQuery, model)
	if err != nil {
		return fmt.Errorf("error inserting user to database: %w", err)
	}

	return nil
}
func (r UserPostgresRepository) UpdateUser(ctx context.Context, userID string, updateFn func(context.Context, userPkg.User) (userPkg.User, error)) error {
	user, err := r.GetUser(ctx, userID)
	if err != nil {
		return fmt.Errorf("cannot get user from database to update: %w", err)
	}

	updatedUser, err := updateFn(ctx, user)
	if err != nil {
		return fmt.Errorf("error updating user: %w", err)
	}

	_, err = r.db.NamedExecContext(ctx, upsertQuery, updatedUser)
	if err != nil {
		return fmt.Errorf("error upserting user: %w", err)
	}

	return nil
}
func (r UserPostgresRepository) DeleteUser(ctx context.Context, userID string) error {
	_, err := r.db.ExecContext(ctx, r.bindQueryVariables(deleteUserQuery), userID)
	if err != nil {
		return fmt.Errorf("error deleting user from database: %w", err)
	}

	return nil
}
func (r UserPostgresRepository) GetUser(ctx context.Context, userID string) (userPkg.User, error) {
	var model UserModel
	err := r.db.QueryRowxContext(ctx, r.bindQueryVariables(getUserQuery), userID).StructScan(&model)
	if err != nil {
		return userPkg.User{}, fmt.Errorf("error getting user %q from database: %w", userID, err)
	}

	return userPkg.User(model), nil
}
func (r UserPostgresRepository) IsUsernameAvailable(ctx context.Context, username string) (bool, error) {
	err := r.db.QueryRowxContext(ctx, r.bindQueryVariables(checkUsernameQuery), username).Scan(&username)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return false, fmt.Errorf("error checking username in database: %w", err)
		}

		return true, nil
	}

	return false, nil
}

func ConnectToPostgres(ctx context.Context, dsn string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("error connecting to postgres: %w", err)
	}

	err = db.PingContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("error pinging postgres: %w", err)
	}

	return db, nil
}
