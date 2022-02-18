package sqlite_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stillwondering/xone"
	"github.com/stillwondering/xone/sqlite"
)

var examplePersons = []xone.Person{
	{
		ID:          1,
		PID:         "1",
		FirstName:   "Harry",
		LastName:    "Potter",
		DateOfBirth: time.Date(1980, time.July, 31, 0, 0, 0, 0, time.UTC),
	},
	{
		ID:          2,
		PID:         "2",
		FirstName:   "Ron",
		LastName:    "Weasley",
		DateOfBirth: time.Date(1980, time.March, 1, 0, 0, 0, 0, time.UTC),
	},
	{
		ID:          3,
		PID:         "3",
		FirstName:   "Hermione",
		LastName:    "Granger",
		DateOfBirth: time.Date(1979, time.September, 19, 0, 0, 0, 0, time.UTC),
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
	})
	if err != nil {
		t.Errorf("Create() error =%v, wantErr nil", err)
	}
	if person.ID != examplePersons[0].ID {
		t.Errorf("find() person.ID = %v, want %v", person.ID, examplePersons[0].ID)
	}

	pid1 := person.PID

	person, err = service.Create(context.Background(), xone.CreatePersonData{
		FirstName:   "Ron",
		LastName:    "Weasley",
		DateOfBirth: time.Date(1980, time.March, 1, 0, 0, 0, 0, time.UTC),
	})
	if err != nil {
		t.Errorf("Create() error =%v, wantErr nil", err)
	}
	if person.ID != examplePersons[1].ID {
		t.Errorf("find() person.ID = %v, want %v", person.ID, examplePersons[1].ID)
	}

	pid2 := person.PID

	err = service.Delete(context.Background(), pid2)
	if err != nil {
		t.Errorf("Delete() error =%v, wantErr nil", err)
	}

	persons, err = service.FindAll(context.Background())
	if err != nil {
		t.Errorf("findAll() error = %v, wantErr nil", err)
	}
	expectedCount := 1
	if expectedCount != len(persons) {
		t.Errorf("findAll() want slice of size %d, got %v", expectedCount, persons)
	}

	person, found, err := service.Find(context.Background(), pid2)
	if err != nil {
		t.Errorf("find() error = %v, wantErr nil", err)
	}
	if found {
		t.Errorf("find() found = %v, want false", found)
	}
	if person.ID != 0 {
		t.Errorf("find() person.ID = %v, want %v", person.ID, 0)
	}

	person, found, err = service.Find(context.Background(), pid1)
	if err != nil {
		t.Errorf("find() error = %v, wantErr nil", err)
	}
	if !found {
		t.Errorf("find() found = %v, want true", found)
	}
	if person.ID != examplePersons[0].ID {
		t.Errorf("find() person.ID = %v, want %v", person.ID, examplePersons[0].ID)
	}
}

func Test_NewUserService(t *testing.T) {
	db := MustOpenDB(t)
	defer MustCloseDB(t, db)

	_, err := sqlite.NewUserService(db)
	if err != nil {
		t.Fatal(err)
	}
}

func TestUserService(t *testing.T) {
	db := MustOpenDB(t)
	defer MustCloseDB(t, db)

	userService, err := sqlite.NewUserService(db)
	if err != nil {
		t.Fatal(err)
	}

	_, found, err := userService.FindByEmail(context.Background(), "albus.dumbledore@hogwarts.co.uk")
	if err != nil {
		t.Fatal(err)
	}
	if found != false {
		t.Errorf("UserService.FindByEmail() gotFound = %v, want %v", found, false)
	}

	user, err := userService.Create(context.Background(), xone.CreateUserData{
		Email:    "albus.dumbledore@hogwarts.co.uk",
		Password: "Harrydidyouputyournameinthegobletoffire",
	})
	if err != nil {
		t.Errorf("UserService.Create() err = %v, want %v", err, nil)
	}
	wantUser := xone.User{Email: "albus.dumbledore@hogwarts.co.uk", Password: "Harrydidyouputyournameinthegobletoffire"}
	if user != wantUser {
		t.Errorf("UserService.Create() want = %v, got %v", wantUser, user)
	}

	_, err = userService.Create(context.Background(), xone.CreateUserData{
		Email:    "albus.dumbledore@hogwarts.co.uk",
		Password: "A different password",
	})
	var e *xone.ErrUserExists
	if !errors.As(err, &e) {
		t.Errorf("UserService.Create() wantErr = %v, got %v", e, err)
	}
}
