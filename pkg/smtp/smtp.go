package smtp

import (
	"ForumWeb/setting"
	"crypto/tls"
	"go.uber.org/zap"
	"net/smtp"
	"strings"
)

// SendEmail 使用SMTP发送邮件
func SendEmail(recipient string, message []byte) (err error) {
	from := setting.Conf.MyEmailConfig.Email
	password := setting.Conf.MyEmailConfig.Password
	// SMTP服务器地址和端口
	smtpHost := setting.Conf.MyEmailConfig.SmtpHost
	smtpPort := setting.Conf.MyEmailConfig.SmtpPort
	// 创建TLS配置
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true, // 或者设置为false，并提供正确的证书链
		ServerName:         smtpHost,
	}
	// 连接到SMTP服务器
	conn, err := tls.Dial("tcp", strings.Join([]string{smtpHost, ":", smtpPort}, ""), tlsconfig)
	if err != nil {
		zap.L().Error("tls.Dial failed", zap.Error(err))
		return
	}
	defer conn.Close()
	// 创建smtp客户端
	client, err := smtp.NewClient(conn, smtpHost)
	if err != nil {
		zap.L().Error("smtp.NewClient failed", zap.Error(err))
		return
	}
	defer client.Close()

	// 认证
	auth := smtp.PlainAuth("", from, password, smtpHost)
	if err = client.Auth(auth); err != nil {
		zap.L().Error("client.Auth failed", zap.Error(err))
		return
	}

	// 设置发送者和接收者
	if err = client.Mail(from); err != nil {
		zap.L().Error("client.Mail failed", zap.Error(err))
		return
	}
	if err = client.Rcpt(recipient); err != nil {
		zap.L().Error("client.Rcpt failed", zap.Error(err))
		return
	}

	// 发送邮件正文
	w, err := client.Data()
	if err != nil {
		zap.L().Error("client.Data failed", zap.Error(err))
		return
	}
	_, err = w.Write(message)
	if err != nil {
		zap.L().Error("w.Write failed", zap.Error(err))
		return
	}
	err = w.Close()
	if err != nil {
		zap.L().Error("w.Close failed", zap.Error(err))
		return
	}
	return
}
