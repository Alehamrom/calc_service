package calculation

import (
	"errors"
	"math"
	"testing"
)

func TestCalc(t *testing.T) {

	testCasesSuccess := []struct {
		name           string
		expression     string
		expectedResult float64
	}{
		{
			name:           "simple_addition",
			expression:     "1+1",
			expectedResult: 2,
		},
		{
			name:           "parentheses_priority",
			expression:     "(2+2)*2",
			expectedResult: 8,
		},
		{
			name:           "multiplication_priority",
			expression:     "2+2*2",
			expectedResult: 6,
		},
		{
			name:           "nested_parentheses",
			expression:     "2+(2*(2+3))",
			expectedResult: 12,
		},
		{
			name:           "simple_division",
			expression:     "1/2",
			expectedResult: 0.5,
		},
	}

	const epsilon = 1e-9

	for _, testCase := range testCasesSuccess {
		t.Run(testCase.name, func(t *testing.T) {
			val, err := Calc(testCase.expression)
			if err != nil {
				t.Fatalf("expression %s returned unexpected error: %v", testCase.expression, err)
			}
			if math.Abs(val-testCase.expectedResult) > epsilon {
				t.Errorf("expected %f, got %f", testCase.expectedResult, val)
			}
		})
	}

	// Тестовые случаи с ожидаемыми ошибками
	testCasesFail := []struct {
		name        string
		expression  string
		expectedErr error
	}{
		{
			name:        "invalid_expression_end_with_operator",
			expression:  "1+1*",
			expectedErr: ErrInvalidExpression,
		},
		{
			name:        "missing_closing_parenthesis",
			expression:  "(2+2*2",
			expectedErr: ErrInvalidParentheses,
		},
		{
			name:        "extra_closing_parenthesis",
			expression:  "2+2*2)-",
			expectedErr: ErrInvalidParentheses,
		},
		{
			name:        "division_by_zero",
			expression:  "2/0",
			expectedErr: ErrDivisionByZero,
		},
		{
			name:        "non_numeric_expression",
			expression:  "hello world 2+2*2",
			expectedErr: ErrInvalidExpression,
		},
	}

	for _, testCase := range testCasesFail {
		t.Run(testCase.name, func(t *testing.T) {
			_, err := Calc(testCase.expression)
			if !errors.Is(err, testCase.expectedErr) {
				t.Fatalf("for expression %s, expected error %v, got %v", testCase.expression, testCase.expectedErr, err)
			}
			// Нет необходимости проверять значение val при ошибке
		})
	}

}
