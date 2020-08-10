package models

import "github.com/astaxie/beego/orm"

type Region struct {
	Id   int
	Name string
}

type Type struct {
	Id   int
	Name string
}

/**
 * 获取频道下地区
 * @param int $channelId
 * @return num, regions, err
 */
func GetChannelRegion(channelId int) (int64, []Region, error) {
	o := orm.NewOrm()
	var regions []Region
	num, err := o.Raw("select id, name from channel_region where status=1 and channel_id=? order by sort desc", channelId).QueryRows(&regions)
	return num, regions, err
}

/**
 * 获取频道下类型
 * @param int $channelId
 * @return rs
 */
func GetChannelType(channelId int) (int64, []Type, error) {
	o := orm.NewOrm()
	var types []Type
	num, err := o.Raw("select id, name from channel_type where status=1 and channel_id=? order by sort desc", channelId).QueryRows(&types)
	return num, types, err
}
