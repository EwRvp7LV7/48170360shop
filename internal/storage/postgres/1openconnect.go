package postgres

import (
	"log"

	"github.com/EwRvp7LV7/48170360shop/internal/config"
	"github.com/jmoiron/sqlx"
)

//db database
var db *sqlx.DB

//OpenConnectDB - открыть соединение.
func OpenConnectDB() {
	var err error
	db, err = sqlx.Open("postgres", config.GetConnectionStringDB())
	if err != nil {
		log.Fatalln("Cant open connection to postgres", err.Error())
	}

	log.Printf("Connected to DB: %s\n", config.GetInfoDB())

	err = db.Ping()
	if err != nil {
		log.Fatalln("Cant ping", err.Error())
	}

}

//CloseConnectionDB disconnect from database
func CloseConnectionDB() {
	if nil != db {
		db.Close()
	}
}
