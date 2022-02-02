package sqlite_test

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/stillwondering/xone"
	"github.com/stillwondering/xone/sqlite"
)

var examplePersons = []xone.Person{
	{
		ID:          1,
		FirstName:   "Harry",
		LastName:    "Potter",
		DateOfBirth: time.Date(1980, time.July, 31, 0, 0, 0, 0, time.UTC),
		Gender:      xone.Male,
	},
	{
		ID:          2,
		FirstName:   "Ron",
		LastName:    "Weasley",
		DateOfBirth: time.Date(1980, time.March, 1, 0, 0, 0, 0, time.UTC),
		Gender:      xone.Male,
	},
	{
		ID:          3,
		FirstName:   "Hermione",
		LastName:    "Granger",
		DateOfBirth: time.Date(1979, time.September, 19, 0, 0, 0, 0, time.UTC),
		Gender:      xone.Female,
	},
}

func Test_NewPersonService(t *testing.T) {
	db := MustOpenDB(t)
	defer MustCloseDB(t, db)

	sqlite.NewPersonService(db)
}

func Test_PersonService_FindAll(t *testing.T) {
	db := MustOpenDB(t)
	defer MustCloseDB(t, db)

	service := sqlite.NewPersonService(db)

	persons, err := service.FindAll(context.Background())
	if err != nil {
		t.Errorf("findAll() error = %v, wantErr nil", err)
	}
	if persons != nil {
		t.Errorf("findAll() persons = %v, want nil", persons)
	}

	person, err := service.Create(context.Background(), xone.CreatePersonData{
		FirstName:   "Harry",
		LastName:    "Potter",
		DateOfBirth: time.Date(1980, time.July, 31, 0, 0, 0, 0, time.UTC),
		Gender:      xone.Male,
	})
	if err != nil {
		t.Errorf("Create() error =%v, wantErr nil", err)
	}
	if person != examplePersons[0] {
		t.Errorf("Create() person = %v, want %v", person, examplePersons[0])
	}

	person, err = service.Create(context.Background(), xone.CreatePersonData{
		FirstName:   "Ron",
		LastName:    "Weasley",
		DateOfBirth: time.Date(1980, time.March, 1, 0, 0, 0, 0, time.UTC),
		Gender:      xone.Male,
	})
	if err != nil {
		t.Errorf("Create() error =%v, wantErr nil", err)
	}
	if person != examplePersons[1] {
		t.Errorf("Create() person = %v, want %v", person, examplePersons[1])
	}

	err = service.Delete(context.Background(), 2)
	if err != nil {
		t.Errorf("Delete() error =%v, wantErr nil", err)
	}

	persons, err = service.FindAll(context.Background())
	if err != nil {
		t.Errorf("findAll() error = %v, wantErr nil", err)
	}
	expected := []xone.Person{examplePersons[0]}
	if !reflect.DeepEqual(expected, persons) {
		t.Errorf("findAll() persons = %v, want %v", persons, expected)
	}

	person, found, err := service.Find(context.Background(), 2)
	if err != nil {
		t.Errorf("find() error = %v, wantErr nil", err)
	}
	if found {
		t.Errorf("find() found = %v, want false", found)
	}
	expectedPerson := xone.Person{}
	if person != expectedPerson {
		t.Errorf("find() person = %v, want %v", person, expectedPerson)
	}

	person, found, err = service.Find(context.Background(), 1)
	if err != nil {
		t.Errorf("find() error = %v, wantErr nil", err)
	}
	if !found {
		t.Errorf("find() found = %v, want true", found)
	}
	expectedPerson = examplePersons[0]
	if person != expectedPerson {
		t.Errorf("find() person = %v, want %v", person, expectedPerson)
	}
}
