package controllers

import (
	"fuck_youku_api/models"

	"github.com/astaxie/beego"
)

type VideoController struct {
	beego.Controller
}

//  频道页 - 获取顶部广告
//  @router /channel/advert [get]
func (c *VideoController) ChannelAdvert() {
	channelId, _ := c.GetInt("channelId")

	if channelId == 0 {
		c.Data["json"] = ReturnError(4001, "必须指定频道")
		c.ServeJSON()
		return
	}

	num, videos, err := models.GetChannelAdvert(channelId)
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

// 频道页 - 正在热播
// @router /channel/hot [get]
func (c *VideoController) ChannelHotList() {
	channelId, _ := c.GetInt("channelId")
	if channelId == 0 {
		c.Data["json"] = ReturnError(4001, "必须指定频道")
		c.ServeJSON()
		return
	}
	num, videos, err := models.GetChannelHotList(channelId)
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

// 获取日漫、国漫推荐
// @router /channel/recommend/region [get]
func (c *VideoController) ChannelRegionRecommendList() {
	channelId, _ := c.GetInt("channelId")
	if channelId == 0 {
		c.Data["json"] = ReturnError(4001, "必须指定频道")
		c.ServeJSON()
		return
	}

	regionId, _ := c.GetInt("regionId")
	if regionId == 0 {
		c.Data["json"] = ReturnError(4001, "必须指定频道地区")
		c.ServeJSON()
		return
	}
	num, videos, err := models.GetChannelRegionRecommend(channelId, regionId)
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

/**
 * 获取少女推荐
 * @param int channelId
 * @param int typeId
 */
// @router /channel/recommend/type [get]
func (c *VideoController) ChannelTypeRecommendList() {
	channelId, _ := c.GetInt("channelId")
	if channelId == 0 {
		c.Data["json"] = ReturnError(4001, "必须指定频道")
		c.ServeJSON()
		return
	}

	typeId, _ := c.GetInt("typeId")
	if typeId == 0 {
		c.Data["json"] = ReturnError(4001, "必须指定频道地区")
		c.ServeJSON()
		return
	}
	num, videos, err := models.GetChannelTypeRecommend(channelId, typeId)
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

// 视频列表接口
// @router /channel/video [get]
func (c *VideoController) ChannelVideo() {
	channelId, _ := c.GetInt("channelId")
	if channelId == 0 {
		c.Data["json"] = ReturnError(4001, "必须指定频道")
		c.ServeJSON()
		return
	}
	///
	regionId, _ := c.GetInt("regionId")
	///
	typeId, _ := c.GetInt("typeId")
	///
	end := c.GetString("end")
	///
	sort := c.GetString("sort")
	///
	limit, _ := c.GetInt("limit")
	offset, _ := c.GetInt("offset")
	if limit == 0 {
		limit = 12
	}

	num, videos, err := models.GetChannelVideoList(channelId, regionId, typeId, end, sort, offset, limit)
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

// 视频详情
// @router /video/info [get]
func (c *VideoController) VideoInfo() {
	videoId, _ := c.GetInt("videoId")
	if videoId == 0 {
		c.Data["json"] = ReturnError(4001, "必须指定视频ID")
		c.ServeJSON()
		return
	}

	video, err := models.RedisGetVideoInfo(videoId)
	if err == nil {
		c.Data["json"] = ReturnSuccess(0, "success", video, 1)
		c.ServeJSON()
		return
	} else {
		c.Data["json"] = ReturnError(4004, "没有相关内容")
		c.ServeJSON()
		return
	}
}

//视频剧集列表
// @router /video/episodes/list [get]
func (c *VideoController) VideoEpisodesList() {
	videoId, _ := c.GetInt("videoId")
	if videoId == 0 {
		c.Data["json"] = ReturnError(4001, "必须指定视频ID")
		c.ServeJSON()
		return
	}
	nums, videos, err := models.RedisGetVideoEpisodesList(videoId)
	if err == nil {
		c.Data["json"] = ReturnSuccess(0, "success", videos, nums)
		c.ServeJSON()
		return
	} else {
		c.Data["json"] = ReturnError(4004, "没有相关内容")
		c.ServeJSON()
		return
	}
}

// 我的视频管理
// @router /user/video [*]
func (c *VideoController) UserVideo() {
	uid := c.GetString("uid")
	if uid == "" {
		c.Data["json"] = ReturnError(4001, "必须指定用户ID")
		c.ServeJSON()
		return
	}

	num, videos, err := models.GetUserVideo(uid)
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

//
// @router /video/save [post]
func (c *VideoController) VideoSave() {
	playUrl := c.GetString("playUrl")
	title := c.GetString("title")
	subTitle := c.GetString("subTitle")
	channelId, _ := c.GetInt("channelId")
	typeId, _ := c.GetInt("typeId")
	regionId, _ := c.GetInt("regionId")
	uid := c.GetString("uid")
	if channelId == 0 {
		c.Data["json"] = ReturnError(4001, "必须指定视频ID")
		c.ServeJSON()
		return
	}

	if uid == "" {
		c.Data["json"] = ReturnError(4001, "必须指定视频ID")
		c.ServeJSON()
		return
	}

	if playUrl == "" {
		c.Data["json"] = ReturnError(4001, "必须指定视频ID")
		c.ServeJSON()
		return
	}

	if regionId == 0 {
		c.Data["json"] = ReturnError(4001, "必须指定视频ID")
		c.ServeJSON()
		return
	}

	if typeId == 0 {
		c.Data["json"] = ReturnError(4001, "必须指定视频ID")
		c.ServeJSON()
		return
	}
	err := models.MVideoSave(title, subTitle, channelId, regionId, typeId, uid, playUrl)
	if err == nil {
		c.Data["json"] = ReturnSuccess(0, "success", nil, 1)
		c.ServeJSON()
		return
	} else {
		c.Data["json"] = ReturnError(8000, err)
		c.ServeJSON()
		return
	}
}
