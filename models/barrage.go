package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

type BarrageInfo struct {
	Content     string `form:"content" valid:"Required"`
	CurrentTime int    `form:"currentTime"`
	Uid         string `form:"uid" valid:"Required"`
	UerId       string `form:"uuid" valid:"Required"`
	EpisodesId  int    `form:"episodesId"`
	VideoId     int    `form:"videoId"`
}

type Barrage struct {
	Id          int
	Content     string
	CurrentTime int
	AddTime     int64
	UserId      string
	Status      int
	EpisodesId  int
	VideoId     int
}

type BarrageData struct {
	Id          int    `json:"id"`
	Content     string `json:"content"`
	CurrentTime int    `json:"currentTime"`
}

func init() {
	orm.RegisterModel(new(Barrage))
}

func BarrageList(episodesId int, startTime int, endTime int) (int64, []BarrageData, error) {
	o := orm.NewOrm()
	var barrages []BarrageData
	num, err := o.Raw("select id, content, `current_time` from barrage where status=1 and episodes_id=? and `current_time`>=? and `current_time`<? order by `current_time` asc", episodesId, startTime, endTime).QueryRows(&barrages)
	return num, barrages, err
}

func SaveBarrage(info BarrageInfo) error {
	o := orm.NewOrm()
	var barrage Barrage
	barrage.Content = info.Content
	barrage.CurrentTime = info.CurrentTime
	barrage.EpisodesId = info.EpisodesId
	barrage.AddTime = time.Now().Unix()
	barrage.Status = 1
	barrage.UserId = info.UerId
	barrage.VideoId = info.VideoId
	_, err := o.Insert(&barrage)
	return err
}
