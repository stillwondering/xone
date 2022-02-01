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

func (ps *PersonService) Create(ctx context.Context, data xone.CreatePersonData) (xone.Person, error) {
	tx, err := ps.db.BeginTx(ctx, nil)
	if err != nil {
		return xone.Person{}, err
	}
	defer tx.Rollback()

	person, err := createPerson(ctx, tx, data)
	if err != nil {
		return xone.Person{}, err
	}

	return person, tx.Commit()
}

func findPersons(ctx context.Context, tx *Tx) ([]xone.Person, error) {
	rows, err := tx.QueryContext(ctx, `
		SELECT
			id,
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
	var firstName, lastName, dobString, genderString string
	for rows.Next() {
		if err := rows.Scan(&id, &firstName, &lastName, &dobString, &genderString); err != nil {
			return nil, err
		}

		p := xone.Person{
			ID:        id,
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

func createPerson(ctx context.Context, tx *Tx, data xone.CreatePersonData) (xone.Person, error) {
	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO person (
			first_name,
			last_name,
			date_of_birth,
			gender
		) VALUES (
			?,
			?,
			?,
			?
		)
	`)
	if err != nil {
		return xone.Person{}, err
	}

	result, err := stmt.ExecContext(
		ctx,
		data.FirstName,
		data.LastName,
		data.DateOfBirth.Format(xone.FormatDateOfBirth),
		formatGender(data.Gender),
	)
	if err != nil {
		return xone.Person{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return xone.Person{}, err
	}

	p := xone.Person{
		ID:          int(id),
		FirstName:   data.FirstName,
		LastName:    data.LastName,
		DateOfBirth: data.DateOfBirth,
		Gender:      data.Gender,
	}

	return p, nil
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

func formatGender(g xone.Gender) string {
	genders := map[xone.Gender]string{
		xone.Female: "f",
		xone.Male:   "m",
		xone.Other:  "o",
	}

	return genders[g]
}
