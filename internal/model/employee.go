package model

import "time"

// Employee ...
type Employee struct {
	ID uint `gorm:"primaryKey"`

	DepartmentID uint       `gorm:"not null;index"`
	Department   Department `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`

	FullName string `gorm:"type:varchar(200);not null"`
	Position string `gorm:"type:varchar(200);not null"`

	HiredAt *time.Time `gorm:"type:date"`

	CreatedAt time.Time `gorm:"not null;default:now()"`
}
