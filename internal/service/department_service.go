// Package service ...
package service

import (
	"context"
	"orgstructure/internal/model"
)

// DepartmentService ...
type DepartmentService interface {
	CreateDepartment(ctx context.Context, departmentName string, pareintID *int) (department model.Department, err error)
	GetDepartment(ctx context.Context, treeDepth int, includEmployees bool, deportmentID int) (department *model.Department, err error)
	UpdateDepartment(ctx context.Context, departmentName *string, pareintID *int, departmentID int) (department *model.Department, err error)
	DelDepartment(ctx context.Context, mode string, reassignToDepartmentID *int, departmentID int) error
}
