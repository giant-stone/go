package gsql

import (
	"fmt"
	"log"
	"strings"

	"github.com/jmoiron/sqlx"
)

// SearchFullText returns query records matched fulltext index.
// This query required created index likes `alter table mytbl add FULLTEXT ft_search (idx_col_a, idx_col_b, ...) WITH PARSER ngram`.
// MySQL *ONLY*.
func (its *GSql) SearchFullText(
	db *sqlx.DB,
	objs interface{},
	columnsQuery *[]string,
	columnsSearch *[]string,
	keywords string,
	limit int) (err error) {
	if db == nil {
		db, err = its.OpenDB()
		if err != nil {
			return
		}
		defer db.Close()
	}

	var columns string
	if columnsQuery != nil && len(*columnsQuery) == 0 {
		columns = strings.Join(*columnsQuery, ",")
	} else if len(its.Columns) > 0 {
		columns = strings.Join(its.Columns, ",")
	} else {
		columns = "*"
	}

	if columnsSearch == nil || len(*columnsSearch) == 0 {
		return ErrQueryOrArgumentIsInvalid
	}

	var args []interface{}
	args = append(args, keywords)
	args = append(args, limit)

	s := fmt.Sprintf("SELECT %s FROM %s WHERE MATCH (%s) AGAINST (?) LIMIT ?",
		columns,
		its.TableName,
		strings.Join(*columnsSearch, ","),
	)

	if its.Debug {
		log.Println("[debug] sql", s, args)
	}

	err = db.Select(objs, s, args...)
	return
}

// Search query records with where EQUAL(=) and LIKE conditions, mysql driver *ONLY*.
func (its *GSql) Search(
	db *sqlx.DB,
	objs interface{},
	columnsQuery *[]string,
	conditionsWhere *map[string]interface{},
	conditionsLike *map[string]interface{},
	limit int) (err error) {
	if db == nil {
		db, err = its.OpenDB()
		if err != nil {
			return
		}
		defer db.Close()
	}

	var columns string
	if columnsQuery != nil && len(*columnsQuery) > 0 {
		columns = strings.Join(*columnsQuery, ",")
	} else if len(its.Columns) > 0 {
		columns = strings.Join(its.Columns, ",")
	} else {
		columns = "*"
	}

	var k string
	wheres := []string{}
	args := []interface{}{}
	if conditionsWhere != nil {
		for k = range *conditionsWhere {
			wheres = append(wheres, fmt.Sprintf("%s= ? ", k))
			value := (*conditionsWhere)[k]
			args = append(args, value)
		}
	}

	if conditionsLike != nil {
		for k = range *conditionsLike {
			wheres = append(wheres, fmt.Sprintf("%s LIKE ?", k))
			value := fmt.Sprintf(`%%%s%%`, (*conditionsLike)[k])
			args = append(args, value)
		}
	}

	s := fmt.Sprintf("SELECT %s FROM %s WHERE %s LIMIT %d",
		columns,
		its.TableName,
		strings.Join(wheres, " OR "),
		limit)

	if its.Debug {
		log.Println("[debug] sql", s, args, conditionsWhere, conditionsLike)
	}

	err = db.Select(objs, s, args...)
	return
}
