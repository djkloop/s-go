package controllers

import (
	"encoding/json"
	"fuck_youku_api/models"
	"net/http"

	"github.com/astaxie/beego/validation"

	"github.com/astaxie/beego"
	"github.com/gorilla/websocket"
)

type BarrageController struct {
	beego.Controller
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WsData struct {
	CurrentTime int
	EpisodesId  int
}

// 获取弹幕websocket
// @router /barrage/ws [*]
func (c *BarrageController) BarrageWs() {
	var (
		conn     *websocket.Conn
		err      error
		data     []byte
		barrages []models.BarrageData
	)

	if conn, err = upgrader.Upgrade(c.Ctx.ResponseWriter, c.Ctx.Request, nil); err != nil {
		goto ERR
	}

	for {
		if _, data, err = conn.ReadMessage(); err != nil {
			goto ERR
		}
		var wsData WsData
		_ = json.Unmarshal([]byte(data), &wsData)
		endTime := wsData.CurrentTime + 60
		// 获取弹幕数据
		_, barrages, err = models.BarrageList(wsData.EpisodesId, wsData.CurrentTime, endTime)
		if err == nil {
			if err := conn.WriteJSON(barrages); err != nil {
				goto ERR
			}
		}
	}

ERR:
	conn.Close()
}

// 保存弹幕
// @router /barrage/save [post]
func (c *BarrageController) Save() {
	barrageInfo := models.BarrageInfo{}
	if err := c.ParseForm(&barrageInfo); err != nil {
		c.Data["json"] = ReturnError(5000, err)
		c.ServeJSON()
		return
	} else {
		valid := validation.Validation{}
		b, err := valid.Valid(&barrageInfo)
		if err != nil {
			c.Data["json"] = ReturnError(5000, err)
			c.ServeJSON()
			return
		}

		if !b {
			// 验证没通过 输出错误信息
			for _, err := range valid.Errors {
				c.Data["json"] = ReturnError(4001, err.Message)
				c.ServeJSON()
			}
		}
	}

	err := models.SaveBarrage(barrageInfo)
	if err == nil {
		c.Data["json"] = ReturnSuccess(0, "", "", 1)
		c.ServeJSON()
		return
	} else {
		c.Data["json"] = ReturnError(5000, err)
		c.ServeJSON()
		return
	}
}
