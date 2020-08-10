package models

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"fuck_youku_api/services/es"
	redisClient "fuck_youku_api/services/redis"

	"github.com/olivere/elastic/v7"

	"github.com/gomodule/redigo/redis"

	"github.com/astaxie/beego/orm"
)

type Video struct {
	Id                 int
	Title              string
	SubTitle           string
	Img                string
	Img1               string
	AddTime            int64
	EpisodesCount      int
	IsEnd              int
	IsHot              int
	IsRecommend        int
	ChannelId          int
	Status             int
	RegionId           int
	TypeId             int
	EpisodesUpdateTime int64
	Comment            int
	UserId             string
}

type VideoData struct {
	Id            int
	Title         string
	SubTitle      string
	Img           string
	Img1          string
	AddTime       int64
	EpisodesCount int
	IsEnd         int
	Comment       int
}

type Episodes struct {
	Id      int
	Title   string
	AddTime int64
	Num     int
	PlayUrl string
	Comment int
}

func init() {
	orm.RegisterModel(new(Video), new(Episodes))
}

// 根据频道ID获取正在热播视频
func GetChannelHotList(channelId int) (int64, []VideoData, error) {
	o := orm.NewOrm()
	var videos []VideoData
	num, err := o.Raw("select id, title, sub_title,img, img1, add_time,episodes_count,is_end from video where channel_id=? and is_hot=1 and status=1 order by episodes_update_time desc limit 9", channelId).QueryRows(&videos)
	return num, videos, err
}

// 根据频道下地区ID获取推荐视频
func GetChannelRegionRecommend(channelId int, regionId int) (int64, []VideoData, error) {
	o := orm.NewOrm()
	var videos []VideoData
	num, err := o.Raw("select id, title, sub_title, img, img1, add_time, episodes_count, is_end from video where status=1 and is_recommend=1 and channel_id=? and region_id=? order by episodes_update_time desc limit 9", channelId, regionId).QueryRows(&videos)
	return num, videos, err
}

// 根据频道下类型ID获取推荐视频
func GetChannelTypeRecommend(channelId int, typeId int) (int64, []VideoData, error) {
	o := orm.NewOrm()
	var videos []VideoData
	num, err := o.Raw("select id, title, sub_title, img, img1, add_time, episodes_count, is_end from video where channel_id=? and type_id=? and is_recommend=1 and status=1 order by episodes_update_time desc limit 9", channelId, typeId).QueryRows(&videos)
	return num, videos, err
}

/**
 * 频道下根据不同条件和排序方式获取视频信息
 * @param int $channelId
 * @param int $regionId
 * @param int $typeId
 * @param string $end
 * @param string $sort
 * @param int $offset
 * @param int $limit
 * @return rs
 */
func GetChannelVideoList(channelId int, regionId int, typeId int, end string, sort string, offset int, limit int) (int64, []orm.Params, error) {
	o := orm.NewOrm()
	var videos []orm.Params

	qs := o.QueryTable("video")
	qs = qs.Filter("channel_id", channelId)
	qs = qs.Filter("status", 1)

	if regionId > 0 {
		qs = qs.Filter("region_id", regionId)
	}
	if typeId > 0 {
		qs = qs.Filter("type_id", typeId)
	}
	if end == "n" {
		qs = qs.Filter("is_end", 0)
	} else if end == "y" {
		qs = qs.Filter("is_end", 1)
	}
	if sort == "episodesUpdateTime" {
		qs = qs.OrderBy("-episodes_update_time")
	} else if sort == "comment" {
		qs = qs.OrderBy("-comment")
	} else {
		qs = qs.OrderBy("-add_time")
	}
	var params = []string{
		"id",
		"title",
		"sub_title",
		"img",
		"img1",
		"add_time",
		"episodes_count",
		"is_end",
	}
	nums, _ := qs.Values(&videos, params...)
	qs.Limit(limit, offset)
	_, err := qs.Values(&videos, params...)
	return nums, videos, err
}

func GetChannelVideoListEs(channelId int, regionId int, typeId int, end string, sort string, offset int, limit int) (int64, []Video, error) {
	source := make(map[string]interface{})
	query := make(map[string]interface{})
	b := make(map[string]interface{})
	var must []map[string]interface{}
	must = append(must, map[string]interface{}{
		"term": map[string]interface{}{
			"channel_id": channelId,
		},
	})
	must = append(must, map[string]interface{}{
		"term": map[string]interface{}{
			"status": 1,
		},
	})
	if regionId > 0 {
		must = append(must, map[string]interface{}{
			"term": map[string]interface{}{
				"region_id": regionId,
			},
		})
	}
	if typeId > 0 {
		must = append(must, map[string]interface{}{
			"term": map[string]interface{}{
				"type_id": typeId,
			},
		})
	}
	if end == "n" {
		must = append(must, map[string]interface{}{
			"term": map[string]interface{}{
				"is_end": 0,
			},
		})
	} else if end == "y" {
		must = append(must, map[string]interface{}{
			"term": map[string]interface{}{
				"is_end": 1,
			},
		})
	}

	b["must"] = must
	query["bool"] = b
	var sortBy []elastic.Sorter
	sortBy = append(sortBy, elastic.NewFieldSort("add_time").Desc())
	if sort == "episodesUpdateTime" {
		sortBy = sortBy[0:]
		sortBy = append(sortBy, elastic.NewFieldSort("episodes_update_time").Desc())
	} else if sort == "comment" {
		sortBy = sortBy[0:]
		sortBy = append(sortBy, elastic.NewFieldSort("comment").Desc())
	}
	source["query"] = query
	res := es.Search("fuck_video", source, offset, limit, sortBy)
	var resData []Video
	for _, hit := range res.Hits.Hits {
		var data Video
		err := json.Unmarshal(hit.Source, &data)
		if err != nil {
			fmt.Println(err)
		}
		resData = append(resData, data)
	}

	return res.TotalHits(), resData, nil

}

/**
 * 获取视频详情
 * @param int $videoId
 * @return \think\response\Json
 */
func GetVideoInfo(videoId int) (Video, error) {
	o := orm.NewOrm()
	var video Video
	err := o.Raw("select * from video where id=? limit 1", videoId).QueryRow(&video)
	return video, err
}

/// 增加redis获取视频缓存
func RedisGetVideoInfo(videoId int) (Video, error) {
	var video Video
	conn := redisClient.PoolConnect()
	defer conn.Close()

	// 定义redis的key
	redisKey := "video:id:" + strconv.Itoa(videoId)
	// 判断redis是否存在缓存
	exists, err := redis.Bool(conn.Do("exists", redisKey))
	if exists {
		res, _ := redis.Values(conn.Do("hgetall", redisKey))
		err = redis.ScanStruct(res, &video)
	} else {
		o := orm.NewOrm()
		var video Video
		err := o.Raw("select * from video where id=? limit 1", videoId).QueryRow(&video)
		if err == nil {
			// 保存redis
			_, err := conn.Do("hmset", redis.Args{redisKey}.AddFlat(video)...)
			if err == nil {
				conn.Do("expire", redisKey, 86400)
			}
		}
	}
	return video, err
}

/**
 * 根据视频ID获取剧集列表
 * @param int $videoId
 * @return rs
 */
func GetVideoEpisodesList(videoId int) (int64, []Episodes, error) {
	o := orm.NewOrm()
	var episodes []Episodes
	nums, err := o.Raw("select id, title, add_time, num, play_url, comment from video_episodes where video_id=? and status=1 order by num asc", videoId).QueryRows(&episodes)
	return nums, episodes, err
}

/// 增加redis接口
func RedisGetVideoEpisodesList(videoId int) (int64, []Episodes, error) {
	var (
		episodes []Episodes
		num      int64
		err      error
	)
	conn := redisClient.PoolConnect()
	defer conn.Close()

	redisKey := "video:episodes:videoId:" + strconv.Itoa(videoId)
	// 判断rediskey是否存在
	exists, err := redis.Bool(conn.Do("exists", redisKey))
	if exists {
		num, err = redis.Int64(conn.Do("llen", redisKey))
		if err == nil {
			values, _ := redis.Values(conn.Do("lrange", redisKey, "0", "-1"))
			var episodesInfo Episodes
			for _, v := range values {
				err = json.Unmarshal(v.([]byte), &episodesInfo)
				if err == nil {
					episodes = append(episodes, episodesInfo)
				}
			}
		}
	} else {
		num, episodes, err = GetVideoEpisodesList(videoId)
		if err == nil {
			for _, v := range episodes {
				jsonValue, err := json.Marshal(v)
				if err == nil {
					conn.Do("rpush", redisKey, jsonValue)
				}
			}
			conn.Do("expire", redisKey, 86400)
		}
	}
	return num, episodes, err
}

/**
 * 根据频道ID获取排行榜
 * @param int $channelId
 * @return \think\response\Json
 */
func GetChannelTop(channelId int) (int64, []VideoData, error) {
	o := orm.NewOrm()
	var videos []VideoData
	num, err := o.Raw("select id,title,sub_title,img,img1,add_time,episodes_count,is_end from video where channel_id=? and status=1 order by comment desc limit 10", channelId).QueryRows(&videos)
	return num, videos, err
}

/// 增加redis缓存 - 频道排行榜
func RedisGetChannelTop(channelId int) (int64, []VideoData, error) {
	var (
		videos []VideoData
		num    int64
		err    error
	)
	conn := redisClient.PoolConnect()
	defer conn.Close()

	redisKey := "video:top:channel:channelId" + strconv.Itoa(channelId)
	// 判断是否存在
	exists, err := redis.Bool(conn.Do("exists", redisKey))
	if exists {
		num = 0
		res, _ := redis.Values(conn.Do("zrevrange", redisKey, "0", "10", "WITHSCORES"))
		for k, v := range res {
			fmt.Println(string(v.([]byte)))
			if k%2 == 0 {
				videoId, err := strconv.Atoi(string(v.([]byte)))
				videoInfo, err := RedisGetVideoInfo(videoId)
				if err == nil {
					var videoDataInfo VideoData
					videoDataInfo.Id = videoInfo.Id
					videoDataInfo.Title = videoInfo.Title
					videoDataInfo.SubTitle = videoInfo.SubTitle
					videoDataInfo.Img = videoInfo.Img
					videoDataInfo.Img1 = videoInfo.Img1
					videoDataInfo.AddTime = videoInfo.AddTime
					videoDataInfo.EpisodesCount = videoInfo.EpisodesCount
					videoDataInfo.Comment = videoInfo.Comment
					videos = append(videos, videoDataInfo)
					num++
				}
			}
		}
	} else {
		o := orm.NewOrm()
		num, err = o.Raw("select id,title,sub_title,img,img1,add_time,episodes_count,is_end from video where channel_id=? and status=1 order by comment desc limit 10", channelId).QueryRows(&videos)
		if err == nil {
			for _, v := range videos {
				conn.Do("zadd", redisKey, v.Comment, v.Id)
			}
			conn.Do("expire", redisKey, 86400*30)
		}
	}
	return num, videos, err
}

/**
 * 获取类型下排行榜
 * @param int $typeId
 * @return rs
 */
func GetTypeTop(typeId int) (int64, []VideoData, error) {
	o := orm.NewOrm()
	var videos []VideoData
	num, err := o.Raw("select id,title,sub_title,img,img1,add_time,episodes_count,is_end from video where type_id=? and status=1 order by comment desc limit 10", typeId).QueryRows(&videos)
	return num, videos, err
}

// 增加redis类型 - 类型排行榜
func RedisGetTypeTop(typeId int) (int64, []VideoData, error) {
	var (
		videos []VideoData
		num    int64
	)
	conn := redisClient.PoolConnect()
	defer conn.Close()

	redisKey := "video:top:type:typeId:" + strconv.Itoa(typeId)
	exists, err := redis.Bool(conn.Do("exists", redisKey))
	if exists {
		num = 0
		res, _ := redis.Values(conn.Do("zrevrange", redisKey, "0", "10", "WITHSCORES"))
		for k, v := range res {
			if k%2 == 0 {
				videoId, err := strconv.Atoi(string(v.([]byte)))
				videoInfo, err := RedisGetVideoInfo(videoId)
				if err == nil {
					var videoDataInfo VideoData
					videoDataInfo.Id = videoInfo.Id
					videoDataInfo.Title = videoInfo.Title
					videoDataInfo.SubTitle = videoInfo.SubTitle
					videoDataInfo.Img = videoInfo.Img
					videoDataInfo.Img1 = videoInfo.Img1
					videoDataInfo.AddTime = videoInfo.AddTime
					videoDataInfo.EpisodesCount = videoInfo.EpisodesCount
					videoDataInfo.Comment = videoInfo.Comment
					videos = append(videos, videoDataInfo)
					num++
				}
			}
		}
	} else {
		o := orm.NewOrm()
		num, err = o.Raw("select id,title,sub_title,img,img1,add_time,episodes_count,is_end from video where type_id=? and status=1 order by comment desc limit 10", typeId).QueryRows(&videos)
		if err == nil {
			for _, v := range videos {
				conn.Do("zadd", redisKey, v.Comment, v.Id)
			}
			conn.Do("expire", redisKey, 86400*30)
		}
	}
	return num, videos, err
}

//
func GetUserVideo(uid string) (int64, []VideoData, error) {
	o := orm.NewOrm()
	var videos []VideoData
	num, err := o.Raw("select id, title, sub_title, img, img1, episodes_count, is_end from video where user_id=? order by add_time desc", uid).QueryRows(&videos)
	return num, videos, err
}

func MVideoSave(title string, subTitle string, channelId int, regionId int, typeId int, uid string, playUrl string) error {
	o := orm.NewOrm()
	var video Video
	nowTime := time.Now().Unix()
	video.Title = title
	video.SubTitle = subTitle
	video.Status = 1
	video.AddTime = nowTime
	video.Img = ""
	video.Img1 = ""
	video.ChannelId = channelId
	video.TypeId = typeId
	video.RegionId = regionId
	video.UserId = uid
	video.EpisodesCount = 1
	video.EpisodesUpdateTime = nowTime
	video.IsEnd = 1
	video.Comment = 0
	videoId, err := o.Insert(&video)
	if err == nil {
		_, _ = o.Raw("insert into video_episodes (title, add_time, num, video_id, play_url, status, comment) values (?,?,?,?,?,?,?)", subTitle, nowTime, 1, videoId, playUrl, 1, 0).Exec()
	}
	return err
}

// 获取所有视频数据
func GetAllList() (int64, []Video, error) {
	o := orm.NewOrm()
	var video []Video
	num, err := o.Raw("select * from video").QueryRows(&video)
	return num, video, err
}
