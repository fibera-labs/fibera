package model

import (
	"fibera/internal/infrastructure/db"
	"time"

	"gorm.io/gorm"
)

type SysUser struct {
	ID        uint       `gorm:"primaryKey;autoIncrement"`
	UserID    int64      `gorm:"column:user_id;not null;comment:用户ID"`
	DeptID    int64      `gorm:"column:dept_id;not null;comment:部门ID"`
	UserName  string     `gorm:"column:user_name;type:varchar(50);not null;default:'';comment:用户名"`
	Password  string     `gorm:"column:password;type:varchar(60);not null;comment:密码"`
	Gender    string     `gorm:"column:gender;type:char(1);not null;default:'1';comment:用户性别（0女 1男 2未知）"`
	Telephone string     `gorm:"column:telephone;type:varchar(20);comment:电话"`
	Email     string     `gorm:"column:email;type:varchar(100);default:'';comment:邮箱"`
	Position  string     `gorm:"column:position;type:varchar(50);comment:职位"`
	Leader    string     `gorm:"column:leader;type:varchar(50);comment:上级"`
	Avatar    string     `gorm:"column:avatar;type:varchar(255);comment:头像路径"`
	LoginIP   string     `gorm:"column:login_ip;type:varchar(50);comment:最后登录IP"`
	LoginDate *time.Time `gorm:"column:login_date;comment:最后登录时间"`
	Status    bool       `gorm:"column:status;type:tinyint(1);not null;default:true;comment:帐号状态（0停用 1正常）"`
	Operator  *int64     `gorm:"column:operator;comment:操作人"`
	Remark    string     `gorm:"column:remark;type:text;comment:备注"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (SysUser) TableName() string {
	return "sys_user"
}

func init() {
	db.Register(new(SysUser))
}
