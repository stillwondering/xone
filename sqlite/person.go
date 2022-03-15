package sqlite

import (
	"context"
	"database/sql"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/stillwondering/xone"

	_ "github.com/mattn/go-sqlite3"
)

var _ xone.PersonRepository = (*PersonService)(nil)

type PersonService struct {
	db         *sql.DB
	GenerateID func() string
}

func NewPersonService(db *sql.DB) *PersonService {
	service := PersonService{
		db: db,
		GenerateID: func() string {
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

	person, err := createPerson(ctx, tx, ps.GenerateID(), data)
	if err != nil {
		return xone.Person{}, err
	}

	_, err = createMembership(ctx, tx, xone.CreateMembershipData{
		PersonID:         person.ID,
		MembershipTypeID: data.MembershipTypeID,
		EffectiveFrom:    data.EffectiveFrom,
	})
	if err != nil {
		return xone.Person{}, err
	}

	if err := attachMemberships(ctx, tx, &person); err != nil {
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

func (ps *PersonService) Update(ctx context.Context, id string, data xone.UpdatePersonData) error {
	tx, err := ps.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := updatePerson(ctx, tx, id, data); err != nil {
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
			email,
			phone,
			mobile,
			street,
			house_number,
			zip_code,
			city
		FROM
			person
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var persons []xone.Person
	var id int
	var pid, firstName, lastName, dobString, email, phone, mobile, street, houseNumber, zipCode, city string
	for rows.Next() {
		if err := rows.Scan(&id, &pid, &firstName, &lastName, &dobString, &email, &phone, &mobile, &street, &houseNumber, &zipCode, &city); err != nil {
			return nil, err
		}

		p := xone.Person{
			ID:          id,
			PID:         pid,
			FirstName:   firstName,
			LastName:    lastName,
			DateOfBirth: time.Time{},
			Email:       email,
			Phone:       phone,
			Mobile:      mobile,
			Street:      street,
			HouseNumber: houseNumber,
			ZipCode:     zipCode,
			City:        city,
		}

		if dobString != "" {
			if p.DateOfBirth, err = parseDateOfBirth(dobString); err != nil {
				return nil, err
			}
		}

		if err := attachMemberships(ctx, tx, &p); err != nil {
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
			email,
			phone,
			mobile,
			street,
			house_number,
			zip_code,
			city
		FROM
			person
		WHERE
			public_id = ?
	`)
	if err != nil {
		return xone.Person{}, false, err
	}

	row := stmt.QueryRowContext(ctx, pid)

	p := xone.Person{}
	var id int
	var firstName, lastName, dobString, email, phone, mobile, street, houseNumber, zipCode, city string

	if err := row.Scan(&id, &firstName, &lastName, &dobString, &email, &phone, &mobile, &street, &houseNumber, &zipCode, &city); err != nil {
		if err == sql.ErrNoRows {
			return xone.Person{}, false, nil
		}

		return xone.Person{}, false, err
	}

	p = xone.Person{
		ID:          id,
		PID:         pid,
		FirstName:   firstName,
		LastName:    lastName,
		DateOfBirth: time.Time{},
		Email:       email,
		Phone:       phone,
		Mobile:      mobile,
		Street:      street,
		HouseNumber: houseNumber,
		ZipCode:     zipCode,
		City:        city,
	}

	if dobString != "" {
		if p.DateOfBirth, err = parseDateOfBirth(dobString); err != nil {
			return xone.Person{}, true, err
		}
	}

	if err := attachMemberships(ctx, tx, &p); err != nil {
		return xone.Person{}, true, err
	}

	return p, true, nil
}

func createPerson(ctx context.Context, tx dbtx, pid string, data xone.CreatePersonData) (xone.Person, error) {
	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO person (
			public_id,
			first_name,
			last_name,
			date_of_birth,
			email,
			phone,
			mobile,
			street,
			house_number,
			zip_code,
			city
		) VALUES (
			?,
			?,
			?,
			?,
			?,
			?,
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

	dob := ""
	if !data.DateOfBirth.IsZero() {
		dob = data.DateOfBirth.Format(xone.FormatDateOfBirth)
	}

	result, err := stmt.ExecContext(
		ctx,
		pid,
		data.FirstName,
		data.LastName,
		dob,
		data.Email,
		data.Phone,
		data.Mobile,
		data.Street,
		data.HouseNumber,
		data.ZipCode,
		data.City,
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
		Email:       data.Email,
		Phone:       data.Phone,
		Mobile:      data.Mobile,
		Street:      data.Street,
		HouseNumber: data.HouseNumber,
		ZipCode:     data.ZipCode,
		City:        data.City,
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

func updatePerson(ctx context.Context, tx dbtx, id string, upd xone.UpdatePersonData) error {
	stmt, err := tx.PrepareContext(ctx, `
		UPDATE
			person
		SET
			first_name = ?,
			last_name = ?,
			date_of_birth = ?,
			email = ?,
			phone = ?,
			mobile = ?,
			street = ?,
			house_number = ?,
			zip_code = ?,
			city = ?
		WHERE
			public_id = ?
	`)
	if err != nil {
		return err
	}

	dob := ""
	if !upd.DateOfBirth.IsZero() {
		dob = upd.DateOfBirth.Format(xone.FormatDateOfBirth)
	}

	_, err = stmt.ExecContext(ctx, upd.FirstName, upd.LastName, dob, upd.Email, upd.Phone, upd.Mobile, upd.Street, upd.HouseNumber, upd.ZipCode, upd.City, id)

	return err
}

func attachMemberships(ctx context.Context, tx dbtx, p *xone.Person) error {
	memberships, err := findMembershipsByPerson(ctx, tx, p.PID)
	if err != nil {
		return err
	}
	p.Memberships = memberships

	return nil
}

func parseDateOfBirth(s string) (time.Time, error) {
	return time.Parse(xone.FormatDateOfBirth, s)
}
