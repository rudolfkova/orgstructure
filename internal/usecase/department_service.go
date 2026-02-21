// Package usecase ...
package usecase

import (
	"context"
	"fmt"
	orgerror "orgstructure/internal/errors"
	"orgstructure/internal/model"
	"orgstructure/internal/repository"
	"strings"
)

// DepartmentService ...
type DepartmentService struct {
	deptRepo repository.DepartmentRepository
}

// NewDepartmentService ...
func NewDepartmentService(repo repository.DepartmentRepository) *DepartmentService {
	return &DepartmentService{deptRepo: repo}
}

// CreateDepartment ...
func (s *DepartmentService) CreateDepartment(ctx context.Context, name string, parentID *uint) (*model.Department, error) {
	name = strings.TrimSpace(name)
	if name == "" {
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

	duplicate, err := s.deptRepo.ExistsName(ctx, name, parentID, nil)
	if err != nil {
		return nil, err
	}
	if duplicate {
		return nil, orgerror.ErrDuplicateName
	}

	dept := &model.Department{
		Name:     name,
		ParentID: parentID,
	}

	if err := s.deptRepo.Create(ctx, dept); err != nil {
		return nil, err
	}

	return dept, nil
}

// GetDepartment ...
func (s *DepartmentService) GetDepartment(ctx context.Context, depth int, includeEmployees bool, id uint) (*model.DepartmentTree, error) {
	if depth > 5 {
		depth = 5
	}

	dept, err := s.deptRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if dept == nil {
		return nil, orgerror.ErrDepartmentNotFound
	}

	return s.buildTree(ctx, dept, depth, includeEmployees)
}

func (s *DepartmentService) buildTree(ctx context.Context, dept *model.Department, depth int, includeEmployees bool) (*model.DepartmentTree, error) {
	node := &model.DepartmentTree{
		ID:        dept.ID,
		Name:      dept.Name,
		CreatedAt: dept.CreatedAt,
	}

	if includeEmployees {
		emps, err := s.deptRepo.GetEmployees(ctx, dept.ID, "created_at")
		if err != nil {
			return nil, err
		}

		for _, e := range emps {
			node.Employees = append(node.Employees, model.EmployeeDTO{
				ID:        e.ID,
				FullName:  e.FullName,
				Position:  e.Position,
				HiredAt:   e.HiredAt,
				CreatedAt: e.CreatedAt,
			})
		}
	}

	if depth <= 0 {
		return node, nil
	}

	children, err := s.deptRepo.GetChildren(ctx, dept.ID)
	if err != nil {
		return nil, err
	}

	for _, child := range children {
		childNode, err := s.buildTree(ctx, &child, depth-1, includeEmployees)
		if err != nil {
			return nil, err
		}
		node.Children = append(node.Children, *childNode)
	}

	return node, nil
}

// UpdateDepartment ...
func (s *DepartmentService) UpdateDepartment(ctx context.Context, id uint, name *string, parentID *uint) (*model.Department, error) {
	dept, err := s.deptRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if dept == nil {
		return nil, orgerror.ErrDepartmentNotFound
	}

	if parentID != nil {
		if *parentID == id {
			return nil, orgerror.ErrCyclicDependency
		}

		subtreeIDs, err := s.deptRepo.GetSubtreeIDs(ctx, id)
		if err != nil {
			return nil, err
		}
		for _, subtreeID := range subtreeIDs {
			if subtreeID == *parentID {
				return nil, orgerror.ErrCyclicDependency
			}
		}

		exists, err := s.deptRepo.Exists(ctx, *parentID)
		if err != nil {
			return nil, err
		}
		if !exists {
			return nil, orgerror.ErrParentNotFound
		}

		dept.ParentID = parentID
	}

	if name != nil {
		trimmed := strings.TrimSpace(*name)
		if trimmed == "" {
			return nil, orgerror.ErrInvalidDepartmentName
		}

		duplicate, err := s.deptRepo.ExistsName(ctx, trimmed, dept.ParentID, &id)
		if err != nil {
			return nil, err
		}
		if duplicate {
			return nil, orgerror.ErrDuplicateName
		}

		dept.Name = trimmed
	}

	if err := s.deptRepo.Update(ctx, dept); err != nil {
		return nil, err
	}

	return dept, nil
}

// DelDepartment ...
func (s *DepartmentService) DelDepartment(ctx context.Context, id uint, mode string, reassignToDepartmentID *uint) error {
	exists, err := s.deptRepo.Exists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return orgerror.ErrDepartmentNotFound
	}

	subtreeIDs, err := s.deptRepo.GetSubtreeIDs(ctx, id)
	if err != nil {
		return err
	}

	switch mode {
	case "cascade":
		return s.deptRepo.DeleteByIDs(ctx, subtreeIDs)

	case "reassign":
		if reassignToDepartmentID == nil {
			return orgerror.ErrInvalidReassignToDepartmentID
		}

		targetExists, err := s.deptRepo.Exists(ctx, *reassignToDepartmentID)
		if err != nil {
			return err
		}
		if !targetExists {
			return orgerror.ErrParentNotFound
		}
		for _, subtreeID := range subtreeIDs {
			if subtreeID == *reassignToDepartmentID {
				return fmt.Errorf("reassign_to_department_id is within the deleted subtree")
			}
		}

		if err := s.deptRepo.MoveEmployees(ctx, subtreeIDs, *reassignToDepartmentID); err != nil {
			return err
		}

		return s.deptRepo.DeleteByIDs(ctx, subtreeIDs)

	default:
		return orgerror.ErrInvalidMode
	}
}
