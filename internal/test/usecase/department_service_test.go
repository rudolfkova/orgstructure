// Package usecase_test ...
package usecase_test

import (
	"context"
	"orgstructure/internal/model"
	"orgstructure/internal/service"
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

	err := s.DelDepartment(ctx, backend.ID, service.CascadeMode, nil)

	require.NoError(t, err)
}
