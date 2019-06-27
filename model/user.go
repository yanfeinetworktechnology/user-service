package model

import "time"

// User 用户表字段
type User struct {
	ID        int64     `json:"id" gorm:"primary_key"`
	UserName  string    `json:"user_name" gorm:"size:15" gorm:"not null"`
	Password  string    `json:"password" gorm:"size:20" gorm:"not null"`
	Role      int       `json:"role"`
	InfoID    int64     `json:"info_id"`
	Wechat    string    `json:"wechat"`
	QQ        string    `json:"qq"`
	UpdatedAt time.Time `json:"updated_at" gorm:"not null"`
}

// LoginRequest 登录字段
type LoginRequest struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
}
