// Package usecase ...
package usecase

import (
	"context"
	"orgstructure/internal/model"
)

// EmployeeRepository определяет контракт для работы с сотрудниками в хранилище.
type EmployeeRepository interface {
	// Create сохраняет нового сотрудника в хранилище.
	Create(ctx context.Context, emp *model.Employee) error

	// Проверяет, существует ли подразделение в БД.
	DepartmentExists(ctx context.Context, departmentID uint) (bool, error)
}
