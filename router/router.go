package router

import (
	"ForumWeb/controller"
	"ForumWeb/log"
	"ForumWeb/pkg/ws"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Init(mode string) *gin.Engine {
	if mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()
	r.Use(log.GinLogger(), log.GinRecovery(true))
	r.Static("/static", "./static")
	// 模板解析
	r.LoadHTMLGlob("templates/*")
	r.GET("/ws_test", func(c *gin.Context) {
		c.HTML(http.StatusOK, "ws_test.html", nil)
	})
	v1 := r.Group("/api/v1")
	// 测试接口连通性
	v1.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})
	v1.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})
	// 为WebSocket请求创建一个新的路由
	v1.GET("/ws", ws.WebsocketHandler) // WebSocket 连接处理器
	// 登录以及注册请求
	v1.POST("/login", controller.LoginHandler)
	v1.POST("/send_code", controller.SendCodeHandler) // 发送验证码的路由
	v1.POST("/register", controller.RegisterHandler)
	v1.POST("/refresh_token", controller.RefreshTokenHandler)
	// 密码找回请求
	v1.POST("/reset_password", controller.ResetPasswordHandler) // 重置密码
	v1.GET("/post", controller.GetListHandler)
	v1.GET("/post/:id", controller.GetDetailHandler)
	// 帖子的分类
	v1.GET("/community", controller.CommunityHandler)
	v1.GET("/community/:id", controller.CommunityDetailHandler)
	v1.GET("/comment", controller.CommentListHandler)
	v1.Use(controller.JWTAuthMiddleware())
	{
		// 帖子的创建以及查询操作
		v1.POST("/post", controller.CreatePostHandler)
		// 帖子的投票
		v1.POST("/vote", controller.VoteHandler)
		// 评论的创建以及查询操作
		v1.POST("/comment", controller.CommentHandler)
	}
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "404 Not Found",
		})
	})
	return r
}
