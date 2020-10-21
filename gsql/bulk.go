package gsql

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/giant-stone/go/gutil"
)

// BulkCreateOrUpdate create or update record(s) in bulk
func (it *GSql) BulkCreateOrUpdate(
	db *sqlx.DB,
	objs []interface{},
	batchsize int,
) (totalWritten int, err error) {
	if db == nil {
		db, err = it.OpenDB()
		if err != nil {
			return
		}
		defer db.Close()
	}

	var columns []string
	var placeholdersValue []string
	var placeholderUpdateOne []string

	if len(objs) == 0 {
		return
	}

	var placeholderInsertOne string
	obj := objs[0]
	_obj, isMap := obj.(*map[string]interface{})
	if !isMap {
		_obj = gutil.Struct2map(obj)
	}
	for k := range *_obj {
		columns = append(columns, k)
		placeholdersValue = append(placeholdersValue, "?")
		placeholderUpdateOne = append(placeholderUpdateOne, fmt.Sprintf("%s=values(%s)", k, k))
	}
	placeholderInsertOne = "( " + strings.Join(placeholdersValue, ", ") + " )"

	total := len(objs)
	for i := 0; i < total; i += batchsize {
		var offset int
		if i+batchsize < total {
			offset = batchsize
		} else {
			offset = total - i
		}

		args := []interface{}{}
		placeholdersInsertAll := []string{}
		for _, obj := range objs[i : i+offset] {
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

		query := fmt.Sprintf("INSERT INTO %s (%s) VALUES %s ON DUPLICATE KEY UPDATE %s",
			it.TableName,
			strings.Join(columns, ", "),
			strings.Join(placeholdersInsertAll, ", "),
			strings.Join(placeholderUpdateOne, ", "),
		)

		ts := time.Now()
		result, errExec := db.Exec(query, args...)
		if errExec != nil {
			err = errExec
			break
		}

		if it.Debug {
			log.Println(fmt.Sprintf("[debug] Writes %d records in %v", offset, time.Since(ts)))
		}

		_totalWritten, err := result.RowsAffected()
		gutil.CheckErr(err)
		totalWritten += int(_totalWritten)
	}

	return
}
