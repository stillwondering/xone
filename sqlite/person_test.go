package sqlite

import (
	"context"
	"io/ioutil"
	"reflect"
	"testing"
	"time"

	"github.com/stillwondering/xone"
)

var persons = []xone.Person{
	{
		PID:         "1",
		FirstName:   "Harry",
		LastName:    "Potter",
		DateOfBirth: time.Date(1980, time.July, 31, 0, 0, 0, 0, time.UTC),
		Gender:      xone.Male,
	},
	{
		PID:         "2",
		FirstName:   "Ron",
		LastName:    "Weasley",
		DateOfBirth: time.Date(1980, time.March, 1, 0, 0, 0, 0, time.UTC),
		Gender:      xone.Male,
	},
	{
		PID:         "3",
		FirstName:   "Hermione",
		LastName:    "Granger",
		DateOfBirth: time.Date(1979, time.September, 19, 0, 0, 0, 0, time.UTC),
		Gender:      xone.Female,
	},
}

func Test_findPersons(t *testing.T) {
	type args struct {
		ctx      context.Context
		testfile string
	}
	tests := []struct {
		name    string
		args    args
		want    []xone.Person
		wantErr bool
	}{
		{
			name: "Empty database",
			args: args{
				ctx: context.Background(),
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "Filled database",
			args: args{
				ctx:      context.Background(),
				testfile: "testdata/Test_findPersons_multiple-people.sql",
			},
			want:    persons,
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
			tx := mustBeginTx(t, db, tt.args.ctx)

			got, err := findPersons(tt.args.ctx, tx)
			if (err != nil) != tt.wantErr {
				t.Errorf("findPersons() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("findPersons() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseDateOfBirth(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    time.Time
		wantErr bool
	}{
		{
			name:    "Invalid",
			args:    args{s: "2021-04-05"},
			want:    time.Date(2021, time.April, 5, 0, 0, 0, 0, time.UTC),
			wantErr: false,
		},
		{
			name:    "Valid",
			args:    args{s: "5.4.2021"},
			want:    time.Time{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseDateOfBirth(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseDatOfBirth() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseDatOfBirth() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseGender(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    xone.Gender
		wantErr bool
	}{
		{
			name:    "Male",
			args:    args{s: "m"},
			want:    xone.Male,
			wantErr: false,
		},
		{
			name:    "Female",
			args:    args{s: "f"},
			want:    xone.Female,
			wantErr: false,
		},
		{
			name:    "Other",
			args:    args{s: "o"},
			want:    xone.Other,
			wantErr: false,
		},
		{
			name:    "Unsupported",
			args:    args{s: "bla"},
			want:    xone.Other,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseGender(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseGender() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseGender() = %v, want %v", got, tt.want)
			}
		})
	}
}

// MustOpenDB returns a new, open DB. Fatal on error.
func mustOpenDB(tb testing.TB) *DB {
	tb.Helper()

	db := NewDB(":memory:")
	if err := db.Open(); err != nil {
		tb.Fatal(err)
	}
	return db
}

// MustCloseDB closes the DB. Fatal on error.
func mustCloseDB(tb testing.TB, db *DB) {
	tb.Helper()
	if err := db.Close(); err != nil {
		tb.Fatal(err)
	}
}

func mustExecuteSQL(tb testing.TB, db *DB, file string) {
	tb.Helper()
	sql, err := ioutil.ReadFile(file)
	if err != nil {
		tb.Fatal(err)
	}

	_, err = db.db.Exec(string(sql))
	if err != nil {
		tb.Fatal(err)
	}
}

func mustBeginTx(tb testing.TB, db *DB, ctx context.Context) *Tx {
	tb.Helper()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		tb.Fatal(err)
	}

	return tx
}
