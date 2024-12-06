package main

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
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
				return 0, errors.New("деление на ноль")
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
			return 0, errors.New("ожидалась закрывающая скобка")
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
	for p.pos < len(p.input) && (unicode.IsDigit(rune(p.input[p.pos])) || p.input[p.pos] == '.') {
		p.pos++
	}
	if start == p.pos {
		return 0, errors.New("ожидалось число")
	}
	numberStr := p.input[start:p.pos]
	number, err := strconv.ParseFloat(numberStr, 64)
	if err != nil {
		return 0, errors.New("некорректное число: " + numberStr)
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
