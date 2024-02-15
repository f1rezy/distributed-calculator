package application

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/f1rezy/distributed-calculator/http/server"
	"github.com/f1rezy/distributed-calculator/pkg/db"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Config struct {
	MaxGoroutines    int    `env:"MAX_GOROUTINES" env-default:"15"`
	PostgresHost     string `env:"POSTGRES_HOST" env-default:"localhost"`
	PostgresPort     string `env:"POSTGRES_PORT" env-default:"5432"`
	PostgresUser     string `env:"POSTGRES_USER" env-default:"admin"`
	PostgresPassword string `env:"POSTGRES_PASSWORD" env-default:"password"`
	PostgresDb       string `env:"POSTGRES_DB" env-default:"database"`
}

type Application struct {
	Cfg Config
}

func New(config Config) *Application {
	return &Application{
		Cfg: config,
	}
}

func (a *Application) Run(ctx context.Context) int {
	logger := setupLogger()

	database, err := db.ConnectToPostgreSQL(
		a.Cfg.PostgresHost,
		a.Cfg.PostgresPort,
		a.Cfg.PostgresUser,
		a.Cfg.PostgresPassword,
		a.Cfg.PostgresDb,
	)
	if err != nil {
		logger.Error(err.Error())

		return 1
	}

	if err = database.AutoMigrate(&db.Expression{}); err != nil {
		logger.Error(err.Error())

		return 1
	}

	shutDownFunc, err := server.Run(ctx, logger, a.Cfg.MaxGoroutines, database)
	if err != nil {
		logger.Error(err.Error())

		return 1
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	_, cancel := context.WithCancel(context.Background())

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	<-c
	cancel()
	shutDownFunc(ctx)

	return 0
}

func setupLogger() *zap.Logger {
	config := zap.NewDevelopmentConfig()

	config.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)

	logger, err := config.Build()
	if err != nil {
		fmt.Printf("Ошибка настройки логгера: %v\n", err)
	}

	return logger
}
