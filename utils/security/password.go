package security

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword 使用 bcrypt 对明文密码加密，返回可直接持久化的密文
func HashPassword(plainPassword string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

// VerifyPassword 校验明文密码与 bcrypt 密文是否匹配，不匹配时返回 bcrypt.ErrMismatchedHashAndPassword
func VerifyPassword(plainPassword, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
}
