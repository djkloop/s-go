package models

import (
	"encoding/json"
	"fuck_youku_api/services/mq"
	"time"

	"github.com/astaxie/beego/orm"
)

type Comment struct {
	Id          int
	Content     string
	AddTime     int64
	UserId      string
	Status      int
	Stamp       int
	PraiseCount int
	EpisodesId  int
	VideoId     int
}

func init() {
	orm.RegisterModel(new(Comment))
}

/**
 * 根据剧集数获取评论列表
 * @param int $episodesId
 * @param int $limit
 * @param int $offset
 * @return rs
 */
func GetCommentList(episodesId int, limit int, offset int) (int64, []Comment, error) {
	o := orm.NewOrm()
	var comments []Comment
	nums, _ := o.Raw("select * from comment where episodes_id=? and status=1", episodesId).QueryRows(&comments)
	_, err := o.Raw("select id, content, add_time, user_id, stamp, praise_count, episodes_id from comment where status=1 and episodes_id=? order by add_time desc limit ?,?", episodesId, limit, offset).QueryRows(&comments)
	return nums, comments, err
}

func SaveComment(content string, uid string, episodesId int, videoId int) error {
	o := orm.NewOrm()
	var comment Comment
	comment.Content = content
	comment.UserId = uid
	comment.EpisodesId = episodesId
	comment.VideoId = videoId
	comment.Stamp = 0
	comment.Status = 1
	comment.AddTime = time.Now().Unix()
	_, err := o.Insert(&comment)
	if err == nil {

		//修改视频的总评论数
		_, _ = o.Raw("UPDATE video SET comment=comment+1 WHERE id=?", videoId).Exec()
		//修改视频剧集的评论数
		_, _ = o.Raw("UPDATE video_episodes SET comment=comment+1 WHERE id=?", episodesId).Exec()

		//更新redis排行榜 - 通过MQ来实现
		//创建一个简单模式的MQ
		//把要传递的数据转换为json字符串
		videoObj := map[string]int{
			"VideoId": videoId,
		}
		videoJson, _ := json.Marshal(videoObj)
		_ = mq.Publish("", "fuck_top", string(videoJson))

		//延迟增加评论数
		videoCountObj := map[string]int{
			"VideoId":    videoId,
			"EpisodesId": episodesId,
		}
		videoCountJson, _ := json.Marshal(videoCountObj)
		_ = mq.PublishDlx("fuck.comment.count", string(videoCountJson))
	}
	return err
}
