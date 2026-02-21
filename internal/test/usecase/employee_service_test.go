// Package usecase_test ...
package usecase_test

import (
	"context"
	"orgstructure/internal/usecase"
	repoMocks "orgstructure/mocks/repository"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEmployeeUseCase_Success(t *testing.T) {
	deptRepo := new(repoMocks.DepartmentRepository)
	empRepo := new(repoMocks.EmployeeRepository)

	s := usecase.NewEmployeeService(
		empRepo,
		deptRepo,
	)

	ctx := context.Background()
	_, _, backend, _ := TwoTrees()
	emp := backend.Employees[0]

	deptRepo.
		On("Exists", ctx, emp.DepartmentID).
		Return(true, nil)
	empRepo.
		On("Create", ctx, &emp).
		Return(nil)
	newEmp, err := s.CreateEmployeeInDepartment(ctx, emp.DepartmentID, emp.FullName, emp.Position, nil)

	require.NoError(t, err)
	assert.Equal(t, emp, *newEmp)

}
