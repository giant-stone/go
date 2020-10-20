// github.com/jmoiron/sqlx CRUD wrapper. DO NOT REPEAT YOURSELF.
package gsql

import (
	"errors"
	"log"

	"github.com/go-sql-driver/mysql"

	"github.com/jmoiron/sqlx"
)

var (
	ErrRecordNotFound           = errors.New("record not found")
	ErrDuplicatedUniqueKey      = errors.New("duplicated unique key")
	ErrQueryOrArgumentIsInvalid = errors.New("query or argument is invalid")
)

type GS interface {
	RawQuery(db *sqlx.DB, objs interface{}, s string, args ...interface{}) error
}

// GSql contains database connection settings and info.
type GSql struct {
	// DriverName Go SQL driver name, such as "mysql"
	DriverName string

	// Dsn is short for Data source name, such as "test:test@tcp(127.0.0.1:3306)/test?charset=utf8mb4,utf8&timeout=2s&writeTimeout=2s&readTimeout=2s&parseTime=true"
	Dsn string

	// Debug set true to print query and query arguments, default is false.
	Debug bool

	// TableName table for read or write
	TableName string

	// Columns columns for query
	Columns []string
}

// OpenDB sqlx.Open wrapper.
func (its *GSql) OpenDB() (db *sqlx.DB, err error) {
	db, err = sqlx.Open(its.DriverName, its.Dsn)

	if err != nil {
		if mysqlError, ok := err.(*mysql.MySQLError); ok {
			errorUnknownDatabase := uint16(1049)
			if mysqlError.Number == errorUnknownDatabase {
				return db, nil
			}
		}

		log.Fatalln("[fatal] db.Ping", err)
	}

	err = db.Ping()
	return
}

// MustOpenDB sqlx.Open wrapper with fatal on error.
func (its *GSql) MustOpenDB() (db *sqlx.DB) {
	db, err := sqlx.Open(its.DriverName, its.Dsn)
	if err != nil {
		log.Fatalln("[fatal] sqlx.Open", its.DriverName, its.Dsn, err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalln("[fatal] db.Ping", err)
	}
	return
}
