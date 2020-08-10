package models

import (
	"encoding/json"
	"time"

	"fuck_youku_api/services/mq"

	"github.com/astaxie/beego/orm"
)

type Message struct {
	Id      int
	Content string
	AddTime int64
}

type MessageUser struct {
	Id        int
	MessageId int64
	AddTime   int64
	Status    int
	UserId    int
}

func init() {
	orm.RegisterModel(new(Message), new(MessageUser))
}

// 发送消息入库
func SendMessageDo(content string) (int64, error) {
	o := orm.NewOrm()
	var message Message
	message.Content = content
	message.AddTime = time.Now().Unix()
	messageId, err := o.Insert(&message)
	return messageId, err
}

// 给某个用户发送消息
func SendMessageUser(userId int, messageId int64) error {
	o := orm.NewOrm()
	var messageUser MessageUser
	messageUser.UserId = userId
	messageUser.MessageId = messageId
	messageUser.Status = 1
	messageUser.AddTime = time.Now().Unix()
	_, err := o.Insert(&messageUser)
	return err
}

// 保存消息到mq队列中
func SendMessageUserMq(userId int, messageId int64) error {
	// 把数据转换成json字符串
	type Data struct {
		UserId    int
		MessageID int64
	}

	var data Data
	data.UserId = userId
	data.MessageID = messageId
	dataJSON, err := json.Marshal(data)
	_ = mq.Publish("", "fuck_send_message_user", string(dataJSON))
	return err
}
