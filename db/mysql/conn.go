package mysql

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func init() {
	db, _ = sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/fileserver?charset=utf8mb4&parseTime=True")
	err := db.Ping()
	if err != nil {
		fmt.Println("Failed to connect to mysql,err = " + err.Error())
		os.Exit(1)
	}
}

// 返回db对象
func DBConn() *sql.DB {
	return db
}
