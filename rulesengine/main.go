package main

import (
	"fmt"

	"github.com/gorules/zen-go"
)

func main() {
	r1, err := zen.EvaluateExpression[int]("1 + 1", nil) // 2
	fmt.Println(r1, err)
	r2, err := zen.EvaluateExpression[int]("a + b", map[string]any{"a": 10, "b": 20}) // 30
	fmt.Println(r2, err)

	r3, err := zen.EvaluateUnaryExpression("> 10", map[string]any{"$": 5}) // false
	fmt.Println(r3, err)
	r4, err := zen.EvaluateUnaryExpression("> 10", map[string]any{"$": 15}) // true
	fmt.Println(r4, err)
}
