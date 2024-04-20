package logic

import (
	"ForumWeb/setting"
	"strings"
)

func emailHeader(email string) string {
	from := setting.Conf.MyEmailConfig.Email
	var header strings.Builder
	header.WriteString(strings.Join([]string{"From: ", from, "\r\n"}, " "))
	header.WriteString(strings.Join([]string{"To: ", email, "\r\n"}, " "))
	header.WriteString("Subject: Verify Your Email\r\n")
	header.WriteString("MIME-Version: 1.0\r\n")
	header.WriteString("Content-Type: text/html; charset=UTF-8\r\n\r\n")
	return header.String()
}

func tokenBody(token string) string {
	var body strings.Builder
	// 写入HTML结构
	body.WriteString("<html><body>")
	body.WriteString("<p>Please use the following token to complete your registration:</p>")
	body.WriteString(strings.Join([]string{"<p><strong>", token, "</strong></p>"}, " "))
	body.WriteString("<p>This token will expire in <strong>10 minutes</strong>.</p>")
	body.WriteString("</body></html>")
	return body.String()
}

func welcomeBody(username string) string {
	var body strings.Builder
	// 写入HTML结构
	body.WriteString("<html><body>")
	body.WriteString(strings.Join([]string{"<p><strong>Hello ", username, ": </strong></p>"}, ""))
	body.WriteString(strings.Join([]string{"<p><h1>Hello, welcome to register for the forum, ",
		"please pay attention to abide by the rules of the forum, ",
		"I wish you a happy use</h1></p>"}, ""))

	body.WriteString("<p>This token will expire in <strong>10 minutes</strong>.</p>")
	body.WriteString("</body></html>")
	return body.String()
}
