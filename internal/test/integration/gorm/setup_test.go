// Package integration_test ...
package integration_test

import (
	"database/sql"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"

	gormdb "orgstructure/internal/infrastructure/gorm"
	gormrepo "orgstructure/internal/infrastructure/gorm/repository"
	orgserver "orgstructure/internal/server"
	"orgstructure/internal/server/handler"
	"orgstructure/internal/server/middleware"
	"orgstructure/internal/usecase"
)

type testApp struct {
	handler http.Handler
}

var app *testApp

func TestMain(m *testing.M) {
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		os.Exit(0)
	}

	sqlDB, err := sql.Open("postgres", dsn)
	if err != nil {
		panic("failed to open sql db: " + err.Error())
	}
	if err := sqlDB.Ping(); err != nil {
		panic("failed to ping test db: " + err.Error())
	}

	if err := goose.SetDialect("postgres"); err != nil {
		panic("goose set dialect: " + err.Error())
	}
	if err := goose.Up(sqlDB, "../../../../migrations"); err != nil {
		panic("goose up: " + err.Error())
	}

	_, _ = sqlDB.Exec("DELETE FROM employees")
	_, _ = sqlDB.Exec("DELETE FROM departments")

	db, err := gormdb.NewPostgresDB(dsn)
	if err != nil {
		panic("failed to connect gorm: " + err.Error())
	}

	deptRepo := gormrepo.NewDepartmentRepository(db)
	empRepo := gormrepo.NewEmployeeRepository(db)

	deptService := usecase.NewDepartmentService(deptRepo)
	empService := usecase.NewEmployeeService(empRepo, deptRepo)

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	srv := orgserver.NewServer(logger)

	deptHandler := handler.NewDepartmentHandler(deptService, srv)
	empHandler := handler.NewEmployeeHandler(empService, srv)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /departments/", deptHandler.HandleCreateDepartment())
	mux.HandleFunc("GET /departments/{id}", deptHandler.HandleGetDepartment())
	mux.HandleFunc("PATCH /departments/{id}", deptHandler.HandleUpdateDepartment())
	mux.HandleFunc("DELETE /departments/{id}", deptHandler.HandleDelDepartment())
	mux.HandleFunc("POST /departments/{id}/employees/", empHandler.HandleCreateEmployeeInDepartment())

	middleware.Use(middleware.RequestID)
	app = &testApp{handler: middleware.Apply(mux)}

	code := m.Run()

	_, _ = sqlDB.Exec("DELETE FROM employees")
	_, _ = sqlDB.Exec("DELETE FROM departments")
	_ = sqlDB.Close()

	os.Exit(code)
}

func (a *testApp) do(req *http.Request) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	a.handler.ServeHTTP(rec, req)
	return rec
}
