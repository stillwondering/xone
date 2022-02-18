package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/stillwondering/xone"
)

type MembershipService struct {
	db *sql.DB
}

func NewMembershipService(db *sql.DB) *MembershipService {
	service := MembershipService{
		db: db,
	}

	return &service
}

func (s *MembershipService) FindAllMembershipTypes(ctx context.Context) ([]xone.MembershipType, error) {
	return findAllMembershipTypes(ctx, s.db)
}

func findAllMembershipTypes(ctx context.Context, db dbtx) ([]xone.MembershipType, error) {
	rows, err := db.QueryContext(ctx, `
		SELECT
			id,
			name
		FROM
			membership_type
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var membershipTypes []xone.MembershipType
	for rows.Next() {
		var mt xone.MembershipType

		if err := rows.Scan(&mt.ID, &mt.Name); err != nil {
			return nil, err
		}

		membershipTypes = append(membershipTypes, mt)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return membershipTypes, nil
}

func findMembershipsByPerson(ctx context.Context, db dbtx, pid string) ([]xone.Membership, error) {
	stmt, err := db.PrepareContext(ctx, `
		SELECT
			membership.id,
			membership.effective_from,
			membership_type.id,
			membership_type.name
		FROM
			membership
			JOIN membership_type ON membership.type_id = membership_type.id
			JOIN person ON membership.person_id = person.id
		WHERE
			person.public_id = ?
		ORDER BY
			membership.id
	`)
	if err != nil {
		return nil, err
	}

	rows, err := stmt.QueryContext(ctx, pid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var memberships []xone.Membership

	for rows.Next() {
		membershipType := xone.MembershipType{}
		membership := xone.Membership{}
		var effectiveFromText string

		if err := rows.Scan(&membership.ID, &effectiveFromText, &membershipType.ID, &membershipType.Name); err != nil {
			return nil, err
		}

		membership.EffectiveFrom, err = time.Parse("2006-01-02", effectiveFromText)
		if err != nil {
			return nil, err
		}

		membership.Type = membershipType

		memberships = append(memberships, membership)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return memberships, nil
}

func findMembership(ctx context.Context, db dbtx, id int) (xone.Membership, bool, error) {
	stmt, err := db.PrepareContext(ctx, `
		SELECT
			membership.id,
			membership.effective_from,
			membership_type.id,
			membership_type.name
		FROM
			membership
			JOIN membership_type ON membership.type_id = membership_type.id
		WHERE
			membership.id = ?
	`)
	if err != nil {
		return xone.Membership{}, false, err
	}

	row := stmt.QueryRowContext(ctx, id)

	membership := xone.Membership{}
	var effectiveFromText string

	if err := row.Scan(&membership.ID, &effectiveFromText, &membership.Type.ID, &membership.Type.Name); err != nil {
		if err == sql.ErrNoRows {
			return xone.Membership{}, false, nil
		}

		return xone.Membership{}, false, err
	}

	membership.EffectiveFrom, err = time.Parse("2006-01-02", effectiveFromText)
	if err != nil {
		return xone.Membership{}, true, err
	}

	return membership, true, nil
}

func createMembership(ctx context.Context, db dbtx, data xone.CreateMembershipData) (xone.Membership, error) {
	stmt, err := db.PrepareContext(ctx, `
		INSERT INTO membership (
			person_id,
			type_id,
			effective_from
		) VALUES (
			?,
			?,
			?
		)
	`)
	if err != nil {
		return xone.Membership{}, err
	}

	res, err := stmt.ExecContext(ctx, data.PersonID, data.MembershipTypeID, data.EffectiveFrom.Format("2006-01-02"))
	if err != nil {
		return xone.Membership{}, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return xone.Membership{}, err
	}

	membership, found, err := findMembership(ctx, db, int(id))
	if err != nil || !found {
		return xone.Membership{}, errors.New("cannot find new membership")
	}

	return membership, nil
}
