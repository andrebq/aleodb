package aleodb

import (
	"github.com/andrebq/alejson"
	"github.com/kode4food/ale/data"
	uuid "github.com/satori/go.uuid"
)

type (
	col struct {
		obj data.Object
		db  *db
		id  uuid.UUID
	}
)

func (c *col) String() string {
	return c.obj.String()
}

func (c *col) initObj() {
	c.obj = data.NewObject(
		data.Keyword("id"), c.id,
		data.Keyword("database"), c.db.obj)
}

func (c *col) putObject(key string, content data.Object) (uuid.UUID, error) {
	oid := uuid.NewV5(c.id, key)
	jsonContent, err := alejson.Marshal(content)
	if err != nil {
		return uuid.UUID{}, err
	}
	_, err = c.db.conn.Exec(`
		insert into fda_dbs.user_objs(id, key, db_id, col_id, content)
		values($1,$2,$3,$4,$5::jsonb)
		on conflict (db_id, id) do update set content = excluded.content;`,
		oid.String(), key, c.db.id, c.id, jsonContent)
	return oid, err
}

func (c *col) getObject(key string) (data.Object, error) {
	oid := uuid.NewV5(c.id, key)
	var jsonContent string
	err := c.db.conn.QueryRow(`
	select content from fda_dbs.user_objs
	where id = $1 and db_id = $2 and col_id = $3
	`, oid, c.db.id, c.id).Scan(&jsonContent)
	if err != nil {
		return nil, err
	}
	obj, err := alejson.Unmarshal(jsonContent)
	if err != nil {
		return nil, err
	}
	return obj.(data.Object), nil
}
