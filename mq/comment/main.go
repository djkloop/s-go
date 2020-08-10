package main

import (
	"encoding/json"
	"fmt"

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

	mq.ConsumerDlx("fuck.comment.count", "fuck_comment_count", "fuck.comment.count.dlx", "fuck_comment_count_dlx", 10, callback)
}

func callback(s string) {
	type Data struct {
		VideoId    int
		EpisodesId int
	}
	var data Data
	err := json.Unmarshal([]byte(s), &data)
	if err == nil {
		o := orm.NewOrm()
		//修改视频的总评论数
		_, _ = o.Raw("UPDATE video SET comment=comment+1 WHERE id=?", data.VideoId).Exec()
		//修改视频剧集的评论数
		_, _ = o.Raw("UPDATE video_episodes SET comment=comment+1 WHERE id=?", data.EpisodesId).Exec()

		//更新redis排行榜 - 通过MQ来实现
		//创建一个简单模式的MQ
		//把要传递的数据转换为json字符串
		videoObj := map[string]int{
			"VideoId": data.VideoId,
		}
		videoJson, _ := json.Marshal(videoObj)
		_ = mq.Publish("", "fuck_top", string(videoJson))
	}
	fmt.Printf("msg is :%s\n", s)
}
