// Package model ...
package model

import "time"

// Department ...
type Department struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"size:200;not null"`

	ParentID *uint
	Parent   *Department  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Children []Department `gorm:"foreignKey:ParentID"`

	Employees []Employee `gorm:"constraint:OnDelete:RESTRICT;"`

	CreatedAt time.Time `gorm:"not null;default:now()"`
}

// DepartmentTree ...
type DepartmentTree struct {
	ID        uint             `json:"id"`
	Name      string           `json:"name"`
	CreatedAt time.Time        `json:"created_at"`
	Employees []EmployeeDTO    `json:"employees,omitempty"`
	Children  []DepartmentTree `json:"children,omitempty"`
}

// DepartmentDTO ...
type DepartmentDTO struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	ParentID  *uint     `json:"parent_id"`
	CreatedAt time.Time `json:"created_at"`
}

// ToDTO ...
func (d *Department) ToDTO() DepartmentDTO {
	return DepartmentDTO{
		ID:        d.ID,
		Name:      d.Name,
		ParentID:  d.ParentID,
		CreatedAt: d.CreatedAt,
	}
}
