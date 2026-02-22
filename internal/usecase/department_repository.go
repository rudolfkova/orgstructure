// Package usecase принимает запросы на обращение в БД, формирует запрос и обращается в БД.
// Возвращает данные, которые лежат в хранилище.
package usecase

// department_repository.go содержит репозиторий для работы с подразделениями в БД.

import (
	"context"
	"orgstructure/internal/model"
)

// DepartmentRepository определяет контракт для работы с подразделениями в хранилище.
type DepartmentRepository interface {
	// Create сохраняет новое подразделение в хранилище.
	Create(ctx context.Context, dept *model.Department) error

	// GetByID возвращает подразделение по ID.
	// Возвращает ErrDepartmentNotFound если подразделение не найдено.
	GetByID(ctx context.Context, id uint) (*model.Department, error)

	// Exists проверяет, существует ли подразделение с данным ID.
	Exists(ctx context.Context, id uint) (bool, error)

	// GetChildren возвращает прямых потомков подразделения (глубина 1).
	GetChildren(ctx context.Context, parentID uint) ([]model.Department, error)

	// GetEmployees возвращает сотрудников подразделения.
	// sortBy - поле сортировки: "created_at" или "full_name".
	GetEmployees(ctx context.Context, departmentID uint, sortBy string) ([]model.Employee, error)

	// Update сохраняет изменения подразделения.
	Update(ctx context.Context, dept *model.Department) error

	// Delete удаляет подразделение по ID.
	Delete(ctx context.Context, id uint) error

	// GetSubtreeIDs возвращает все ID подразделений в поддереве (включая корень)
	GetSubtreeIDs(ctx context.Context, rootID uint) ([]uint, error)

	// MoveEmployees переводит сотрудников из одного подразделения в другое
	MoveEmployees(ctx context.Context, fromDepartmentIDs []uint, toDepartmentID uint) error

	// DeleteByIDs удаляет подразделения по списку ID
	DeleteByIDs(ctx context.Context, ids []uint) error

	// ExistsName проверяет, существует ли подразделение с таким именем у данного родителя
	ExistsName(ctx context.Context, name string, parentID *uint, excludeID *uint) (bool, error)
}
