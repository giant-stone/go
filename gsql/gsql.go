// github.com/jmoiron/sqlx CRUD wrapper. DO NOT REPEAT YOURSELF.
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

	"github.com/giant-stone/go/gstr"
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

// GetsWhere query records with where conditions, all conditions are concatenate with " AND " operator by default.
// You have to put all "OR" conditions at the end if mix with "AND" ones.
// case
//   NO "select * from mytbl where a=b or c=d or z=y and e=f"
//   YES "select * from mytbl where a=b and e=f and (c=d or z=y)"
func (its *GSql) GetsWhere(
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
	wheres := []string{}
	args := []interface{}{}

	orCount := 0

	for _, item := range *conditionsWhere {
		cond, okCond := item["cond"].(string)

		if !okCond {
			cond = "AND"
		} else {
			cond = strings.ToUpper(cond)
			if !gstr.StrInSlice([]string{"AND", "OR"}, cond) {
				cond = "AND"
			}
		}
		if cond == "OR" {
			orCount += 1
		}

		// hard-coded fix pass `is/is not null` condition
		v, ok := item["value"].(string)
		if ok && v == "null" {
			cond := fmt.Sprintf("%v %v null %s",
				item["key"],
				item["op"],
				cond,
			)
			if orCount == 1 {
				cond = "(" + cond
			}
			wheres = append(wheres, cond)
		} else {
			cond := fmt.Sprintf("%v %v ? %s",
				item["key"],
				item["op"],
				cond,
			)
			if orCount == 1 {
				cond = "(" + cond
			}
			wheres = append(wheres, cond)
			args = append(args, item["value"])
		}
	}

	if orCount > 0 {
		wheres = append(wheres, "0)")
	} else {
		wheres = append(wheres, "1=1")
	}

	var s string
	s = fmt.Sprintf("SELECT %s FROM %s WHERE %s LIMIT %d",
		columnsQuery,
		its.TableName,
		strings.Join(wheres, " "),
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

// CreateOrUpdateFromStruct insert a record or update record(s) by a struct.
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
		tname := reflect.TypeOf(item["value"]).String()
		switch tname {
		case "string":
			{
				// hard-coded fix pass `is/is not null` condition
				if tname == "null" {
					cond := fmt.Sprintf("%v %v null",
						item["key"],
						item["op"],
					)
					wheres = append(wheres, cond)
					break
				}
			}
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
				break
			}
		default:
			{
				cond := fmt.Sprintf("%v %v ?",
					item["key"],
					item["op"],
				)
				wheres = append(wheres, cond)
				args = append(args, item["value"])
				break
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
