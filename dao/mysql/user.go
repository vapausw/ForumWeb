package mysql

import (
	"ForumWeb/model"
	"database/sql"
	"errors"
	"github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
)

func GetUserByUserName(username string) (*model.User, error) {
	sqlStr := `SELECT user_id, username, password FROM user WHERE username = ?`
	user := new(model.User)
	err := db.QueryRow(sqlStr, username).Scan(&user.UserID, &user.UserName, &user.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			zap.L().Error("No user found", zap.String("username", username), zap.Error(err))
			return nil, ErrUserNotFound
		}
		zap.L().Error("db.QueryRow failed", zap.String("username", username), zap.Error(err))
		return nil, ErrServiceBusy
	}
	return user, nil
}

func InsertUser(user *model.User) error {
	sqlStr := `INSERT INTO user(user_id, username, password, email) VALUES (?, ?, ?, ?)`
	_, err := db.Exec(sqlStr, user.UserID, user.UserName, user.Password, user.Email)
	if err != nil {
		zap.L().Error("Insert user failed", zap.String("username", user.UserName), zap.Error(err))
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			return ErrUserExists
		}
		return ErrServiceBusy
	}
	return nil
}

func GetUserByID(uid uint64) (*model.User, error) {
	sqlStr := `SELECT user_id, username, password FROM user WHERE user_id = ?`
	user := new(model.User)
	err := db.QueryRow(sqlStr, uid).Scan(&user.UserID, &user.UserName, &user.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			zap.L().Error("No user found", zap.Uint64("user_id", uid), zap.Error(err))
			return nil, ErrUserNotFound
		}
		zap.L().Error("db.QueryRow failed", zap.Uint64("user_id", uid), zap.Error(err))
		return nil, ErrServiceBusy
	}
	return user, nil
}

func UpdatePassword(re *model.RegisterForm) (err error) {
	// 准备更新密码的SQL语句
	sqlStr := `UPDATE user SET password = ? WHERE email = ?`
	// 执行更新操作
	_, err = db.Exec(sqlStr, re.Password, re.Email)
	if err != nil {
		zap.L().Error("Failed to update password", zap.Error(err))
		return err
	}
	return
}
