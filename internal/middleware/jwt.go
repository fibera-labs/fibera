package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"fibera/internal/common"
	"fibera/internal/config"
	"fibera/internal/infrastructure/redis"
	"fibera/utils/security"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
	goredis "github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type JWTConfig struct {
	Secret         []byte
	Issuer         string
	SigningMethod  jwt.SigningMethod
	ExpirationTime time.Duration
}

type CustomClaims struct {
	UserID   int    `json:"userID"`
	UserName string `json:"userName"`
	jwt.RegisteredClaims
}

var (
	jwtConfig *JWTConfig
	mutex     sync.Mutex
)

// Jwt JWT中间件
func Jwt() fiber.Handler {
	excludedPaths := map[string]bool{
		"/login":  true,
		"/public": true,
	}

	return func(c fiber.Ctx) error {
		// 检查排除路由
		if excludedPaths[c.Path()] {
			return c.Next()
		}

		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return common.Error(c, fiber.StatusUnauthorized, "授权头缺失")
		}

		tokenString, err := parseBearerToken(authHeader)
		if err != nil {
			return common.Error(c, fiber.StatusUnauthorized, "无效的token格式")
		}

		claims, err := validateJWT(tokenString)
		if err != nil {
			return handleJWTError(c, err)
		}

		// 存储claims到上下文
		c.Locals("jwt_claims", claims)
		return c.Next()
	}
}

// GenerateToken 生成token
func GenerateToken(userID int, userName string) (string, error) {
	conf, _ := config.LoadConfig()
	jwtConf, err := loadJwtConfig(conf)
	if err != nil {
		return "", err
	}

	claims := CustomClaims{
		UserID:   userID,
		UserName: userName,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(jwtConf.ExpirationTime)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    jwtConf.Issuer,
		},
	}

	token := jwt.NewWithClaims(jwtConf.SigningMethod, claims)
	tokenString, err := token.SignedString(jwtConf.Secret)
	if err != nil {
		return "", err
	}

	// 启用 Redis 缓存时，将 token 写入 Redis
	if conf.Redis.Enable {
		redisCli, _ := redis.GetRedis()
		claimsJSON, mErr := json.Marshal(claims)
		if mErr != nil {
			zap.L().Error("failed to marshal claims", zap.Error(mErr))
			return "", mErr
		}

		err = redisCli.Set(context.Background(), tokenString, claimsJSON, jwtConf.ExpirationTime).Err()
		if err != nil {
			zap.L().Error("failed to set token to redis", zap.Error(err))
			return "", err
		}
	}
	return tokenString, nil
}

// validateJWT 验证并解析token
func validateJWT(tokenString string) (*CustomClaims, error) {
	conf, _ := config.LoadConfig()

	// 如果启用redis缓存，则从redis中获取claims
	if conf.Redis.Enable {
		redisCli, _ := redis.GetRedis()
		var claims CustomClaims
		claimsJSON, err := redisCli.Get(context.Background(), tokenString).Bytes()
		if err != nil {
			if errors.Is(err, goredis.Nil) {
				return nil, jwt.ErrTokenInvalidClaims
			}
			zap.L().Error("failed to get token from redis", zap.Error(err))
			return nil, err
		}
		if err = json.Unmarshal(claimsJSON, &claims); err != nil {
			zap.L().Error("failed to unmarshal claims", zap.Error(err))
			return nil, err
		}
		return &claims, nil
	}

	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwtConfig.SigningMethod {
			return nil, fmt.Errorf("不支持的签名方法: %v", token.Header["alg"])
		}
		return jwtConfig.Secret, nil
	})
	if err != nil {
		zap.L().Warn("failed to validate token", zap.Error(err))
		return nil, err
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, jwt.ErrTokenInvalidClaims
}

// handleJWTError 处理JWT错误
func handleJWTError(c fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, jwt.ErrTokenExpired):
		ResetJWTConfig()
		return common.Error(c, fiber.StatusUnauthorized, "token过期")
	case errors.Is(err, jwt.ErrTokenMalformed):
		return common.Error(c, fiber.StatusUnauthorized, "无效的token格式")
	case errors.Is(err, jwt.ErrTokenSignatureInvalid):
		return common.Error(c, fiber.StatusUnauthorized, "无效的token签名")
	default:
		return common.Error(c, fiber.StatusUnauthorized, "认证失败")
	}
}

// parseBearerToken 解析Bearer token
func parseBearerToken(header string) (string, error) {
	const bearerPrefix = "Bearer "
	if len(header) < len(bearerPrefix) || !strings.HasPrefix(header, bearerPrefix) {
		return "", fmt.Errorf("invalid authorization header format")
	}
	return header[len(bearerPrefix):], nil
}

// loadJwtConfig 加载JWT配置
func loadJwtConfig(conf *config.Config) (*JWTConfig, error) {
	mutex.Lock()
	defer mutex.Unlock()

	if jwtConfig != nil {
		return jwtConfig, nil
	}

	randSecret, err := security.CryptoRandSecret(32)
	if err != nil {
		return nil, fmt.Errorf("failed to generate random secret: %v", err)
	}
	secret := []byte(randSecret)

	signingMethod := jwt.GetSigningMethod(conf.Jwt.SigningMethod)
	if signingMethod == nil {
		return nil, fmt.Errorf("invalid signing method: %s", conf.Jwt.SigningMethod)
	}

	expirationTime, err := time.ParseDuration(conf.Jwt.ExpirationTime)
	if err != nil {
		return nil, fmt.Errorf("invalid expiration time format: %s", conf.Jwt.ExpirationTime)
	}

	jwtConfig = &JWTConfig{
		Secret:         secret,
		Issuer:         conf.Jwt.Issuer,
		SigningMethod:  signingMethod,
		ExpirationTime: expirationTime,
	}
	return jwtConfig, nil
}

// ResetJWTConfig 重置jwt配置
func ResetJWTConfig() {
	mutex.Lock()
	defer mutex.Unlock()
	jwtConfig = nil
}
