# About

GSql is github.com/jmoiron/sqlx CRUD wrapper:

 - Gets - query record(s) with where conditions
 - Create - insert a record
 - Creates - insert records in bulk
 - CreateOrUpdate - create or update record(s)
 - Update update record(s)
 - Del - delete record(s)

Search, MySQL *ONLY*

 - Search - query records with where EQUAL(=) and LIKE conditions
 - SearchFullText - query records with MySQL fulltext index

Misc

 - RawQuery - custom query SQL and arguments
 - RawExec
 - GetColumns - compose xx in `SELECT xx from ...`


For more detail about example, see `gsql_test.go` .

Install 

    go get -v -u github.com/lib/pq
    go get -v -u github.com/go-sql-driver/mysql
    go get -v -u github.com/jmoiron/sqlx
    go get -v -u github.com/giant-stone/go/gsql


## See also

- http://go-database-sql.org/index.html
- https://github.com/go-sql-driver/mysql
- https://github.com/jmoiron/sqlx
- https://github.com/lib/pq



