package xone

import (
	"time"
)

// Person contains the personal data of a organization member.
type Person struct {
	ID          int
	PID         string
	FirstName   string
	LastName    string
	DateOfBirth time.Time
	Email       string
	Phone       string
	Mobile      string
	Street      string
	HouseNumber string
	ZipCode     string
	City        string
	Memberships []Membership
}

func (p Person) CurrentAge() int {
	return p.Age(time.Now())
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

func (p Person) CurrentMembership() *Membership {
	return p.Membership(time.Now())
}

func (p Person) Membership(today time.Time) *Membership {
	if len(p.Memberships) == 0 {
		return nil
	}

	m := p.Memberships[0]

	for i := 1; i < len(p.Memberships); i++ {
		if p.Memberships[i].EffectiveFrom.IsZero() {
			m = p.Memberships[i]
			continue
		}

		if p.Memberships[i].EffectiveFrom.After(today) {
			continue
		}

		if p.Memberships[i].EffectiveFrom.After(m.EffectiveFrom) {
			m = p.Memberships[i]
		}
	}

	return &m
}

// CreatePersonData contains all data which is necessary to create a new Person entry
// in any kind of repository.
type CreatePersonData struct {
	FirstName        string
	LastName         string
	DateOfBirth      time.Time
	Email            string
	Phone            string
	Mobile           string
	Street           string
	HouseNumber      string
	ZipCode          string
	City             string
	MembershipTypeID int
	EffectiveFrom    time.Time
}

// UpdatePersonData contains a person's data points which can be updated.
type UpdatePersonData struct {
	FirstName   string
	LastName    string
	DateOfBirth time.Time
	Email       string
	Phone       string
	Mobile      string
	Street      string
	HouseNumber string
	ZipCode     string
	City        string
}

// ToUpdateData returns a struct that can be used as a starting point when
// updating a person's data.
func (p Person) ToUpdateData() UpdatePersonData {
	return UpdatePersonData{
		FirstName:   p.FirstName,
		LastName:    p.LastName,
		DateOfBirth: p.DateOfBirth,
		Email:       p.Email,
		Phone:       p.Phone,
		Mobile:      p.Mobile,
		Street:      p.Street,
		HouseNumber: p.HouseNumber,
		ZipCode:     p.ZipCode,
		City:        p.City,
	}
}
