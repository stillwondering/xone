package sqlite_test

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/stillwondering/xone"
	"github.com/stillwondering/xone/sqlite"
)

func Test(t *testing.T) {
	db := MustOpenDB(t)
	defer MustCloseDB(t, db)

	personService := sqlite.NewPersonService(db)
	personService.GenerateID = func() string {
		return "id"
	}

	membershipService := sqlite.NewMembershipService(db)

	mt, err := membershipService.CreateMembershipType(context.Background(), "active")
	if err != nil {
		t.Fatalf("MembershipService.CreateMembershipType() error = %v, wantErr nil", err)
	}
	if !reflect.DeepEqual(mt, xone.MembershipType{ID: 1, Name: "active"}) {
		t.Fatalf("MembershipService.CreateMembershipType() = %v, want %v", mt, xone.MembershipType{ID: 1, Name: "active"})
	}

	p, err := personService.Create(context.Background(), xone.CreatePersonData{
		FirstName:        "Harry",
		LastName:         "Potter",
		DateOfBirth:      time.Date(1980, time.July, 31, 0, 0, 0, 0, time.UTC),
		MembershipTypeID: mt.ID,
		EffectiveFrom:    time.Date(1998, time.July, 31, 0, 0, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatalf("PersonService.Create() error = %v, wantErr nil", err)
	}

	expectedPerson := xone.Person{
		ID:          1,
		PID:         "id",
		FirstName:   "Harry",
		LastName:    "Potter",
		DateOfBirth: time.Date(1980, time.July, 31, 0, 0, 0, 0, time.UTC),
		Memberships: []xone.Membership{
			{
				ID:            1,
				Type:          mt,
				EffectiveFrom: time.Date(1998, time.July, 31, 0, 0, 0, 0, time.UTC),
			},
		},
	}
	if !reflect.DeepEqual(expectedPerson, p) {
		t.Fatalf("PersonService.Create() = %v, want %v", p, expectedPerson)
	}

	findPerson, found, err := personService.Find(context.Background(), expectedPerson.PID)
	if err != nil {
		t.Fatalf("PersonService.Find() error = %v, wantErr false", err)
	}
	if !found {
		t.Fatalf("PersonService.Find() found = %v, wantFound true", found)
	}
	if !reflect.DeepEqual(expectedPerson, findPerson) {
		t.Fatalf("PersonService.Create() = %v, want %v", findPerson, expectedPerson)
	}

	if err := personService.Delete(context.Background(), expectedPerson.PID); err != nil {
		t.Fatalf("PersonService.Delete() error = %v, wantErr false", err)
	}

	_, found, err = personService.Find(context.Background(), expectedPerson.PID)
	if err != nil {
		t.Fatalf("PersonService.Find() error = %v, wantErr false", err)
	}
	if found {
		t.Fatalf("PersonService.Find() found = %v, wantFound false", found)
	}
}

func Test_NewPersonService(t *testing.T) {
	db := MustOpenDB(t)
	defer MustCloseDB(t, db)

	sqlite.NewPersonService(db)
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
