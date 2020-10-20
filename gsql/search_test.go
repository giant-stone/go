package gsql_test

import (
	"testing"

	"github.com/giant-stone/go/gutil"
)

func TestSearch(t *testing.T) {
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

	var objs []account
	err = mgr.Search(db, &objs, nil, nil, &map[string]interface{}{"mobileno": "1380"}, 1)
	cnt := len(objs)
	if err != nil || cnt != 1 {
		t.Errorf("expected Search err=nil cnt=1, got %v %d", err, cnt)
	}

	obj := objs[0]

	if obj.Mobileno != mobilenoExpected {
		t.Errorf("expected mobileno=%s, got %s", mobilenoExpected, obj.Mobileno)
	}

	var accountsMiss []account
	err = mgr.Search(db, &accountsMiss, nil, nil, &map[string]interface{}{"mobileno": "8888"}, 1)
	if err != nil {
		t.Errorf("expected Search err=nil, got %v", err)
	}

	cnt = len(accountsMiss)
	if cnt != 0 {
		t.Errorf("expected Search len(objs)=0, got %d", cnt)
	}

	tearDown(db)
}
