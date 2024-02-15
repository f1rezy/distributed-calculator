package main

import (
	"context"
	"os"

	"github.com/f1rezy/distributed-calculator/internal/application"
	"github.com/ilyakaznacheev/cleanenv"
)

func main() {
	ctx := context.Background()
	os.Exit(mainWithExitCode(ctx))
}

func mainWithExitCode(ctx context.Context) int {
	var cfg application.Config
	cleanenv.ReadEnv(&cfg)

	app := application.New(cfg)

	return app.Run(ctx)
}
