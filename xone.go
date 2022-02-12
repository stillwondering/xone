package xone

import (
	"context"
	"fmt"
)

const (
	FormatDateOfBirth = "2006-01-02"
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

//go:generate stringer -type=Gender
type Gender int

const (
	Other Gender = iota
	Female
	Male
)

func ParseGender(s string) (Gender, error) {
	genderMap := map[string]Gender{
		Other.String():  Other,
		Female.String(): Female,
		Male.String():   Male,
	}

	gender, ok := genderMap[s]
	if !ok {
		return Other, fmt.Errorf("%s is not a valid gender", s)
	}

	return gender, nil
}
