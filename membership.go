package xone

import "time"

type Membership struct {
	ID            int
	Type          MembershipType
	EffectiveFrom time.Time
}

type MembershipType struct {
	ID   int
	Name string
}

type CreateMembershipData struct {
	PersonID         int
	MembershipTypeID int
	EffectiveFrom    time.Time
}

type UpdateMembershipData struct {
	MembershipTypeID int
	EffectiveFrom    time.Time
}

func (m Membership) ToUpdateData() UpdateMembershipData {
	return UpdateMembershipData{
		MembershipTypeID: m.Type.ID,
		EffectiveFrom:    m.EffectiveFrom,
	}
}
