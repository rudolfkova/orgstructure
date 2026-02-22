// Package main ...
package main

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"orgstructure/internal/config"
	gormdb "orgstructure/internal/repository/gorm"
	gormrepo "orgstructure/internal/repository/gorm/repository"
	orgserver "orgstructure/internal/server"
	"orgstructure/internal/server/handler"
	"orgstructure/internal/server/middleware"
	"orgstructure/internal/usecase"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg, err := config.ParseConfig()
	if err != nil {
		log.Fatal(err)
	}
	logger := config.NewLogger(&cfg)

	db, err := gormdb.NewPostgresDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}

	deptRepo := gormrepo.NewDepartmentRepository(db)
	empRepo := gormrepo.NewEmployeeRepository(db)

	deptService := usecase.NewDepartmentService(deptRepo)
	empService := usecase.NewEmployeeService(empRepo, deptRepo)

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
	wrappedMux := middleware.Apply(mux)

	httpServer := &http.Server{
		Addr:         cfg.BindAddr,
		Handler:      wrappedMux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		logger.Info("starting server", slog.String("addr", cfg.BindAddr))
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}()

	<-quit
	logger.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		log.Fatalf("server forced to shutdown: %v", err)
	}

	logger.Info("server stopped")

}
