// Simple tests.
//
// Issue following grants before run this script
//
//   reset root password
//     mysqladmin -hlocalhost -uroot password
//
//   grants root user to create database
//     create user if not exists 'root'@'127.0.0.1' identified by 'root';
//     grant all on *.* to 'root'@'127.0.0.1';
//
package gsql_test

import (
	"log"
	"os"
	"testing"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"

	"github.com/giant-stone/go/gsql"
)

var (
	sqlCreateDb    = `CREATE DATABASE IF NOT EXISTS testgsql;`
	sqlUseDb       = `USE testgsql;`
	sqlCreateTable = `CREATE TABLE IF NOT EXISTS testusers (
	id int not null AUTO_INCREMENT,
	mobileno varchar(255) default '',
	password varchar(255) default '',
	UNIQUE KEY mobileno (mobileno),
	PRIMARY KEY (id)
); `

	sqlDropDb = `DROP DATABASE IF EXISTS testgsql;`
)

func tearDown(db *sqlx.DB) {
	_, err := db.Exec(sqlDropDb)
	if err != nil {
		if mysqlError, ok := err.(*mysql.MySQLError); ok {
			errorUnknownDatabase := uint16(1049)
			if mysqlError.Number == errorUnknownDatabase {
				return
			}
		}

		log.Fatalln("[fatal] db.Exec", err)
	}
}

func setUp(db *sqlx.DB) {
	_, err := db.Exec(sqlCreateDb)
	if err != nil {
		log.Fatalln("[fatal] db.Exec", err)
	}

	_, err = db.Exec(sqlUseDb)
	if err != nil {
		log.Fatalln("[fatal] db.Exec", err)
	}

	_, err = db.Exec(sqlCreateTable)
	if err != nil {
		log.Fatalln("[fatal] db.Exec", err)
	}
}

type account struct {
	Id       int    `json:"id" db:"id"`
	Mobileno string `json:"mobileno" db:"mobileno"`
	Password string `json:"password" db:"password"`
}

type accountProxy struct {
	gsql.GSql
}

func newAccountProxy() *accountProxy {
	p := accountProxy{}
	p.DriverName = "mysql"
	p.Debug = true
	p.Dsn = "root:root@tcp(127.0.0.1:3306)/?charset=utf8mb4,utf8&timeout=2s&writeTimeout=2s&readTimeout=2s&parseTime=true"
	p.TableName = "testusers"
	p.Columns = p.GetColumns(&account{})
	return &p
}

func TestMain(m *testing.M) {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	os.Exit(m.Run())
}
