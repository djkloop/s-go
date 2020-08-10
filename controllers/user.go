package controllers

import (
	"fmt"
	"fuck_youku_api/models"
	"regexp"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
)

// Operations about Users
type UserController struct {
	beego.Controller
}

// 用户注册功能
// @router /register/save [post]
func (c *UserController) RegisterUser() {
	var (
		mobile   string
		password string
		err      error
	)

	mobile = c.GetString("mobile")
	password = c.GetString("password")

	if mobile == "" {
		c.Data["json"] = ReturnError(4001, "手机号不能为空")
		c.ServeJSON()
		return
	}

	isPhoneError, _ := regexp.MatchString(`^(?:\+?86)?1(?:3\d{3}|5[^4\D]\d{2}|8\d{3}|7(?:[235-8]\d{2}|4(?:0\d|1[0-2]|9\d))|9[0-35-9]\d{2}|66\d{2})\d{6}$`, mobile)
	if !isPhoneError {
		c.Data["json"] = ReturnError(4002, "手机号格式不正确")
		c.ServeJSON()
		return
	}

	if password == "" {
		c.Data["json"] = ReturnError(4003, "密码不能为空")
		c.ServeJSON()
		return
	}

	/// 判断手机时候已经被注册过了
	isRegister := models.IsUserMobileRegister(mobile)
	if isRegister {
		c.Data["json"] = ReturnError(4005, "该手机号已经被注册了")
		c.ServeJSON()
		return
	} else {
		err = models.UserSave(mobile, MD5V(password))
		if err != nil {
			c.Data["json"] = ReturnError(5000, err)
		} else {
			c.Data["json"] = ReturnSuccess(0, "注册成功", nil, 0)
		}
		c.ServeJSON()
		return
	}
}

/// 用户登录接口
// @router /login/do [post]
func (c *UserController) UserLogin() {
	var (
		mobile   string
		password string
	)

	mobile = c.GetString("mobile")
	password = c.GetString("password")

	if mobile == "" {
		c.Data["json"] = ReturnError(4001, "手机号不能为空")
		c.ServeJSON()
		return
	}

	isPhoneError, _ := regexp.MatchString(`^(?:\+?86)?1(?:3\d{3}|5[^4\D]\d{2}|8\d{3}|7(?:[235-8]\d{2}|4(?:0\d|1[0-2]|9\d))|9[0-35-9]\d{2}|66\d{2})\d{6}$`, mobile)
	if !isPhoneError {
		c.Data["json"] = ReturnError(4002, "手机号格式不正确")
		c.ServeJSON()
		return
	}

	if password == "" {
		c.Data["json"] = ReturnError(4003, "密码不能为空")
		c.ServeJSON()
		return
	}

	uid, name := models.UserLogin(mobile, MD5V(password))
	if uid != "" {
		c.Data["json"] = ReturnSuccess(0, "登录成功", map[string]interface{}{"uid": uid, "username": name}, 1)
		c.ServeJSON()
		return
	} else {
		c.Data["json"] = ReturnError(4004, "手机号或密码不正确")
		c.ServeJSON()
		return
	}
}

//// 批量发送通知消息
//// @router /send/message [post]
//func (c *UserController) SendMessageDo() {
//	uids := c.GetString("uids")
//	if uids == "" {
//		c.Data["json"] = ReturnError(4001, "填写接收人")
//		c.ServeJSON()
//		return
//	}
//
//	content := c.GetString("content")
//	if content == "" {
//		c.Data["json"] = ReturnError(4001, "填写接收信息")
//		c.ServeJSON()
//		return
//	}
//
//	messageId, err := models.SendMessageDo(content)
//	if err == nil {
//		uidConfig := strings.Split(uids, ",")
//		for _, v := range uidConfig {
//			userId, _ := strconv.Atoi(v)
//			_ = models.SendMessageUserMq(userId, messageId)
//		}
//		c.Data["json"] = ReturnSuccess(0, "发送成功", "", 1)
//		c.ServeJSON()
//		return
//	} else {
//		c.Data["json"] = ReturnError(5000, "发送失败，请联系客服")
//		c.ServeJSON()
//		return
//	}
//}

type SendData struct {
	UserId    int
	MessageId int64
}

// 批量发送通知消息
// @router /send/message [post]
func (c *UserController) SendMessageDo() {
	uids := c.GetString("uids")
	if uids == "" {
		c.Data["json"] = ReturnError(4001, "填写接收人")
		c.ServeJSON()
		return
	}

	content := c.GetString("content")
	if content == "" {
		c.Data["json"] = ReturnError(4001, "填写接收信息")
		c.ServeJSON()
		return
	}

	messageId, err := models.SendMessageDo(content)
	if err == nil {
		uidConfig := strings.Split(uids, ",")
		count := len(uidConfig)

		sendChan := make(chan SendData, count)
		closeChan := make(chan bool, count)

		go func() {
			var data SendData
			for _, v := range uidConfig {
				userId, _ := strconv.Atoi(v)
				data.UserId = userId
				data.MessageId = messageId
				sendChan <- data
			}
			close(sendChan)
		}()

		for i := 0; i < 5; i++ {
			go sendMessageFunc(sendChan, closeChan)
		}

		for i := 0; i < 5; i++ {
			<-closeChan
		}
		close(closeChan)

		c.Data["json"] = ReturnSuccess(0, "发送成功", "", 1)
		c.ServeJSON()
		return
	} else {
		c.Data["json"] = ReturnError(5000, "发送失败，请联系客服")
		c.ServeJSON()
		return
	}
}

func sendMessageFunc(sendChan chan SendData, closeChan chan bool) {
	for t := range sendChan {
		fmt.Println(t)
		_ = models.SendMessageUserMq(t.UserId, t.MessageId)
	}
	closeChan <- true
}
