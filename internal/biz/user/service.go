package user

import (
	"fibera/internal/infrastructure/db"
	"fibera/internal/model"

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
func (s *Service) GetUserList() ([]*model.SysUser, error) {
	var users []*model.SysUser
	err := s.db.Find(&users).Error
	return users, err
}

// GetUserByID 根据主键查询用户
func (s *Service) GetUserByID(id uint) (*model.SysUser, error) {
	var user model.SysUser
	if err := s.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// CreateUser 创建用户
func (s *Service) CreateUser(user *model.SysUser) error {
	return s.db.Create(user).Error
}

// UpdateUser 根据主键更新用户
func (s *Service) UpdateUser(user *model.SysUser) error {
	return s.db.Save(user).Error
}

// DeleteUser 根据主键删除用户
func (s *Service) DeleteUser(id uint) error {
	return s.db.Delete(&model.SysUser{}, id).Error
}
