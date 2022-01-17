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
	FirstName   string
	LastName    string
	DateOfBirth time.Time
	Gender      Gender
}
