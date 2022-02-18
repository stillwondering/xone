package xone

import (
	"fmt"
)

type ErrUserExists struct {
	Data CreateUserData
}

func (e *ErrUserExists) Error() string {
	return fmt.Sprintf(`user with email "%s" already exists`, e.Data.Email)
}
