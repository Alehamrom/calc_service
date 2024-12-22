package evaluator

import (
	"errors"
	"math"
	"strconv"
	"strings"
	"unicode"
)

type Operator struct {
	Precedence    int
	Associativity string
}

var operators = map[string]Operator{
	"+": {Precedence: 2, Associativity: "left"},
	"-": {Precedence: 2, Associativity: "left"},
	"*": {Precedence: 3, Associativity: "left"},
	"/": {Precedence: 3, Associativity: "left"},
	"^": {Precedence: 4, Associativity: "right"},
}

func Evaluate(expression string) (float64, error) {
	tokens, err := tokenize(expression)
	if err != nil {
		return 0, err
	}

	rpn, err := toRPN(tokens)
	if err != nil {
		return 0, err
	}

	result, err := evalRPN(rpn)
	if err != nil {
		return 0, err
	}

	return result, nil
}

func tokenize(expr string) ([]string, error) {
	var tokens []string
	var number strings.Builder

	for i, ch := range expr {
		if unicode.IsSpace(ch) {
			continue
		}

		if unicode.IsDigit(ch) || ch == '.' {
			number.WriteRune(ch)
		} else if isOperator(string(ch)) || ch == '(' || ch == ')' {
			if number.Len() > 0 {
				tokens = append(tokens, number.String())
				number.Reset()
			}
			tokens = append(tokens, string(ch))
		} else {
			return nil, errors.New("invalid character in expression")
		}

		if ch == '.' {
			start := i - 1
			if start >= 0 && unicode.IsDigit(rune(expr[start])) {
			}
		}
	}

	if number.Len() > 0 {
		tokens = append(tokens, number.String())
	}

	return tokens, nil
}

func isOperator(op string) bool {
	_, exists := operators[op]
	return exists
}

func toRPN(tokens []string) ([]string, error) {
	var outputQueue []string
	var operatorStack []string

	for _, token := range tokens {
		if isNumber(token) {
			outputQueue = append(outputQueue, token)
		} else if isOperator(token) {
			for len(operatorStack) > 0 {
				top := operatorStack[len(operatorStack)-1]
				if isOperator(top) {
					curOp := operators[token]
					topOp := operators[top]
					if (curOp.Associativity == "left" && curOp.Precedence <= topOp.Precedence) ||
						(curOp.Associativity == "right" && curOp.Precedence < topOp.Precedence) {
						outputQueue = append(outputQueue, top)
						operatorStack = operatorStack[:len(operatorStack)-1]
						continue
					}
				}
				break
			}
			operatorStack = append(operatorStack, token)
		} else if token == "(" {
			operatorStack = append(operatorStack, token)
		} else if token == ")" {
			foundLeftParen := false
			for len(operatorStack) > 0 {
				top := operatorStack[len(operatorStack)-1]
				if top == "(" {
					foundLeftParen = true

					operatorStack = operatorStack[:len(operatorStack)-1]
					break
				}
				outputQueue = append(outputQueue, top)
				operatorStack = operatorStack[:len(operatorStack)-1]
			}
			if !foundLeftParen {
				return nil, errors.New("mismatched parentheses")
			}
		} else {
			return nil, errors.New("unknown token")
		}
	}

	for len(operatorStack) > 0 {
		top := operatorStack[len(operatorStack)-1]
		if top == "(" || top == ")" {
			return nil, errors.New("mismatched parentheses")
		}
		outputQueue = append(outputQueue, top)
		operatorStack = operatorStack[:len(operatorStack)-1]
	}

	return outputQueue, nil
}

func isNumber(token string) bool {
	_, err := strconv.ParseFloat(token, 64)
	return err == nil
}

func evalRPN(rpn []string) (float64, error) {
	var stack []float64

	for _, token := range rpn {
		if isNumber(token) {
			num, _ := strconv.ParseFloat(token, 64)
			stack = append(stack, num)
		} else if isOperator(token) {
			if len(stack) < 2 {
				return 0, errors.New("invalid expression")
			}
			b := stack[len(stack)-1]
			a := stack[len(stack)-2]
			stack = stack[:len(stack)-2]
			var res float64
			switch token {
			case "+":
				res = a + b
			case "-":
				res = a - b
			case "*":
				res = a * b
			case "/":
				if b == 0 {
					return 0, errors.New("division by zero")
				}
				res = a / b
			case "^":
				res = math.Pow(a, b)
			default:
				return 0, errors.New("unknown operator")
			}
			stack = append(stack, res)
		} else {
			return 0, errors.New("invalid token in RPN")
		}
	}

	if len(stack) != 1 {
		return 0, errors.New("invalid RPN expression")
	}

	return stack[0], nil
}
