package models

import (
	redisClient "fuck_youku_api/services/redis"
	"time"

	"github.com/garyburd/redigo/redis"

	"github.com/astaxie/beego/orm"
	"github.com/bwmarrin/snowflake"
)

type User struct {
	Id       int
	Uuid     string
	Name     string
	Password string
	Mobile   string
	Avatar   string
	Status   int
	AddTime  int64
}

type UserInfo struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	AddTime int64  `json:"addTime"`
	Avatar  string `json:"avatar"`
	Uuid    string `json:"uuid"`
}

func init() {
	orm.RegisterModel(new(User))
}

/// 根据手机号判断用户是否存在
func IsUserMobileRegister(mobile string) bool {
	o := orm.NewOrm()
	user := User{
		Mobile: mobile,
	}
	err := o.Read(&user, "Mobile")
	if err == orm.ErrNoRows {
		return false
	} else {
		return true
	}
}

/// 保存用户
func UserSave(mobile string, password string) error {
	var (
		user User
		err  error
	)
	node, _ := snowflake.NewNode(1)

	id := node.Generate().Base64()

	o := orm.NewOrm()
	user.Name = ""
	user.Uuid = id
	user.Password = password
	user.Mobile = mobile
	user.Status = 1
	user.AddTime = time.Now().Unix()
	_, err = o.Insert(&user)
	return err
}

/// 登录功能
func UserLogin(mobile string, password string) (uid string, name string) {
	o := orm.NewOrm()
	var (
		user User
		err  error
	)
	err = o.QueryTable("user").Filter("mobile", mobile).Filter("password", password).One(&user)
	if err == orm.ErrNoRows {
		return "", ""
	} else if err == orm.ErrMissPK {
		return "", ""
	}
	return user.Uuid, user.Name
}

// 根据用户id获取信息
func GetUserInfo(uid string) (UserInfo, error) {
	o := orm.NewOrm()
	var user UserInfo
	err := o.Raw("select id, name, add_time, avatar, uuid from user where uuid=? limit 1", uid).QueryRow(&user)
	return user, err
}

func RedisGetUserInfo(uid string) (UserInfo, error) {
	var (
		user UserInfo
		err  error
	)
	conn := redisClient.PoolConnect()
	defer conn.Close()

	redisKey := "user:id:" + uid
	// 判断redis是否存在
	exists, err := redis.Bool(conn.Do("exists", redisKey))
	if exists {
		res, _ := redis.Values(conn.Do("hgetall", redisKey))
		err = redis.ScanStruct(res, &user)
	} else {
		user, err = GetUserInfo(uid)
		if err == nil {
			_, err := conn.Do("hmset", redis.Args{redisKey}.AddFlat(user)...)
			if err == nil {
				conn.Do("expire", redisKey, 86400)
			}
		}
	}
	return user, err
}
