package logic

import "errors"

var (
	ErrInvalidParams   = errors.New("参数错误")
	ErrInvalidPassword = errors.New("密码错误")
	ErrUserExist       = errors.New("用户已存在")
	ErrorCaptcha       = errors.New("验证码错误")
)
