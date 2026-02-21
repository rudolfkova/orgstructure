// Package repository ...
package repository

import (
	"context"
	"orgstructure/internal/model"
)

// DepartmentRepository ...
type DepartmentRepository interface {
	Create(ctx context.Context, dept *model.Department) error
	GetByID(ctx context.Context, id uint) (*model.Department, error)
	Exists(ctx context.Context, id uint) (bool, error)
}
