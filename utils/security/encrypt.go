package security

import (
	"crypto/rand"
	"encoding/base64"
)

// CryptoRandSecret 生成随机密钥
func CryptoRandSecret(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	// 使用 Base64 编码扩展字符集（包含大小写字母、数字、+、/、=）
	return base64.URLEncoding.EncodeToString(b), nil
}
