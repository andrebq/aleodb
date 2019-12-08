package aleodb

import (
	"database/sql"
	"testing"

	"github.com/kode4food/ale/core/bootstrap"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/eval"
)

func TestAleInteraction(t *testing.T) {
	manager := bootstrap.TopLevelManager()
	bootstrap.Into(manager)
	RegisterQualified(manager)

	ns := manager.GetAnonymous()

	cleanupDB(t)
	val := eval.String(ns, data.String(`
	(define server (odb/connect "localhost:5432" "fda_owner" "fda_owner"))

	(define user (odb/create-user server "little bob tables"))

	(define db (odb/create-db server user "mydb"))

	(define stuff (odb/create-collection db "stuff"))

	(define obj-to-save {:say "Hello world"})

	(odb/put stuff "my/object" obj-to-save)

	(odb/get stuff "my/object")
	`))

	_, ok := val.(data.Object)
	if !ok {
		t.Errorf("Should be a *db got %#v", val)
	}
}

func cleanupDB(t *testing.T) {
	conn, err := sql.Open("postgres", "host=localhost user=fda_owner password=fda_owner dbname=fda_db sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}
	_, err = conn.Exec("TRUNCATE TABLE fda_users.users CASCADE;")
	if err != nil {
		t.Fatal(err)
	}
}
