package user

import (
	"fibera/internal/model"
	"time"
)

// CreateUserRequest 创建用户接口请求参数
type CreateUserRequest struct {
	UserID    int64  `json:"userId"`
	DeptID    int64  `json:"deptId"`
	UserName  string `json:"userName"`
	Password  string `json:"password"`
	Gender    string `json:"gender"`
	Telephone string `json:"telephone"`
	Email     string `json:"email"`
	Position  string `json:"position"`
	Leader    string `json:"leader"`
	Avatar    string `json:"avatar"`
	Status    bool   `json:"status"`
	Operator  *int64 `json:"operator"`
	Remark    string `json:"remark"`
}

// UpdateUserRequest 更新用户接口请求参数
type UpdateUserRequest struct {
	DeptID    int64  `json:"deptId"`
	UserName  string `json:"userName"`
	Password  string `json:"password"`
	Gender    string `json:"gender"`
	Telephone string `json:"telephone"`
	Email     string `json:"email"`
	Position  string `json:"position"`
	Leader    string `json:"leader"`
	Avatar    string `json:"avatar"`
	Status    bool   `json:"status"`
	Operator  *int64 `json:"operator"`
	Remark    string `json:"remark"`
}

// UserResponse 用户接口返回数据，不包含密码等敏感字段
type UserResponse struct {
	ID        uint       `json:"id"`
	UserID    int64      `json:"userId"`
	DeptID    int64      `json:"deptId"`
	UserName  string     `json:"userName"`
	Gender    string     `json:"gender"`
	Telephone string     `json:"telephone"`
	Email     string     `json:"email"`
	Position  string     `json:"position"`
	Leader    string     `json:"leader"`
	Avatar    string     `json:"avatar"`
	LoginIP   string     `json:"loginIp"`
	LoginDate *time.Time `json:"loginDate"`
	Status    bool       `json:"status"`
	Operator  *int64     `json:"operator"`
	Remark    string     `json:"remark"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
}

// toModel 将创建请求转换为用户实体
func (r *CreateUserRequest) toModel() *model.SysUser {
	return &model.SysUser{
		UserID:    r.UserID,
		DeptID:    r.DeptID,
		UserName:  r.UserName,
		Password:  r.Password,
		Gender:    r.Gender,
		Telephone: r.Telephone,
		Email:     r.Email,
		Position:  r.Position,
		Leader:    r.Leader,
		Avatar:    r.Avatar,
		Status:    r.Status,
		Operator:  r.Operator,
		Remark:    r.Remark,
	}
}

// toModel 将更新请求转换为用户实体
func (r *UpdateUserRequest) toModel(id uint) *model.SysUser {
	return &model.SysUser{
		ID:        id,
		DeptID:    r.DeptID,
		UserName:  r.UserName,
		Password:  r.Password,
		Gender:    r.Gender,
		Telephone: r.Telephone,
		Email:     r.Email,
		Position:  r.Position,
		Leader:    r.Leader,
		Avatar:    r.Avatar,
		Status:    r.Status,
		Operator:  r.Operator,
		Remark:    r.Remark,
	}
}

// newUserResponse 由用户实体组装响应数据
func newUserResponse(user *model.SysUser) *UserResponse {
	return &UserResponse{
		ID:        user.ID,
		UserID:    user.UserID,
		DeptID:    user.DeptID,
		UserName:  user.UserName,
		Gender:    user.Gender,
		Telephone: user.Telephone,
		Email:     user.Email,
		Position:  user.Position,
		Leader:    user.Leader,
		Avatar:    user.Avatar,
		LoginIP:   user.LoginIP,
		LoginDate: user.LoginDate,
		Status:    user.Status,
		Operator:  user.Operator,
		Remark:    user.Remark,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
