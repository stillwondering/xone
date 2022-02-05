package xone

import (
	"fmt"
	"time"
)

const (
	FormatDateOfBirth = "2006-01-02"
)

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

// Person contains the personal data of a organization member.
type Person struct {
	ID          int
	FirstName   string
	LastName    string
	DateOfBirth time.Time
	Gender      Gender
}

// CreatePersonData contains all data which is necessary to create a new Person entry
// in any kind of repository.
type CreatePersonData struct {
	FirstName   string
	LastName    string
	DateOfBirth time.Time
	Gender      Gender
}

type User struct {
	Email    string
	Password string
}
