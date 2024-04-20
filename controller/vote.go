package controller

import (
	"ForumWeb/dao/redis"
	"ForumWeb/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// VoteHandler 帖子的投票
func VoteHandler(c *gin.Context) {
	// 给哪个文章投什么票
	// 如果没有投过票，则投票
	// 如果投过票了，判断和之前的投票方向是否一致，不一致则更改
	// 如果之前的投票方向和现在的一致，则取消此次投票
	var vote model.VoteData
	if err := c.ShouldBindJSON(&vote); err != nil {
		zap.L().Info("VoteHandler", zap.Error(err))
		ResponseErrorWithMsg(c, CodeInvalidParams, err.Error())
		return
	}
	// 从请求中获取当前用户的ID
	userID, err := getCurrentUserID(c)
	if err != nil {
		zap.L().Error("getCurrentUserID() failed", zap.Error(err))
		ResponseError(c, CodeNotLogin)
		return
	}
	// 投票
	if err := redis.PostVote(vote.PostID, fmt.Sprint(userID), vote.Direction); err != nil {
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, nil)
}
