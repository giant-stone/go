package gsql_test

import (
	"testing"

	"github.com/giant-stone/go/gutil"
)

func TestGSql_GetsWhere(t *testing.T) {
	mgr := newAccountProxy()
	db, err := mgr.OpenDB()
	gutil.ExitOnErr(err)
	defer db.Close()

	tearDown(db)
	setUp(db)

	// insert samples
	changes := []interface{}{
		account{Id: 1, Mobileno: "12345", Password: "12345"},
		account{Id: 2, Mobileno: "10000", Password: "10000"},
		account{Id: 3, Mobileno: "9999", Password: "9999"},
	}
	_, err = mgr.BulkCreateOrUpdate(db, changes)
	gutil.ExitOnErr(err)

	// query samples
	var objsGot []account
	columns := mgr.GetColumns(&account{})
	limit := 10000
	err = mgr.GetsWhere(db, &objsGot, &columns, nil, limit)
	if err != nil {
		t.Errorf("want GetsWhere err=nil, got %v", err)
	}

	if len(changes) != len(objsGot) {
		t.Errorf("want len(objsGot)=%d, got %d", len(changes), len(objsGot))
	}

	objsMapGot := map[int]account{}
	for _, obj := range objsGot {
		objsMapGot[obj.Id] = obj
	}
	for _, objI := range changes {
		obj := objI.(account)
		objGot, ok := objsMapGot[obj.Id]
		if !ok {
			t.Errorf("want id=%d found, got not found", obj.Id)
		} else {
			gutil.CmpExpectedGot(t, "id", obj.Id, objGot.Id)
			gutil.CmpExpectedGot(t, "mobileno", obj.Mobileno, objGot.Mobileno)
			gutil.CmpExpectedGot(t, "password", obj.Password, objGot.Password)
		}

	}

	tearDown(db)
}
