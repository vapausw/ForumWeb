package bcrypt

import "golang.org/x/crypto/bcrypt"

// Encrypt 密码加密
func Encrypt(s string) string {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	return string(hashedPassword)
}

// Compare 密码比较第一个参数是数据库中的密码，第二个参数是用户输入的密码
func Compare(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
