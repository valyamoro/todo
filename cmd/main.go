package main

import (
	"log"
	"net/http"
	"time"
	"github.com/valyamoro/TODO/internal/repository/psql"
	"github.com/valyamoro/TODO/internal/service"
	"github.com/valyamoro/TODO/internal/transport/rest"
	"github.com/valyamoro/TODO/pkg/database"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"os"
	"strconv"
)

func main() {
	logger, err := initializeLogger()

	if err != nil {
		log.Fatal("failed to create zap logger", err)
	}

	defer logger.Sync()

	host := os.Getenv("DB_HOST")
	port, _ := strconv.Atoi(os.Getenv("DB_PORT"))
	username := os.Getenv("DB_USERNAME")
	dbName := os.Getenv("DB_NAME")
	sslMode := os.Getenv("DB_SSLMODE")
	password := os.Getenv("DB_PASSWORD")

	db, err := database.NewPostgresConnection(database.ConnectionInfo{
		Host:     host,
		Port:     port,
		Username: username,
		DBName:   dbName,
		SSLMode:  sslMode,
		Password: password,
	})

	if err != nil {
		logger.Fatal("failed to connnect to database", zap.Error(err))
	}

	defer db.Close()

	tasksRepo := psql.NewTasks(db)
	tasksService := service.NewTasks(tasksRepo)
	handler:= rest.NewHandler(tasksService)

	r := handler.InitRouter(logger)

	srv := &http.Server{
		Addr: ":8080",
		Handler: r,
	}

	logger.Info("SERVER STARTED", zap.String("time", time.Now().Format(time.RFC3339)))

	if err := srv.ListenAndServe(); err != nil {
		logger.Fatal("server failed", zap.Error(err))
	}
}

func initializeLogger() (*zap.Logger, error) {
	var logger *zap.Logger 
	var err error 

	if os.Getenv("APP_ENV") == "development" {
		logger, err = zap.NewDevelopment() 
	} else {
		logger, err = zap.NewProduction()
	}

	return logger, err
}
