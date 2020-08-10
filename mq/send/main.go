package main

import (
	"encoding/json"
	"fmt"

	"fuck_youku_api/models"
	"fuck_youku_api/services/mq"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	_ = beego.LoadAppConfig("ini", "../../conf/app.conf")
	defaultDB := beego.AppConfig.String("defaultDB")
	_ = orm.RegisterDriver("mysql", orm.DRMySQL)
	_ = orm.RegisterDataBase("default", "mysql", defaultDB, 30, 30)

	mq.Consumer("", "fuck_send_message_user", callback)
}

func callback(s string) {
	type Data struct {
		UserId    int
		MessageId int64
	}
	var data Data
	err := json.Unmarshal([]byte(s), &data)
	if err == nil {
		_ = models.SendMessageUser(data.UserId, data.MessageId)
	}
	fmt.Printf("msg is :%s\n", s)
}
