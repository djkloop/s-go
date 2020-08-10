package controllers

import (
	"fuck_youku_api/models"

	"github.com/astaxie/beego"
)

type BaseControllers struct {
	beego.Controller
}

// 获取频道地区列表
// @router /channel/region [get]
func (c *BaseControllers) ChannelRegion() {
	channelId, _ := c.GetInt("channelId")
	if channelId == 0 {
		c.Data["json"] = ReturnError(4001, "必须返回channelId")
		c.ServeJSON()
		return
	}

	nums, regions, err := models.GetChannelRegion(channelId)
	if err == nil {
		c.Data["json"] = ReturnSuccess(0, "success", regions, nums)
		c.ServeJSON()
		return
	} else {
		c.Data["json"] = ReturnError(4004, "没有相关内容")
		c.ServeJSON()
		return
	}
}

// 获取频道下类型信息
// @router /channel/type [get]
func (c *BaseControllers) ChannelType() {
	channelId, _ := c.GetInt("channelId")
	if channelId == 0 {
		c.Data["json"] = ReturnError(4001, "必须返回channelId")
		c.ServeJSON()
		return
	}
	nums, types, err := models.GetChannelType(channelId)
	if err == nil {
		c.Data["json"] = ReturnSuccess(0, "success", types, nums)
		c.ServeJSON()
		return
	} else {
		c.Data["json"] = ReturnError(4004, "没有相关内容")
		c.ServeJSON()
		return
	}

}
