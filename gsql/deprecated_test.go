package gsql_test

import (
	"testing"

	"github.com/giant-stone/go/gsql"
	"github.com/giant-stone/go/gutil"
)

func TestUpdate(t *testing.T) {
	mgr := newAccountProxy()
	db, err := mgr.OpenDB()
	gutil.ExitOnErr(err)
	defer db.Close()

	tearDown(db)
	setUp(db)

	mobilenoExpected := "13800138000"

	// Test Create
	result, err := mgr.Create(db, &map[string]interface{}{"mobileno": mobilenoExpected})
	if err != nil {
		t.Errorf("expected CreateOrUpdate err=nil, got %v", err)
	}
	lastInsertID, err := result.LastInsertId()
	if err != nil || lastInsertID <= 0 {
		t.Errorf("expected LastInsertId lastInsertID>0 err=nil, got %d %v", lastInsertID, err)
	}

	// Test Update
	conditionsWhere := []map[string]interface{}{
		{"key": "mobileno", "op": "like", "value": "%3800%"},
		{"key": "mobileno", "op": "in", "value": []string{"13800138000"}},
		{"key": "mobileno", "op": "is not", "value": "null"},
	}
	result, err = mgr.Update(db, &conditionsWhere, &map[string]interface{}{"password": "1111"})
	if err != nil {
		t.Errorf("expected Update err=nil, got %v", err)
	}
	row, err := result.RowsAffected()
	if err != nil || row <= 0 {
		t.Errorf("expected RowsAffected RowsAffected>0 err=nil, got %d %v", row, err)
	}

	var objs []account
	err = mgr.Gets(db, &objs, nil, &[]map[string]interface{}{{"key": "id", "op": "=", "value": lastInsertID}}, 1)
	cnt := len(objs)
	if err == gsql.ErrRecordNotFound || cnt == 0 {
		t.Errorf("expected Gets cnt>0 err=gsql.ErrRecordNotFound, got %d %v", cnt, err)
	}
	objGot := objs[0]
	if objGot.Mobileno != mobilenoExpected || objGot.Password != "1111" {
		t.Errorf(`expected Gets password="", got %v`, objGot.Password)
	}

	tearDown(db)
}

func TestCreate(t *testing.T) {
	mgr := newAccountProxy()
	db, err := mgr.OpenDB()
	gutil.ExitOnErr(err)
	defer db.Close()

	tearDown(db)
	setUp(db)

	var objs []account
	mobilenoExpected := "13800138000"

	err = mgr.Gets(db, &objs, nil, &[]map[string]interface{}{{"key": "id", "op": "=", "value": 1}}, 1)
	cnt := len(objs)
	if err != nil || cnt > 0 {
		t.Errorf("expected Gets err=nil cnt=0, got %v cnt=%d", err, cnt)
	}

	result, err := mgr.Create(db, &map[string]interface{}{"mobileno": mobilenoExpected})
	if err != nil {
		t.Errorf("expected Create err=nil, got %v", err)
	}
	lastInsertID, err := result.LastInsertId()
	if err != nil || lastInsertID <= 0 {
		t.Errorf("expected LastInsertId lastInsertID>0 err=nil, got %d err=%v", lastInsertID, err)
	}

	err = mgr.Gets(db, &objs, nil, &[]map[string]interface{}{{"key": "id", "op": "=", "value": lastInsertID}}, 1)
	cnt = len(objs)
	if err == gsql.ErrRecordNotFound || cnt == 0 {
		t.Errorf("expected Gets cnt>0 err=gsql.ErrRecordNotFound, got %d %v", cnt, err)
	}

	conditionsWhere := []map[string]interface{}{
		{"key": "mobileno", "op": "like", "value": "%3800%"},
		{"key": "mobileno", "op": "in", "value": []string{"13800138000"}},
		{"key": "mobileno", "op": "is not", "value": "null"},
	}
	err = mgr.Gets(db, &objs, nil, &conditionsWhere, 1)
	cnt = len(objs)
	if err == gsql.ErrRecordNotFound || cnt == 0 {
		t.Errorf("expected Gets err=nil cnt=0, got %v cnt=%d", err, cnt)
	}

	objGot := objs[0]
	if objGot.Mobileno != mobilenoExpected || objGot.Password != "" {
		t.Errorf(`expected Gets password="", got %v`, objGot.Password)
	}

	_, err = mgr.Create(db, &map[string]interface{}{"mobileno": mobilenoExpected})
	if err != gsql.ErrDuplicatedUniqueKey {
		t.Errorf("expected Create err=ErrDuplicatedUniqueKey, got %v", err)
	}

	tearDown(db)
}

func TestGets(t *testing.T) {
	TestCreate(t)
}

func TestCreateOrUpdate(t *testing.T) {
	mgr := newAccountProxy()
	db, err := mgr.OpenDB()
	gutil.ExitOnErr(err)
	defer db.Close()

	tearDown(db)
	setUp(db)

	var objs []account

	m := map[string]interface{}{"mobileno": "13800138000"}
	result, err := mgr.CreateOrUpdate(db, &m)
	if err != nil {
		t.Errorf("expected CreateOrUpdate err=nil, got %v", err)
	}
	lastInsertID, err := result.LastInsertId()
	if err != nil || lastInsertID <= 0 {
		t.Errorf("expected LastInsertId lastInsertID>0, got lastInsertId=%d err=%v", lastInsertID, err)
	}

	passwordNew := "secret"
	mUpdates := map[string]interface{}{
		"id":       lastInsertID,
		"password": passwordNew,
	}
	result, err = mgr.CreateOrUpdate(db, &mUpdates)
	if err != nil {
		t.Errorf("expected CreateOrUpdate err=nil, got%v", err)
	}

	lastInsertID, err = result.LastInsertId()
	if err != nil || lastInsertID <= 0 {
		t.Errorf(`expected CreateOrUpdate lastInsertID>0 err=nil, got %d %v`, lastInsertID, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected <= 0 {
		t.Errorf(`expected CreateOrUpdate rowsAffected>0 err=nil, got %d %v`, rowsAffected, err)
	}

	err = mgr.Gets(db, &objs, nil, &[]map[string]interface{}{{"key": "id", "op": "=", "value": lastInsertID}}, 1)
	cnt := len(objs)
	if err != nil || cnt != 1 {
		t.Errorf("expected Gets err=nil cnt=1, got %v %d", err, cnt)
	}
	obj := objs[0]
	if obj.Password != passwordNew {
		t.Errorf("expected password=%s, got %s", passwordNew, obj.Password)
	}

	tearDown(db)
}
