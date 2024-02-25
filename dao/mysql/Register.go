package mysql

import (
	"errors"
	"github.com/go-sql-driver/mysql"
	"log"
	"server-client-alarm/modules"
)

func Register(c *modules.Client) (err error) {
	sqlStr := "insert into register(`id`,`ip`,`name`) values(?,?,?)"
	_, err = db.Exec(sqlStr, c.Id, c.Ip, c.Name)
	if err != nil {
		if driverErr, ok := err.(*mysql.MySQLError); ok {
			if driverErr.Number == 1062 {
				return errors.New("数据已存在")
			}
		}
		log.Println("注册client-写入数据库失败", err)
		return
	}
	return nil
}

func ClientInit() (clist []*modules.Client, err error) {
	sqlStr := "select `id`,`name`,`ip` from register"
	err = db.Select(&clist, sqlStr)
	if err != nil {
		log.Println("数据库查询失败", err)
		return
	}
	return
}
