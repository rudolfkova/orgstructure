// Package repository ...
package repository

import (
	"context"
	"orgstructure/internal/model"
	"time"
)

// EmployeeRepository ...
type EmployeeRepository interface {
	CreateEmployeeInDepartment(ctx *context.Context, employeeFullName string, position string, hiredAt time.Time) (employee model.Employee, err error)
}
