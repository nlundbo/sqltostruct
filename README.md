 # sqltostruct 

reads stdin for create-table definitions, makes silly assumptions and prints
a go-struct with `db` struct-tags.

Panics if structure does not follow silly assumptions.

## example
```
$ go run main.go << EOF 
create table example(
id serial,
c1 varchar(255)
);
EOF
```
```
type Example struct{
	Id  `db:"id"`
	C1 string `db:"c1"`
}
```
