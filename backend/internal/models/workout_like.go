package models

import "gorm.io/gorm"

type WorkoutLike struct {
	gorm.Model
	UserID   uint `gorm:"not null;index;uniqueIndex:ux_user_record"`
	RecordID uint `gorm:"not null;index;uniqueIndex:ux_user_record"`

	User   User          `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Record WorkoutRecord `gorm:"foreignKey:RecordID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
