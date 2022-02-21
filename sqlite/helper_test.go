package sqlite

import (
	"database/sql"
	"io/ioutil"
	"testing"
)

// mustOpenDB is a helper function that creates a temporary file and opens a
// sqlite database. Please note that this function takes care of adding cleanup
// functions (closing database, removing file) to the testcase as well.
func mustOpenDB(t *testing.T) *sql.DB {
	t.Helper()

	db, err := Open("file::memory:?cache=shared")
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		db.Close()
	})

	return db
}

func mustMigrateFile(t *testing.T, db *sql.DB, file string) {
	t.Helper()

	migration, err := ioutil.ReadFile(file)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := db.Exec(string(migration)); err != nil {
		t.Fatal(err)
	}
}
