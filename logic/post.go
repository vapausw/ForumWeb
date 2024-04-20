package logic

import (
	"ForumWeb/dao/mysql"
	"ForumWeb/dao/redis"
	"ForumWeb/model"
	"ForumWeb/pkg/snowflake"
	"fmt"
	"go.uber.org/zap"
	"strconv"
)

func CreatePost(p *model.Post) (err error) {
	// 1.生成帖子ID
	postID := snowflake.GenID()
	p.PostID = uint64(postID)
	// 2.存储帖子到数据库
	// 创建帖子
	if err := mysql.CreatePost(p); err != nil {
		zap.L().Error("mysql.CreatePost(&post) failed", zap.Error(err))
		return err
	}
	// 3.根据帖子的社区ID存储到redis
	community, err := mysql.GetCommunityByID(strconv.Itoa(int(p.CommunityID)))
	if err != nil {
		zap.L().Error("mysql.GetCommunityByID failed", zap.Error(err))
		return err
	}
	if err := redis.CreatePost(
		fmt.Sprint(p.PostID),
		fmt.Sprint(p.AuthorId),
		p.Title,
		TruncateByWords(p.Content, 120),
		community.CommunityName); err != nil {
		zap.L().Error("redis.CreatePost failed", zap.Error(err))
		return err
	}
	return
}

func GetPost(postID string) (data interface{}, err error) {
	// 1.查询帖子详情
	post, err := mysql.GetPostByID(postID)
	if err != nil {
		zap.L().Error("mysql.GetPostByID(postID) failed", zap.String("post_id", postID), zap.Error(err))
		return nil, err
	}
	// 获取作者信息
	user, err := mysql.GetUserByID(post.AuthorId)
	if err != nil {
		zap.L().Error("mysql.GetUserByID() failed", zap.String("author_id", fmt.Sprint(post.AuthorId)), zap.Error(err))
		return
	}
	post.AuthorName = user.UserName
	community, err := mysql.GetCommunityByID(fmt.Sprint(post.CommunityID))
	if err != nil {
		zap.L().Error("mysql.GetCommunityByID() failed", zap.String("community_id", fmt.Sprint(post.CommunityID)), zap.Error(err))
		return
	}
	post.CommunityName = community.CommunityName
	return post, nil
}
