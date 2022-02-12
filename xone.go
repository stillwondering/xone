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

// Age calculates a person's age based on their date of birth and with respect
// to the given date.
func (p Person) Age(today time.Time) int {
	if !p.HasDateOfBirth() {
		return 0
	}

	today = today.In(p.DateOfBirth.Location())
	ty, tm, td := today.Date()
	today = time.Date(ty, tm, td, 0, 0, 0, 0, time.UTC)

	by, bm, bd := p.DateOfBirth.Date()
	dob := time.Date(by, bm, bd, 0, 0, 0, 0, time.UTC)

	if today.Before(dob) {
		return 0
	}

	age := ty - by
	anniversary := dob.AddDate(age, 0, 0)
	if anniversary.After(today) {
		age--
	}

	return age
}

func (p Person) HasDateOfBirth() bool {
	return !p.DateOfBirth.IsZero()
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

type CreateUserData struct {
	Email    string
	Password string
}
