package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"time"
)

type UserRole string

const (
	RoleAdmin UserRole = "admin"
	RoleUser  UserRole = "user"
)

type User struct {
	// 主键
	ID uint `gorm:"primaryKey" json:"id"`

	// 用户名 (唯一索引，长度限制50)
	Username string `gorm:"uniqueIndex;size:50;not null" json:"username" validate:"required,min=3,max=50"`

	// 密码 (存储bcrypt哈希值，不序列化到JSON)
	Password string `gorm:"size:255;not null" json:"-" validate:"required,min=8"`

	// 邮箱 (唯一索引，长度限制100)
	Email string `gorm:"uniqueIndex;size:100;not null" json:"email" validate:"required,email"`

	// 角色 (枚举值)
	Role UserRole `gorm:"type:enum('admin','user');default:'user'" json:"role" validate:"oneof=admin user"`

	// 创建时间
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}

// 密码加密钩子
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.Password != "" {
		hashed, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		u.Password = string(hashed)
	}
	return nil
}

// 密码验证方法
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

func (User) TableName() string {
	return "users"
}
