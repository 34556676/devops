package libUtils

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"golang.org/x/crypto/bcrypt"
)

func EncryptPassword(password, salt string) (string, error) {
	saltedPassword := salt + password
	hashed, err := bcrypt.GenerateFromPassword([]byte(saltedPassword), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

// GetClientIp 获取客户端IP
func GetClientIp(ctx context.Context) string {
	return g.RequestFromCtx(ctx).GetClientIp()
}

// GetUserAgent 获取user-agent
func GetUserAgent(ctx context.Context) string {
	return ghttp.RequestFromCtx(ctx).Header.Get("User-Agent")
}

// HashPassword 使用 Bcrypt 算法生成密码哈希值
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// ComparePasswords 比较输入的密码与哈希值是否匹配
func ComparePasswords(hashedPassword, password, userSalt string) bool {
	saltedPassword := userSalt + password
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(saltedPassword))
	return err == nil
}

// 生成安全的盐值
func GenerateSecureSalt() (string, error) {
	salt := make([]byte, 8)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(salt), nil
}
