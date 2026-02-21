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
	GetChildren(ctx context.Context, parentID uint) ([]model.Department, error)
	GetEmployees(ctx context.Context, departmentID uint, sortBy string) ([]model.Employee, error)
	Update(ctx context.Context, dept *model.Department) error
	Delete(ctx context.Context, id uint) error
	// GetSubtreeIDs возвращает все ID подразделений в поддереве (включая корень)
	GetSubtreeIDs(ctx context.Context, rootID uint) ([]uint, error)
	// MoveEmployees переводит сотрудников из одного подразделения в другое
	MoveEmployees(ctx context.Context, fromDepartmentIDs []uint, toDepartmentID uint) error
	// DeleteByIDs удаляет подразделения по списку ID
	DeleteByIDs(ctx context.Context, ids []uint) error
	// ExistsName проверяет, существует ли подразделение с таким именем у данного родителя
	ExistsName(ctx context.Context, name string, parentID *uint, excludeID *uint) (bool, error)
}
