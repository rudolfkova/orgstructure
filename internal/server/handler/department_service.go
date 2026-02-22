// Package handler транспортный слой. Обрабатывает запрос, передаёт готовые валидированные данные в service.
// Также обрабатывает ошибки, подготавливает ответ пользователю.
package handler

import (
	"context"
	"orgstructure/internal/model"
)

// DepartmentService определяет контракт для обработки запросов. Работает с подразделениями.
// Внутри получает данные из БД, обрабатывает их и возвращает данные, необходимые для ответа.
type DepartmentService interface {
	// CreateDepartment создаёт подразделение. Если указать родителя - создаёт дерево.
	// Если не указать родителя, делает созданное подразделение корнем в БД.
	// Возвращает созданное подразделение.
	CreateDepartment(ctx context.Context, name string, parentID *uint) (*model.Department, error)

	// GetDepartment возвращает дерево подразделений.
	// dept - глубина возвращаемого дерева.
	// includeEmployees - если true, возвращает дерево и сотрудников в подразделениях.
	// Корень дерева - подразделение с deptID = id.
	GetDepartment(ctx context.Context, depth int, includeEmployees bool, id uint) (*model.DepartmentTree, error)

	// UpdateDepartment редактирует сотрудника. Изменяет имя и родителя. Оба поля не обязательны.
	// Если и name и paretnID - nil, тогда метод ничего не делает.
	UpdateDepartment(ctx context.Context, id uint, name *string, parentID *uint) (*model.Department, error)
	DelDepartment(ctx context.Context, id uint, mode string, reassignToDepartmentID *uint) error
}
