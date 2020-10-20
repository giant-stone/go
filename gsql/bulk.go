package gsql

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"

	"github.com/giant-stone/go/gutil"
)

// BulkCreateOrUpdate create or update record(s) in bulk
func (it *GSql) BulkCreateOrUpdate(
	db *sqlx.DB,
	// objs map[string]*map[string]interface{},
	objs map[string]interface{},
) (rowsAffected int64, err error) {
	if db == nil {
		db, err = it.OpenDB()
		if err != nil {
			return
		}
		defer db.Close()
	}

	var columns []string
	var args []interface{}

	var placeholdersValue []string
	var placeholdersInsertAll []string
	var placeholderUpdateOne []string

	if len(objs) == 0 {
		return
	}

	var placeholderInsertOne string

	var skipPk string
	for pk, obj := range objs {
		_obj, isMap := obj.(*map[string]interface{})
		if !isMap {
			_obj = gutil.Struct2map(obj)
		}

		skipPk = pk
		for k, v := range *_obj {
			columns = append(columns, k)
			args = append(args, v)
			placeholdersValue = append(placeholdersValue, "?")
			placeholderUpdateOne = append(placeholderUpdateOne, fmt.Sprintf("%s=values(%s)", k, k))
		}
		break
	}

	placeholderInsertOne = "( " + strings.Join(placeholdersValue, ", ") + " )"
	placeholdersInsertAll = append(placeholdersInsertAll, placeholderInsertOne)

	for pk, obj := range objs {
		if pk == skipPk {
			continue
		}

		_obj, isMap := obj.(*map[string]interface{})
		if !isMap {
			_obj = gutil.Struct2map(obj)
		}

		for _, k := range columns {
			v, ok := (*_obj)[k]
			if !ok {
				log.Printf("[warn] skip for key not found, key=%s obj=%v", k, _obj)
				continue
			}
			args = append(args, v)
		}

		placeholdersInsertAll = append(placeholdersInsertAll, placeholderInsertOne)
	}

	//  INSERT INTO mytbl (name, cnt, ver) VALUES ('foo',11, 1),('bar',222, 2) ON DUPLICATE KEY UPDATE cnt=values(cnt), cnt=values(cnt), ver=values(ver);
	s := fmt.Sprintf("INSERT INTO %s (%s) VALUES %s ON DUPLICATE KEY UPDATE %s",
		it.TableName,
		strings.Join(columns, ", "),
		strings.Join(placeholdersInsertAll, ", "),
		strings.Join(placeholderUpdateOne, ", "),
	)

	ts := time.Now()
	result, err := db.Exec(s, args...)
	if it.Debug {
		log.Println(fmt.Sprintf("[debug] Writes %d records in %v", len(objs), time.Since(ts)))
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

	if result != nil {
		rowsAffected, err = result.RowsAffected()
	}

	return
}
