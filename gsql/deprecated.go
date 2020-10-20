// DEPRECATED.
package gsql

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// Gets query records with where conditions, all conditions are concatenate with " AND " operator.
// DEPRECATED.
func (its *GSql) Gets(
	db *sqlx.DB,
	objs interface{},
	columns *[]string,
	conditionsWhere *[]map[string]interface{},
	limit int) (err error) {
	if db == nil {
		db, err = its.OpenDB()
		if err != nil {
			return
		}
		defer db.Close()
	}

	var columnsQuery string
	if columns != nil && len(*columns) > 0 {
		columnsQuery = strings.Join(*columns, ",")
	} else if len(its.Columns) > 0 {
		columnsQuery = strings.Join(its.Columns, ",")
	} else {
		columnsQuery = "*"
	}

	// make WHERE works with empty conditionsWhere
	wheres := []string{
		"1 = 1",
	}
	args := []interface{}{}
	for _, item := range *conditionsWhere {
		tname := reflect.TypeOf(item["value"]).String()
		switch tname {
		case "[]string":
			{
				// hard-coded fix pass `in/not in (arg1, arg2, ...)` condition
				cond := fmt.Sprintf("%v %v (?)",
					item["key"],
					item["op"],
				)
				wheres = append(wheres, cond)
				v := item["value"].([]string)
				args = append(args, strings.Join(v, ","))
			}
		default:
			{
				// hard-coded fix pass `is/is not null` condition
				if tname == "null" {
					cond := fmt.Sprintf("%v %v null",
						item["key"],
						item["op"],
					)
					wheres = append(wheres, cond)
				} else {
					cond := fmt.Sprintf("%v %v ?",
						item["key"],
						item["op"],
					)
					wheres = append(wheres, cond)
					args = append(args, item["value"])
				}
			}
		}
	}

	var s string
	s = fmt.Sprintf("SELECT %s FROM %s WHERE %s LIMIT %d",
		columnsQuery,
		its.TableName,
		strings.Join(wheres, " AND "),
		limit)

	if its.Debug {
		log.Println("[debug] sql", s, args)
	}

	err = db.Select(objs, s, args...)
	return
}

// CreateOrUpdate insert a record or update record(s).
// DEPRECATED.
func (its *GSql) CreateOrUpdate(
	db *sqlx.DB,
	m *map[string]interface{},
) (result sql.Result, err error) {
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

	for k := range *m {
		createKeys = append(createKeys, k)
		createValuesPlaceholder = append(createValuesPlaceholder, fmt.Sprintf(":%s", k))
		updates = append(updates, fmt.Sprintf("%s=:%s", k, k))

	}

	s := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s) ON DUPLICATE KEY UPDATE %s",
		its.TableName,
		strings.Join(createKeys, ","),
		strings.Join(createValuesPlaceholder, ","),
		strings.Join(updates, ","),
	)
	if its.Debug {
		log.Println("[debug] sql", s, m)
	}
	result, err = db.NamedExec(s, *m)
	if err != nil {
		if mysqlError, ok := err.(*mysql.MySQLError); ok {
			if mysqlError.Number == 1062 {
				err = ErrDuplicatedUniqueKey
			}
		}
	}

	return
}

// CreateOrUpdateFromStruct insert a record or update record(s) by a struct.
// DEPRECATED.
func (its *GSql) CreateOrUpdateFromStruct(
	db *sqlx.DB,
	t interface{},
) (result sql.Result, err error) {
	if db == nil {
		db, err = its.OpenDB()
		if err != nil {
			return
		}
		defer db.Close()
	}

	m := map[string]interface{}{}
	b, _ := json.Marshal(t)
	_ = json.Unmarshal(b, &m)

	createKeys := []string{}
	createValuesPlaceholder := []string{}
	updates := []string{}

	for k := range m {
		createKeys = append(createKeys, k)
		createValuesPlaceholder = append(createValuesPlaceholder, fmt.Sprintf(":%s", k))
		updates = append(updates, fmt.Sprintf("%s=:%s", k, k))

	}

	s := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s) ON DUPLICATE KEY UPDATE %s",
		its.TableName,
		strings.Join(createKeys, ","),
		strings.Join(createValuesPlaceholder, ","),
		strings.Join(updates, ","),
	)
	if its.Debug {
		log.Println("[debug] sql", s, m)
	}
	result, err = db.NamedExec(s, m)
	if err != nil {
		if mysqlError, ok := err.(*mysql.MySQLError); ok {
			if mysqlError.Number == 1062 {
				err = ErrDuplicatedUniqueKey
			}
		}
	}

	return
}

// Creates insert records in bulk.
// DEPRECATED.
func (its *GSql) Creates(db *sqlx.DB, items *[]map[string]interface{}) (result sql.Result, err error) {
	if db == nil {
		db, err = its.OpenDB()
		if err != nil {
			return
		}
		defer db.Close()
	}

	createKeys := []string{}
	recordsPlaceholder := []string{}
	args := []interface{}{}

	itemMap := (*items)[0]
	for k := range itemMap {
		createKeys = append(createKeys, k)
	}
	totalKeys := len(itemMap)

	i := 1
	for _, itemMap := range *items {
		if len(itemMap) != totalKeys {
			err = errors.New("count of keys must be equal in bulk insert")
			return
		}
		placeholders := []string{}

		for _, k := range createKeys {
			v := itemMap[k]

			switch v.(type) {
			default:
				{
					args = append(args, v)
				}
			}

			if its.DriverName == "mysql" {
				placeholders = append(placeholders, "?")
			} else {
				err = errors.New("got unsupport driver " + its.DriverName)
				return
			}

			i++
		}

		l := fmt.Sprintf("(%s)", strings.Join(placeholders, ","))
		recordsPlaceholder = append(recordsPlaceholder, l)
	}

	s := fmt.Sprintf("INSERT INTO %s (%s) VALUES %s",
		its.TableName,
		strings.Join(createKeys, ","),
		strings.Join(recordsPlaceholder, ","),
	)

	//
	// Too many to write, just mute it.
	//
	//	if its.Debug {
	//		log.Println("sql", s)
	//		log.Println(" Args", args)
	//	}

	ts := time.Now()
	result, err = db.Exec(s, args...)
	if its.Debug {
		log.Println(fmt.Sprintf("[debug] Writes %d records in %v", len(*items), time.Since(ts)))
	}

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

// Update update records with where conditions.
// DEPRECATED.
func (its *GSql) Update(
	db *sqlx.DB,
	conditionsWhere *[]map[string]interface{},
	updatesMap *map[string]interface{},
) (result sql.Result, err error) {
	limit := 1

	if db == nil {
		db, err = its.OpenDB()
		if err != nil {
			return
		}
		defer db.Close()
	}

	args := []interface{}{}

	updates := []string{}
	for key, value := range *updatesMap {
		update := fmt.Sprintf("%v=?", key)
		args = append(args, value)
		updates = append(updates, update)
	}

	// make WHERE always works
	wheres := []string{
		"1 = 1",
	}

	for _, item := range *conditionsWhere {
		tname := reflect.TypeOf(item["value"]).String()
		switch tname {
		case "[]string":
			{
				// hard-coded fix pass `in/not in (arg1, arg2, ...)` condition
				cond := fmt.Sprintf("%v %v (?)",
					item["key"],
					item["op"],
				)
				wheres = append(wheres, cond)
				v := item["value"].([]string)
				args = append(args, strings.Join(v, ","))
			}
		default:
			{
				// hard-coded fix pass `is/is not null` condition
				if item["value"] == "null" {
					cond := fmt.Sprintf("%v %v null",
						item["key"],
						item["op"],
					)
					wheres = append(wheres, cond)
				} else {
					cond := fmt.Sprintf("%v %v ?",
						item["key"],
						item["op"],
					)
					wheres = append(wheres, cond)
					args = append(args, item["value"])
				}
			}
		}
	}

	var s string
	s = fmt.Sprintf("UPDATE %s SET %s WHERE %s LIMIT %d",
		its.TableName,
		strings.Join(updates, ","),
		strings.Join(wheres, " AND "),
		limit)

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
