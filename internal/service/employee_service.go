// Package service ...
package service

import (
	"context"
	"orgstructure/internal/model"
	"time"
)

// EmployeeService ...
type EmployeeService interface {
	CreateEmployeeInDepartment(ctx *context.Context, employeeFullName string, position string, hiredAt *time.Time, id int) (employee model.Employee, err error)
}
