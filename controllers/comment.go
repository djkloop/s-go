package controllers

import (
	"fuck_youku_api/models"

	"github.com/astaxie/beego"
)

type CommentController struct {
	beego.Controller
}

type CommentInfo struct {
	Id           int             `json:"id"`
	Content      string          `json:"content"`
	AddTime      int64           `json:"addTime"`
	AddTimeTitle string          `json:"addTimeTitle"`
	UserId       string          `json:"userId"`
	Stamp        int             `json:"stamp"`
	PraiseCount  int             `json:"praiseCount"`
	UserInfo     models.UserInfo `json:"userinfo"`
}

// 获取评论列表
// @router /comment/list [get]
//func (c *CommentController) List() {
//	episodesId, _ := c.GetInt("episodesId")
//	if episodesId == 0 {
//		c.Data["json"] = ReturnError(4001, "必须指定视频ID")
//		c.ServeJSON()
//		return
//	}
//
//	limit, _ := c.GetInt("limit")
//	offset, _ := c.GetInt("offset")
//	if limit == 0 {
//		limit = 12
//	}
//
//	nums, comments, err := models.GetCommentList(episodesId, limit, offset)
//	if err == nil {
//		var data []CommentInfo
//		var commentInfo CommentInfo
//		for _, v := range comments {
//			commentInfo.Id = v.Id
//			commentInfo.Content = v.Content
//			commentInfo.AddTime = v.AddTime
//			commentInfo.AddTimeTitle = DateFormat(v.AddTime)
//			commentInfo.UserId = v.UserId
//			commentInfo.Stamp = v.Stamp
//			commentInfo.PraiseCount = v.PraiseCount
//			// 获取用户信息
//			commentInfo.UserInfo, _ = models.RedisGetUserInfo(v.UserId)
//			data = append(data, commentInfo)
//		}
//
//		c.Data["json"] = ReturnSuccess(0, "success", data, nums)
//		c.ServeJSON()
//		return
//	} else {
//		c.Data["json"] = ReturnError(4004, "没有相关内容")
//		c.ServeJSON()
//		return
//	}
//}

// 获取评论列表
// @router /comment/list [get]
func (c *CommentController) List() {
	episodesId, _ := c.GetInt("episodesId")
	if episodesId == 0 {
		c.Data["json"] = ReturnError(4001, "必须指定视频ID")
		c.ServeJSON()
		return
	}

	limit, _ := c.GetInt("limit")
	offset, _ := c.GetInt("offset")
	if limit == 0 {
		limit = 12
	}

	nums, comments, err := models.GetCommentList(episodesId, limit, offset)
	if err == nil {
		var data []CommentInfo
		var commentInfo CommentInfo

		// 获取uid channel
		uidChan := make(chan string, 12)
		closeChan := make(chan bool, 5)
		resChan := make(chan models.UserInfo, 12)
		/// 把获取到的uid放到channel中
		go func() {
			for _, v := range comments {
				uidChan <- v.UserId
			}
			close(uidChan)
		}()

		/// 处理uidChannel中的信息
		for i := 0; i < 5; i++ {
			go channelGetUserInfo(uidChan, resChan, closeChan)
		}

		/// 判断是否执行完成，信息聚合
		go func() {
			for i := 0; i < 5; i++ {
				<-closeChan
			}
			close(resChan)
			close(closeChan)
		}()

		userInfoMap := make(map[string]models.UserInfo)
		for r := range resChan {
			userInfoMap[r.Uuid] = r
		}

		for _, v := range comments {
			commentInfo.Id = v.Id
			commentInfo.Content = v.Content
			commentInfo.AddTime = v.AddTime
			commentInfo.AddTimeTitle = DateFormat(v.AddTime)
			commentInfo.UserId = v.UserId
			commentInfo.Stamp = v.Stamp
			commentInfo.PraiseCount = v.PraiseCount
			// 获取用户信息
			commentInfo.UserInfo, _ = userInfoMap[v.UserId]
			data = append(data, commentInfo)
		}

		c.Data["json"] = ReturnSuccess(0, "success", data, nums)
		c.ServeJSON()
		return
	} else {
		c.Data["json"] = ReturnError(4004, "没有相关内容")
		c.ServeJSON()
		return
	}
}

func channelGetUserInfo(uidChan chan string, resChan chan models.UserInfo, closeChan chan bool) {
	for uid := range uidChan {
		res, err := models.RedisGetUserInfo(uid)
		if err == nil {
			resChan <- res
		}
	}
	closeChan <- true
}

//保存评论
// @router /comment/save [*]
func (c *CommentController) Save() {
	content := c.GetString("content")
	uid := c.GetString("uid")
	episodesId, _ := c.GetInt("episodesId")
	videoId, _ := c.GetInt("videoId")

	if content == "" {
		c.Data["json"] = ReturnError(4001, "内容不能为空")
		c.ServeJSON()
	}
	if uid == "" {
		c.Data["json"] = ReturnError(4002, "请先登录")
		c.ServeJSON()
	}
	if episodesId == 0 {
		c.Data["json"] = ReturnError(4003, "必须指定评论剧集ID")
		c.ServeJSON()
	}
	if videoId == 0 {
		c.Data["json"] = ReturnError(4005, "必须指定视频ID")
		c.ServeJSON()
	}
	err := models.SaveComment(content, uid, episodesId, videoId)
	if err == nil {
		c.Data["json"] = ReturnSuccess(0, "success", "", 1)
		c.ServeJSON()
	} else {
		c.Data["json"] = ReturnError(5000, err)
		c.ServeJSON()
	}
}
