# About

GSql is github.com/jmoiron/sqlx CRUD wrapper:
 - GetsWhere - query record(s) with mixed AND/OR where conditions
 - Create - insert a record
 - BulkCreateOrUpdate - create or update record(s) in bulk
 - Del - delete record(s)

Search
 - Search - query records with where EQUAL(=) and LIKE conditions
 - SearchFullText - query records with MySQL fulltext index

Misc
 - RawQuery - custom query SQL and arguments
 - RawExec
 - GetColumns - auto compose xx in `SELECT xx from ...`


DEPRECATED.
 - Gets - query record(s) with where AND conditions
 - Creates - insert records in bulk
 - CreateOrUpdate - create or update record(s)
 - Update update record(s)

For more detail about example, see `gsql_test.go` `bulk_test.go` .

Install 

    go get -v -u github.com/giant-stone/go


## See also

- http://go-database-sql.org/index.html
- https://github.com/go-sql-driver/mysql
- https://github.com/jmoiron/sqlx


