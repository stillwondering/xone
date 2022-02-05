package sqlite

import (
	"context"
	"database/sql"
	"errors"

	"github.com/stillwondering/xone"
)

type UserService struct {
	db *DB
}

func NewUserService(db *DB) (*UserService, error) {
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

func findUserByEmail(ctx context.Context, tx *Tx, email string) (xone.User, bool, error) {
	stmt, err := tx.PrepareContext(ctx, `
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
