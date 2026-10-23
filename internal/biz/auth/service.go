package auth

import (
	"errors"
	"strings"
	"time"

	"fibera/internal/infrastructure/db"
	"fibera/internal/middleware"
	"fibera/internal/model"
	"fibera/utils/security"

	"gorm.io/gorm"
)

// 认证相关错误
var (
	ErrInvalidCredentials = errors.New("用户名或密码错误")
	ErrUserDisabled       = errors.New("账号已被停用")
)

// Service 认证业务服务
type Service struct {
	db *gorm.DB
}

// NewService 创建认证业务服务
func NewService() *Service {
	database, _ := db.GetDB()
	return &Service{db: database}
}

// Login 校验账号密码，校验通过后签发 token 并记录登录信息
func (s *Service) Login(req *LoginRequest, loginIP string) (*LoginResponse, error) {
	userName := strings.TrimSpace(req.UserName)
	if req.Password == "" {
		return nil, ErrInvalidCredentials
	}

	// 优先按用户ID登录，其次按用户名登录
	var (
		user model.SysUser
		err  error
	)
	switch {
	case req.UserID != nil && *req.UserID > 0:
		err = s.db.Where("user_id = ?", *req.UserID).First(&user).Error
	case userName != "":
		err = s.db.Where("user_name = ?", userName).First(&user).Error
	default:
		return nil, ErrInvalidCredentials
	}
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	// 使用 bcrypt 校验明文密码与数据库中的密文
	if err := security.VerifyPassword(req.Password, user.Password); err != nil {
		return nil, ErrInvalidCredentials
	}

	if !user.Status {
		return nil, ErrUserDisabled
	}

	token, err := middleware.GenerateToken(user.UserID, user.UserName)
	if err != nil {
		return nil, err
	}

	// 更新最后登录 IP 与时间
	now := time.Now()
	if err := s.db.Model(&user).Updates(map[string]any{
		"login_ip":   loginIP,
		"login_date": now,
	}).Error; err != nil {
		return nil, err
	}

	user.LoginIP = loginIP
	user.LoginDate = &now

	return newLoginResponse(token, &user), nil
}
