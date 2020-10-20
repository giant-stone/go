package gsql_test

import (
	"testing"

	"github.com/giant-stone/go/gutil"
)

func TestDel(t *testing.T) {
	mgr := newAccountProxy()
	db, err := mgr.OpenDB()
	gutil.ExitOnErr(err)
	defer db.Close()

	tearDown(db)
	setUp(db)

	mobilenoExpected := "13800138000"
	result, err := mgr.Create(db, &map[string]interface{}{"mobileno": mobilenoExpected})
	if err != nil {
		t.Errorf("expected Create err=nil, got %v", err)
	}
	lastInsertID, err := result.LastInsertId()
	if err != nil || lastInsertID <= 0 {
		t.Errorf("expected LastInsertId lastInsertID>0 err=nil, got %d err=%v", lastInsertID, err)
	}

	err = mgr.Del(db, &[]map[string]interface{}{{"key": "id", "op": "=", "value": lastInsertID}})
	if err != nil {
		t.Errorf("expected Mgr.Del() err=nil, got %v", err)
	}

	var objs []account
	err = mgr.Gets(db, &objs, nil, &[]map[string]interface{}{{"key": "id", "op": "=", "value": lastInsertID}}, 1)
	cnt := len(objs)
	if err != nil || cnt != 0 {
		t.Errorf("expected Get err=nil cnt=0, got %v %d", err, cnt)
	}

	tearDown(db)
}
