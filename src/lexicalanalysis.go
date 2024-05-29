package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode"
)

type position struct {
	line   int
	column int
}

type lexer struct {
	pos    position
	reader *bufio.Reader
}

type token int

func LexicalAnalysis(shouldprint bool) ([]node, []string) {
	var includes []string
	file, err := os.Open("build/build.ap")
	if err != nil {
		panic(err)
	}
	lexer := NewLexer(file)
	var nodes []node
	for {
		pos, tok, lit := lexer.Lex()
		if tok == EOF {
			break
		}
		var node node
		node.token = tok
		node.value = lit
		node.size = 0
		node.linenumber = pos.line
		node.columnnumber = pos.column
		if tok != INCLUDE {
			nodes = append(nodes, node)
		} else {
			includes = append(includes, lit)
		}
		if shouldprint {
			fmt.Printf("%d:%d\t%s\t%s\n", pos.line, pos.column, tok, lit)
		}
	}
	return nodes, includes
}

func NewLexer(reader io.Reader) *lexer {
	return &lexer{
		pos:    position{line: 1, column: 0},
		reader: bufio.NewReader(reader),
	}
}

func StandardizeSpaces(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

func (l *lexer) Lex() (position, token, string) {
	for {
		r, _, err := l.reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				return l.pos, EOF, ""
			}

			panic(err)
		}
		// updates the column to the position of the neqly read in rune
		l.pos.column++
		switch r {
		case '\n':
			l.ResetPosition()
		case ';':
			return l.pos, SEMI, ";"
		case '+':
			return l.pos, ADD, "+"
		case '-':
			return l.pos, SUB, "-"
		case '*':
			return l.pos, IMUL, "*"
		case '/':
			return l.pos, IDIV, "/"
		case '=':
			l.pos.column++
			r, _, err = l.reader.ReadRune()
			switch r {
			case '=':
				return l.pos, ISEQUAL, "=="
			default:
				l.pos.column--
				l.Backup()
				return l.pos, ASSIGN, "="
			}
		case '>':
			l.pos.column++
			r, _, err = l.reader.ReadRune()
			switch r {
			case '=':
				return l.pos, ISEQUAL, "=="
			default:
				l.pos.column--
				l.Backup()
				return l.pos, GREATER, ">"
			}
		case '!':
			l.pos.column++
			r, _, err = l.reader.ReadRune()
			switch r {
			case '=':
				return l.pos, ISNOTEQUAL, "!="
			default:
				l.pos.column--
				l.Backup()
				return l.pos, ILLEGAL, "!"
			}
		case '.':
			return l.pos, POINT, "."
		case '(':
			return l.pos, OPENBRACKET, "("
		case ')':
			return l.pos, CLOSEBRACKET, ")"
		case '{':
			return l.pos, OPENCBRACKET, "{"
		case '}':
			return l.pos, CLOSECBRACKET, "}"
		case '|':
			return l.pos, OR, "|"
		case '&':
			return l.pos, AND, "&"
		case '%':
			return l.pos, MOD, "%"
		case ',':
			return l.pos, COMMA, ","
		case '<':
			l.pos.column++
			r, _, err = l.reader.ReadRune()
			if r != '<' && r != '=' {
				l.pos.column--
				l.Backup()
				return l.pos, LESS, "<"
			} else if r == '=' {
				l.pos.column--
				l.Backup()
				return l.pos, LESSOREQUAL, "<="
			}
			l.pos.column++
			r, _, err = l.reader.ReadRune()
			if r != '<' {
				l.pos.column--
				l.Backup()
			}
			return l.pos, PIPEIN, "<<<"
		case '#':
			returnstring := ""
			include := "include"
			for _, val := range include {
				l.pos.column++
				r, _, err = l.reader.ReadRune()
				if r != val {
					break
				}
			}
			for {
				l.pos.column++
				r, _, err = l.reader.ReadRune()
				if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '.' || r == '/' {
					l.pos.column--
					l.Backup()
					break
				}
			}
			for {
				l.pos.column++
				r, _, err = l.reader.ReadRune()
				if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '.' || r == '/' {
					returnstring += string(r)
				} else {
					l.pos.column--
					l.Backup()
					break
				}
			}
			return l.pos, INCLUDE, returnstring
		default:
			if unicode.IsSpace(r) {
				continue
			} else if unicode.IsDigit(r) {
				// backup and let lexInt rescan the beginning of the int
				startPos := l.pos
				l.Backup()
				lit, isfloat := l.LexInt()
				if !isfloat {
					return startPos, INT, lit
				} else {
					return startPos, FLOAT, lit
				}
			} else if unicode.IsLetter(r) {
				startPos := l.pos
				l.Backup()
				lit := l.LexIdent()
				switch lit {
				case "func":
					return startPos, FUNC, lit
				case "return":
					return startPos, RETURN, lit
				case "true":
					return startPos, BOOL, lit
				case "false":
					return startPos, BOOL, lit
				case "if":
					return startPos, IF, lit
				case "exit":
					return startPos, EXIT, lit
				case "do":
					return startPos, DO, lit
				case "break":
					return startPos, BREAK, lit
				default:
					return startPos, IDENT, lit
				}
			} else {
				return l.pos, ILLEGAL, string(r)
			}
		}
	}
}

// lexInt scans the input until the end of an integer and then returns the literal
func (l *lexer) LexInt() (string, bool) {
	var lit string
	isfloat := false
	for {
		r, _, err := l.reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				// at the end of the int
				return lit, isfloat
			}
		}

		l.pos.column++
		if unicode.IsDigit(r) || string(r) == "." {
			if string(r) == "." {
				isfloat = true
			}
			lit = lit + string(r)
		} else {
			// scanned something not in the integer
			l.Backup()
			return lit, isfloat
		}
	}
}

// lexIdent scans the input until the end of an identifier and then returns the literal.
func (l *lexer) LexIdent() string {
	var lit string
	for {
		r, _, err := l.reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				// at the end of the identifier
				return lit
			}
		}

		l.pos.column++
		if unicode.IsLetter(r) {
			lit = lit + string(r)
		} else {
			// scanned something not in the identifier
			l.Backup()
			return lit
		}
	}
}

func (l *lexer) Backup() {
	if err := l.reader.UnreadRune(); err != nil {
		panic(err)
	}

	l.pos.column--
}

func (l *lexer) ResetPosition() {
	l.pos.line++
	l.pos.column = 0
}
