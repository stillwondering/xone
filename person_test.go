package xone

import (
	"reflect"
	"testing"
	"time"
)

func TestPerson_Membership(t *testing.T) {
	type fields struct {
		ID          int
		PID         string
		FirstName   string
		LastName    string
		DateOfBirth time.Time
		Email       string
		Phone       string
		Mobile      string
		Memberships []Membership
	}
	type args struct {
		today time.Time
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Membership
	}{
		{
			name:   "No memberships",
			fields: fields{},
			args: args{
				today: time.Date(2022, time.February, 23, 12, 0, 0, 0, time.UTC),
			},
			want: nil,
		},
		{
			name: "One membership without effective from date",
			fields: fields{
				Memberships: []Membership{
					{
						ID: 1,
						Type: MembershipType{
							ID:   1,
							Name: "active",
						},
						EffectiveFrom: time.Time{},
					},
				},
			},
			args: args{
				today: time.Date(2022, time.February, 23, 12, 0, 0, 0, time.UTC),
			},
			want: &Membership{
				ID: 1,
				Type: MembershipType{
					ID:   1,
					Name: "active",
				},
				EffectiveFrom: time.Time{},
			},
		},
		{
			name: "Two memberships without effective from date",
			fields: fields{
				Memberships: []Membership{
					{
						ID: 1,
						Type: MembershipType{
							ID:   1,
							Name: "active",
						},
						EffectiveFrom: time.Time{},
					},
					{
						ID: 2,
						Type: MembershipType{
							ID:   2,
							Name: "passive",
						},
						EffectiveFrom: time.Time{},
					},
				},
			},
			args: args{
				today: time.Date(2022, time.February, 23, 12, 0, 0, 0, time.UTC),
			},
			want: &Membership{
				ID: 2,
				Type: MembershipType{
					ID:   2,
					Name: "passive",
				},
				EffectiveFrom: time.Time{},
			},
		},
		{
			name: "One membership without effective from, one in the future",
			fields: fields{
				Memberships: []Membership{
					{
						ID: 1,
						Type: MembershipType{
							ID:   1,
							Name: "active",
						},
						EffectiveFrom: time.Time{},
					},
					{
						ID: 2,
						Type: MembershipType{
							ID:   2,
							Name: "passive",
						},
						EffectiveFrom: time.Date(2022, time.April, 1, 0, 0, 0, 0, time.UTC),
					},
				},
			},
			args: args{
				today: time.Date(2022, time.February, 23, 12, 0, 0, 0, time.UTC),
			},
			want: &Membership{
				ID: 1,
				Type: MembershipType{
					ID:   1,
					Name: "active",
				},
				EffectiveFrom: time.Time{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Person{
				ID:          tt.fields.ID,
				PID:         tt.fields.PID,
				FirstName:   tt.fields.FirstName,
				LastName:    tt.fields.LastName,
				DateOfBirth: tt.fields.DateOfBirth,
				Email:       tt.fields.Email,
				Phone:       tt.fields.Phone,
				Mobile:      tt.fields.Mobile,
				Memberships: tt.fields.Memberships,
			}
			if got := p.Membership(tt.args.today); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Person.Membership() = %v, want %v", got, tt.want)
			}
		})
	}
}
