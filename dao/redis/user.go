package redis

import "time"

func SetToken(email, token string) error {
	return rdb.Set(email, token, time.Minute*10).Err()
}

func GetToken(email string) string {
	return rdb.Get(email).Val()
}
