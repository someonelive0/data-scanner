package testing

import (
	"database/sql"
	"flag"
	"log"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

// RUN: go test -v ./testing/ --db mysql --dsn "rroot:<PASSWORD>@tcp(localhost:3306)/db"
var arg_db = flag.String("db", "mysql", "database type, such as mysql, postgresql")
var arg_dsn = flag.String("dsn", "root@tcp(localhost:3306)/db", "input mysql dsn, such as root:<PASSWORD>@tcp(localhost:3306)/db")

func Test_mysql_tablenames(t *testing.T) {
	flag.Parse()

	db, err := sql.Open(*arg_db, *arg_dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Printf("ping failed: %s\n", err)
	}

	// get schema of tables with sql sentence
	const sentence = "select table_name, table_schema from information_schema.tables where table_schema = ?"
	rows, err := db.Query(sentence, "fakedb")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var table_name string
	var table_schema string
	for rows.Next() {
		err := rows.Scan(&table_name, &table_schema)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(table_name, table_schema)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
}
