package xone

import "time"

const (
	FormatDateOfBirth = "2006-01-02"
)

// Person contains the personal data of a organization member.
type Person struct {
	FirstName   string
	LastName    string
	DateOfBirth time.Time
}
