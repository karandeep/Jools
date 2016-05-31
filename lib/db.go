package lib

import (
	"../config"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

var mysqlDbConn *sql.DB
var hotOrNotDbConn *sql.DB
var uploadDbConn *sql.DB
var newsDbConn *sql.DB

func GetDBConnection() *sql.DB {
	if mysqlDbConn == nil {
		configData := config.GetConfig()
		var err error
		mysqlDbConn, err = sql.Open("mysql", "jools:@ceVentur@P#one!x@tcp("+configData.DB_HOST+":3306)/jools")
		if err != nil {
			log.Fatal("Oh noez, could not connect to database", err)
		}
	}
	return mysqlDbConn
}

func GetHotOrNotDBConnection() *sql.DB {
	if hotOrNotDbConn == nil {
		configData := config.GetConfig()
		var err error
		hotOrNotDbConn, err = sql.Open("mysql", "jools:@ceVentur@P#one!x@tcp("+configData.DB_HOST+":3306)/hotornot")
		if err != nil {
			log.Fatal("Oh noez, could not connect to database", err)
		}
	}
	return hotOrNotDbConn
}

func GetUploadDBConnection() *sql.DB {
	if uploadDbConn == nil {
		configData := config.GetConfig()
		var err error
		uploadDbConn, err = sql.Open("mysql", "jools:@ceVentur@P#one!x@tcp("+configData.DB_HOST+":3306)/Uploads")
		if err != nil {
			log.Fatal("Oh noez, could not connect to database", err)
		}
	}
	return uploadDbConn
}

func GetNewsDBConnection() *sql.DB {
	if newsDbConn == nil {
		configData := config.GetConfig()
		var err error
		newsDbConn, err = sql.Open("mysql", "jools:@ceVentur@P#one!x@tcp("+configData.DB_HOST+":3306)/news")
		if err != nil {
			log.Fatal("Oh noez, could not connect to database", err)
		}
	}
	return newsDbConn
}
