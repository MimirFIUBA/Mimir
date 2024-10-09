package config

import (
	"fmt"
	"mimir/triggers"
	"strconv"
	"strings"
	"unicode"
)

//Para parsear las condiciones vamos a construir un Árbol de Sintaxis Abstracta (AST)

type TokenType int

const (
	TOKEN_AND TokenType = iota
	TOKEN_OR
	TOKEN_COMPARE
	TOKER_AVG
	TOKEN_LPAREN // (
	TOKEN_RPAREN // )
	TOKEN_IDENT  // $(sensorId)
	TOKEN_OP     // >, <, ==, etc.
	TOKEN_NUMBER // 10, 5, etc.
	TOKEN_END    // End of input
)

const (
	OR_STRING      = "OR"
	AND_STRING     = "AND"
	AVERAGE_STRING = "AVG"
)

type Token struct {
	Type  TokenType
	Value string
}

func (t *Token) getValue() string {
	switch t.Type {
	case TOKEN_IDENT:
		return t.Value[2 : len(t.Value)-1]
	default:
		return t.Value
	}
}

func Tokenize(input string) []Token {
	var tokens []Token
	input = strings.TrimSpace(input)
	i := 0
	for i < len(input) {
		switch {
		case input[i] == '(':
			tokens = append(tokens, Token{Type: TOKEN_LPAREN, Value: "("})
			i++
		case input[i] == ')':
			tokens = append(tokens, Token{Type: TOKEN_RPAREN, Value: ")"})
			i++
		case unicode.IsDigit(rune(input[i])):
			// Parse number
			start := i
			for i < len(input) && (unicode.IsDigit(rune(input[i])) || unicode.IsPunct(rune(input[i]))) {
				i++
			}
			tokens = append(tokens, Token{Type: TOKEN_NUMBER, Value: input[start:i]})
		case input[i] == '>':
			tokens = append(tokens, Token{Type: TOKEN_OP, Value: ">"})
			i++
		case input[i] == '<':
			tokens = append(tokens, Token{Type: TOKEN_OP, Value: "<"})
			i++
		case input[i] == '=' && i+1 < len(input) && input[i+1] == '=':
			tokens = append(tokens, Token{Type: TOKEN_OP, Value: "=="})
			i += 2
		case strings.HasPrefix(input[i:], AND_STRING):
			tokens = append(tokens, Token{Type: TOKEN_AND, Value: AND_STRING})
			i += 3
		case strings.HasPrefix(input[i:], OR_STRING):
			tokens = append(tokens, Token{Type: TOKEN_OR, Value: OR_STRING})
			i += 2
		case strings.HasPrefix(input[i:], AVERAGE_STRING):
			tokens = append(tokens, Token{Type: TOKER_AVG, Value: AVERAGE_STRING})
			i += 3
		case input[i] == '$' && input[i+1] == '(':
			// Parse identifier like $(sensorId)
			start := i
			i += 2
			for i < len(input) && input[i] != ')' {
				i++
			}
			tokens = append(tokens, Token{Type: TOKEN_IDENT, Value: input[start : i+1]})
			i++ // consume the closing parenthesis
		default:
			i++
		}
	}
	tokens = append(tokens, Token{Type: TOKEN_END})
	return tokens
}

// ParserState holds the current state of the parser
type ParserState struct {
	tokens []Token
	pos    int
}

// Current returns the current token
func (p *ParserState) Current() Token {
	if p.pos >= len(p.tokens) {
		return Token{Type: TOKEN_END}
	}
	return p.tokens[p.pos]
}

// Advance moves to the next token
func (p *ParserState) Advance() {
	if p.pos < len(p.tokens) {
		p.pos++
	}
}

// ParseCondition parses the tokens into an AST
func ParseCondition(tokens []Token) (triggers.Condition, error) {
	state := &ParserState{tokens: tokens, pos: 0}
	return parseExpression(state)
}

// parseExpression parses AND/OR expressions
func parseExpression(state *ParserState) (triggers.Condition, error) {
	fmt.Println("parseExpression", state)
	left, err := parsePrimary(state)
	if err != nil {
		fmt.Println("err ", err)
		return nil, err
	}

	for state.Current().Type == TOKEN_AND || state.Current().Type == TOKEN_OR {
		operator := state.Current()
		state.Advance()
		right, err := parsePrimary(state)
		if err != nil {
			return nil, err
		}
		if operator.Type == TOKEN_AND {
			left = triggers.NewAndCondition([]triggers.Condition{left, right})
		} else if operator.Type == TOKEN_OR {
			left = triggers.NewOrCondition([]triggers.Condition{left, right})
		}
	}
	return left, nil
}

// parsePrimary parses individual conditions like $(sensorId) > 10
func parsePrimary(state *ParserState) (triggers.Condition, error) {
	fmt.Println("parsePrimary", state)
	token := state.Current()
	fmt.Println("token", token)
	switch token.Type {
	case TOKEN_IDENT:
		state.Advance()
		operator := state.Current()
		if operator.Type != TOKEN_OP {
			return nil, fmt.Errorf("expected operator after identifier")
		}
		state.Advance()
		right := state.Current()
		if right.Type != TOKEN_NUMBER {
			return nil, fmt.Errorf("expected number after operator")
		}
		state.Advance()
		value, err := strconv.ParseFloat(right.Value, 64)
		if err != nil {
			panic("Cannot convert string to float")
		}

		compareCondition := triggers.NewCompareCondition(operator.Value, value)
		compareCondition.SetSenderId(token.getValue())
		return compareCondition, nil

	case TOKEN_LPAREN:
		state.Advance()
		node, err := parseExpression(state)
		if err != nil {
			return nil, err
		}
		if state.Current().Type != TOKEN_RPAREN {
			return nil, fmt.Errorf("expected closing parenthesis")
		}
		state.Advance()
		return node, nil
	default:
		return nil, fmt.Errorf("unexpected token: %v", token)
	}
}
