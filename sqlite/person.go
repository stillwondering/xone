package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/stillwondering/xone"

	_ "github.com/mattn/go-sqlite3"
)

var _ xone.PersonRepository = (*PersonService)(nil)

type PersonService struct {
	db         *sql.DB
	generateID func() string
}

func NewPersonService(db *sql.DB) *PersonService {
	service := PersonService{
		db: db,
		generateID: func() string {
			return uuid.NewV4().String()
		},
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

func (ps *PersonService) Find(ctx context.Context, id string) (xone.Person, bool, error) {
	tx, err := ps.db.BeginTx(ctx, nil)
	if err != nil {
		return xone.Person{}, false, err
	}
	defer tx.Rollback()

	person, found, err := findPerson(ctx, tx, id)
	if err != nil {
		return xone.Person{}, false, err
	}

	return person, found, tx.Commit()
}

func (ps *PersonService) Create(ctx context.Context, data xone.CreatePersonData) (xone.Person, error) {
	tx, err := ps.db.BeginTx(ctx, nil)
	if err != nil {
		return xone.Person{}, err
	}
	defer tx.Rollback()

	person, err := createPerson(ctx, tx, ps.generateID(), data)
	if err != nil {
		return xone.Person{}, err
	}

	return person, tx.Commit()
}

func (ps *PersonService) Delete(ctx context.Context, id string) error {
	tx, err := ps.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := deletePerson(ctx, tx, id); err != nil {
		return err
	}

	return tx.Commit()
}

func findPersons(ctx context.Context, tx dbtx) ([]xone.Person, error) {
	rows, err := tx.QueryContext(ctx, `
		SELECT
			id,
			public_id,
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
			ID:        id,
			PID:       pid,
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

func findPerson(ctx context.Context, tx dbtx, pid string) (xone.Person, bool, error) {
	stmt, err := tx.PrepareContext(ctx, `
		SELECT
			id,
			first_name,
			last_name,
			date_of_birth,
			gender
		FROM
			person
		WHERE
			public_id = ?
	`)
	if err != nil {
		return xone.Person{}, false, err
	}

	rows, err := stmt.QueryContext(ctx, pid)
	if err != nil {
		return xone.Person{}, false, err
	}
	defer rows.Close()

	p := xone.Person{}
	found := false
	for rows.Next() {
		var id int
		var firstName, lastName, dobString, genderString string

		if err := rows.Scan(&id, &firstName, &lastName, &dobString, &genderString); err != nil {
			return xone.Person{}, false, err
		}

		p = xone.Person{
			ID:        id,
			PID:       pid,
			FirstName: firstName,
			LastName:  lastName,
		}

		if p.DateOfBirth, err = parseDateOfBirth(dobString); err != nil {
			return xone.Person{}, false, err
		}

		if p.Gender, err = parseGender(genderString); err != nil {
			return xone.Person{}, false, err
		}

		found = true
	}
	if err := rows.Err(); err != nil {
		return xone.Person{}, false, err
	}

	return p, found, nil
}

func createPerson(ctx context.Context, tx dbtx, pid string, data xone.CreatePersonData) (xone.Person, error) {
	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO person (
			public_id,
			first_name,
			last_name,
			date_of_birth,
			gender
		) VALUES (
			?,
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
		pid,
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
		PID:         pid,
		FirstName:   data.FirstName,
		LastName:    data.LastName,
		DateOfBirth: data.DateOfBirth,
		Gender:      data.Gender,
	}

	return p, nil
}

func deletePerson(ctx context.Context, tx dbtx, id string) error {
	stmt, err := tx.PrepareContext(ctx, `
		DELETE FROM
			person
		WHERE
			public_id = ?
	`)
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(ctx, id)

	return err
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
