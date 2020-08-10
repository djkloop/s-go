package controllers

import (
	"fuck_youku_api/models"

	"github.com/astaxie/beego"
)

type TopController struct {
	beego.Controller
}

// 频道排行榜
// @router /channel/top [get]
func (c *TopController) ChannelTop() {
	channelId, _ := c.GetInt("channelId")
	if channelId == 0 {
		c.Data["json"] = ReturnError(4001, "必须返回channelId")
		c.ServeJSON()
		return
	}

	num, videos, err := models.RedisGetChannelTop(channelId)
	if err == nil {
		c.Data["json"] = ReturnSuccess(0, "success", videos, num)
		c.ServeJSON()
		return
	} else {
		c.Data["json"] = ReturnError(4004, "没有相关内容")
		c.ServeJSON()
		return
	}
}

// 类型排行榜
// @router /type/top [get]
func (c *TopController) TypeTop() {
	typeId, _ := c.GetInt("typeId")
	if typeId == 0 {
		c.Data["json"] = ReturnError(4001, "必须返回typeId")
		c.ServeJSON()
		return
	}

	num, videos, err := models.RedisGetTypeTop(typeId)
	if err == nil {
		c.Data["json"] = ReturnSuccess(0, "success", videos, num)
		c.ServeJSON()
		return
	} else {
		c.Data["json"] = ReturnError(4004, "没有相关内容")
		c.ServeJSON()
		return
	}
}
