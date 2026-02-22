// Package usecase_test ...
package usecase_test

import (
	"context"
	orgerror "orgstructure/internal/errors"
	"orgstructure/internal/model"
	"orgstructure/internal/usecase"
	repoMocks "orgstructure/mocks/repository"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func uintPtr(i int) *uint {
	if i < 0 {
		return nil
	}
	uintVal := uint(i)
	uintPtr := &uintVal
	return uintPtr
}

func TwoTrees() (developers, frontend, backend, offchain *model.Department) {
	frontend = &model.Department{
		ID:   uint(11),
		Name: "frontend",

		ParentID: uintPtr(10),
		Parent:   developers,
		Children: make([]model.Department, 0),

		Employees: make([]model.Employee, 0),

		CreatedAt: time.Now(),
	}

	developers = &model.Department{
		ID:   uint(10),
		Name: "developers",

		ParentID: nil,
		Parent:   nil,
		Children: []model.Department{
			*frontend,
		},

		Employees: make([]model.Employee, 0),

		CreatedAt: time.Now(),
	}

	offchain = &model.Department{
		ID:   uint(21),
		Name: "offchain",

		ParentID: uintPtr(20),
		Parent:   backend,
		Children: make([]model.Department, 0),

		Employees: make([]model.Employee, 0),

		CreatedAt: time.Now(),
	}

	backend = &model.Department{
		ID:   uint(20),
		Name: "backend",

		ParentID: nil,
		Parent:   nil,
		Children: []model.Department{
			*offchain,
		},

		Employees: []model.Employee{
			{
				DepartmentID: uint(20),
				FullName:     "Anton Denisov",
				Position:     "Lead",
				HiredAt:      nil,
			},
		},

		CreatedAt: time.Now(),
	}

	return
}

// --- CreateDepartment ---

func TestDepartmentUsecase_CreateDepartment_Success(t *testing.T) {
	deptRepo := new(repoMocks.DepartmentRepository)

	s := usecase.NewDepartmentService(
		deptRepo,
	)

	ctx := context.Background()
	name := "backend"
	parentID := uintPtr(1)
	var nilUint *uint

	dept := &model.Department{
		Name:     name,
		ParentID: parentID,
	}

	deptRepo.
		On("Exists", ctx, *parentID).
		Return(true, nil)
	deptRepo.
		On("ExistsName", ctx, name, parentID, nilUint).
		Return(false, nil)
	deptRepo.
		On("Create", ctx, dept).
		Return(nil)

	respDept, err := s.CreateDepartment(ctx, name, parentID)

	require.NoError(t, err)
	assert.Equal(t, name, respDept.Name)
	assert.Equal(t, parentID, respDept.ParentID)
}

func TestDepartmentUsecase_CreateDepartment_EmptyName(t *testing.T) {
	deptRepo := new(repoMocks.DepartmentRepository)
	s := usecase.NewDepartmentService(deptRepo)

	_, err := s.CreateDepartment(context.Background(), "   ", nil)

	require.ErrorIs(t, err, orgerror.ErrInvalidDepartmentName)
}

func TestDepartmentUsecase_CreateDepartment_ParentNotFound(t *testing.T) {
	deptRepo := new(repoMocks.DepartmentRepository)
	s := usecase.NewDepartmentService(deptRepo)

	ctx := context.Background()
	parentID := uintPtr(99)

	deptRepo.
		On("Exists", ctx, *parentID).
		Return(false, nil)

	_, err := s.CreateDepartment(ctx, "backend", parentID)

	require.ErrorIs(t, err, orgerror.ErrParentNotFound)
}

func TestDepartmentUsecase_CreateDepartment_DuplicateName(t *testing.T) {
	deptRepo := new(repoMocks.DepartmentRepository)
	s := usecase.NewDepartmentService(deptRepo)

	ctx := context.Background()
	name := "backend"
	var nilUint *uint

	deptRepo.
		On("ExistsName", ctx, name, nilUint, nilUint).
		Return(true, nil)

	_, err := s.CreateDepartment(ctx, name, nil)

	require.ErrorIs(t, err, orgerror.ErrDuplicateName)
}

// --- GetDepartment ---

func TestDepartmentUsecase_GetDepartment_Success(t *testing.T) {
	deptRepo := new(repoMocks.DepartmentRepository)

	s := usecase.NewDepartmentService(
		deptRepo,
	)

	ctx := context.Background()

	dept := &model.Department{
		ID:   uint(1),
		Name: "backend",

		ParentID: nil,
		Parent:   nil,
		Children: make([]model.Department, 0),

		Employees: make([]model.Employee, 0),

		CreatedAt: time.Now(),
	}

	deptRepo.
		On("GetByID", ctx, dept.ID).
		Return(dept, nil)

	respDept, err := s.GetDepartment(ctx, 0, false, dept.ID)

	require.NoError(t, err)
	assert.Equal(t, dept.ID, respDept.ID)
	assert.Equal(t, dept.Name, respDept.Name)
}

func TestDepartmentUsecase_GetDepartment_NotFound(t *testing.T) {
	deptRepo := new(repoMocks.DepartmentRepository)
	s := usecase.NewDepartmentService(deptRepo)

	ctx := context.Background()

	deptRepo.
		On("GetByID", ctx, uint(99)).
		Return(nil, orgerror.ErrDepartmentNotFound)

	_, err := s.GetDepartment(ctx, 1, false, 99)

	require.ErrorIs(t, err, orgerror.ErrDepartmentNotFound)
}

// --- UpdateDepartment ---

func TestDepartmentUsecase_UpdateDepartment_Success(t *testing.T) {
	deptRepo := new(repoMocks.DepartmentRepository)

	s := usecase.NewDepartmentService(
		deptRepo,
	)

	ctx := context.Background()

	developers, _, backend, offchain := TwoTrees()

	// Перемещаем backend в отдел developers.

	deptRepo.
		On("GetByID", ctx, backend.ID).
		Return(backend, nil)
	deptRepo.
		On("GetSubtreeIDs", ctx, backend.ID).
		Return([]uint{backend.ID, offchain.ID}, nil)
	deptRepo.
		On("Exists", ctx, developers.ID).
		Return(true, nil)
	deptRepo.
		On("Update", ctx, backend).
		Return(nil)

	respDept, err := s.UpdateDepartment(ctx, backend.ID, nil, &developers.ID)

	require.NoError(t, err)
	assert.Equal(t, backend.ID, respDept.ID)
	assert.Equal(t, developers.ID, *respDept.ParentID)
}

func TestDepartmentUsecase_UpdateDepartment_CyclicDependency_Self(t *testing.T) {
	deptRepo := new(repoMocks.DepartmentRepository)
	s := usecase.NewDepartmentService(deptRepo)

	ctx := context.Background()
	_, _, backend, _ := TwoTrees()

	deptRepo.
		On("GetByID", ctx, backend.ID).
		Return(backend, nil)

	// Пытаемся сделать backend родителем самого себя
	_, err := s.UpdateDepartment(ctx, backend.ID, nil, &backend.ID)

	require.ErrorIs(t, err, orgerror.ErrCyclicDependency)
}

func TestDepartmentUsecase_UpdateDepartment_CyclicDependency_Subtree(t *testing.T) {
	deptRepo := new(repoMocks.DepartmentRepository)
	s := usecase.NewDepartmentService(deptRepo)

	ctx := context.Background()
	_, _, backend, offchain := TwoTrees()

	deptRepo.
		On("GetByID", ctx, backend.ID).
		Return(backend, nil)
	deptRepo.On("GetSubtreeIDs", ctx, backend.ID).Return([]uint{backend.ID, offchain.ID}, nil)

	// Пытаемся переместить backend внутрь его дочернего offchain
	_, err := s.UpdateDepartment(ctx, backend.ID, nil, &offchain.ID)

	require.ErrorIs(t, err, orgerror.ErrCyclicDependency)
}

func TestDepartmentUsecase_UpdateDepartment_NotFound(t *testing.T) {
	deptRepo := new(repoMocks.DepartmentRepository)
	s := usecase.NewDepartmentService(deptRepo)

	ctx := context.Background()

	deptRepo.
		On("GetByID", ctx, uint(99)).
		Return(nil, orgerror.ErrDepartmentNotFound)

	_, err := s.UpdateDepartment(ctx, 99, nil, nil)

	require.ErrorIs(t, err, orgerror.ErrDepartmentNotFound)
}

// --- DelDepartment ---

func TestDepartmentUsecase_DelDepartment_Success_Cascade(t *testing.T) {
	deptRepo := new(repoMocks.DepartmentRepository)

	s := usecase.NewDepartmentService(
		deptRepo,
	)

	ctx := context.Background()

	_, _, backend, offchain := TwoTrees()

	deptRepo.
		On("Exists", ctx, backend.ID).
		Return(true, nil)
	deptRepo.
		On("GetSubtreeIDs", ctx, backend.ID).
		Return([]uint{backend.ID, offchain.ID}, nil)
	deptRepo.
		On("DeleteByIDs", ctx, []uint{backend.ID, offchain.ID}).
		Return(nil)

	err := s.DelDepartment(ctx, backend.ID, usecase.CascadeMode, nil)

	require.NoError(t, err)
}

func TestDepartmentUsecase_DelDepartment_NotFound(t *testing.T) {
	deptRepo := new(repoMocks.DepartmentRepository)
	s := usecase.NewDepartmentService(deptRepo)

	ctx := context.Background()

	deptRepo.
		On("Exists", ctx, uint(99)).
		Return(false, nil)

	err := s.DelDepartment(ctx, 99, usecase.CascadeMode, nil)

	require.ErrorIs(t, err, orgerror.ErrDepartmentNotFound)
}

func TestDepartmentUsecase_DelDepartment_Reassign_Success(t *testing.T) {
	deptRepo := new(repoMocks.DepartmentRepository)
	s := usecase.NewDepartmentService(deptRepo)

	ctx := context.Background()
	developers, _, backend, offchain := TwoTrees()

	subtreeIDs := []uint{backend.ID, offchain.ID}

	deptRepo.
		On("Exists", ctx, backend.ID).
		Return(true, nil)
	deptRepo.
		On("GetSubtreeIDs", ctx, backend.ID).
		Return(subtreeIDs, nil)
	deptRepo.
		On("Exists", ctx, developers.ID).
		Return(true, nil)
	deptRepo.
		On("MoveEmployees", ctx, subtreeIDs, developers.ID).
		Return(nil)
	deptRepo.
		On("DeleteByIDs", ctx, subtreeIDs).
		Return(nil)

	err := s.DelDepartment(ctx, backend.ID, usecase.ReassignMode, &developers.ID)

	require.NoError(t, err)
}

func TestDepartmentUsecase_DelDepartment_Reassign_TargetInSubtree(t *testing.T) {
	deptRepo := new(repoMocks.DepartmentRepository)
	s := usecase.NewDepartmentService(deptRepo)

	ctx := context.Background()
	_, _, backend, offchain := TwoTrees()

	subtreeIDs := []uint{backend.ID, offchain.ID}

	deptRepo.
		On("Exists", ctx, backend.ID).
		Return(true, nil)
	deptRepo.
		On("GetSubtreeIDs", ctx, backend.ID).
		Return(subtreeIDs, nil)
	deptRepo.
		On("Exists", ctx, offchain.ID).
		Return(true, nil)

	// Пытаемся reassign в offchain, который сам входит в удаляемое поддерево
	err := s.DelDepartment(ctx, backend.ID, usecase.ReassignMode, &offchain.ID)

	assert.Error(t, err)
}
