package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"unicode"
)

var (
	ErrUnexpectedSymbol     = errors.New("неожиданный символ")
	ErrExpectedClosingParen = errors.New("ожидалась закрывающая скобка")
	ErrExpectedNumber       = errors.New("ожидалось число")
	ErrInvalidNumber        = errors.New("некорректное число")
	ErrDivisionByZero       = errors.New("деление на ноль")
)

func Calc(expression string) (float64, error) {
	parser := &Parser{input: strings.TrimSpace(expression), pos: 0}
	result, err := parser.parseExpression()
	if err != nil {
		return 0, err
	}
	parser.skipSpaces()
	if parser.pos < len(parser.input) {
		return 0, errors.New("неожиданный символ: " + string(parser.input[parser.pos]))
	}
	return result, nil
}

type Parser struct {
	input string
	pos   int
}

func (p *Parser) parseExpression() (float64, error) {
	result, err := p.parseTerm()
	if err != nil {
		return 0, err
	}
	for {
		p.skipSpaces()
		if p.match('+') {
			p.pos++
			term, err := p.parseTerm()
			if err != nil {
				return 0, err
			}
			result += term
		} else if p.match('-') {
			p.pos++
			term, err := p.parseTerm()
			if err != nil {
				return 0, err
			}
			result -= term
		} else {
			break
		}
	}
	return result, nil
}

func (p *Parser) parseTerm() (float64, error) {
	result, err := p.parseFactor()
	if err != nil {
		return 0, err
	}
	for {
		p.skipSpaces()
		if p.match('*') {
			p.pos++
			factor, err := p.parseFactor()
			if err != nil {
				return 0, err
			}
			result *= factor
		} else if p.match('/') {
			p.pos++
			factor, err := p.parseFactor()
			if err != nil {
				return 0, err
			}
			if factor == 0 {
				return 0, ErrDivisionByZero
			}
			result /= factor
		} else {
			break
		}
	}
	return result, nil
}

func (p *Parser) parseFactor() (float64, error) {
	p.skipSpaces()
	if p.match('(') {
		p.pos++
		expr, err := p.parseExpression()
		if err != nil {
			return 0, err
		}
		p.skipSpaces()
		if !p.match(')') {
			return 0, ErrExpectedClosingParen
		}
		p.pos++
		return expr, nil
	}
	return p.parseNumber()
}

func (p *Parser) parseNumber() (float64, error) {
	p.skipSpaces()
	start := p.pos
	if p.match('+') || p.match('-') {
		p.pos++
	}
	dotCount := 0
	for p.pos < len(p.input) {
		ch := p.input[p.pos]
		if unicode.IsDigit(rune(ch)) {
			p.pos++
		} else if ch == '.' {
			if dotCount == 1 {
				break
			}
			dotCount++
			p.pos++
		} else {
			break
		}
	}
	if start == p.pos {
		return 0, ErrExpectedNumber
	}
	numberStr := p.input[start:p.pos]
	number, err := strconv.ParseFloat(numberStr, 64)
	if err != nil {
		return 0, ErrInvalidNumber
	}
	return number, nil
}

func (p *Parser) skipSpaces() {
	for p.pos < len(p.input) && unicode.IsSpace(rune(p.input[p.pos])) {
		p.pos++
	}
}

func (p *Parser) match(ch byte) bool {
	return p.pos < len(p.input) && p.input[p.pos] == ch
}

type Request struct {
	Expression string `json:"expression"`
}

type Response struct {
	Result string `json:"result,omitempty"`
	Error  string `json:"error,omitempty"`
}

func calculateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	var req Request
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	result, err := Calc(req.Expression)
	if err != nil {
		if errors.Is(err, ErrExpectedNumber) || errors.Is(err, ErrInvalidNumber) || errors.Is(err, ErrExpectedClosingParen) || strings.HasPrefix(err.Error(), ErrUnexpectedSymbol.Error()) {
			w.WriteHeader(http.StatusUnprocessableEntity)
			json.NewEncoder(w).Encode(Response{Error: "Expression is not valid"})
			return
		}

		if errors.Is(err, ErrDivisionByZero) {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(Response{Error: "Division by zero"})
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{Error: "Internal server error"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Response{Result: strconv.FormatFloat(result, 'f', -1, 64)})
}

func main() {
	http.HandleFunc("/api/v1/calculate", calculateHandler)
	log.Println("Server is running on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
