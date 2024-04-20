package jwt

import (
	"errors"
	"github.com/golang-jwt/jwt"
	"go.uber.org/zap"
	"time"
)

// MyClaims 自定义声明结构体并内嵌jwt.StandardClaims
// jwt包自带的jwt.StandardClaims只包含了官方字段
// 我们这里需要额外记录一个UserID字段，所以要自定义结构体
// 如果想要保存更多信息，都可以添加到这个结构体中
type MyClaims struct {
	UserID uint64 `json:"user_id"`
	jwt.StandardClaims
}

var mySecret = []byte("VvsdD/?65VSAvdsdcc,.,.64cas")

func keyFunc(_ *jwt.Token) (i interface{}, err error) {
	return mySecret, nil
}

const TokenExpireDuration = time.Minute * 10

// GenToken 生成access token 和 refresh token
func GenToken(userID uint64) (aToken, rToken string, err error) {
	// 创建一个我们自己的声明
	c := MyClaims{
		userID, // 自定义字段
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(TokenExpireDuration).Unix(), // 过期时间
			Issuer:    "bluebell",                                 // 签发人
		},
	}
	// 加密并获得完整的编码后的字符串token
	aToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, c).SignedString(mySecret)

	// refresh token 不需要存任何自定义数据
	rToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Hour * 24 * 7).Unix(), // 过期时间
		Issuer:    "bluebell",                                // 签发人
	}).SignedString(mySecret)
	// 使用指定的secret签名并获得完整的编码后的字符串token
	return
}

// ParseToken 解析JWT
func ParseToken(tokenString string) (claims *MyClaims, err error) {
	// 解析token
	var token *jwt.Token
	claims = new(MyClaims)
	token, err = jwt.ParseWithClaims(tokenString, claims, keyFunc)
	if err != nil {
		return
	}
	if !token.Valid { // 校验token
		err = errors.New("invalid token")
	}
	return
}

// RefreshToken 刷新AccessToken
func RefreshToken(aToken, rToken string) (newAToken, newRToken string, err error) {
	// 验证刷新令牌的有效性
	if _, err = jwt.Parse(rToken, keyFunc); err != nil {
		zap.L().Error("jwt.Parse(rToken, keyFunc) failed", zap.Error(err))
		return "", "", err
	}

	// 从旧的访问令牌中解析出claims数据
	var claims MyClaims
	zap.L().Info("Parsing aToken", zap.String("aToken", aToken))
	_, err = jwt.ParseWithClaims(aToken, &claims, keyFunc)

	if err != nil {
		var v *jwt.ValidationError
		if errors.As(err, &v) {
			// 检查错误是否为令牌过期
			if v.Errors&jwt.ValidationErrorExpired != 0 {
				// 生成新的访问令牌和刷新令牌
				zap.L().Info("Access token is expired, generating new tokens...")
				return GenToken(claims.UserID)
			}
		}
		zap.L().Error("Error parsing aToken with claims", zap.Error(err))
		return "", "", err
	}

	zap.L().Info("Access token parsed and claims verified successfully")
	// 如果访问令牌没有问题，则不需要更新令牌
	return aToken, rToken, nil
}
