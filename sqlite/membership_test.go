package sqlite

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/stillwondering/xone"
)

func Test_findMembershipsByPerson(t *testing.T) {
	type args struct {
		ctx context.Context
		pid string
	}
	tests := []struct {
		name    string
		args    args
		want    []xone.Membership
		wantErr bool
	}{
		{
			name: "Nonexistant person ID",
			args: args{
				ctx: context.Background(),
				pid: "nonexistant",
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "Multiple memberships",
			args: args{
				ctx: context.Background(),
				pid: "1",
			},
			want: []xone.Membership{
				{
					ID: 1,
					Type: xone.MembershipType{
						ID:   1,
						Name: "active",
					},
					EffectiveFrom: time.Date(1998, time.July, 31, 0, 0, 0, 0, time.UTC),
				},
				{
					ID: 2,
					Type: xone.MembershipType{
						ID:   2,
						Name: "passive",
					},
					EffectiveFrom: time.Date(2045, time.July, 31, 0, 0, 0, 0, time.UTC),
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := mustOpenDB(t)
			mustMigrateFile(t, db, "testdata/Test_findMembershipsByPerson.sql")

			got, err := findMembershipsByPerson(tt.args.ctx, db, tt.args.pid)
			if (err != nil) != tt.wantErr {
				t.Errorf("findMembershipsByPerson() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("findMembershipsByPerson() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_findAllMembershipTypes(t *testing.T) {
	tests := []struct {
		name    string
		want    []xone.MembershipType
		wantErr bool
	}{
		{
			name: "Two membership types",
			want: []xone.MembershipType{
				{
					ID:   1,
					Name: "active",
				},
				{
					ID:   2,
					Name: "passive",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := mustOpenDB(t)
			mustMigrateFile(t, db, "testdata/Test_findMembershipsByPerson.sql")

			got, err := findAllMembershipTypes(context.Background(), db)
			if (err != nil) != tt.wantErr {
				t.Errorf("findAllMembershipTypes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("findAllMembershipTypes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_findMembership(t *testing.T) {
	type args struct {
		id int
	}
	tests := []struct {
		name      string
		args      args
		want      xone.Membership
		wantFound bool
		wantErr   bool
	}{
		{
			name: "Nonexistant membership",
			args: args{
				id: 123,
			},
			want:      xone.Membership{},
			wantFound: false,
			wantErr:   false,
		},
		{
			name: "Existing membership",
			args: args{
				id: 1,
			},
			want: xone.Membership{
				ID: 1,
				Type: xone.MembershipType{
					ID:   1,
					Name: "active",
				},
				EffectiveFrom: time.Date(1998, time.July, 31, 0, 0, 0, 0, time.UTC),
			},
			wantFound: true,
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := mustOpenDB(t)
			mustMigrateFile(t, db, "testdata/Test_findMembership.sql")

			got, found, err := findMembership(context.Background(), db, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("findMembership() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("findMembership() got = %v, want %v", got, tt.want)
			}
			if found != tt.wantFound {
				t.Errorf("findMembership() found = %v, want %v", found, tt.wantFound)
			}
		})
	}
}

func Test_createMembership(t *testing.T) {
	type args struct {
		data xone.CreateMembershipData
	}
	tests := []struct {
		name    string
		args    args
		want    xone.Membership
		wantErr bool
	}{
		{
			name: "Invalid membership type",
			args: args{
				data: xone.CreateMembershipData{
					PersonID:         3,
					MembershipTypeID: 3,
					EffectiveFrom:    time.Date(1997, time.September, 19, 0, 0, 0, 0, time.UTC),
				},
			},
			want:    xone.Membership{},
			wantErr: true,
		},
		{
			name: "Nonexistant person",
			args: args{
				data: xone.CreateMembershipData{
					PersonID:         123,
					MembershipTypeID: 1,
					EffectiveFrom:    time.Date(1997, time.September, 19, 0, 0, 0, 0, time.UTC),
				},
			},
			want:    xone.Membership{},
			wantErr: true,
		},
		{
			name: "Valid data",
			args: args{
				data: xone.CreateMembershipData{
					PersonID:         3,
					MembershipTypeID: 1,
					EffectiveFrom:    time.Date(1997, time.September, 19, 0, 0, 0, 0, time.UTC),
				},
			},
			want: xone.Membership{
				ID: 3,
				Type: xone.MembershipType{
					ID:   1,
					Name: "active",
				},
				EffectiveFrom: time.Date(1997, time.September, 19, 0, 0, 0, 0, time.UTC),
			},
			wantErr: false,
		},
		{
			name: "Creating additional membership",
			args: args{
				data: xone.CreateMembershipData{
					PersonID:         1,
					MembershipTypeID: 2,
					EffectiveFrom:    time.Date(2045, time.July, 31, 0, 0, 0, 0, time.UTC),
				},
			},
			want: xone.Membership{
				ID: 3,
				Type: xone.MembershipType{
					ID:   2,
					Name: "passive",
				},
				EffectiveFrom: time.Date(2045, time.July, 31, 0, 0, 0, 0, time.UTC),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := mustOpenDB(t)
			mustMigrateFile(t, db, "testdata/Test_createMembership.sql")

			got, err := createMembership(context.Background(), db, tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("createMembership() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("createMembership() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_createMembershipType(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		args    args
		want    xone.MembershipType
		wantErr bool
	}{
		{
			name: "Already existing membership type",
			args: args{
				name: "active",
			},
			want:    xone.MembershipType{},
			wantErr: true,
		},
		{
			name: "New membership type",
			args: args{
				name: "passive",
			},
			want: xone.MembershipType{
				ID:   2,
				Name: "passive",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := mustOpenDB(t)
			mustMigrateFile(t, db, "testdata/Test_createMembershipType.sql")

			got, err := createMembershipType(context.Background(), db, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("createMembershipType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("createMembershipType() = %v, want %v", got, tt.want)
			}
		})
	}
}
