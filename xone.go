package xone

import "time"

// Person contains the personal data of a organization member.
type Person struct {
	FirstName   string
	LastName    string
	DateOfBirth time.Time
}
