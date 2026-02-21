// Package service ...
package service

import (
	"context"
	"orgstructure/internal/model"
)

// DepartmentService ...
type DepartmentService interface {
	CreateDepartment(ctx context.Context, name string, parentID *uint) (*model.Department, error)
	GetDepartment(ctx context.Context, depth int, includeEmployees bool, id uint) (*model.DepartmentTree, error)
	UpdateDepartment(ctx context.Context, id uint, name *string, parentID *uint) (*model.Department, error)
	DelDepartment(ctx context.Context, id uint, mode string, reassignToDepartmentID *uint) error
}
