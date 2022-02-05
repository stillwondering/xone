package sqlite

import (
	"context"
	"reflect"
	"testing"

	"github.com/stillwondering/xone"
)

func Test_findUserByEmail(t *testing.T) {
	type args struct {
		ctx      context.Context
		email    string
		testfile string
	}
	tests := []struct {
		name      string
		args      args
		want      xone.User
		wantFound bool
		wantErr   bool
	}{
		{
			name: "No users in database",
			args: args{
				ctx:   context.Background(),
				email: "albus.dumbledore@hogwarts.co.uk",
			},
			want:      xone.User{},
			wantFound: false,
			wantErr:   false,
		},
		{
			name: "User in database",
			args: args{
				ctx:      context.Background(),
				email:    "albus.dumbledore@hogwarts.co.uk",
				testfile: "testdata/Test_findUserByEmail_prefill.sql",
			},
			want: xone.User{
				Email:    "albus.dumbledore@hogwarts.co.uk",
				Password: "Harrydidyouputyournameinthegobletoffire",
			},
			wantFound: true,
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := mustOpenDB(t)
			defer mustCloseDB(t, db)
			if tt.args.testfile != "" {
				mustExecuteSQL(t, db, tt.args.testfile)
			}
			tx := mustBeginTx(t, db, context.Background())

			got, gotFound, err := findUserByEmail(tt.args.ctx, tx, tt.args.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("findUserByEmail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("findUserByEmail() got = %v, want %v", got, tt.want)
			}
			if gotFound != tt.wantFound {
				t.Errorf("findUserByEmail() gotFound = %v, want %v", gotFound, tt.wantFound)
			}
		})
	}
}

func Test_createUser(t *testing.T) {
	type args struct {
		ctx      context.Context
		data     xone.CreateUserData
		testfile string
	}
	tests := []struct {
		name    string
		args    args
		want    xone.User
		wantErr bool
	}{
		{
			name: "Empty database",
			args: args{
				ctx: context.Background(),
				data: xone.CreateUserData{
					Email:    "albus.dumbledore@hogwarts.co.uk",
					Password: "Harrydidyouputyournameinthegobletoffire",
				},
			},
			want: xone.User{
				Email:    "albus.dumbledore@hogwarts.co.uk",
				Password: "Harrydidyouputyournameinthegobletoffire",
			},
			wantErr: false,
		},
		{
			name: "User already exists",
			args: args{
				ctx: context.Background(),
				data: xone.CreateUserData{
					Email:    "albus.dumbledore@hogwarts.co.uk",
					Password: "Harrydidyouputyournameinthegobletoffire",
				},
				testfile: "testdata/Test_createUser_prefill.sql",
			},
			want:    xone.User{},
			wantErr: true,
		},
		{
			name: "Non-empty database",
			args: args{
				ctx: context.Background(),
				data: xone.CreateUserData{
					Email:    "severus.snape@hogwarts.co.uk",
					Password: "Detention!",
				},
				testfile: "testdata/Test_createUser_prefill.sql",
			},
			want: xone.User{
				Email:    "severus.snape@hogwarts.co.uk",
				Password: "Detention!",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := mustOpenDB(t)
			defer mustCloseDB(t, db)
			if tt.args.testfile != "" {
				mustExecuteSQL(t, db, tt.args.testfile)
			}
			tx := mustBeginTx(t, db, context.Background())

			got, err := createUser(tt.args.ctx, tx, tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("createUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("createUser() = %v, want %v", got, tt.want)
			}
		})
	}
}
