package config

import (
	"fmt"
	"mimir/triggers"
	"strconv"
	"strings"
	"time"
	"unicode"
)

//Para parsear las condiciones vamos a construir un Árbol de Sintaxis Abstracta (AST)

type TokenType int

const (
	TOKEN_AND TokenType = iota
	TOKEN_OR
	TOKEN_COMPARE
	TOKEN_AVG
	TOKEN_LPAREN // (
	TOKEN_RPAREN // )
	TOKEN_LBRACE // [
	TOKEN_RBRACE // ]
	TOKEN_IDENT  // $(topic/subtopic)
	TOKEN_OP     // >, <, ==, etc.
	TOKEN_NUMBER // 10, 5, etc.
	TOKEN_STRING // hello etc. TODO: see if we need to add ""
	TOKEN_END    // End of input
)

// Reserved words
const (
	OR_STRING      = "OR"
	AND_STRING     = "AND"
	AVERAGE_STRING = "AVG"
)

const (
	AVG_COND_MIN_AMOUNT_DEFAULT        = 0
	AVG_COND_MAX_AMOUNT_DEFAULT        = 10
	AVG_COND_TIMEFRAME_SECONDS_DEFAULT = 10
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
		case input[i] == '[':
			tokens = append(tokens, Token{Type: TOKEN_LBRACE, Value: "["})
			i++
		case input[i] == ']':
			tokens = append(tokens, Token{Type: TOKEN_RBRACE, Value: "]"})
			i++
		case unicode.IsDigit(rune(input[i])):
			start := i
			for i < len(input) && (unicode.IsDigit(rune(input[i])) || input[i] == '.') {
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
			tokens = append(tokens, Token{Type: TOKEN_AVG, Value: AVERAGE_STRING})
			i += 3
		case input[i] == '$' && input[i+1] == '(':
			start := i
			i += 2
			for i < len(input) && input[i] != ')' {
				i++
			}
			tokens = append(tokens, Token{Type: TOKEN_IDENT, Value: input[start : i+1]})
			i++
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
		//TODO: call parse condition function
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
	case TOKEN_AVG:
		averageCondition, err := parseAverageCondition(state)
		state.Advance()
		return averageCondition, err
	default:
		return nil, fmt.Errorf("unexpected token: %v", token)
	}
}

func parseAverageCondition(state *ParserState) (triggers.Condition, error) {
	state.Advance()
	params, err := parseParameters(state)
	if err != nil {
		return nil, err
	}
	fmt.Println("params: ", params)

	state.Advance()
	metadata, err := parseMetadata(state)
	if err != nil {
		return nil, err
	}
	fmt.Println("metadata: ", metadata)

	state.Advance()
	condition, err := parseConditionForExpression(state)
	if err != nil {
		fmt.Println(err)
	}

	avgCondition := buildAverageCondition(params, metadata, condition)
	return avgCondition, nil
}

func parseConditionForExpression(state *ParserState) (triggers.Condition, error) {
	operator := state.Current()
	switch operator.Type {
	case TOKEN_OP:
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
		return compareCondition, nil
	default:
		return nil, fmt.Errorf("expected operator after identifier")
	}
}

func parseParameters(state *ParserState) ([]Token, error) {
	fmt.Println("parse parameters")
	currentToken := state.Current()
	tokens := make([]Token, 0)
	for currentToken.Type != TOKEN_RPAREN {
		fmt.Println("current token ", currentToken)
		state.Advance()
		currentToken = state.Current()
		switch currentToken.Type {
		case TOKEN_LPAREN:

			if currentToken.Type != TOKEN_IDENT {
				return nil, fmt.Errorf("wrong format for parameter, expecting identity")
			}
		case TOKEN_RPAREN:
			return tokens, nil
		case TOKEN_IDENT:
			tokens = append(tokens, currentToken)
		default:
			return nil, fmt.Errorf("wrong format for parameter, expecting identity")
		}
	}
	return nil, nil
}

func parseMetadata(state *ParserState) ([]string, error) {
	currentToken := state.Current()
	params := make([]string, 0)
	fmt.Println("parse metadata ", currentToken.Type)
	for currentToken.Type != TOKEN_RBRACE {
		switch currentToken.Type {
		case TOKEN_LBRACE:
			state.Advance()
		case TOKEN_NUMBER, TOKEN_STRING:
			params = append(params, currentToken.Value)
			state.Advance()
		default:
			return nil, fmt.Errorf("expected number or string for metadata")
		}

		currentToken = state.Current()
	}
	return params, nil
}

func buildAverageCondition(parameters []Token, metadata []string, condition triggers.Condition) triggers.Condition {
	minAmount, err := strconv.Atoi(metadata[0])
	if err != nil {
		minAmount = AVG_COND_MIN_AMOUNT_DEFAULT
	}

	maxAmount, err := strconv.Atoi(metadata[1])
	if err != nil {
		maxAmount = AVG_COND_MAX_AMOUNT_DEFAULT
	}

	timeFrame, err := strconv.Atoi(metadata[2])
	if err != nil {
		maxAmount = AVG_COND_TIMEFRAME_SECONDS_DEFAULT
	}

	//TODO: we should always expect one param
	senderId := parameters[0]

	avgCondition := triggers.NewAverageCondition(minAmount, maxAmount, time.Duration(timeFrame)*time.Second)
	avgCondition.SetSenderId(senderId.getValue())
	avgCondition.SetCondition(condition)
	return avgCondition
}
