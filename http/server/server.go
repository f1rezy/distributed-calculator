package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/f1rezy/distributed-calculator/http/server/handler"
	"github.com/f1rezy/distributed-calculator/internal/work"
	"github.com/f1rezy/distributed-calculator/pkg/db"
	"github.com/f1rezy/distributed-calculator/pkg/evaluator"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func new(ctx context.Context,
	logger *zap.Logger,
	workerPool *work.Pool,
	db *gorm.DB,
) (http.Handler, error) {
	muxHandler, err := handler.New(ctx, workerPool, db)
	if err != nil {
		return nil, fmt.Errorf("handler initialization error: %w", err)
	}
	muxHandler = handler.Decorate(muxHandler, loggingMiddleware(logger))

	return muxHandler, nil
}

func Run(
	ctx context.Context,
	logger *zap.Logger,
	maxGoroutines int,
	database *gorm.DB,
) (func(context.Context) error, error) {
	workerPool := work.New(maxGoroutines)

	muxHandler, err := new(ctx, logger, workerPool, database)
	if err != nil {
		return nil, err
	}

	go func() {
		var expressions []db.Expression

		database.Find(&expressions)

		for _, expression := range expressions {
			if expression.Status != "in_progress" {
				continue
			}
			go func(expression db.Expression) {
				res, err := evaluator.Evaluate(expression.Expression, workerPool)
				if err != nil {
					expression.Status = "error"
				} else {
					expression.Status = "ok"
					expression.Result = fmt.Sprint(res)
				}
				expression.EvaluatedAt = time.Now()
				database.Save(&expression)
			}(expression)
		}
	}()

	srv := &http.Server{Addr: ":8080", Handler: muxHandler}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			logger.Error("ListenAndServe",
				zap.String("err", err.Error()))
		}
	}()

	return srv.Shutdown, nil
}

func loggingMiddleware(logger *zap.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			next.ServeHTTP(w, r)

			duration := time.Since(start)
			logger.Info("HTTP request",
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.Duration("duration", duration),
			)
		})
	}
}
