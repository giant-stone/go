package gsql

import (
	"database/sql"
	"log"
	"reflect"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// RawQuery query records custom SQL and arguments.
func (its *GSql) RawQuery(db *sqlx.DB, objs interface{}, s string, args ...interface{}) (err error) {
	if db == nil {
		db, err = its.OpenDB()
		if err != nil {
			return
		}
		defer db.Close()
	}

	if its.Debug {
		log.Println("[debug] sql", s, args)
	}

	err = db.Select(objs, s, args...)
	return

}

// RawExec db.Exec wrapper.
func (its *GSql) RawExec(db *sqlx.DB, s string, args ...interface{}) (result sql.Result, err error) {
	if db == nil {
		db, err = its.OpenDB()
		if err != nil {
			return
		}
		defer db.Close()
	}

	if its.Debug {
		log.Println("[debug] sql", s, args)
	}

	result, err = db.Exec(s, args...)
	if err != nil {
		if mysqlError, ok := err.(*mysql.MySQLError); ok {
			if mysqlError.Number == 1062 {
				err = ErrDuplicatedUniqueKey
			}
		}
	}
	return
}

func parseStructTags(s interface{}, tagsDb *[]string, tagName string) {
	t := reflect.TypeOf(s)
	te := t.Elem()
	parseStructFieldTag(te, tagsDb, tagName)
}

func parseStructFieldTag(tt reflect.Type, tags *[]string, tagName string) {
	for i := 0; i < tt.NumField(); i++ {
		subT := tt.Field(i).Type
		subTName := subT.Kind().String()
		if subTName == "struct" {
			parseStructFieldTag(subT, tags, tagName)
			continue
		}

		field := tt.Field(i).Tag.Get(tagName)
		if field != "" && field != "-" {
			*tags = append(*tags, field)
		}
	}
}

// GetColumns returns query columns from tag `db` in strutt.
// Example GetColumns(&myObj{})
func (its *GSql) GetColumns(obj interface{}) []string {
	var tagsDb []string
	parseStructTags(obj, &tagsDb, "db")
	return tagsDb

}
