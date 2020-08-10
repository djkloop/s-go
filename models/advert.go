package models

import (
	"github.com/astaxie/beego/orm"
)

type Advert struct {
	Id       int
	Title    string
	SubTitle string
	AddTime  int64
	Img      string
	Url      string
	Status   int
}

func init() {
	orm.RegisterModel(new(Advert))
}

// 根据频道ID获取顶部广告
func GetChannelAdvert(channelId int) (int64, []Advert, error) {
	var (
		num     int64
		adverts []Advert
		err     error
	)
	o := orm.NewOrm()
	num, err = o.Raw("select id, title, sub_title, img, add_time, url from advert where status=1 and channel_id=? order by sort desc limit 1", channelId).QueryRows(&adverts)
	return num, adverts, err
}
