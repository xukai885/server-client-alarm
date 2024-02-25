package mysql

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"log"
	"server-client-alarm/settings"
)

var db *sqlx.DB

func Init(cfg *settings.Db) (err error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=Local",
		cfg.User,
		cfg.Pass,
		cfg.Host,
		cfg.Port,
		cfg.Dbname,
	)
	db, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		log.Println("mysql连接报错", err)
		return
	}
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	return
}

func Close() {
	_ = db.Close()
}
