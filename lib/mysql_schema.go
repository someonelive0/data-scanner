package lib

import (
	"database/sql"
	"fmt"

	"github.com/jimsmart/schema"
)

func Get_tablenames(db *sql.DB) error {
	// Fetch names of all tables
	tnames, err := schema.TableNames(db)
	if err != nil {
		print("schema.TableNames failed: %s", err)
		return err
	}

	// tnames is [][2]string
	for i := range tnames {
		fmt.Println("Table:", tnames[i][1])
	}

	return nil
}

func Get_columns(db *sql.DB, tablename string) error {
	// Fetch column metadata for given table
	tcols, err := schema.ColumnTypes(db, "", tablename)
	if err != nil {
		print("schema.ColumnTypes failed: %s", err)
		return err
	}

	// tcols is []*sql.ColumnType
	for i := range tcols {
		fmt.Println("Column:", tcols[i].Name(), tcols[i].DatabaseTypeName())
	}

	return nil
}

func Get_views(db *sql.DB, tablename string) error {
	// Fetch names of all views
	vnames, err := schema.ViewNames(db)
	if err != nil {
		print("schema.ViewNames failed: %s", err)
		return err
	}
	// vnames is [][2]string
	for i := range vnames {
		fmt.Println("View:", vnames[i][1])
	}

	// Fetch column metadata for given view
	vcols, err := schema.ColumnTypes(db, "", "monthly_sales_view")
	// tcols is []*sql.ColumnType
	for i := range vcols {
		fmt.Println("Column:", vcols[i].Name(), vcols[i].DatabaseTypeName())
	}

	// Fetch column metadata for all views
	views, err := schema.Views(db)
	// views is [][2]string
	for i := range views {
		fmt.Println("View:", views[i][1])
	}

	return nil
}

func Get_pk(db *sql.DB, tablename string) error {
	// Fetch primary key for given table
	pks, err := schema.PrimaryKey(db, "", "employee_tbl")
	if err != nil {
		print("schema.PrimaryKey failed: %s", err)
		return err
	}

	// pks is []string
	for i := range pks {
		fmt.Println("Primary Key:", pks[i])
	}

	// Output:
	// Primary Key: employee_id

	return nil
}
