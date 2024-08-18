package main

import (
	"fmt"
	"strconv"
)

var keywords = map[string]TokenType{
	"and":    AND,
	"class":  CLASS,
	"else":   ELSE,
	"false":  FALSE,
	"for":    FOR,
	"fun":    FUN,
	"if":     IF,
	"nil":    NIL,
	"or":     OR,
	"print":  PRINT,
	"return": RETURN,
	"super":  SUPER,
	"this":   THIS,
	"true":   TRUE,
	"var":    VAR,
	"while":  WHILE,
}

type Scanner struct {
	source  string
	tokens  []Token
	start   int
	current int
	line    int
}

func (s *Scanner) ScanTokens() []Token {
	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken()
	}

	s.tokens = append(s.tokens, Token{
		Type:    EOF,
		lexeme:  "",
		literal: nil,
		line:    s.line,
	})
	return s.tokens
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func (s *Scanner) scanToken() {
	c := s.advance()
	switch c {
	case '(':
		s.addSingleCharacterToken(LEFT_PAREN)
		break
	case ')':
		s.addSingleCharacterToken(RIGHT_PAREN)
		break
	case '{':
		s.addSingleCharacterToken(LEFT_BRACE)
		break
	case '}':
		s.addSingleCharacterToken(RIGHT_BRACE)
		break
	case '.':
		s.addSingleCharacterToken(DOT)
		break
	case ',':
		s.addSingleCharacterToken(COMMA)
		break
	case '-':
		s.addSingleCharacterToken(MINUS)
		break
	case '+':
		s.addSingleCharacterToken(PLUS)
		break
	case ';':
		s.addSingleCharacterToken(SEMICOLON)
		break
	case '*':
		s.addSingleCharacterToken(STAR)
		break
	case '!':
		if s.match('=') {
			s.addSingleCharacterToken(BANG_EQUAL)
		} else {
			s.addSingleCharacterToken(EQUAL)
		}
	case '=':
		if s.match('=') {
			s.addSingleCharacterToken(EQUAL_EQUAL)
		} else {
			s.addSingleCharacterToken(EQUAL)
		}
	case '<':
		if s.match('=') {
			s.addSingleCharacterToken(LESS_EQUAL)
		} else {
			s.addSingleCharacterToken(LESS)
		}
	case '>':
		if s.match('=') {
			s.addSingleCharacterToken(GREATER_EQUAL)
		} else {
			s.addSingleCharacterToken(GREATER)
		}
	case '/':
		if s.match('/') {
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		} else if s.peek() == '*' {
			s.advance()
			for {
				if s.peek() == '*' {
					s.advance()
					if s.peek() == '/' {
						s.advance()
						break
					}
					continue
				}
				if s.isAtEnd() {
					break
				}
				s.advance()
			}
		} else {
			s.addSingleCharacterToken(SLASH)
		}
	case ' ':
	case '\r':
	case '\t':
		break
	case 'n':
		s.line++
		break
	case '"':
		s.string()
		break
	default:
		if s.isDigit(c) {
			s.number()
		} else if s.isAlpha(c) {
			s.identifier()
		} else {
			Error(s.line, "unexpected character")
		}
		break
	}
}

func (s *Scanner) isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func (s *Scanner) number() {
	for s.isDigit(s.peek()) {
		s.advance()
	}

	if s.peek() == '.' && s.isDigit(s.peekNext()) {
		s.advance()

		for s.isDigit(s.peek()) {
			s.advance()
		}
	}

	value, _ := strconv.ParseFloat(s.source[s.start:s.current], 10)

	s.addToken(NUMBER, value)
}

func (s *Scanner) identifier() {
	for s.isAlphaNumeric(s.peek()) {
		s.advance()
	}

	text := s.source[s.start:s.current]
	tokenType, ok := keywords[text]
	if !ok {
		tokenType = IDENTIFIER
	}

	s.addSingleCharacterToken(tokenType)
}

func (s *Scanner) isAlpha(c byte) bool {
	return (c >= 'a' && c <= 'z') ||
		(c >= 'A' && c <= 'Z') ||
		c == '_'
}

func (s *Scanner) isAlphaNumeric(c byte) bool {
	return s.isAlpha(c) || s.isDigit(c)
}

func (s *Scanner) peekNext() byte {
	if s.current+1 >= len(s.source) {
		return byte(rune(0))
	}
	return s.source[s.current+1]
}

func (s *Scanner) peek() byte {
	if s.isAtEnd() {
		return byte(rune(0))
	}
	return s.source[s.current]
}

func (s *Scanner) addSingleCharacterToken(tokenType TokenType) {
	s.addToken(tokenType, nil)
}

func (s *Scanner) addToken(tokenType TokenType, literal any) {
	text := s.source[s.start:s.current]
	s.tokens = append(s.tokens, Token{
		Type:    tokenType,
		lexeme:  text,
		literal: literal,
		line:    s.line,
	})
}

func (s *Scanner) advance() byte {
	b := s.source[s.current]
	s.current++
	return b
}

func (s *Scanner) match(expected byte) bool {
	if s.isAtEnd() {
		return false
	}
	if s.source[s.current] != expected {
		return false
	}

	s.current++
	return true
}

func (s *Scanner) string() {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}

	if s.isAtEnd() {
		Error(s.line, "unterminated string.")
		return
	}

	s.advance()

	s.addToken(STRING, s.source[s.start+1:s.current-1])
}

type TokenType int

var tokenTypeStringMap = map[TokenType]string{
	LEFT_PAREN:    "LEFT_PAREN",
	RIGHT_PAREN:   "RIGHT_PAREN",
	LEFT_BRACE:    "LEFT_BRACE",
	RIGHT_BRACE:   "RIGHT_BRACE",
	COMMA:         "COMMA",
	DOT:           "DOT",
	MINUS:         "MINUS",
	PLUS:          "PLUS",
	SEMICOLON:     "SEMICOLON",
	SLASH:         "SLASH",
	STAR:          "STAR",
	BANG:          "BANG",
	BANG_EQUAL:    "BANG_EQUAL",
	EQUAL:         "EQUAL",
	EQUAL_EQUAL:   "EQUAL_EQUAL",
	GREATER:       "GREATER",
	GREATER_EQUAL: "GREATER_EQUAL",
	LESS:          "LESS",
	LESS_EQUAL:    "LESS_EQUAL",
	IDENTIFIER:    "IDENTIFIER",
	STRING:        "STRING",
	NUMBER:        "NUMBER",
	AND:           "AND",
	CLASS:         "CLASS",
	ELSE:          "ELSE",
	FALSE:         "FALSE",
	FUN:           "FUN",
	FOR:           "FOR",
	IF:            "IF",
	NIL:           "NIL",
	OR:            "OR",
	PRINT:         "PRINT",
	RETURN:        "RETURN",
	SUPER:         "SUPER",
	THIS:          "THIS",
	TRUE:          "TRUE",
	VAR:           "VAR",
	WHILE:         "WHILE",
	EOF:           "EOF",
}

func (t TokenType) String() string {
	return tokenTypeStringMap[t]
}

const (
	// Single-character tokens.
	LEFT_PAREN = iota
	RIGHT_PAREN
	LEFT_BRACE
	RIGHT_BRACE
	COMMA
	DOT
	MINUS
	PLUS
	SEMICOLON
	SLASH
	STAR

	// One or two character tokens.
	BANG
	BANG_EQUAL
	EQUAL
	EQUAL_EQUAL
	GREATER
	GREATER_EQUAL
	LESS
	LESS_EQUAL

	// Literals.
	IDENTIFIER
	STRING
	NUMBER

	// Keywords.
	AND
	CLASS
	ELSE
	FALSE
	FUN
	FOR
	IF
	NIL
	OR
	PRINT
	RETURN
	SUPER
	THIS
	TRUE
	VAR
	WHILE

	EOF
)

type Token struct {
	Type    TokenType
	lexeme  string
	literal any
	line    int
}

func (t *Token) String() string {
	return fmt.Sprintf("%s %s %s", t.Type, t.lexeme, t.literal)
}

func NewScanner(source string) *Scanner {
	return &Scanner{
		source: source,
	}
}
