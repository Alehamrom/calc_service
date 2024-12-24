package calculation

import "errors"

var (
	ErrInvalidExpression  = errors.New("invalid expression")
	ErrDivisionByZero     = errors.New("division by zero")
	ErrInvalidParentheses = errors.New("invalid parentheses")
)
