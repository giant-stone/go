// Simple tests.
package gsql_test

import (
	"testing"

	"github.com/giant-stone/go/gutil"
)

func TestGSql_BulkCreateOrUpdate(t *testing.T) {
	mgr := newAccountProxy()
	db, err := mgr.OpenDB()
	gutil.ExitOnErr(err)
	defer db.Close()

	tearDown(db)
	setUp(db)

	// insertSamples
	changes := []interface{}{
		account{Id: 1, Mobileno: "12345", Password: "12345"},
		account{Id: 2, Mobileno: "12345", Password: "12345"},
		account{Id: 3, Mobileno: "12345", Password: "12345"},
	}

	_, err = mgr.BulkCreateOrUpdate(db, changes, 2)
	gutil.ExitOnErr(err)

	// update samples
	obj1 := account{Id: 1, Mobileno: "1111", Password: "111"}
	obj2 := account{Id: 2, Mobileno: "222", Password: "2222"}
	changes = []interface{}{obj1, obj2}

	_, err = mgr.BulkCreateOrUpdate(db, changes, 5)
	if err != nil {
		t.Errorf("want BulkCreateOrUpdate err=nil, got %v", err)
	}

	// query samples
	var objsGot []account
	columns := mgr.GetColumns(&account{})
	limit := 10000
	where := []map[string]interface{}{
		map[string]interface{}{"key": "id", "op": "in", "value": []interface{}{"1", "2"}},
	}
	err = mgr.GetsWhere(db, &objsGot, &columns, &where, limit)
	if err != nil {
		t.Errorf("want GetsWhere err=nil, got %v", err)
	}

	objsMapGot := map[int]account{}
	for _, obj := range objsGot {
		objsMapGot[obj.Id] = obj
	}

	// compare samples
	total := 2
	if total != len(objsGot) {
		t.Errorf("want total=%d, got %d", total, len(objsGot))
	}

	obj1Got, ok := objsMapGot[1]
	if !ok {
		t.Error("id=1 not found")
	}
	if obj1Got.Mobileno != obj1.Mobileno {
		t.Errorf("want Mobileno=%v got %v", obj1.Mobileno, obj1Got.Mobileno)
	}

	if obj1Got.Password != obj1.Password {
		t.Errorf("want Password=%v got %v", obj1.Mobileno, obj1Got.Mobileno)
	}

	obj2Got, ok := objsMapGot[2]
	if !ok {
		t.Error("id=1 not found")
	}
	if obj2Got.Mobileno != obj2.Mobileno {
		t.Errorf("want Mobileno=%v got %v", obj2.Mobileno, obj2Got.Mobileno)
	}

	if obj2Got.Password != obj2.Password {
		t.Errorf("want Password=%v got %v", obj2.Mobileno, obj2Got.Mobileno)
	}

	tearDown(db)
}
