package gsql

import (
	"fmt"
	"log"
	"strings"

	"github.com/jmoiron/sqlx"
)

// Del delete record(s).
func (its *GSql) Del(
	db *sqlx.DB,
	conditionsWhere *[]map[string]interface{},
) (err error) {
	if db == nil {
		db, err = its.OpenDB()
		if err != nil {
			return
		}
		defer db.Close()
	}

	// make WHERE works with empty conditionsWhere
	argsWhere := []string{
		"1 = 1",
	}
	var args []interface{}
	for _, item := range *conditionsWhere {
		// hard-coded fix pass `is/is not null` condition
		v, ok := item["value"].(string)
		if ok && v == "null" {
			cond := fmt.Sprintf("%v %v null",
				item["key"],
				item["op"],
			)
			argsWhere = append(argsWhere, cond)
		} else {
			cond := fmt.Sprintf("%v %v ?",
				item["key"],
				item["op"],
			)
			argsWhere = append(argsWhere, cond)
			args = append(args, item["value"])
		}
	}

	s := fmt.Sprintf("DELETE FROM %s WHERE %s LIMIT 1", its.TableName, strings.Join(argsWhere, " AND "))
	if its.Debug {
		log.Println("[debug] sql", s, args)
	}
	_, err = db.Exec(s, args...)
	return
}
