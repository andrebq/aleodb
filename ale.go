package aleodb

import (
	"github.com/kode4food/ale/compiler/arity"
	"github.com/kode4food/ale/core/builtin"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/namespace"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

// RegisterQualified register all couchdb functions to manage the databases
func RegisterQualified(manager *namespace.Manager) namespace.Type {
	ns := manager.GetQualified(data.Name("odb"))
	ns.Declare("connect").Bind(data.MakeApplicative(aleConnectServer, arity.MakeFixedChecker(3)))
	ns.Declare("create-user").Bind(data.MakeApplicative(aleCreateUser, arity.MakeFixedChecker(2)))
	ns.Declare("create-db").Bind(data.MakeApplicative(aleCreateDatabase, arity.MakeFixedChecker(3)))
	ns.Declare("create-collection").Bind(data.MakeApplicative(aleCreateCollection, arity.MakeFixedChecker(2)))
	ns.Declare("random-uuid").Bind(data.MakeApplicative(aleRandomUUID, arity.MakeFixedChecker(0)))
	ns.Declare("put").Bind(data.MakeApplicative(alePutObject, arity.MakeFixedChecker(3)))
	ns.Declare("get").Bind(data.MakeApplicative(aleGetObject, arity.MakeFixedChecker(2)))
	return ns
}

func aleRandomUUID(args ...data.Value) data.Value {
	return data.String(uuid.NewV4().String())
}

func aleGetObject(args ...data.Value) data.Value {
	col := expectCol(args[0])
	id := expectString(args[1])

	obj, err := col.getObject(id)
	if err != nil {
		raiseString(err.Error())
	}
	return obj
}

func alePutObject(args ...data.Value) data.Value {
	col := expectCol(args[0])
	id := expectString(args[1])
	obj := expectObject(args[2])
	var oid uuid.UUID
	var err error
	if oid, err = col.putObject(id, obj); err != nil {
		raiseString(err.Error())
	}
	return obj.Merge(data.NewObject(data.Keyword("id"), data.String(oid.String())))
}

func aleCreateCollection(args ...data.Value) data.Value {
	db := expectDB(args[0])
	name := expectString(args[1])

	col, err := db.createCollection(name)
	if err != nil {
		db.logentry().
			WithField("action", "create-collection").
			WithField("name", name).
			WithError(err).
			Error("Unable to create collection")
		raiseString(err.Error())
	}
	return col
}

func aleConnectServer(args ...data.Value) data.Value {
	host := expectString(args[0])
	user := expectString(args[1])
	passwd := expectString(args[2])

	server, err := connectToServer(host, user, passwd)
	if err != nil {
		logrus.WithField("operation", "connect-server").
			WithField("user", user).
			WithField("host", host).
			WithError(err).
			Error("Unable to connect to server")
		raiseString(err.Error())
	}
	return server
}

func aleCreateUser(args ...data.Value) data.Value {
	s := expectServer(args[0])
	username := expectString(args[1])

	uid, err := s.createUser(username)
	if err != nil {
		s.logentry().
			WithField("ale", true).
			WithField("username", username).
			WithField("operation", "create-user").
			WithError(err).
			Error("Unable to create user")
		raiseString("unable to create user " + username)
	}
	return uid
}

func aleCreateDatabase(args ...data.Value) data.Value {
	s := expectServer(args[0])
	ownerid := expectUUID(args[1])
	dbname := expectString(args[2])

	db, err := s.createDatabase(ownerid, dbname)
	if err != nil {
		s.logentry().
			WithField("ale", true).
			WithField("operation", "create-database").
			WithError(err).
			Error("Unable to create database")
		raiseString("unable to create database " + dbname)
	}
	return db
}

func expectServer(v data.Value) *server {
	s, ok := v.(*server)
	if !ok {
		raiseString("expected a server object")
	}
	return s
}

func expectDB(v data.Value) *db {
	d, ok := v.(*db)
	if !ok {
		raiseString("expected a db object")
	}
	return d
}

func expectCol(v data.Value) *col {
	c, ok := v.(*col)
	if !ok {
		raiseString("expected a col object")
	}
	return c
}

func expectObject(v data.Value) data.Object {
	o, ok := v.(data.Object)
	if !ok {
		raiseString("expected a data.Object")
	}
	return o
}

func expectString(v data.Value) string {
	return v.String()
}

func expectUUID(v data.Value) uuid.UUID {
	u, ok := v.(uuid.UUID)
	if !ok {
		str := v.String()
		u, err := uuid.FromString(str)
		if err != nil {
			raiseString("expecting a uuid.UUID/string that parses to uuid.UUID")
		}
		return u
	}
	return u
}

func raiseString(str string) {
	builtin.Raise(data.String(str))
}
