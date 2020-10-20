package gsql

import (
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/jmoiron/sqlx"

	"github.com/giant-stone/go/gstr"
)

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

	defaultConditionOperator := "AND"

	if conditionsWhere != nil {
		for _, item := range *conditionsWhere {
			conditionOperator, okCond := item["cond"].(string)

			if !okCond {
				conditionOperator = defaultConditionOperator
			} else {
				conditionOperator = strings.ToUpper(conditionOperator)
				if !gstr.StrInSlice([]string{"AND", "OR"}, conditionOperator) {
					conditionOperator = defaultConditionOperator
				}
			}
			if conditionOperator == "OR" {
				orCount += 1
			}

			// hard-coded fix pass `is/is not null` condition
			v, ok := item["value"].(string)
			if ok && v == "null" {
				condExpr := fmt.Sprintf("%v %v null %s",
					item["key"],
					item["op"],
					conditionOperator,
				)
				if orCount == 1 {
					condExpr = "(" + condExpr
				}
				wheres = append(wheres, condExpr)
			} else {
				tname := reflect.TypeOf(item["value"]).String()
				switch tname {
				case "[]interface {}":
					{
						values := item["value"].([]interface{})
						// hard-coded fix pass `in/not in (arg1, arg2, ...)` condition
						cond := fmt.Sprintf("%v %v (%s) %s",
							item["key"],
							item["op"],
							strings.Join(composeNPlaceholders(len(values)), ","),
							conditionOperator,
						)
						if orCount == 1 {
							cond = "(" + cond
						}
						wheres = append(wheres, cond)
						args = append(args, values...)
						break
					}
				default:
					{
						// hard-coded fix pass `is/is not null` condition
						if tname == "null" {
							cond := fmt.Sprintf("%v %v null %s",
								item["key"],
								item["op"],
								conditionOperator,
							)
							if orCount == 1 {
								cond = "(" + cond
							}
							wheres = append(wheres, cond)
							break
						} else {
							cond := fmt.Sprintf("%v %v ? %s",
								item["key"],
								item["op"],
								conditionOperator,
							)
							if orCount == 1 {
								cond = "(" + cond
							}
							wheres = append(wheres, cond)
							args = append(args, item["value"])
						}
					}
				}
			}
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

func composeNPlaceholders(n int) []string {
	var rs []string
	for i := 0; i < n; i += 1 {
		rs = append(rs, "?")
	}
	return rs
}
