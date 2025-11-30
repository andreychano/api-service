package main

import (
	"log/slog"
	"net/http"
	"os"

	transport "github.com/andreychano/api-service/internal/adapters/http"
	"github.com/andreychano/api-service/internal/adapters/storage/postgres"
	"github.com/andreychano/api-service/internal/config"
	"github.com/andreychano/api-service/internal/service"
)

func main() {
	// Setup structured logging (JSON format)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	cfg, err := config.LoadConfig()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	db, err := postgres.ConnectDB(cfg.FormatDSN())
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	slog.Info("database connected successfully", "host", cfg.DBHost)

	// Init Adapters
	questionRepo := postgres.NewQuestionRepo(db)
	answerRepo := postgres.NewAnswerRepo(db)

	// Init Services
	questionService := service.NewQuestionService(questionRepo)
	answerService := service.NewAnswerService(answerRepo, questionRepo)

	// Init Router
	router := transport.NewRouter(questionService, answerService)

	serverAddr := ":8080"
	slog.Info("starting server", "address", serverAddr)

	if err := http.ListenAndServe(serverAddr, router); err != nil {
		slog.Error("server failed", "error", err)
		os.Exit(1)
	}
}
