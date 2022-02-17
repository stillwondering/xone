package sqlite

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/stillwondering/xone"
)

var persons = []xone.Person{
	{
		ID:          1,
		PID:         "1",
		FirstName:   "Harry",
		LastName:    "Potter",
		DateOfBirth: time.Date(1980, time.July, 31, 0, 0, 0, 0, time.UTC),
		Email:       "harry.potter@hogwarts.co.uk",
	},
	{
		ID:          2,
		PID:         "2",
		FirstName:   "Ron",
		LastName:    "Weasley",
		DateOfBirth: time.Time{},
		Email:       "ron.weasley@hogwarts.co.uk",
	},
	{
		ID:          3,
		PID:         "3",
		FirstName:   "Hermione",
		LastName:    "Granger",
		DateOfBirth: time.Date(1979, time.September, 19, 0, 0, 0, 0, time.UTC),
		Email:       "hermione.granger@hogwarts.co.uk",
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
				testfile: "testdata/people.sql",
			},
			want:    persons,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := mustOpenDB(t)

			if tt.args.testfile != "" {
				mustMigrateFile(t, db, tt.args.testfile)
			}

			got, err := findPersons(tt.args.ctx, db)
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

func Test_findPerson(t *testing.T) {
	type args struct {
		ctx      context.Context
		testfile string
		id       string
	}
	tests := []struct {
		name    string
		args    args
		want    xone.Person
		found   bool
		wantErr bool
	}{
		{
			name: "Empty database",
			args: args{
				ctx: context.Background(),
				id:  "5",
			},
			want:    xone.Person{},
			found:   false,
			wantErr: false,
		},
		{
			name: "Filled database",
			args: args{
				ctx:      context.Background(),
				testfile: "testdata/people.sql",
				id:       "1",
			},
			want:    persons[0],
			found:   true,
			wantErr: false,
		},
		{
			name: "Empty date of birth",
			args: args{
				ctx:      context.Background(),
				testfile: "testdata/people.sql",
				id:       "2",
			},
			want:    persons[1],
			found:   true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := mustOpenDB(t)

			if tt.args.testfile != "" {
				mustMigrateFile(t, db, tt.args.testfile)
			}

			got, found, err := findPerson(tt.args.ctx, db, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("findPerson() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if found != tt.found {
				t.Errorf("findPerson() found = %v, want %v", found, tt.found)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("findPerson() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_createPerson(t *testing.T) {
	type args struct {
		ctx      context.Context
		testfile string
		id       string
		data     xone.CreatePersonData
	}
	tests := []struct {
		name    string
		args    args
		want    xone.Person
		wantErr bool
	}{
		{
			name: "Empty database",
			args: args{
				ctx: context.Background(),
				id:  "1",
				data: xone.CreatePersonData{
					FirstName:   "Harry",
					LastName:    "Potter",
					DateOfBirth: time.Date(1980, time.July, 31, 0, 0, 0, 0, time.UTC),
					Email:       "harry.potter@hogwarts.co.uk",
				},
			},
			want:    persons[0],
			wantErr: false,
		},
		{
			name: "Filled database",
			args: args{
				ctx:      context.Background(),
				testfile: "testdata/Test_createPerson_prefill.sql",
				id:       "3",
				data: xone.CreatePersonData{
					FirstName:   "Hermione",
					LastName:    "Granger",
					DateOfBirth: time.Date(1979, time.September, 19, 0, 0, 0, 0, time.UTC),
					Email:       "hermione.granger@hogwarts.co.uk",
				},
			},
			want:    persons[2],
			wantErr: false,
		},
		{
			name: "Empty date of birth",
			args: args{
				ctx:      context.Background(),
				testfile: "testdata/Test_createPerson_prefill.sql",
				id:       "3",
				data: xone.CreatePersonData{
					FirstName: "Hermione",
					LastName:  "Granger",
				},
			},
			want: xone.Person{
				ID:          3,
				PID:         "3",
				FirstName:   "Hermione",
				LastName:    "Granger",
				DateOfBirth: time.Time{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := mustOpenDB(t)

			if tt.args.testfile != "" {
				mustMigrateFile(t, db, tt.args.testfile)
			}

			got, err := createPerson(tt.args.ctx, db, tt.args.id, tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("createPerson() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("createPerson() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_deletePerson(t *testing.T) {
	type args struct {
		ctx      context.Context
		testfile string
		id       string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Empty database",
			args: args{
				ctx: context.Background(),
				id:  "1",
			},
			wantErr: false,
		},
		{
			name: "Filled database",
			args: args{
				ctx:      context.Background(),
				testfile: "testdata/Test_deletePerson_prefill.sql",
				id:       "1",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := mustOpenDB(t)

			if tt.args.testfile != "" {
				mustMigrateFile(t, db, tt.args.testfile)
			}

			err := deletePerson(tt.args.ctx, db, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("createPerson() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_updatePerson(t *testing.T) {
	type args struct {
		id  string
		upd xone.UpdatePersonData
	}
	tests := []struct {
		name       string
		args       args
		wantErr    bool
		wantPerson xone.Person
	}{
		{
			name: "Update existing person",
			args: args{
				id: "2",
				upd: xone.UpdatePersonData{
					FirstName:   "Ronald",
					LastName:    "Weasley",
					DateOfBirth: time.Date(1980, time.March, 1, 0, 0, 0, 0, time.UTC),
					Phone:       "1234",
				},
			},
			wantErr: false,
			wantPerson: xone.Person{
				ID:          2,
				PID:         "2",
				FirstName:   "Ronald",
				LastName:    "Weasley",
				DateOfBirth: time.Date(1980, time.March, 1, 0, 0, 0, 0, time.UTC),
				Email:       "",
				Phone:       "1234",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := mustOpenDB(t)
			mustMigrateFile(t, db, "testdata/people.sql")

			if err := updatePerson(context.Background(), db, tt.args.id, tt.args.upd); (err != nil) != tt.wantErr {
				t.Errorf("updatePerson() error = %v, wantErr %v", err, tt.wantErr)
			}

			p, _, err := findPerson(context.Background(), db, tt.args.id)
			if err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(p, tt.wantPerson) {
				t.Errorf("updatePerson() = %v, want %v", p, tt.wantPerson)
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
