// Package service ...
package service

import (
	"context"
	"orgstructure/internal/model"
	"time"
)

// EmployeeService ...
type EmployeeService interface {
	CreateEmployeeInDepartment(ctx context.Context, departmentID uint, fullName, position string, hiredAt *time.Time) (*model.Employee, error)
}
