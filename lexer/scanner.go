package lexer

import (
	"errors"
	"fmt"
	"strings"

	"github.com/ruiconti/gojs/internal"
	gojs "github.com/ruiconti/gojs/internal"
)

const EOF rune = -1

var (
	ReservedKeywords       = internal.MapInvert(ReservedWordNames)
	TokenUnknown     Token = Token{TUnknown, "", "", 0, 0}
)

// Errors
var (
	errEOF                  = fmt.Errorf("EOF")
	errNoLiteralAfterNumber = fmt.Errorf("no literal after number")
	errUnexpectedToken      = fmt.Errorf("unexpected token")
	// TODO: below are untested
	errUnterminatedStringLiteral = fmt.Errorf("unterminated string literal")
	errInvalidEscapedSequence    = errors.New("invalid escaped sequence")
)

type Lexer struct {
	// source string being scanned
	src string
	// token slice
	tokens []Token
	// reference to the logging mechanism
	logger *gojs.SimpleLogger
	// offset of the beginning of the current char
	srcCursor int
	// offset of the current index of the current char
	srcCursorHead int
	// whether the cursor is out of bounds
	srcCursorOOB bool
	// last index of the source string
	srcEnd int
	// store errors found while scanning
	errors []error

	// TODO: implement proper positioning
	line   uint
	column uint
}

func NewLexer(src string, logger *gojs.SimpleLogger) *Lexer {
	if logger == nil {
		logger = gojs.NewSimpleLogger(gojs.ModeDebug)
	}

	return &Lexer{
		src:           src,
		logger:        logger,
		srcCursor:     -2,
		srcCursorHead: 0,
		srcCursorOOB:  false,
		srcEnd:        len(src) - 1,
		tokens:        []Token{},
		line:          0,
		column:        0,
	}
}

// PeekLoop wraps the usual for { Peek(); Next(); } loop
// in a way to prevent infinite loops in a coordinated fashion
func (s *Lexer) PeekLoop(callback func(rune) bool) {
	var (
		ch      rune
		lastPos int
	)
	for {
		ch = s.Peek()
		lastPos = s.srcCursorHead

		s.logger.Debug("%d: loop:%c", s.srcCursorHead, ch)
		cont := callback(ch) // side-effect

		if !cont || s.srcCursorOOB {
			break
		}
		if lastPos == s.srcCursorHead {
			panic("infinite loop detected")
		}
	}
}

// Advance the cursor by 1 char
func (s *Lexer) Next() {
	s.Jump(1)
}

// Advance the cursor by N char
func (s *Lexer) Jump(offset uint) {
	width := s.srcCursorHead + int(offset)

	if width > s.srcEnd {
		s.srcCursorOOB = true
		return
	}

	s.srcCursorHead = width
	s.logger.Debug("%d: jump(offset:%d):%c", s.srcCursorHead, offset, s.Peek())
	s.PrettyPrintSrc()
}

// 1-char look-ahead
func (s *Lexer) Peek() rune {
	return s.PeekN(0)
}

// N-char look-ahead
func (s *Lexer) PeekN(offset uint) rune {
	lookAhead := s.srcCursorHead + int(offset)
	if lookAhead > s.srcEnd {
		return EOF
	}
	return rune(s.src[lookAhead])
}

// CreateLiteralToken abstracts the common task of creating
// a token for a literal (eg bool, string, number)
func (s *Lexer) CreateLiteralToken(typ TokenType) Token {
	var candidate string
	if s.srcCursorHead == s.srcEnd && s.srcCursorOOB {
		// if we are at the end of the source, we can't use srcCursorHead for trimming
		// because that would be out of bounds
		candidate = s.src[s.srcCursor:]
	} else {
		candidate = s.src[s.srcCursor:s.srcCursorHead]
	}

	// try to parse it as a reserved word
	if typ, ok := ReservedKeywords[candidate]; ok {
		return Token{
			Lexeme: candidate,
			Type:   typ,
		}
	}

	// regular literal
	return Token{
		Type:    typ,
		Literal: candidate,
		Lexeme:  candidate,
	}
}

func isWhitespace(r rune) bool { return r == ' ' || r == '\t' || r == '\r' }
func isNewline(r rune) bool    { return r == '\n' }
func isStr(r rune) bool        { return r == '\'' || r == '"' }
func isNumeric(r rune) bool    { return r == '.' || isDec(r) }

// Scan only the next token
func (s *Lexer) Scan() Token {
	var token Token

	ch := s.Peek()
	switch {
	case isId(ch):
		token = s.scanIdentifier()
	case isStr(ch):
		token = s.scanStringLiteral()
	case isNumeric(ch):
		token = s.scanNumericLiteral()
	case isPunctuation(ch):
		token = s.scanPunctuation()
	case isWhitespace(ch):
		token = Token{Type: TWhitespace, Lexeme: " ", Literal: " "}
	case isNewline(ch):
		s.line++
		s.Next()
		token = Token{Type: TWhitespace, Lexeme: " ", Literal: " "}
	default:
		token = TokenUnknown
	}
	return token
}

// Scan up until src's EOF
func (s *Lexer) ScanAll() ([]Token, []error) {
	s.logger.Debug("SRC:\n%s\n\n", s.src)
	defer func() {
		stack := recover()
		if stack != nil {
			s.PrettyPrintSrc()
			s.logger.DumpLogs()
			panic(stack)
		}
	}()

mainloop:
	for s.srcCursorHead <= s.srcEnd {
		if s.srcCursorHead == s.srcCursor {
			panic("infinite loop found, aborting")
		}
		s.srcCursor = s.srcCursorHead

		tok := s.Scan()
		switch tok.Type {
		case TUnknown:
			break mainloop
		case TWhitespace:
			s.Next()
		default:
			s.tokens = append(s.tokens, tok)
		}

		if s.srcCursorOOB {
			// we advanced past the end of the source
			break
		}
	}

	s.logger.Debug("\nTOKENS:\n%s\n", s.Tokens())
	return s.tokens, s.errors
}

// Printing utilities: all tokens
func (s *Lexer) Tokens() string {
	var sb strings.Builder
	sb.WriteByte('(')
	for i, token := range s.tokens {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.Write([]byte(token.Lexeme))
	}
	sb.WriteByte(')')
	return sb.String()
}

// Printing utilities: where the cursor is
func (s *Lexer) PrettyPrintSrc() {
	s.logger.Debug("%v", s.src)
	cursor := []byte{}
	for i := 0; i < s.srcCursorHead; i++ {
		cursor = append(cursor, ' ')
	}
	cursor = append(cursor, '^')
	s.logger.Debug("%s", cursor)
}

// Errorf is a convenience function to log errors
func (s *Lexer) Errorf(format string, values ...any) {
	formatted := fmt.Sprintf(format, values...)
	serr := fmt.Sprintf("%s COL:%d CH:%c", formatted, s.srcCursorHead, s.Peek())

	s.errors = append(s.errors, errors.New(serr))
	s.PrettyPrintSrc()
	s.logger.Error(serr + "\n")
}
