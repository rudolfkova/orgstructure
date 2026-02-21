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

// EmployeeDTO ...
type EmployeeDTO struct {
	ID        uint       `json:"id"`
	FullName  string     `json:"full_name"`
	Position  string     `json:"position"`
	HiredAt   *time.Time `json:"hired_at,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
}

// ToDTO ...
func (e *Employee) ToDTO() EmployeeDTO {
	return EmployeeDTO{
		ID:        e.ID,
		FullName:  e.FullName,
		Position:  e.Position,
		HiredAt:   e.HiredAt,
		CreatedAt: e.CreatedAt,
	}
}
