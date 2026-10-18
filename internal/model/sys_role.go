package model

import (
	"fibera/internal/infrastructure/db"
	"time"

	"gorm.io/gorm"
)

type SysRole struct {
	ID        uint   `gorm:"primaryKey;autoIncrement"`
	RoleID    int64  `gorm:"column:role_id;not null;comment:角色ID"`
	RoleName  string `gorm:"column:role_name;type:varchar(50);not null;default:'';comment:角色名称"`
	Status    bool   `gorm:"column:status;type:tinyint(1);not null;default:true;comment:角色状态（0停用 1正常）"`
	Operator  *int64 `gorm:"column:operator;comment:操作人"`
	Remark    string `gorm:"column:remark;type:text;comment:备注"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (SysRole) TableName() string {
	return "sys_role"
}

func init() {
	db.Register(new(SysRole))
}
