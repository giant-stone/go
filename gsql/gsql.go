// github.com/jmoiron/sqlx CRUD wrapper. DO NOT REPEAT YOURSELF.
package gsql

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"

	"github.com/jmoiron/sqlx"
)

var (
	ErrRecordNotFound      = errors.New("record not found")
	ErrDuplicatedUniqueKey = errors.New("duplicated unique key")
	ErrQueryOrArgumentIsInvalid = errors.New("query or argument is invalid")
)

type GS interface {
	RawQuery(db *sqlx.DB, objs interface{}, s string, args ...interface{}) error
}

// GSql contains database connection settings and info.
type GSql struct {
	// DriverName Go SQL driver name, such as "mysql" "postgres"
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

// NewGSql setup DSN(data source name) and table, sub-class have to override it.
func NewGSql() *GSql {
	w := new(GSql)
	w.DriverName = "mysql"
	w.Dsn = "test:test@tcp(127.0.0.1:3306)/test?charset=utf8mb4,utf8&timeout=2s&writeTimeout=2s&readTimeout=2s&parseTime=true"
	w.TableName = "test"
	return w
}

// OpenDB sqlx.Open wrapper.
func (its *GSql) OpenDB() (db *sqlx.DB, err error) {
	db, err = sqlx.Open(its.DriverName, its.Dsn)
	if err != nil {
		return
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

// Gets query records with where conditions, all conditions are concatenate with " AND " operator.
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
		// hard-coded fix pass `is/is not null` condition
		v, ok := item["value"].(string)
		if ok && v == "null" {
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

// CreateOrUpdate insert a record or update record(s).
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

// Creates insert records in bulk.
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
			case map[string]interface{}:
				{

					var j JSONB
					for key, value := range v.(map[string]interface{}) {
						j[key] = value
					}

					args = append(args, j)
				}
			default:
				{
					args = append(args, v)
				}
			}

			if its.DriverName == "mysql" {
				placeholders = append(placeholders, "?")
			} else if its.DriverName == "postgres" {
				placeholders = append(placeholders, fmt.Sprintf("$%d", i))
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

// GetColumns returns query columns from tag `db` in strutt.
func (its *GSql) GetColumns(obj interface{}) []string {
	te := reflect.TypeOf(obj).Elem()
	columns := []string{}
	for i := 0; i < te.NumField(); i++ {
		field := te.Field(i).Tag.Get("db")
		if field != "" && field != "-" {
			columns = append(columns, field)
		}
	}
	return columns

}

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

// JSONB maps PostgreSQL JSONB type into `map` in Go.
// See also http://coussej.github.io/2016/02/16/Handling-JSONB-in-Go-Structs/
type JSONB map[string]interface{}

// Value convert value from map[string]interface{} into []byte for PostgreSQL JSONB
func (p JSONB) Value() (driver.Value, error) {
	j, err := json.Marshal(p)
	return j, err
}

// Scan convert value from PostgreSQL JSONB into map[string]interface{}
func (p *JSONB) Scan(src interface{}) error {
	source, ok := src.([]byte)
	if !ok {
		return errors.New("Type assertion .([]byte) failed.")
	}

	// walk around issue `NULL`
	if len(source) == 4 && string(source) == "null" {
		return nil
	}

	var i interface{}
	err := json.Unmarshal(source, &i)
	if err != nil {
		return err
	}

	*p, ok = i.(map[string]interface{})
	if !ok {
		return errors.New("Type assertion .(map[string]interface{}) failed.")
	}

	return nil
}

// Update update records with where conditions.
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
		// hard-coded fix pass `is/is not null` condition
		v, ok := item["value"].(string)
		if ok && v == "null" {
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
