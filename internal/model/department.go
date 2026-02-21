// Package model ...
package model

import "time"

// Department ...
type Department struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"size:200;not null"`

	ParentID *uint
	Parent   *Department `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Children []Department `gorm:"foreignKey:ParentID"`

	Employees []Employee `gorm:"constraint:OnDelete:RESTRICT;"`

	CreatedAt time.Time `gorm:"not null;default:now()"`
}
