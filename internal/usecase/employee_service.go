// Package usecase ...
package usecase

import (
	"context"
	orgerror "orgstructure/internal/errors"
	"orgstructure/internal/model"
	"orgstructure/internal/validation"
	"strings"
	"time"
)

// EmployeeService ...
type EmployeeService struct {
	empRepo  EmployeeRepository
	deptRepo DepartmentRepository
}

// NewEmployeeService ...
func NewEmployeeService(empRepo EmployeeRepository, deptRepo DepartmentRepository) *EmployeeService {
	return &EmployeeService{empRepo: empRepo, deptRepo: deptRepo}
}

// CreateEmployeeInDepartment ...
func (s *EmployeeService) CreateEmployeeInDepartment(ctx context.Context, departmentID uint, fullName, position string, hiredAt *time.Time) (*model.Employee, error) {
	fullName = strings.TrimSpace(fullName)
	if err := validation.ValidateStr(fullName, 200); err != nil {
		return nil, orgerror.ErrInvalidFullName
	}
	position = strings.TrimSpace(position)
	if err := validation.ValidateStr(position, 200); err != nil {
		return nil, orgerror.ErrInvalidPosition
	}

	exists, err := s.deptRepo.Exists(ctx, departmentID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, orgerror.ErrDepartmentNotFound
	}

	emp := &model.Employee{
		DepartmentID: departmentID,
		FullName:     fullName,
		Position:     position,
		HiredAt:      hiredAt,
	}

	if err := s.empRepo.Create(ctx, emp); err != nil {
		return nil, err
	}

	return emp, nil
}
