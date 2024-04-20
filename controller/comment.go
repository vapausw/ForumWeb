package controller

import (
	"ForumWeb/dao/mysql"
	"ForumWeb/model"
	"ForumWeb/pkg/snowflake"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// CommentHandler 评论的创建
func CommentHandler(c *gin.Context) {
	// 获取请求的参数，只需要保证评论内容不为空即可
	// 评论内容的判断可以加一些过滤的规则
	var comment model.Comment
	if err := c.ShouldBindJSON(&comment); err != nil {
		zap.L().Error("CommentHandler with invalid params", zap.Error(err))
		ResponseError(c, CodeInvalidParams)
		return
	}
	comment.CommentID = uint64(snowflake.GenID())
	// 获取当前用户的ID
	userID, err := getCurrentUserID(c)
	if err != nil {
		zap.L().Error("CommentHandler with getCurrentUserID failed", zap.Error(err))
		ResponseError(c, CodeNotLogin)
		return
	}
	comment.AuthorID = userID
	// 创建评论
	// 创建帖子
	if err := mysql.CreateComment(&comment); err != nil {
		zap.L().Error("mysql.CreatePost(&post) failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, nil)
}

// CommentListHandler 评论的查询
func CommentListHandler(c *gin.Context) {
	// 获取请求的参数
	ids, ok := c.GetQueryArray("ids")
	if !ok {
		zap.L().Error("CommentListHandler with invalid params")
		ResponseError(c, CodeInvalidParams)
		return
	}
	// 查询评论列表
	comments, err := mysql.GetCommentListByIDs(ids)
	if err != nil {
		zap.L().Error("CommentListHandler with GetCommentListByIDs failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, comments)
}
