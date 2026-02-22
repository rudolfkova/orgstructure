// Package model создаёт структуры сотрудников и подразделений. Есть структуры с gorm тэгами для работы с БД,
// а также DTO для красивого ответа на запрос.
package model

// Файл employee.go - модель сотрудника.

import (
	"time"
)

// Employee структура отдельного сотрудника с gorm тэгами. Нужна для работы с БД.
// Нельзя отдавать пользователю, т.к. даёт пользователю подробности работы БД и выглядит в ответе ужасно.
// Пользователю лучше отдавать EmployeeDTO.
type Employee struct {
	ID uint `gorm:"primaryKey"`

	DepartmentID uint       `gorm:"not null;index"`
	Department   Department `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`

	FullName string `gorm:"type:varchar(200);not null"`

	// Должность сотрудника.
	Position string `gorm:"type:varchar(200);not null"`

	// Дата найма.
	HiredAt *time.Time `gorm:"type:date"`

	CreatedAt time.Time `gorm:"not null;default:now()"`
}

// NewEmployee конструктор для создания структуры с gorm тэгами.
func NewEmployee(departmentID uint, fullName, position string, hiredAt *time.Time) *Employee {
	return &Employee{
		DepartmentID: departmentID,
		FullName:     fullName,
		Position:     position,
		HiredAt:      hiredAt,
	}
}

// EmployeeDTO структура для сотрудника для ответа на запрос.
type EmployeeDTO struct {
	ID       uint   `json:"id"`
	FullName string `json:"full_name"`

	// Должность сотрудника.
	Position string `json:"position"`

	// Дата найма.
	HiredAt   *time.Time `json:"hired_at,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
}

// ToDTO конвертирует
func (e *Employee) ToDTO() EmployeeDTO {
	return EmployeeDTO{
		ID:        e.ID,
		FullName:  e.FullName,
		Position:  e.Position,
		HiredAt:   e.HiredAt,
		CreatedAt: e.CreatedAt,
	}
}
