package sqlite

import (
	"context"
	"database/sql"
	"errors"

	"github.com/stillwondering/xone"
)

type UserService struct {
	db *sql.DB
}

func NewUserService(db *sql.DB) (*UserService, error) {
	service := UserService{
		db: db,
	}

	return &service, nil
}

func (us *UserService) FindByEmail(ctx context.Context, email string) (xone.User, bool, error) {
	tx, err := us.db.BeginTx(ctx, nil)
	if err != nil {
		return xone.User{}, false, err
	}
	defer tx.Rollback()

	user, found, err := findUserByEmail(ctx, tx, email)
	if err != nil {
		return xone.User{}, false, err
	}

	return user, found, tx.Commit()
}

func (us *UserService) Create(ctx context.Context, data xone.CreateUserData) (xone.User, error) {
	tx, err := us.db.BeginTx(ctx, nil)
	if err != nil {
		return xone.User{}, err
	}
	defer tx.Rollback()

	if _, found, err := findUserByEmail(ctx, tx, data.Email); err != nil {
		return xone.User{}, err
	} else if found {
		return xone.User{}, &xone.ErrUserExists{Data: data}
	}

	user, err := createUser(ctx, tx, data)
	if err != nil {
		return xone.User{}, err
	}

	return user, tx.Commit()
}

func findUserByEmail(ctx context.Context, db dbtx, email string) (xone.User, bool, error) {
	stmt, err := db.PrepareContext(ctx, `
		SELECT
			email,
			password
		FROM
			users
		WHERE
			email = ?
	`)
	if err != nil {
		return xone.User{}, false, err
	}

	user := xone.User{}
	row := stmt.QueryRowContext(ctx, email)
	if err := row.Scan(&user.Email, &user.Password); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user, false, nil
		}

		return user, false, err
	}

	return user, true, nil
}

func createUser(ctx context.Context, db dbtx, data xone.CreateUserData) (xone.User, error) {
	stmt, err := db.PrepareContext(ctx, `
		INSERT INTO users (
			email,
			password
		) VALUES (
			?,
			?
		)
	`)
	if err != nil {
		return xone.User{}, err
	}

	_, err = stmt.ExecContext(ctx, data.Email, data.Password)
	if err != nil {
		return xone.User{}, err
	}

	return xone.User(data), nil
}
