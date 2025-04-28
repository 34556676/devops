package libUtils

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type CustomClaims struct {
	Data interface{}
	jwt.RegisteredClaims
}

// 生成token
func GenerateToken(ctx context.Context, key string, data interface{}) (keys string, err error) {

	//1 生成claims
	customClaims := CustomClaims{
		data,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)), // 定义过期时间
			Issuer:    "devops",                                           // 签发人
		},
	}

	//2 生成token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, customClaims)

	// 3 签名
	signedToken, err := token.SignedString([]byte(key))
	if err != nil {
		fmt.Println("generate jwt token error:", err.Error())
		return "", err

	}

	return signedToken, nil
}

// 解析jwt
func ParseToken(tokenString, key string) (*CustomClaims, error) {
	// 解析token
	var mc = new(CustomClaims)
	token, err := jwt.ParseWithClaims(tokenString, mc, func(token *jwt.Token) (i interface{}, err error) {
		return []byte(key), nil
	})
	if err != nil {
		return nil, err
	}
	// 对token对象中的Claim进行类型断言
	if token.Valid { // 校验token
		return mc, nil
	}
	return nil, errors.New("invalid token")
}
