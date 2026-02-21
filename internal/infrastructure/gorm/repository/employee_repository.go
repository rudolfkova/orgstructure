// Package gormrepo ...
package gormrepo

import (
	"context"
	"orgstructure/internal/model"

	"gorm.io/gorm"
)

// EmployeeRepository ...
type EmployeeRepository struct {
	db *gorm.DB
}

// NewEmployeeRepository ...
func NewEmployeeRepository(db *gorm.DB) *EmployeeRepository {
	return &EmployeeRepository{db: db}
}

// Create ...
func (r *EmployeeRepository) Create(ctx context.Context, emp *model.Employee) error {
	return r.db.WithContext(ctx).Create(emp).Error
}

// DepartmentExists ...
func (r *EmployeeRepository) DepartmentExists(ctx context.Context, departmentID uint) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Department{}).Where("id = ?", departmentID).Count(&count).Error
	return count > 0, err
}
