// Package usecase ...
package usecase

import (
	"context"
	orgerror "orgstructure/internal/errors"
	"orgstructure/internal/model"
	"orgstructure/internal/repository"
	"strings"
)

// DepartmentService ...
type DepartmentService struct {
	deptRepo repository.DepartmentRepository
}

// CreateDepartment ...
func (s *DepartmentService) CreateDepartment(ctx context.Context, name string, parentID *uint) (*model.Department, error) {
	if strings.TrimSpace(name) == "" {
		return nil, orgerror.ErrInvalidDepartmentName
	}

	if parentID != nil {
		exists, err := s.deptRepo.Exists(ctx, *parentID)
		if err != nil {
			return nil, err
		}
		if !exists {
			return nil, orgerror.ErrParentNotFound
		}
	}

	dept := &model.Department{
		Name:     strings.TrimSpace(name),
		ParentID: parentID,
	}

	if err := s.deptRepo.Create(ctx, dept); err != nil {
		return nil, err
	}

	return dept, nil
}
