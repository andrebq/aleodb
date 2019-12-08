package aleodb

import (
	"database/sql"

	"github.com/kode4food/ale/data"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

type (
	db struct {
		obj     data.Object
		conn    *sql.DB
		name    string
		id      uuid.UUID
		ownerID uuid.UUID
	}
)

func (d *db) initObj() *db {
	d.obj = data.NewObject(
		data.Keyword("id"), d.id,
		data.Keyword("name"), data.String(d.name),
		data.Keyword("owner-id"), d.ownerID)
	return d
}

func (d *db) createCollection(name string) (*col, error) {
	colid := uuid.NewV5(d.id, name)
	err := d.conn.QueryRow(`
	insert into fda_dbs.user_cols(id, db_id, name)
	values ($1, $2, $3)
	returning id`, colid, d.id, name).Scan(&colid)
	if err != nil {
		return nil, err
	}

	return &col{id: colid, db: d}, nil
}

func (d *db) logentry() *logrus.Entry {
	return logrus.WithField("subsys", "couchdb-server").
		WithField("db", d.id).
		WithField("ownerID", d.ownerID)
}

// String implements data.Value
func (d *db) String() string {
	return d.obj.String()
}
