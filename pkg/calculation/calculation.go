package calculation

import (
	"strconv"
)

// Calc вычисляет арифметическое выражение, заданное в виде строки, и возвращает результат.
// Поддерживаются операции +, -, *, / и круглые скобки.
func Calc(expression string) (float64, error) {
	var operations []rune // Список операторов
	var numbers []float64 // Список чисел
	var count int         // Счётчик для проверки скобок
	var i int             // Индекс для перебора символов выражения

	// Проверяем правильность расстановки скобок в выражении.
	for _, ch := range expression {
		switch ch {
		case '(':
			count++
		case ')':
			count--
		}
		if count < 0 {
			return 0.0, ErrInvalidParentheses
		}
	}
	if count != 0 {
		return 0.0, ErrInvalidParentheses
	}

	// Парсинг выражения.
	for i < len(expression) {
		ch := expression[i]

		// Пропускаем пробелы.
		if ch == ' ' {
			i++
			continue
		}

		// Если цифра или точка, считываем число (поддержка многозначных чисел и чисел с плавающей запятой).
		if (ch >= '0' && ch <= '9') || ch == '.' {
			start := i
			for i < len(expression) && ((expression[i] >= '0' && expression[i] <= '9') || expression[i] == '.') {
				i++
			}
			numStr := expression[start:i]
			num, err := strconv.ParseFloat(numStr, 64)
			if err != nil {
				return 0.0, ErrInvalidExpression
			}
			numbers = append(numbers, num)
			continue
		}

		// Если оператор, добавляем его в список операторов.
		if ch == '+' || ch == '-' || ch == '*' || ch == '/' {
			operations = append(operations, rune(ch))
			i++
			continue
		}

		// Если открывающая скобка, находим соответствующую закрывающую и рекурсивно вычисляем выражение внутри скобок.
		if ch == '(' {
			count = 1
			start := i + 1
			i++
			for i < len(expression) && count > 0 {
				if expression[i] == '(' {
					count++
				} else if expression[i] == ')' {
					count--
				}
				i++
			}
			if count != 0 {
				return 0.0, ErrInvalidParentheses
			}
			// Рекурсивный вызов для вычисления выражения внутри скобок.
			result, err := Calc(expression[start : i-1])
			if err != nil {
				return 0.0, err
			}
			numbers = append(numbers, result)
			continue
		}

		// Если символ не распознан, возвращаем ошибку.
		return 0.0, ErrInvalidExpression
	}

	// Проверяем корректность количества чисел и операторов.
	if len(numbers) != len(operations)+1 {
		return 0.0, ErrInvalidExpression
	}

	// Вычисляем операции с высоким приоритетом (* и /).
	for i := 0; i < len(operations); {
		if operations[i] == '*' || operations[i] == '/' {
			var result float64
			switch operations[i] {
			case '*':
				result = numbers[i] * numbers[i+1]
			case '/':
				if numbers[i+1] == 0 {
					return 0.0, ErrDivisionByZero
				}
				result = numbers[i] / numbers[i+1]
			}
			// Обновляем списки чисел и операторов после вычисления.
			numbers = append(numbers[:i], append([]float64{result}, numbers[i+2:]...)...)
			operations = append(operations[:i], operations[i+1:]...)
		} else {
			i++
		}
	}

	// Вычисляем операции с низким приоритетом (+ и -).
	for i := 0; i < len(operations); {
		var result float64
		switch operations[i] {
		case '+':
			result = numbers[i] + numbers[i+1]
		case '-':

			result = numbers[i] - numbers[i+1]
		}
		// Обновляем списки чисел и операторов после вычисления.
		numbers = append(numbers[:i], append([]float64{result}, numbers[i+2:]...)...)
		operations = append(operations[:i], operations[i+1:]...)
	}

	// Возвращаем окончательный результат.
	return numbers[0], nil
}
