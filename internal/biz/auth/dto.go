package auth

import (
	"fibera/internal/model"
	"time"
)

// LoginRequest 登录接口请求参数，用户ID与用户名二选一，同时提供时以用户ID为准
type LoginRequest struct {
	UserID   *int64 `json:"userId"`
	UserName string `json:"userName"`
	Password string `json:"password"`
}

// LoginResponse 登录接口返回数据
type LoginResponse struct {
	Token    string        `json:"token"`
	UserInfo LoginUserInfo `json:"userInfo"`
}

// LoginUserInfo 登录成功后返回的用户信息，不包含密码等敏感字段
type LoginUserInfo struct {
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
	CreatedAt time.Time  `json:"createdAt"`
}

// newLoginResponse 由用户实体组装登录响应数据
func newLoginResponse(token string, user *model.SysUser) *LoginResponse {
	return &LoginResponse{
		Token: token,
		UserInfo: LoginUserInfo{
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
			CreatedAt: user.CreatedAt,
		},
	}
}
