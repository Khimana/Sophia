package lexer

import (
	"fmt"
	"log"
	"sophia/core/token"
	"strings"
	"unicode"
)

type Lexer struct {
	input    []byte
	pos      int
	chr      byte
	line     int
	linepos  int
	HasError bool
}

func New(input []byte) Lexer {
	if len(input) == 0 || len(strings.TrimSpace(string(input))) == 0 {
		log.Println("err: input is empty, stopping")
		return Lexer{
			HasError: true,
		}
	}
	return Lexer{
		input:    input,
		pos:      0,
		chr:      input[0],
		line:     0,
		linepos:  0,
		HasError: false,
	}
}

func (l *Lexer) Lex() []token.Token {
	t := make([]token.Token, 0)
	for l.chr != 0 {
		ttype := token.UNKNOWN

		switch l.chr {
		case '+':
			if l.peek() == '+' {
				ttype = token.MERGE
				l.advance()
			} else {
				ttype = token.ADD
			}
		case '-':
			ttype = token.SUB
		case '/':
			ttype = token.DIV
		case '*':
			ttype = token.MUL
		case '%':
			ttype = token.MOD
		case '(':
			ttype = token.LEFT_BRACE
		case ')':
			ttype = token.RIGHT_BRACE
		case '{':
			ttype = token.LEFT_CURLY
		case '}':
			ttype = token.RIGHT_CURLY
		case '[':
			ttype = token.LEFT_BRACKET
		case ']':
			ttype = token.RIGHT_BRACKET
		case ':':
			ttype = token.COLON
		case '.':
			ttype = token.DOT
		case '_':
			ttype = token.PARAM
		case '\'':
			t = append(t, l.templateString()...)
			continue
		case ' ', '\t', '\r', '\n':
			if l.chr == '\n' {
				l.linepos = 0
				l.line++
			}
			l.advance()
			continue
		case '"':
			t = append(t, l.string())
			continue
		case ';':
			if l.peek() == ';' {
				for l.chr != '\n' && l.chr != 0 {
					l.advance()
				}
				continue
			}
		default:
			if unicode.IsLetter(rune(l.chr)) {
				t = append(t, l.ident())
				continue
			} else if unicode.IsDigit(rune(l.chr)) {
				if tok, err := l.float(); err == nil {
					t = append(t, tok)
				} else {
					l.error(3, tok.Raw)
				}
				continue
			}
		}

		if ttype == token.UNKNOWN {
			l.error(0, "")
		}

		t = append(t, token.Token{
			Pos:  l.pos,
			Type: ttype,
			Line: l.line,
			Raw:  string(l.chr),
		})

		l.advance()
	}
	if l.HasError {
		return []token.Token{
			{Type: token.EOF, Line: l.line},
		}
	}
	t = append(t, token.Token{
		Type: token.EOF, Line: l.line,
	})
	return t
}

func (l *Lexer) error(errType uint, ident string) {
	pos := l.linepos
	iLen := len(ident)

	switch errType {
	case 3:
		pos = pos - iLen
		if pos < 0 {
			pos = 0
		}
		log.Printf("err: Invalid floating point integer '%s' at [l %d:%d]", ident, l.line+1, pos)
	case 2:
		pos = pos - iLen
		if pos < 0 {
			pos = 0
		}
		log.Printf("err: Unterminated String at [l %d:%d]", l.line+1, pos)
	case 1:
		pos = pos - iLen
		if pos < 0 {
			pos = 0
		}
		log.Printf("err: Unknown identifier '%s' at [l %d:%d]", ident, l.line+1, pos)
	default:
		log.Printf("err: Unknown token '%c' at [l %d:%d]", l.chr, l.line+1, pos)
	}

	lines := strings.Split(string(l.input), "\n")

	spaces := pos

	// if no identifier given, print one ^
	if iLen == 0 {
		iLen += 1
	}
	if l.line-1 > -1 {
		if errType == 2 {
			spaces -= 1
		}
		spaces -= 1
		if spaces < 0 {
			spaces = 0
		}
		fmt.Printf("\n%.3d |\t%s\n%.3d |\t%s\n\t%s%s\n\n", l.line, lines[l.line-1], l.line+1, lines[l.line], strings.Repeat(" ", spaces), strings.Repeat("^", iLen))
	} else {
		fmt.Printf("\n%.3d |\t%s\n\t%s%s\n\n", l.line+1, lines[l.line], strings.Repeat(" ", spaces), strings.Repeat("^", iLen))
	}
	l.HasError = true
}

func (l *Lexer) templateString() []token.Token {
	el := make([]token.Token, 0)
	el = append(el, token.Token{
		Type: token.TEMPLATE_STRING,
		Pos:  l.pos,
		Line: l.line,
		Raw:  "",
	})
	b := strings.Builder{}

	l.advance() // skip '

	for {
		if l.chr == '}' {
			l.advance()
			continue
		} else if l.chr == '{' {
			l.advance()
			if b.Len() != 0 {
				el = append(el, token.Token{
					Pos:  l.pos - (len(b.String()) + 2),
					Type: token.STRING,
					Raw:  b.String(),
					Line: l.line,
				})
				b.Reset()
			}
			el = append(el, l.ident())
			continue
		} else if l.chr == '\n' || l.chr == 0 {
			log.Printf("err: Unexpected Newline or EOF at [l %d:%d]", l.line+1, l.pos)
			l.HasError = true
			return []token.Token{}
		} else if l.chr == '\'' {
			if b.Len() != 0 {
				el = append(el, token.Token{
					Pos:  l.pos - (len(b.String()) + 2),
					Type: token.STRING,
					Raw:  b.String(),
					Line: l.line,
				})
				b.Reset()
			}
			break
		}
		b.WriteByte(l.chr)
		l.advance()
	}

	if l.chr == '\'' {
		l.advance()
	}

	el = append(el, token.Token{
		Type: token.TEMPLATE_STRING,
		Pos:  l.pos,
		Line: l.line,
		Raw:  "",
	})
	return el
}

func (l *Lexer) string() token.Token {
	l.advance()
	b := strings.Builder{}
	for l.chr != '"' && l.chr != '\n' && l.chr != 0 {
		b.WriteByte(l.chr)
		l.advance()
	}
	str := b.String()
	if l.chr != '"' {
		l.error(2, str)
	} else {
		l.advance()
	}

	return token.Token{
		Pos:  l.pos - (len(str) + 2),
		Type: token.STRING,
		Raw:  str,
		Line: l.line,
	}
}

func (l *Lexer) ident() token.Token {
	builder := strings.Builder{}
	for unicode.IsLetter(rune(l.chr)) || l.chr == '_' || unicode.IsDigit(rune(l.chr)) {
		builder.WriteByte(l.chr)
		l.advance()
	}
	str := builder.String()
	ttype := token.UNKNOWN
	switch str {
	case "true", "false":
		ttype = token.BOOL
	default:
		ttype = token.IDENT
	}
	if tokenType, ok := token.KEYWORD_MAP[str]; ok {
		ttype = tokenType
	}
	return token.Token{
		Pos:  l.pos - len(str),
		Type: ttype,
		Raw:  str,
		Line: l.line,
	}
}

func (l *Lexer) float() (token.Token, error) {
	builder := strings.Builder{}
	for unicode.IsDigit(rune(l.chr)) || l.chr == '.' || l.chr == '_' || l.chr == 'e' || l.chr == '-' {
		builder.WriteByte(l.chr)
		l.advance()
	}
	str := builder.String()
	return token.Token{
		Pos:  l.pos - len(str),
		Type: token.FLOAT,
		Raw:  str,
		Line: l.line,
	}, nil
}

func (l *Lexer) peek() byte {
	if l.pos+1 < len(l.input) {
		return l.input[l.pos+1]
	}
	return 0
}

func (l *Lexer) advance() {
	if l.pos+1 < len(l.input) {
		l.pos++
		l.linepos++
		l.chr = l.input[l.pos]
		return
	}
	l.chr = 0
}
