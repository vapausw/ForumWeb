package controller

import (
	"ForumWeb/logic"
	"ForumWeb/model"
	"ForumWeb/pkg/jwt"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strings"
	"sync"
	"time"
)

// 全局map和锁，用于限流存储和同步
var (
	ipRateLimit = make(map[string]time.Time)
	lock        sync.Mutex
)

// LoginHandler 处理登录请求
func LoginHandler(c *gin.Context) {
	// 1.获取请求的参数，并且对其进行校验
	// 参数基本只有用户名以及密码
	var u model.User
	if err := c.ShouldBindJSON(&u); err != nil {
		zap.L().Error("invalid params", zap.Error(err))
		// 返回错误信息
		ResponseError(c, CodeInvalidParams)
		return
	}
	// 2. 业务处理
	if err := logic.Login(&u); err != nil {
		zap.L().Error("logic.Login(u) failed", zap.Error(err))
		ResponseErrorWithMsg(c, CodeInvalidPassword, err.Error())
		return
	}
	// 3. 返回响应
	// 生成Token
	aToken, rToken, _ := jwt.GenToken(uint64(u.UserID))
	ResponseSuccess(c, gin.H{
		"accessToken":  aToken,
		"refreshToken": rToken,
		"userID":       u.UserID,
		"username":     u.UserName,
	})
}

// SendCodeHandler 处理注册时发送验证码请求
func SendCodeHandler(c *gin.Context) {
	// 获取客户端IP地址
	clientIP := c.ClientIP()

	// 使用锁保证并发安全
	lock.Lock()
	if lastTime, ok := ipRateLimit[clientIP]; ok {
		// 如果在1分钟内
		if time.Since(lastTime) < time.Minute {
			lock.Unlock()
			ResponseError(c, CodeRateLimiting)
			return
		}
	}
	// 更新当前IP的时间戳
	ipRateLimit[clientIP] = time.Now()
	lock.Unlock()

	// 获取请求参数
	var re model.RegisterSend
	if err := c.ShouldBindJSON(&re); err != nil {
		zap.L().Error("invalid params", zap.Error(err))
		ResponseErrorWithMsg(c, CodeInvalidParams, err.Error())
		return
	}

	// 发送验证码
	if err := logic.SendCode(re.Email); err != nil {
		zap.L().Error("logic.SendCode(re.Email) failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	// 返回响应
	ResponseSuccess(c, nil)
}

// RegisterHandler 处理注册请求
func RegisterHandler(c *gin.Context) {
	// 1.获取请求参数 // 2.校验数据有效性
	var re model.RegisterForm
	if err := c.ShouldBindJSON(&re); err != nil {
		zap.L().Error("invalid params", zap.Error(err))
		ResponseErrorWithMsg(c, CodeInvalidParams, err.Error())
		return
	}

	zap.L().Info("register form", zap.Any("re", re))
	// 3.注册用户
	if err := logic.Register(&model.User{
		Email:    re.Email,
		UserName: re.UserName,
		Password: re.Password,
	}, re.Token); err != nil {
		zap.L().Error("logic.Register(&model.User{}) failed", zap.Error(err))
		if errors.Is(err, logic.ErrUserExist) {
			ResponseError(c, CodeUserExist)
			return
		}
		ResponseErrorWithMsg(c, CodeServerBusy, ErrToken.Error())
		return
	}
	// 注册成功后通过kafka异步发送欢迎邮件
	if err := logic.SendWelcome(re.Email, re.UserName); err != nil {
		zap.L().Error("logic.SendWelcome(re.Email) failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	// 4.返回响应
	ResponseSuccess(c, nil)
}

// RefreshTokenHandler 处理刷新token请求
func RefreshTokenHandler(c *gin.Context) {
	// 客户端携带Token有三种方式 1.放在请求头 2.放在请求体 3.放在URI
	// 这里假设Token放在Header的Authorization中，并使用Bearer开头
	// 这里的具体实现方式要依据你的实际业务情况决定
	rt := c.Query("refresh_token")
	zap.L().Info("refresh token", zap.String("rt", rt))
	authHeader := c.Request.Header.Get("Authorization")
	if authHeader == "" {
		ResponseErrorWithMsg(c, CodeInvalidToken, ErrFormat.Error())
		c.Abort()
		return
	}
	// 按空格分割
	// 按空格分割
	parts := strings.SplitN(authHeader, " ", 2)
	if !(len(parts) == 2 && parts[0] == "Bearer") {
		ResponseErrorWithMsg(c, CodeInvalidToken, ErrFormat.Error())
		c.Abort()
		return
	}
	aToken, rToken, err := jwt.RefreshToken(parts[1], rt)
	fmt.Println(err)
	c.JSON(http.StatusOK, gin.H{
		"access_token":  aToken,
		"refresh_token": rToken,
	})
}

// ResetPasswordHandler 处理重置密码请求
func ResetPasswordHandler(c *gin.Context) {
	var re = new(model.RegisterForm)
	if err := c.ShouldBindJSON(re); err != nil {
		zap.L().Error("invalid params", zap.Error(err))
		ResponseErrorWithMsg(c, CodeInvalidParams, err.Error())
		return
	}
	if err := logic.ResetPassword(re); err != nil {
		zap.L().Error("logic.ResetPassword failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, nil)
}
