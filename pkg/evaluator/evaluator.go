package evaluator

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/f1rezy/distributed-calculator/internal/work"
	"go.nanasi880.dev/rpn"
)

var (
	AddTime int
	SubTime int
	MulTime int
	DivTime int
)

type Expression struct {
	value1   int
	value2   int
	operator string
	ch       chan int
}

func (e *Expression) Task() {
	switch e.operator {
	case "+":
		time.Sleep(time.Duration(AddTime) * time.Second)
		e.ch <- e.value1 + e.value2
	case "-":
		time.Sleep(time.Duration(SubTime) * time.Second)
		e.ch <- e.value1 - e.value2
	case "*":
		time.Sleep(time.Duration(MulTime) * time.Second)
		e.ch <- e.value1 * e.value2
	case "/":
		time.Sleep(time.Duration(DivTime) * time.Second)
		e.ch <- e.value1 / e.value2
	}
}

func Evaluate(expr string, pool *work.Pool) (int, error) {
	r, err := rpn.Parse(expr)
	if err != nil {
		return 0, err
	}
	stack := make([]int, 0)

	for _, token := range r.Tokens {
		if num, err := strconv.Atoi(token); err == nil {
			stack = append(stack, num)
		} else {
			if len(stack) < 2 {
				return 0, fmt.Errorf("invalid RPN expression")
			}

			value2 := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			value1 := stack[len(stack)-1]
			stack = stack[:len(stack)-1]

			ch := make(chan int)
			pool.Run(&Expression{value1: value1, value2: value2, operator: token, ch: ch})
			select {
			case <-context.Background().Done():
				return 0, nil
			case res := <-ch:
				stack = append(stack, res)
			}
		}
	}

	return stack[0], nil
}
