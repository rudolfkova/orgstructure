// Package model ...
package model

import "time"

// Department ...
type Department struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"size:200;not null"`
	ParentID  *uint  `gorm:"uniqueIndex:idx_name_parent"`
	Parent    *Department
	Children  []Department `gorm:"foreignKey:ParentID"`
	Employees []Employee
	CreatedAt time.Time
}
