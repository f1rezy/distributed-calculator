package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/f1rezy/distributed-calculator/internal/work"
	"github.com/f1rezy/distributed-calculator/pkg/db"
	"github.com/f1rezy/distributed-calculator/pkg/evaluator"
	"gorm.io/gorm"
)

type Decorator func(http.Handler) http.Handler

type Worker struct {
	pool *work.Pool
	db   *gorm.DB
}

type ExpressionJSON struct {
	Expression string `json:"expression"`
}

type Operation struct {
	Operator      string `json:"operator"`
	ExecutionTime int    `json:"execution_time"`
}

type OperationsTime struct {
	AddTime int `json:"add_time"`
	SubTime int `json:"sub_time"`
	MulTime int `json:"mul_time"`
	DivTime int `json:"div_time"`
}

func New(ctx context.Context,
	workerPool *work.Pool,
	db *gorm.DB,
) (http.Handler, error) {
	serveMux := http.NewServeMux()

	worker := Worker{
		workerPool,
		db,
	}

	serveMux.HandleFunc("/expression", worker.expression)
	serveMux.HandleFunc("/expressions", worker.getAllEpressions)
	serveMux.HandleFunc("/operations", worker.operations)

	return serveMux, nil
}

func Decorate(next http.Handler, ds ...Decorator) http.Handler {
	decorated := next
	for d := len(ds) - 1; d >= 0; d-- {
		decorated = ds[d](decorated)
	}

	return decorated
}

func (wr *Worker) expression(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		expressionId, err := strconv.Atoi(id)
		if err != nil {
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		var expression db.Expression

		result := wr.db.First(&expression, expressionId)
		if result.Error != nil {
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(expression)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	case http.MethodPost:
		decoder := json.NewDecoder(r.Body)
		var expr ExpressionJSON
		err := decoder.Decode(&expr)
		if err != nil {
			panic(err)
		}
		expression := db.Expression{
			Expression:  expr.Expression,
			CreatedAt:   time.Now(),
			EvaluatedAt: time.Time{},
		}
		result := wr.db.Create(&expression)
		if result.Error != nil {
			fmt.Fprintln(w, result.Error.Error())
		}
		go func() {
			res, err := evaluator.Evaluate(expr.Expression, wr.pool)
			if err != nil {
				expression.Status = "error"
			} else {
				expression.Status = "ok"
				expression.Result = fmt.Sprint(res)
			}
			expression.EvaluatedAt = time.Now()
			wr.db.Save(&expression)
		}()
		fmt.Fprintf(w, "Приняли к обработке\nid: %d", expression.ID)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (wr *Worker) getAllEpressions(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		var expressions []db.Expression

		wr.db.Find(&expressions)

		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(expressions)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (wr *Worker) operations(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		operations := []Operation{
			{Operator: "+", ExecutionTime: evaluator.AddTime},
			{Operator: "-", ExecutionTime: evaluator.SubTime},
			{Operator: "*", ExecutionTime: evaluator.MulTime},
			{Operator: "/", ExecutionTime: evaluator.DivTime},
		}

		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(operations)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	case http.MethodPost:
		decoder := json.NewDecoder(r.Body)
		var operationsTime OperationsTime
		err := decoder.Decode(&operationsTime)
		if err != nil {
			panic(err)
		}

		evaluator.AddTime = operationsTime.AddTime
		evaluator.SubTime = operationsTime.SubTime
		evaluator.MulTime = operationsTime.MulTime
		evaluator.DivTime = operationsTime.DivTime
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
