package controller

import (
	"ForumWeb/dao/mysql"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// CommunityHandler 获取所有社区的信息
func CommunityHandler(c *gin.Context) {
	communityList, err := mysql.GetCommunityList()
	if err != nil {
		zap.L().Error("mysql.GetCommunityList() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, communityList)
}

// CommunityDetailHandler 获取某一个社区的详细信息
func CommunityDetailHandler(c *gin.Context) {
	communityID := c.Param("id")
	communityList, err := mysql.GetCommunityByID(communityID)
	if err != nil {
		zap.L().Error("mysql.GetCommunityByID() failed", zap.Error(err))
		ResponseErrorWithMsg(c, CodeServerBusy, err.Error())
		return
	}
	ResponseSuccess(c, communityList)
}
