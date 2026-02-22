// Package model создаёт структуры сотрудников и подразделений. Есть структуры с gorm тэгами для работы с БД,
// а также DTO для красивого ответа на запрос.
package model

// Файл department.go - модель подразделения.

import "time"

// Department структура отдельного подразделения с gorm тэгами. Нужна для работы с БД.
// Нельзя отдавать пользователю, т.к. даёт пользователю подробности работы БД и выглядит в ответе ужасно.
// Пользователю лучше отдавать EmployeeDTO.
type Department struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"size:200;not null"`

	// Родительское подразделение в дереве.
	ParentID *uint
	Parent   *Department `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

	// Дочерние подразделения в дереве.
	Children []Department `gorm:"foreignKey:ParentID"`

	Employees []Employee `gorm:"constraint:OnDelete:RESTRICT;"`

	CreatedAt time.Time `gorm:"not null;default:now()"`
}

// DepartmentTree структура дерева подразделений для формирования ответа на запрос.
type DepartmentTree struct {
	// ID корневой ноды.
	ID uint `json:"id"`

	Name      string        `json:"name"`
	CreatedAt time.Time     `json:"created_at"`
	Employees []EmployeeDTO `json:"employees,omitempty"`

	// Дочерние подразделения в дереве.
	Children []DepartmentTree `json:"children,omitempty"`
}

// DepartmentDTO структура подразделения для формирования ответа на запрос.
type DepartmentDTO struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`

	// ID родительского подразделения.
	ParentID *uint `json:"parent_id"`

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
