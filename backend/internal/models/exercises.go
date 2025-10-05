package models

import "gorm.io/gorm"

type Exercise struct {
	gorm.Model
	Name string `gorm:"unique;not null"`
}
