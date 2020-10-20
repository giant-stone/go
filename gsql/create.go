package gsql

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// Create insert one record.
func (its *GSql) Create(db *sqlx.DB, creates *map[string]interface{}) (result sql.Result, err error) {
	if db == nil {
		db, err = its.OpenDB()
		if err != nil {
			return
		}
		defer db.Close()
	}

	createKeys := []string{}
	createValuesPlaceholder := []string{}
	updates := []string{}

	for k := range *creates {
		createKeys = append(createKeys, k)
		createValuesPlaceholder = append(createValuesPlaceholder, fmt.Sprintf(":%s", k))
		updates = append(updates, fmt.Sprintf("%s=:%s", k, k))

	}

	s := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		its.TableName,
		strings.Join(createKeys, ","),
		strings.Join(createValuesPlaceholder, ","),
	)
	if its.Debug {
		log.Println("[debug] sql", s, creates)
	}
	result, err = db.NamedExec(s, *creates)
	if err != nil {
		if errMysql, ok := err.(*mysql.MySQLError); ok {
			// duplicated record
			if errMysql.Number == 1062 {
				err = ErrDuplicatedUniqueKey
				return
			}
		}
	}

	return
}
