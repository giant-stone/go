package gutil

import (
	"reflect"
	"testing"
)

func CmpExpectedGot(t *testing.T, key, expI, gotI interface{}) {
	if expI == nil {
		if expI != gotI {
			t.Errorf("%+q expected=%v got=%v", key, expI, gotI)
		}
		return
	}

	tt := reflect.TypeOf(expI).String()
	switch tt {
	case "string":
		{
			exp := expI.(string)
			got := gotI.(string)

			if exp != got {
				t.Errorf("%+q expected=-%s- got=-%s-", key, exp, got)
			}
			break
		}
	case "float32":
	case "float64":
		{
			exp := expI.(float64)
			got := gotI.(float64)

			if (exp - got) < 0.001 {
				t.Errorf("%+q expected=%f got=%f", key, exp, got)
			}
			break
		}

	case "bool":
		{
			exp := expI.(bool)
			got := gotI.(bool)

			if exp != got {
				t.Errorf("%+q expected=%t got=%t", key, exp, got)
			}
			break
		}

	case "uint8":
		{
			exp := expI.(uint8)
			got := gotI.(uint8)

			if exp != got {
				t.Errorf("%+q expected=%d got=%d", key, exp, got)
			}
			break
		}

	case "uint16":
		{
			exp := expI.(uint16)
			got := gotI.(uint16)

			if exp != got {
				t.Errorf("%+q expected=%d got=%d", key, exp, got)
			}
			break
		}

	case "int":
	case "int64":
		{
			exp := expI.(int64)
			got := gotI.(int64)

			if exp != got {
				t.Errorf("%+q expected=%d got=%d", key, exp, got)
			}
			break
		}

	case "uint32":
		{
			exp := expI.(uint32)
			got := gotI.(uint32)

			if exp != got {
				t.Errorf("%+q expected=%d got=%d", key, exp, got)
			}
			break
		}

	case "*errors.errorString":
		{
			exp := expI.(error)
			got := gotI.(error)

			if exp != got {
				t.Errorf("%+q expected=%v got=%v", key, exp, got)
			}
			break
		}
	default:
		{
			t.Errorf("field -%+q- got unknown type %s", key, tt)
		}
	}
}
