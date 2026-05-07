package mysql

import (
	"bluebell/settings"
	"fmt"

	_ "github.com/go-sql-driver/mysql" // 触发驱动注册,将MySql的连接逻辑挂载到标准库database/sql引擎上
	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB

func Init() (err error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=Local",
		settings.Conf.Databaseconfig.User,
		settings.Conf.Databaseconfig.Password,
		settings.Conf.Databaseconfig.Host,
		settings.Conf.Databaseconfig.Port,
		settings.Conf.Databaseconfig.Dbname,
	)
	db, err = sqlx.Open("mysql", dsn)

	if err != nil {
		return err
	}

	if err := db.Ping(); err != nil {
		fmt.Printf("Error: %v", err)
		return err
	} else {
		fmt.Println("connect to mysql suceessfully")
	}
	return nil
}

func Close() {
	_ = db.Close()
}
