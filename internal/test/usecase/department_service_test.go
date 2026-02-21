// Package usecase_test ...
package usecase_test

import (
	"context"
	"orgstructure/internal/model"
	"orgstructure/internal/usecase"
	repoMocks "orgstructure/mocks/repository"
	"testing"

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
