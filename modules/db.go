package modules

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	sqlx "github.com/jmoiron/sqlx"
	"log"
	"os"
)

var DB *sqlx.DB

func InitDB(login, password, host, port, database string) {
	s := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8",
		login,
		password,
		host,
		port,
		database,
	)
	var err error
	DB, err = sqlx.Connect("mysql", s)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)

		err = DB.Ping()
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
	}
	log.Printf("Connected to DB %s:%s / %s", host, port, database)
}
