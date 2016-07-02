// Package database - general database access
// This is a wrapper on github.com/go-xorm/xorm
package database

import (
	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"

	// supported database drivers
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

// -----------------------------------------------------------------------------

// DB stores database handle
type DB struct {
	*xorm.Engine
	InDebug bool
}

// -----------------------------------------------------------------------------

// Flags is a package flags sample
// in form ready for use with github.com/jessevdk/go-flags
type Flags struct {
	Driver  string `long:"db_driver"  default:"sqlite3"   description:"Database driver"`
	Connect string `long:"db_connect" default:"./test.db" description:"Database connect string"`
	Debug   bool   `long:"db_debug"   description:"Print database debug info"`
}

// -----------------------------------------------------------------------------
// Functional options

// Debug sets sql tracing to on when "on" argument is true
func Debug(on bool) func(db *DB) error {
	return func(db *DB) error {
		return db.setDebug(on)
	}
}

// -----------------------------------------------------------------------------
// Internal setters

func (db *DB) setDebug(on bool) error {
	if on {
		db.InDebug = true
		db.Engine.ShowSQL() // = true
		db.Engine.Logger().SetLevel(core.LOG_DEBUG)
	}
	return nil
}

// -----------------------------------------------------------------------------

// New creates db engine object
// Configuration should be set via functional options
func New(driver, connect string, options ...func(db *DB) error) (*DB, error) {

	engine, err := xorm.NewEngine(driver, connect)
	if err != nil {
		return nil, err
	}

	db := DB{Engine: engine, InDebug: false}

	for _, option := range options {
		err := option(&db)
		if err != nil {
			return nil, err
		}
	}

	return &db, nil
}
