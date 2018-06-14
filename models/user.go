package models

import "github.com/jinzhu/gorm"

type User struct {
	gorm.Model
	MobilePhone string `gorm:"column:user_name; type:varchar(11)"`
	Password    string `gorm:"column:password; type:varchar(16)"`
}
