package dbutil

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const (
	DataBaseUser     = "root"
	DataBasePassword = "LXY1019XYXYZ"
	DataBaseHost     = "127.0.0.1"
	DataBasePort     = "3306"
	DataBaseName     = "web_bms"
)

func dbInit() *sql.DB {
	connString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8",
		DataBaseUser,
		DataBasePassword,
		DataBaseHost,
		DataBasePort,
		DataBaseName)
	db, err := sql.Open("mysql", connString)
	if err != nil {
		log.Fatal(err)
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	return db
}

// Prepare statement for reading data
func Query(sqlString string, args ...any) []map[string]any {
	db := dbInit()
	defer db.Close()

	// Execute the query
	rows, err := db.Query(sqlString, args...)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	// Make a slice for the values
	values := make([]sql.RawBytes, len(columns))

	// rows.Scan wants '[]interface{}' as an argument, so we must copy the
	// references into such a slice
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	dataSlice := make([]map[string]any, 0)

	// Fetch rows
	for rows.Next() {
		// get RawBytes from data
		err = rows.Scan(scanArgs...)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}

		// Now do something with the data.
		// Here we just print each column as a string.
		var value string
		oneData := make(map[string]any, len(columns))
		for i, col := range values {
			// Here we can check if the value is nil (NULL value)
			if col == nil {
				value = "NULL"
			} else {
				value = string(col)
			}
			oneData[columns[i]] = value
			// fmt.Println(columns[i], ": ", value)
		}
		dataSlice = append(dataSlice, oneData)
		// fmt.Println("-----------------------------------")
	}
	if err = rows.Err(); err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	return dataSlice
}

// update 更新 【更新，插入，删除都是exec方法】
func Update(sqlString string, params ...any) error {
	db := dbInit()
	defer db.Close()
	result, err := db.Exec(sqlString, params...)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("%s\n", "update error!")
	}
	return nil
}
