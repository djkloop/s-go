package controllers

import (
	"crypto/md5"
	"encoding/hex"
	_ "encoding/json"
	"time"

	"github.com/astaxie/beego"
)

type CommonController struct {
	beego.Controller
}

type JsonStruct struct {
	Code  int         `json:"code"`
	Msg   interface{} `json:"msg"`
	Items interface{} `json:"items"`
	Count int64       `json:"count"`
}

func ReturnSuccess(code int, msg interface{}, items interface{}, count int64) *JsonStruct {
	var json *JsonStruct
	json = &JsonStruct{
		Code:  code,
		Msg:   msg,
		Items: items,
		Count: count,
	}
	return json
}

func ReturnError(code int, msg interface{}) *JsonStruct {
	var json *JsonStruct
	json = &JsonStruct{
		Code: code,
		Msg:  msg,
	}
	return json
}

/// 用户密码加密
func MD5V(password string) string {
	h := md5.New()
	h.Write([]byte(password + beego.AppConfig.String("md5code")))
	return hex.EncodeToString(h.Sum(nil))
}

// 格式化时间
func DateFormat(t int64) string {
	videoTime := time.Unix(t, 0)
	return videoTime.Format("2006-01-02")
}
