package xone

import (
	"context"
)

type PersonRepository interface {
	FindAll(context.Context) ([]Person, error)
	Find(context.Context, string) (Person, bool, error)
	Create(context.Context, CreatePersonData) (Person, error)
	Delete(context.Context, string) error
}

type UserService interface {
	FindByEmail(context.Context, string) (User, bool, error)
	Create(context.Context, CreateUserData) (User, error)
}
