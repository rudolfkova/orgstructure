// Package repository ...
package repository

import (
	"context"
	"orgstructure/internal/model"
)

// DepartmentRepository ...
type DepartmentRepository interface {
	CreateDepartment(ctx *context.Context, departmentName string, pareintID int) error
	GetDepartment(ctx *context.Context, treeDepth int, includEmployees bool) (department *model.Department, err error)
	MoveDepartment(ctx *context.Context, departmentName string, pareintID int) error
	DelDepartment(ctx *context.Context, departmentName string, pareintID int) error
}
