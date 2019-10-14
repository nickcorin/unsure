package db

import (
	"database/sql"
	"flag"
	"io/ioutil"
	"os/exec"
	"runtime"
	"strings"

	"github.com/go-sql-driver/mysql"
	"github.com/luno/jettison/errors"
	"github.com/luno/jettison/j"
)

var playerDB = flag.String("player_db", "", "Player DB URI")
var recreateSchema = flag.Bool("recreate_schema", true,
	"Recreate the Player DB")

const defaultOptions = "?parseTime=true"

// Connect attempts to create a connection to the MySQL database provided
// by the "player_db" flag.
func Connect() (*sql.DB, error) {
	// Always attempt to clean the current schema.
	err := maybeRecreateSchema()
	if err != nil {
		return nil, errors.Wrap(err, "failed to recreate schema")
	}

	// Connect to the db schema.
	db, err := sql.Open("mysql", *playerDB+defaultOptions)
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to player database")
	}

	return db, nil
}

// cat schema.sql | mysql -u root -p
func maybeRecreateSchema() error {
	if !*recreateSchema {
		return nil
	}

	// Parse the URI into an accessible config.
	conf, err := mysql.ParseDSN(*playerDB)
	if err != nil {
		return errors.Wrap(err, "failed to parse player db uri")
	}

	// Drop current database.
	err = executeQuery("drop database if exists "+conf.DBName+";", conf)
	if err != nil {
		return errors.Wrap(err, "failed to drop database")
	}

	// Create new database.
	err = executeQuery("create database "+conf.DBName+";", conf)
	if err != nil {
		return errors.Wrap(err, "failed to drop database")
	}

	// Read schema.sql file.
	schema, err := ioutil.ReadFile(getSchemaPath())
	if err != nil {
		return errors.Wrap(err, "failed to read schema")
	}

	// Create schema.
	err = executeQuery(string(schema), conf)
	if err != nil {
		return errors.Wrap(err, "failed to drop database")
	}

	return errors.New("not implemented")
}

func executeQuery(query string, conf *mysql.Config) error {
	var args []string
	if conf.User != "" {
		args = append(args, "-u", conf.User)
	}
	if conf.Passwd != "" {
		args = append(args, "-p", conf.Passwd)
	}

	queryCmd := exec.Command("mysql", args...)
	queryCmd.Stdin = strings.NewReader(query)
	out, err := queryCmd.CombinedOutput()
	if err != nil {
		return errors.Wrap(err, "failed to execute query", j.KV("output", out))
	}

	return nil
}

func getSchemaPath() string {
	_, filename, _, _ := runtime.Caller(0)
	return strings.Replace(filename, "connect.go", "schema.sql", 1)
}
