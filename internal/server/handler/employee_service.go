// Package handler транспортный слой. Обрабатывает запрос, передаёт готовые валидированные данные в service.
// Также обрабатывает ошибки, подготавливает ответ пользователю.
package handler

import (
	"context"
	"orgstructure/internal/model"
	"time"
)

// EmployeeService определяет контракт для обработки запросов. Работает с сотрудниками.
// Внутри получает данные из БД, обрабатывает их и возвращает данные, необходимые для ответа.
type EmployeeService interface {
	// CreateEmployeeInDepartment создаёт сотрудника внутри подразделения. Возвращает созданного сотрудника.
	CreateEmployeeInDepartment(ctx context.Context, departmentID uint, fullName, position string, hiredAt *time.Time) (*model.Employee, error)
}
