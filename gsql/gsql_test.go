// Simple tests.
// Setup database
//   create database test;
//   create user 'test'@'127.0.0.1' identified by 'test';
//   grant all privileges on `test`.* to  'test'@'127.0.0.1';
//   flush privileges;
package gsql_test

import (
	"log"
	"os"
	"strings"
	"testing"

	"github.com/giant-stone/go/gsql"
)

var (
		sqlCreateTest = `CREATE TABLE IF NOT EXISTS test_gsql (
	id int not null AUTO_INCREMENT,
	mobileno varchar(255) default '',
	password varchar(255) default '',
	UNIQUE KEY mobileno (mobileno),
	PRIMARY KEY (id)
); `
	
	sqlDropTest = `DROP TABLE IF EXISTS test_gsql;`
)

func tearDown(mgr *AccountProxy) {
	db, err := mgr.OpenDB()
	if err != nil {
		log.Fatalln("[fatal] mgr.OpenDB", err)
	}
	defer db.Close()
	_, err = db.Exec(sqlDropTest)
	if err != nil {
		if strings.Index(err.Error(), "Unknown table") != -1 {
			// drop a table not exists, safe to skip
		} else {
			log.Fatalln("[fatal] db.Exec", err)
		}
	}
}

func setUp(mgr *AccountProxy) {
	db, err := mgr.OpenDB()
	if err != nil {
		log.Fatalln("[fatal] mgr.OpenDB", err)
	}
	defer db.Close()
	_, err = db.Exec(sqlCreateTest)
	if err != nil {
		log.Fatalln("[fatal] db.Exec", err)
	}
}

type Account struct {
	ID           uint64    `json:"id" db:"id"`
	Mobileno     string    `json:"mobileno" db:"mobileno"`
	Password     string    `json:"password" db:"password"`
}

type AccountProxy struct {
	gsql.GSql
}

func NewAccountProxy() *AccountProxy {
	p := AccountProxy{}
	p.DriverName = "mysql"
	p.Debug = true
	p.Dsn = "test:test@tcp(127.0.0.1:3306)/test?charset=utf8mb4,utf8&timeout=2s&writeTimeout=2s&readTimeout=2s&parseTime=true"
	p.TableName = "test_gsql"
	p.Columns = p.GetColumns(&Account{})
	return &p
}

func TestCreate(t *testing.T) {
	mgr := NewAccountProxy()
	db := mgr.MustOpenDB()
	defer db.Close()

	tearDown(mgr)
	setUp(mgr)

	var accounts []Account
	mobilenoExpected := "13800138000"


	err := mgr.Gets(db, &accounts, nil, &[]map[string]interface{}{{"key":"id", "op":"=","value":1}}, 1)
	cnt := len(accounts)
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

	err = mgr.Gets(db, &accounts, nil, &[]map[string]interface{}{{"key":"id", "op":"=","value":lastInsertID}}, 1)
	cnt = len(accounts)
	if err == gsql.ErrRecordNotFound || cnt == 0 {
		t.Errorf("expected Gets cnt>0 err=gsql.ErrRecordNotFound, got %d %v", cnt, err)
	}
	
	accountGot := accounts[0]
	if accountGot.Mobileno != mobilenoExpected || accountGot.Password != "" {
		t.Errorf(`expected Gets password="", got %v`, accountGot.Password)
	}

	_, err = mgr.Create(db, &map[string]interface{}{"mobileno": mobilenoExpected})
	if err != gsql.ErrDuplicatedUniqueKey {
		t.Errorf("expected Create err=ErrDuplicatedUniqueKey, got %v", err)
	}

	tearDown(mgr)
}


func TestGets(t *testing.T) {
	TestCreate(t)
}

func TestCreateOrUpdate(t *testing.T) {
	mgr := NewAccountProxy()
	db := mgr.MustOpenDB()
	defer db.Close()

	tearDown(mgr)
	setUp(mgr)

	var accounts []Account
	
	m := map[string]interface{}{"mobileno":"13800138000"}
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

	err = mgr.Gets(db, &accounts, nil, &[]map[string]interface{}{{"key":"id","op":"=","value":lastInsertID}}, 1)
	cnt := len(accounts)
	if err != nil || cnt != 1 {
		t.Errorf("expected Gets err=nil cnt=1, got %v %d", err, cnt)
	}
	account := accounts[0]
	if account.Password != passwordNew {
		t.Errorf("expected password=%s, got %s", passwordNew, account.Password)
	}

	tearDown(mgr)
}


func TestDel(t *testing.T) {
	mgr := NewAccountProxy()
	db := mgr.MustOpenDB()
	defer db.Close()

	tearDown(mgr)
	setUp(mgr)

	mobilenoExpected := "13800138000"
	result, err := mgr.Create(db, &map[string]interface{}{"mobileno": mobilenoExpected})
	if err != nil {
		t.Errorf("expected Create err=nil, got %v", err)
	}
	lastInsertID, err := result.LastInsertId()
	if err != nil || lastInsertID <= 0 {
		t.Errorf("expected LastInsertId lastInsertID>0 err=nil, got %d err=%v", lastInsertID, err)
	}
	
	err = mgr.Del(db, &[]map[string]interface{}{{"key":"id", "op":"=", "value":lastInsertID}})
	if err != nil {
		t.Errorf("expected Mgr.Del() err=nil, got %v", err)
	}

	var accounts []Account
	err = mgr.Gets(db, &accounts, nil, &[]map[string]interface{}{{"key":"id", "op":"=", "value":lastInsertID}}, 1)
	cnt := len(accounts)
	if err != nil || cnt != 0 {
		t.Errorf("expected Get err=nil cnt=0, got %v %d", err, cnt)
	}
	
	tearDown(mgr)
}



func TestSearch(t *testing.T) {
	mgr := NewAccountProxy()
	db := mgr.MustOpenDB()
	defer db.Close()

	tearDown(mgr)
	setUp(mgr)
	
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
	
	accounts := []Account{}
	err = mgr.Search(db, &accounts, nil, nil, &map[string]interface{}{"mobileno": "1380"}, 1)
	cnt := len(accounts)
	if err != nil || cnt != 1 {
		t.Errorf("expected Search err=nil cnt=1, got %v %d", err, cnt)
	}

	account := accounts[0]

	if account.Mobileno != mobilenoExpected {
		t.Errorf("expected mobileno=%s, got %s", mobilenoExpected, account.Mobileno)
	}

	var accountsMiss = []Account{}
	err = mgr.Search(db, &accountsMiss, nil, nil, &map[string]interface{}{"mobileno": "8888"}, 1)
	if err != nil {
		t.Errorf("expected Search err=nil, got %v", err)
	}

	cnt = len(accountsMiss)
	if cnt != 0 {
		t.Errorf("expected Search len(objs)=0, got %d", cnt)
	}

	tearDown(mgr)
}

func TestMain(m *testing.M) {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	os.Exit(m.Run())
}
