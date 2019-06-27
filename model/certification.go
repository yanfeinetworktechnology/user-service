package model

import "time"

// Person 个人表字段
type Person struct {
	ID        int64     `json:"id" gorm:"primary_key"`
	RealName  string    `json:"real_name" gorm:"size:15" gorm:"not null"`
	Sex       string    `json:"sex" gorm:"size:5" gorm:"not null"`
	HomeTown  string    `json:"hometown" gorm:"not null"`
	Phone     string    `json:"phone" gorm:"not null"`
	UpdatedAt time.Time `json:"updated_at" gorm:"not null"`
}

// PersonRequest 个人认证信息请求字段
type PersonRequest struct {
	RealName string `json:"real_name"`
	Sex      string `json:"sex"`
	HomeTown string `json:"hometown"`
	Phone    string `json:"phone"`
}
