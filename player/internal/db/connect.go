package db

import (
	"database/sql"
	"flag"
	"runtime"
	"strings"

	"github.com/corverroos/unsure"
	"github.com/luno/jettison/log"
)

var (
	playerDB = flag.String("player_db", "", "Database name for player")
)

// Connect attempts to create a connection to the MySQL database provided
// by the "player_db" flag.
func Connect() (*sql.DB, error) {
	dbURI := "mysql://root@unix(" + unsure.SockFile() + ")/" + *playerDB + "?"
	ok, err := unsure.MaybeRecreateSchema(dbURI, getSchemaPath())
	if err != nil {
		return nil, err
	} else if ok {
		log.Info(nil, "recreated schema")
	}

	dbc, err := unsure.Connect(dbURI)
	if err != nil {
		return nil, err
	}

	return dbc, nil
}

func getSchemaPath() string {
	_, filename, _, _ := runtime.Caller(0)
	return strings.Replace(filename, "connect.go", "schema.sql", 1)
}
