// Package gormrepo ...
package gormrepo

import (
	"context"
	"errors"
	orgerror "orgstructure/internal/errors"
	"orgstructure/internal/model"

	"gorm.io/gorm"
)

// DepartmentRepository ...
type DepartmentRepository struct {
	db *gorm.DB
}

// NewDepartmentRepository ...
func NewDepartmentRepository(db *gorm.DB) *DepartmentRepository {
	return &DepartmentRepository{db: db}
}

// Create ...
func (r *DepartmentRepository) Create(ctx context.Context, dept *model.Department) error {
	return r.db.WithContext(ctx).Create(dept).Error
}

// GetByID ...
func (r *DepartmentRepository) GetByID(ctx context.Context, id uint) (*model.Department, error) {
	var dept model.Department
	err := r.db.WithContext(ctx).First(&dept, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, orgerror.ErrDepartmentNotFound
	}
	return &dept, err
}

// Exists ...
func (r *DepartmentRepository) Exists(ctx context.Context, id uint) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Department{}).Where("id = ?", id).Count(&count).Error
	return count > 0, err
}

// GetChildren ...
func (r *DepartmentRepository) GetChildren(ctx context.Context, parentID uint) ([]model.Department, error) {
	var children []model.Department
	err := r.db.WithContext(ctx).Where("parent_id = ?", parentID).Find(&children).Error
	return children, err
}

// GetEmployees ...
func (r *DepartmentRepository) GetEmployees(ctx context.Context, departmentID uint, sortBy string) ([]model.Employee, error) {
	var employees []model.Employee

	allowedSorts := map[string]bool{"created_at": true, "full_name": true}
	if !allowedSorts[sortBy] {
		sortBy = "created_at"
	}

	err := r.db.WithContext(ctx).
		Where("department_id = ?", departmentID).
		Order(sortBy).
		Find(&employees).Error
	return employees, err
}

// Update ...
func (r *DepartmentRepository) Update(ctx context.Context, dept *model.Department) error {
	return r.db.WithContext(ctx).Save(dept).Error
}

// Delete ...
func (r *DepartmentRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.Department{}, id).Error
}

// GetSubtreeIDs ...
func (r *DepartmentRepository) GetSubtreeIDs(ctx context.Context, rootID uint) ([]uint, error) {
	query := `
		WITH RECURSIVE subtree AS (
			SELECT id FROM departments WHERE id = ?
			UNION ALL
			SELECT d.id FROM departments d
			INNER JOIN subtree s ON d.parent_id = s.id
		)
		SELECT id FROM subtree
	`

	var ids []uint
	err := r.db.WithContext(ctx).Raw(query, rootID).Scan(&ids).Error
	return ids, err
}

// MoveEmployees ...
func (r *DepartmentRepository) MoveEmployees(ctx context.Context, fromDepartmentIDs []uint, toDepartmentID uint) error {
	if len(fromDepartmentIDs) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).
		Model(&model.Employee{}).
		Where("department_id IN ?", fromDepartmentIDs).
		Update("department_id", toDepartmentID).Error
}

// DeleteByIDs ...
func (r *DepartmentRepository) DeleteByIDs(ctx context.Context, ids []uint) error {
	if len(ids) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).
		Where("id IN ?", ids).
		Delete(&model.Department{}).Error
}

// ExistsName ...
func (r *DepartmentRepository) ExistsName(ctx context.Context, name string, parentID *uint, excludeID *uint) (bool, error) {
	var count int64
	q := r.db.WithContext(ctx).Model(&model.Department{}).Where("name = ?", name)

	if parentID == nil {
		q = q.Where("parent_id IS NULL")
	} else {
		q = q.Where("parent_id = ?", *parentID)
	}

	if excludeID != nil {
		q = q.Where("id != ?", *excludeID)
	}

	err := q.Count(&count).Error
	return count > 0, err
}
