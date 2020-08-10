package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	"fuck_youku_api/models"
	"fuck_youku_api/services/mq"
	redisClient "fuck_youku_api/services/redis"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	_ = beego.LoadAppConfig("ini", "../../conf/app.conf")
	defaultDB := beego.AppConfig.String("defaultDB")
	_ = orm.RegisterDriver("mysql", orm.DRMySQL)
	_ = orm.RegisterDataBase("default", "mysql", defaultDB, 30, 30)

	mq.Consumer("", "fuck_top", callback)
}

func callback(s string) {
	type Data struct {
		VideoId int
	}
	var data Data
	err := json.Unmarshal([]byte(s), &data)
	videoInfo, err := models.RedisGetVideoInfo(data.VideoId)
	if err == nil {
		conn := redisClient.PoolConnect()
		defer conn.Close()
		// 更新排行榜
		redisChannelKey := "video:top:channel:channelId:" + strconv.Itoa(videoInfo.ChannelId)
		redisTypeKey := "video:top:type:typeId:" + strconv.Itoa(videoInfo.TypeId)
		conn.Do("zincrby", redisChannelKey, 1, data.VideoId)
		conn.Do("zincrby", redisTypeKey, 1, data.VideoId)
		fmt.Printf("msg is :%s\n", s)
	}
}
