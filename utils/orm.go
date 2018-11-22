package utils

import (
	"github.com/go-xorm/xorm"
	"github.com/luoxiaojun1992/http-dns/models"
	"log"
	"os"
)

var Orm *xorm.Engine

func InitOrm() {
	//Init ORM
	var err error
	Orm, err = xorm.NewEngine("mysql", os.Getenv("DB_USER")+":"+os.Getenv("DB_PWD")+"@/"+os.Getenv("DB_NAME")+"?charset=utf8mb4")
	if err != nil {
		log.Fatal(err)
	}

	//Sync Tables
	err = Orm.Sync2(new(models.IpList))
	if err != nil {
		log.Fatal(err)
	}
}
