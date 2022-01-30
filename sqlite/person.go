package sqlite

import (
	"context"
	"fmt"
	"time"

	"github.com/stillwondering/xone"
)

type PersonService struct {
	db *DB
}

func NewPersonService(db *DB) *PersonService {
	service := PersonService{
		db: db,
	}

	return &service
}

func (ps *PersonService) FindAll(ctx context.Context) ([]xone.Person, error) {
	tx, err := ps.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	persons, err := findPersons(ctx, tx)
	if err != nil {
		return nil, err
	}

	return persons, tx.Commit()
}

func findPersons(ctx context.Context, tx *Tx) ([]xone.Person, error) {
	rows, err := tx.QueryContext(ctx, `
		SELECT
			id,
			pid,
			first_name,
			last_name,
			date_of_birth,
			gender
		FROM
			person
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var persons []xone.Person
	var id int
	var pid, firstName, lastName, dobString, genderString string
	for rows.Next() {
		if err := rows.Scan(&id, &pid, &firstName, &lastName, &dobString, &genderString); err != nil {
			return nil, err
		}

		p := xone.Person{
			PID:       xone.PersonID(pid),
			FirstName: firstName,
			LastName:  lastName,
		}

		if p.DateOfBirth, err = parseDateOfBirth(dobString); err != nil {
			return nil, err
		}

		if p.Gender, err = parseGender(genderString); err != nil {
			return nil, err
		}

		persons = append(persons, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return persons, nil
}

func parseDateOfBirth(s string) (time.Time, error) {
	return time.Parse(xone.FormatDateOfBirth, s)
}

func parseGender(s string) (xone.Gender, error) {
	genders := map[string]xone.Gender{
		"f": xone.Female,
		"m": xone.Male,
		"o": xone.Other,
	}

	g, found := genders[s]
	if !found {
		return xone.Other, fmt.Errorf("unsupported gender %s", s)
	}

	return g, nil
}
