package mysql

import "errors"

var (
	ErrServiceBusy    = errors.New("服务繁忙，请稍后再试")
	ErrUserExists     = errors.New("用户已存在")
	ErrUserNotFound   = errors.New("用户不存在")
	ErrInsertFailed   = errors.New("插入数据失败")
	ErrRecordNotFound = errors.New("社区分类未找到")
)
