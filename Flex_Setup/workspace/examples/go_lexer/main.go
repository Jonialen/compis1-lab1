package main

import (
	"fmt"
	"os"
	"strings"
	"unicode"
	"unicode/utf8"
)

//  Tipos de token

type TokenType int

const (
	TEOF TokenType = iota
	TILLEGAL
	TIDENT
	TKEYWORD
	TINT
	TFLOAT
	TSCIENTIFIC
	THEX
	TSTRING
	TOP
	TDELIM
	TCOMMENT
)

var tokenName = [...]string{
	"EOF", "ILLEGAL", "IDENT", "KEYWORD",
	"INT", "FLOAT", "SCIENTIFIC", "HEX",
	"STRING", "OP", "DELIM", "COMMENT",
}

func (t TokenType) String() string { return tokenName[t] }

//  Token

type Token struct {
	Type  TokenType
	Value string
}

func (t Token) String() string {
	return fmt.Sprintf("Token(%-12s %q)", t.Type.String()+",", t.Value)
}

//  Lexer — corazón de la máquina de estados

type stateFn func(*Lexer) stateFn

const eof rune = -1

type Lexer struct {
	input  string
	pos    int // posición actual
	start  int // inicio del token en curso (en bytes)
	tokens chan Token
}

func Lex(src string) <-chan Token {
	l := &Lexer{input: src, tokens: make(chan Token, 64)}
	go func() {
		for state := lexRoot; state != nil; {
			state = state(l) // transición: ejecutar estado → obtener siguiente
		}
		close(l.tokens)
	}()
	return l.tokens
}

// Primitivas del lexer

func (l *Lexer) next() rune {
	if l.pos >= len(l.input) {
		return eof
	}
	r, sz := utf8.DecodeRuneInString(l.input[l.pos:])
	l.pos += sz
	return r
}

func (l *Lexer) peek() rune {
	if l.pos >= len(l.input) {
		return eof
	}
	r, _ := utf8.DecodeRuneInString(l.input[l.pos:])
	return r
}

func (l *Lexer) backup() {
	if l.pos > l.start {
		_, sz := utf8.DecodeLastRuneInString(l.input[:l.pos])
		l.pos -= sz
	}
}

func (l *Lexer) emit(t TokenType) {
	l.tokens <- Token{t, l.input[l.start:l.pos]}
	l.start = l.pos
}

func (l *Lexer) discard() { l.start = l.pos }

func (l *Lexer) lexeme() string { return l.input[l.start:l.pos] }

//  Estados del autómata (nodos del grafo de transiciones)

var goKeywords = map[string]bool{
	"break": true, "case": true, "chan": true, "const": true,
	"continue": true, "default": true, "defer": true, "else": true,
	"fallthrough": true, "for": true, "func": true, "go": true,
	"goto": true, "if": true, "import": true, "interface": true,
	"map": true, "package": true, "range": true, "return": true,
	"select": true, "struct": true, "switch": true, "type": true,
	"var": true,
}

func lexRoot(l *Lexer) stateFn {
	r := l.next()
	switch {
	case r == eof:
		l.emit(TEOF)
		return nil

	case r == ' ' || r == '\t' || r == '\n' || r == '\r':
		l.discard()
		return lexRoot

	case r == '"':
		return lexInterpString

	case r == '`':
		return lexRawString

	case r == '/':
		return lexSlash

	// Hex literal: 0x…
	case r == '0' && (l.peek() == 'x' || l.peek() == 'X'):
		l.next() // consume 'x'/'X'
		return lexHex

	case unicode.IsDigit(r):
		return lexNumber

	case unicode.IsLetter(r) || r == '_':
		return lexIdent

	default:
		return lexPunct(r) // retorna una stateFn capturando r
	}
}

func lexInterpString(l *Lexer) stateFn {
	for {
		switch l.next() {
		case '\\':
			l.next() // saltar el carácter escapado (\n, \t, \", \\, …)
		case '"', eof:
			l.emit(TSTRING)
			return lexRoot
		}
	}
}

func lexRawString(l *Lexer) stateFn {
	for {
		if r := l.next(); r == '`' || r == eof {
			l.emit(TSTRING)
			return lexRoot
		}
	}
}

func lexSlash(l *Lexer) stateFn {
	switch l.peek() {

	case '/': // comentario de línea
		l.next()
		for r := l.next(); r != '\n' && r != eof; r = l.next() {
		}
		l.backup()
		l.emit(TCOMMENT)
		return lexRoot

	case '*': // comentario de bloque
		l.next()
		for {
			r := l.next()
			if r == eof {
				l.emit(TCOMMENT)
				return nil
			}
			if r == '*' && l.peek() == '/' {
				l.next()
				break
			}
		}
		l.emit(TCOMMENT)
		return lexRoot

	case '=': // operador /=
		l.next()
		l.emit(TOP)
		return lexRoot

	default:
		l.emit(TOP)
		return lexRoot
	}
}

// Números

func lexHex(l *Lexer) stateFn {
	for {
		r := l.peek()
		isHexDigit := (r >= '0' && r <= '9') ||
			(r >= 'a' && r <= 'f') ||
			(r >= 'A' && r <= 'F') ||
			r == '_'
		if !isHexDigit {
			break
		}
		l.next()
	}
	l.emit(THEX)
	return lexRoot
}

func lexNumber(l *Lexer) stateFn {
	for unicode.IsDigit(l.peek()) || l.peek() == '_' {
		l.next()
	}
	if l.peek() == '.' {
		l.next()
		for unicode.IsDigit(l.peek()) {
			l.next()
		}
		return lexExpSuffix(TFLOAT)
	}
	return lexExpSuffix(TINT)
}

// lexExpSuffix retorna un estado que detecta el sufijo exponencial (e/E±dd…).
// Al retornar una stateFn en lugar de ejecutarse directamente, se integra
// con naturalidad al ciclo del despachador sin romper el patrón.
func lexExpSuffix(base TokenType) stateFn {
	return func(l *Lexer) stateFn {
		if p := l.peek(); p == 'e' || p == 'E' {
			l.next()
			if s := l.peek(); s == '+' || s == '-' {
				l.next()
			}
			for unicode.IsDigit(l.peek()) {
				l.next()
			}
			l.emit(TSCIENTIFIC)
			return lexRoot
		}
		l.emit(base)
		return lexRoot
	}
}

// Identificadores / palabras clave

func lexIdent(l *Lexer) stateFn {
	for r := l.peek(); unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_'; r = l.peek() {
		l.next()
	}
	if goKeywords[l.lexeme()] {
		l.emit(TKEYWORD)
	} else {
		l.emit(TIDENT)
	}
	return lexRoot
}

// Operadores y delimitadores

// lexPunct retorna una stateFn que captura r en su clausura.
// Esto permite manejar operadores multi-carácter (==, :=, <=, &&, …)
// sin consumir lookahead prematuramente en lexRoot.
func lexPunct(r rune) stateFn {
	return func(l *Lexer) stateFn {
		switch r {
		case '(', ')', '{', '}', '[', ']', ',', ';':
			l.emit(TDELIM)

		case '.':
			if l.peek() == '.' { // operador variádico ...
				l.next()
				if l.peek() == '.' {
					l.next()
				}
			}
			l.emit(TDELIM)

		case ':':
			if l.peek() == '=' {
				l.next() // :=
			}
			l.emit(TOP)

		case '=', '!', '<', '>':
			if l.peek() == '=' {
				l.next() // ==, !=, <=, >=
			}
			l.emit(TOP)

		case '+':
			if p := l.peek(); p == '=' || p == '+' {
				l.next() // +=, ++
			}
			l.emit(TOP)

		case '-':
			if p := l.peek(); p == '=' || p == '-' {
				l.next() // -=, --
			}
			l.emit(TOP)

		case '*', '%', '^':
			if l.peek() == '=' {
				l.next() // *=, %=, ^=
			}
			l.emit(TOP)

		case '&':
			if p := l.peek(); p == '&' || p == '=' || p == '^' {
				l.next() // &&, &=, &^
			}
			l.emit(TOP)

		case '|':
			if p := l.peek(); p == '|' || p == '=' {
				l.next() // ||, |=
			}
			l.emit(TOP)

		default:
			l.emit(TILLEGAL)
		}
		return lexRoot
	}
}

//  Entrada de muestra

const sampleSrc = `package main

import "fmt"

/* Fibonacci con memoización — demuestra todos los tipos de token */
func fib(n int, memo map[int]int) int {
	if n <= 1 {
		return n
	}
	if v, ok := memo[n]; ok {
		return v
	}
	result := fib(n-1, memo) + fib(n-2, memo)
	memo[n] = result
	return result
}

func main() {
	// Literales numéricos
	hex      := 0xFF_A3        // hexadecimal con separador
	pi       := 3.14159        // flotante
	avogadro := 6.022e+23      // notación científica
	year     := 2026

	// Operadores relacionales y lógicos
	_ = (hex != 0) && (pi >= 3.0) || !(year == 2025)

	// String con secuencias de escape
	msg := "Hola,\n\"Mundo\""

	// Raw string (sin escapes)
	raw := ` + "`cadena\ncruda\t sin escapes`" + `

	_ = hex; _ = avogadro; _ = year; _ = raw
	fmt.Println(msg)
}
`

//  main

func main() {
	src := sampleSrc
	if len(os.Args) > 1 {
		data, err := os.ReadFile(os.Args[1])
		if err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			os.Exit(1)
		}
		src = string(data)
	}

	fmt.Println("Lexer concurrente de Go")
	fmt.Println()

	counts := make(map[TokenType]int)
	for tok := range Lex(src) {
		if tok.Type == TEOF {
			break
		}
		fmt.Println(tok)
		counts[tok.Type]++
	}

	fmt.Println()
	fmt.Println(strings.Repeat("─", 50))
	fmt.Printf("%-14s  %s\n", "Tipo de token", "Cantidad")
	fmt.Println(strings.Repeat("─", 50))
	order := []TokenType{TKEYWORD, TIDENT, TINT, TFLOAT, TSCIENTIFIC, THEX, TSTRING, TOP, TDELIM, TCOMMENT, TILLEGAL}
	for _, tt := range order {
		if c := counts[tt]; c > 0 {
			fmt.Printf("%-14s  %d\n", tt, c)
		}
	}
}
