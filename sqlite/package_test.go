package sqlite_test

import (
	"context"
	"testing"

	"github.com/stillwondering/xone/sqlite"
)

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
}
