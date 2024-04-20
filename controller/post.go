package controller

import (
	"ForumWeb/dao/redis"
	"ForumWeb/logic"
	"ForumWeb/model"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
)

// GetListHandler 获取帖子列表
func GetListHandler(c *gin.Context) {
	// 不用检查是否登录根据当前所处的页面去获取一定数量的帖子
	// 帖子内部可以进行排序，例如热度排序，时间排序等等，后续进行优化
	// 1.获取page, size
	order, _ := c.GetQuery("order")
	pageStr, ok := c.GetQuery("page")
	if !ok {
		pageStr = "1"
	}
	pageNum, err := strconv.ParseInt(pageStr, 10, 64)
	if err != nil {
		pageNum = 1
	}
	// 2.获取帖子列表
	posts := redis.GetPost(order, pageNum)
	zap.L().Info("redis.GetPost(order, pageNum)", zap.Any("posts", posts))
	// 3.返回响应
	ResponseSuccess(c, posts)
}

// CreatePostHandler 创建帖子
func CreatePostHandler(c *gin.Context) {
	// 1.获取请求信息，进行参数校验，使用反序列化进行
	var p model.Post
	if err := c.ShouldBindJSON(&p); err != nil {
		zap.L().Error("GetDetailHandler with invalid params", zap.Error(err))
		ResponseError(c, CodeInvalidParams)
		return
	}
	// 获取当前作者的ID
	userID, err := getCurrentUserID(c)
	if err != nil {
		zap.L().Error("GetCurrentUserID() failed", zap.Error(err))
		ResponseError(c, CodeNotLogin)
		return
	}
	p.AuthorId = userID
	// 2.业务处理,将帖子进行存储
	zap.L().Info("CreatePostHandler", zap.Any("p", p))
	if err := logic.CreatePost(&p); err != nil {
		zap.L().Error("logic.CreatePost failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	// 3.返回响应
	ResponseSuccess(c, nil)
}

// GetDetailHandler 获取帖子详情
func GetDetailHandler(c *gin.Context) {
	postId := c.Param("id")

	post, err := logic.GetPost(postId)
	if err != nil {
		zap.L().Error("logic.GetPost(postID) failed", zap.String("postId", postId), zap.Error(err))
	}

	ResponseSuccess(c, post)
}
