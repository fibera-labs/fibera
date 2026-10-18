package model

import (
	"fibera/internal/infrastructure/db"
	"time"

	"gorm.io/gorm"
)

type SysDept struct {
	ID        uint   `gorm:"primaryKey;autoIncrement"`
	DeptID    int64  `gorm:"column:dept_id;not null;comment:部门ID"`
	ParentID  int64  `gorm:"column:parent_id;not null;comment:上级部门ID"`
	DeptName  string `gorm:"column:dept_name;type:varchar(50);not null;default:'';comment:部门名称"`
	Leader    string `gorm:"column:leader;type:varchar(50);comment:负责人"`
	Telephone string `gorm:"column:telephone;type:varchar(20);comment:电话"`
	Email     string `gorm:"column:email;type:varchar(100);default:'';comment:邮箱"`
	Status    bool   `gorm:"column:status;type:tinyint(1);not null;default:true;comment:部门状态（0停用 1正常）"`
	Operator  *int64 `gorm:"column:operator;comment:操作人"`
	Remark    string `gorm:"column:remark;type:text;comment:备注"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (SysDept) TableName() string {
	return "sys_dept"
}

func init() {
	db.Register(new(SysDept))
}
