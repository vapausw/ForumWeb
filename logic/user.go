package logic

import (
	"ForumWeb/dao/mysql"
	"ForumWeb/dao/redis"
	"ForumWeb/model"
	"ForumWeb/pkg/bcrypt"
	"ForumWeb/pkg/kafka"
	"ForumWeb/pkg/smtp"
	"ForumWeb/pkg/snowflake"
	"crypto/rand"
	"encoding/base32"
	"encoding/json"
	"go.uber.org/zap"
	"strings"
)

func Login(u *model.User) error {
	if len(u.UserName) == 0 || len(u.Password) == 0 {
		return ErrInvalidParams
	}
	// 从数据库中查找该用户
	v, err := mysql.GetUserByUserName(u.UserName)
	if err != nil {
		return err
	}
	// 判断密码是否正确
	if !bcrypt.Compare(v.Password, u.Password) {
		return ErrInvalidPassword
	}
	u.UserID = v.UserID
	return nil
}

func Register(u *model.User, Token string) error {
	// 从数据库中查找该用户
	_, err := mysql.GetUserByUserName(u.UserName)
	if err == nil {
		return ErrUserExist
	}
	// 从redis中获取token
	token := redis.GetToken(u.Email)
	zap.L().Info("redis.GetToken", zap.String("token", token))
	zap.L().Info("token", zap.String("Token", Token))
	if token != Token {
		return ErrorCaptcha
	}
	// 注册用户
	u.UserID = snowflake.GenID()
	u.Password = bcrypt.Encrypt(u.Password)
	if err := mysql.InsertUser(u); err != nil {
		zap.L().Error("mysql.InsertUser failed", zap.Error(err))
		return err
	}
	return nil
}

// GenerateSecureToken 生成一个安全的随机令牌，长度为 16 位，由大小写字母和数字组成
func GenerateSecureToken() string {
	// 定义令牌的字节长度，base32 编码每5个比特表示一个字符，因此对于8个字符，我们需要5*16/8=10个字节
	tokenLength := 10
	b := make([]byte, tokenLength)
	_, err := rand.Read(b)
	if err != nil {
		zap.L().Error("Failed to generate random token", zap.Error(err))
		return ""
	}
	// 使用 base32 编码生成令牌，然后取前16位作为结果
	// 注意：这里使用了 base32，因为它比 base64 更容易产生人类可读的字符（大小写字母和数字）
	// 但由于输出会更长，我们只取前16个字符
	token := base32.StdEncoding.EncodeToString(b)[:16]
	return token
}

func SendCode(email string) error {
	// 生成验证码
	token := GenerateSecureToken()
	// 设置邮件内容
	// 初始化邮件头部
	header := emailHeader(email)
	// 邮件正文，使用HTML格式
	body := tokenBody(token)
	// 组合邮件头部和正文
	message := []byte(strings.Join([]string{header, body}, ""))
	// 发送邮件
	err := smtp.SendEmail(email, message)
	if err != nil {
		zap.L().Error("smtp.SendEmail failed", zap.Error(err))
		return err
	}
	// 没有错误，将token加一个过期时间存入redis
	if err = redis.SetToken(email, token); err != nil {
		zap.L().Error("redis.SetToken failed", zap.Error(err))
		return err
	}
	return nil
}

func ResetPassword(re *model.RegisterForm) error {
	// 1.从redis中获取token
	token := redis.GetToken(re.Email)
	if token != re.Token {
		return ErrorCaptcha
	}
	// 2.更新用户密码
	re.Password = bcrypt.Encrypt(re.Password)
	if err := mysql.UpdatePassword(re); err != nil {
		zap.L().Error("mysql.UpdatePassword failed", zap.Error(err))
		return err
	}
	return nil
}

func SendWelcome(email, username string) error {
	// 设置邮件内容
	// 初始化邮件头部
	header := emailHeader(email)
	// 邮件正文，使用HTML格式
	body := welcomeBody(username)
	// 组合邮件头部和正文
	message := []byte(strings.Join([]string{header, body}, ""))
	//通过kafka异步发送欢迎邮件
	// 发送消息到 Kafka
	// 通过 Kafka 异步发送欢迎邮件
	msg := model.WelcomeEmailMessage{
		Email:   email,
		Message: string(message),
	}
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		zap.L().Error("json.Marshal failed", zap.Error(err))
		return err
	}
	kafka.StartEmailProducer(msgBytes)
	return nil
}
