package main

import (
	"database/sql"
	"flag"
	"log"

	_ "github.com/go-sql-driver/mysql"

	mylib "data-scanner/lib"
)

// RUN: go run ./myscanner --db mysql --dsn root:<PASSWORD>@tcp(localhost:3306)/db
var arg_db = flag.String("db", "mysql", "database type, such as mysql, postgresql")
var arg_dsn = flag.String("dsn", "root@tcp(localhost:3306)/db", "input mysql dsn, such as root:<PASSWORD>@tcp(localhost:3306)/db")

func init() {
	flag.Parse()
}

func main() {
	db, err := sql.Open(*arg_db, *arg_dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Printf("ping failed: %s\n", err)
	}

	mylib.Get_tablenames(db)

	mylib.Get_columns(db, "works_6796")
}
