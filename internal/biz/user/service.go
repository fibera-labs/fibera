package user

import (
	"fibera/internal/infrastructure/db"
	"fibera/internal/model"
	"fibera/utils/security"

	"gorm.io/gorm"
)

// Service 系统用户业务服务
type Service struct {
	db *gorm.DB
}

// NewService 创建系统用户业务服务
func NewService() *Service {
	database, _ := db.GetDB()
	return &Service{db: database}
}

// GetUserList 查询用户列表
func (s *Service) GetUserList() ([]*UserResponse, error) {
	var users []*model.SysUser
	if err := s.db.Find(&users).Error; err != nil {
		return nil, err
	}

	responses := make([]*UserResponse, 0, len(users))
	for _, user := range users {
		responses = append(responses, newUserResponse(user))
	}
	return responses, nil
}

// GetUserByID 根据主键查询用户
func (s *Service) GetUserByID(id uint) (*UserResponse, error) {
	var user model.SysUser
	if err := s.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return newUserResponse(&user), nil
}

// CreateUser 创建用户，明文密码统一使用 bcrypt 加密后入库
func (s *Service) CreateUser(req *CreateUserRequest) (*UserResponse, error) {
	user := req.toModel()
	if user.Password != "" {
		hashed, err := security.HashPassword(user.Password)
		if err != nil {
			return nil, err
		}
		user.Password = hashed
	}

	if err := s.db.Create(user).Error; err != nil {
		return nil, err
	}
	return newUserResponse(user), nil
}

// UpdateUser 根据主键更新用户，未传密码时保留原密码
func (s *Service) UpdateUser(id uint, req *UpdateUserRequest) (*UserResponse, error) {
	user := req.toModel(id)

	query := s.db
	if user.Password != "" {
		hashed, err := security.HashPassword(user.Password)
		if err != nil {
			return nil, err
		}
		user.Password = hashed
	} else {
		query = query.Omit("password")
	}

	if err := query.Save(user).Error; err != nil {
		return nil, err
	}
	return newUserResponse(user), nil
}

// DeleteUser 根据主键删除用户
func (s *Service) DeleteUser(id uint) error {
	return s.db.Delete(&model.SysUser{}, id).Error
}
