// Package usecase_test ...
package usecase_test

import (
	"context"
	orgerror "orgstructure/internal/errors"
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

func TestEmployeeUsecase_CreateEmployee_EmptyFullName(t *testing.T) {
	deptRepo := new(repoMocks.DepartmentRepository)
	empRepo := new(repoMocks.EmployeeRepository)
	s := usecase.NewEmployeeService(empRepo, deptRepo)

	_, err := s.CreateEmployeeInDepartment(context.Background(), 1, "   ", "Developer", nil)

	require.ErrorIs(t, err, orgerror.ErrInvalidFullName)
}

func TestEmployeeUsecase_CreateEmployee_EmptyPosition(t *testing.T) {
	deptRepo := new(repoMocks.DepartmentRepository)
	empRepo := new(repoMocks.EmployeeRepository)
	s := usecase.NewEmployeeService(empRepo, deptRepo)

	_, err := s.CreateEmployeeInDepartment(context.Background(), 1, "Ivan Petrov", "   ", nil)

	require.ErrorIs(t, err, orgerror.ErrInvalidPosition)
}

func TestEmployeeUsecase_CreateEmployee_DepartmentNotFound(t *testing.T) {
	deptRepo := new(repoMocks.DepartmentRepository)
	empRepo := new(repoMocks.EmployeeRepository)
	s := usecase.NewEmployeeService(empRepo, deptRepo)

	ctx := context.Background()

	deptRepo.
		On("Exists", ctx, uint(99)).
		Return(false, nil)

	_, err := s.CreateEmployeeInDepartment(ctx, 99, "Ivan Petrov", "Developer", nil)

	require.ErrorIs(t, err, orgerror.ErrDepartmentNotFound)
}
