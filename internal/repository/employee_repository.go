// Package repository ...
package repository

import (
	"context"
	"orgstructure/internal/model"
	"time"
)

// EmployeeRepository ...
type EmployeeRepository interface {
	Create(ctx context.Context, emp *model.Employee) error
	DepartmentExists(ctx context.Context, departmentID uint) (bool, error)
}

// NewEmployee создаёт модель сотрудника из переданных данных
func NewEmployee(departmentID uint, fullName, position string, hiredAt *time.Time) *model.Employee {
	return &model.Employee{
		DepartmentID: departmentID,
		FullName:     fullName,
		Position:     position,
		HiredAt:      hiredAt,
	}
}
