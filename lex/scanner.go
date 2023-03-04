package lex

import (
	"fmt"
	"strings"

	gojs "github.com/ruiconti/gojs/internal"
)

var errEOF = fmt.Errorf("EOF")
var errOOB = fmt.Errorf("out of bounds")
var errNoLiteralAfterNumber = fmt.Errorf("no literal after number")
var errUnexpectedToken = fmt.Errorf("unexpected token")
var errDigitExpected = fmt.Errorf("digit expected")
var errInfiniteLoop = fmt.Errorf("infinite loop detected")

// TODO: below are untested
var errUnterminatedStringLiteral = fmt.Errorf("unterminated string literal")
var errBadEscapeSequence = fmt.Errorf("bad character escape sequence")
var errInvalidNumericEscape = fmt.Errorf("invalid numeric escape")

type Scanner struct {
	// source string being scanned
	src string
	// offset of the current index of the current lexeme
	idxHead int
	// offset of the beginning of the current lexeme
	idxHeadStart int
	// token slice
	tokens []Token
	// reference to the logging mechanism
	logger *gojs.SimpleLogger
	// errors
	errors []error
}

func (s *Scanner) hasErrors() bool {
	return len(s.errors) > 0
}

func NewScanner(src string, logger *gojs.SimpleLogger) *Scanner {
	if logger == nil {
		logger = gojs.NewSimpleLogger(gojs.ModeDebug)
	}
	return &Scanner{
		src:          src,
		logger:       logger,
		idxHead:      0,
		idxHeadStart: 0,
		tokens:       []Token{},
	}
}

func (s *Scanner) advanceBy(n int) {
	s.idxHead += n
}

func (s *Scanner) peek() rune {
	return rune(s.src[s.idxHead])
}

func (s *Scanner) peekN(n int) (rune, error) {
	if s.idxHead+n < 0 {
		return 0, errOOB
	} else if s.idxHead+n >= len(s.src) {
		return 0, errEOF
	}
	return rune(s.src[s.idxHead+n]), nil
}

func (s *Scanner) addTokenSafe(t TokenType) {
	var lexeme string
	if s.idxHead == s.idxHeadStart {
		if t == TStringLiteral_DoubleQuote || t == TStringLiteral_SingleQuote {
			// only string literals need this treatment
			// see in strings parser for more context
			lexeme = ""
		} else {
			lexeme = s.src[s.idxHeadStart : s.idxHead+1]
		}
	} else {
		lexeme = s.src[s.idxHeadStart:s.idxHead]
	}
	s.logger.Debug("[%d] addTokenSafe -> %v (lex: %v, len: %d)", s.idxHead, ResolveName(t), lexeme, len(lexeme))
	s.tokens = append(s.tokens, Token{
		T:       t,
		Lexeme:  lexeme,
		Literal: lexeme,
		// TODO: implement positioning
		Line:   0,
		Column: 0,
	})
}

// Checks whether the sequence is found next, if it is, advance headIdx e.g.
//
// >> src := "()===!";
// >> headIdx := 2;
// >> head := '=';
// >> candidates := []rune{'=', '='}
// >> seekMatchSequence(candidates)
// true
// >> fmt.Println(headIdx)
// 4
func (s *Scanner) seekMatchSequence(sequence []rune) bool {
	if len(sequence) == 0 {
		panic("sequence must not be empty")
	}

	i, j := s.idxHead+1, 0
	for j < len(sequence) {
		if i > len(s.src)-1 || j > len(sequence)-1 {
			// out of bounds
			return false
		}
		cursorGot := rune(s.src[i])
		cursorExpected := rune(sequence[j])

		if cursorGot != cursorExpected {
			return false
		}
		i++
		j++
	}

	return true
}

func isWhitespace(r rune) bool {
	return r == ' ' || r == '\t' || r == '\r'
}

func isNewline(r rune) bool {
	return r == '\n'
}

func (s *Scanner) Scan() ([]Token, error) {
	s.logger.Debug("SRC:\n%s\n\n", s.src)

	var lastHead rune
	for s.idxHead < len(s.src) {
		s.logger.Debug("[%d] mainloop: %c", s.idxHead, s.peek())
		s.prettyPrintScan()
		s.idxHeadStart = s.idxHead
		head := s.peek()

		// check whether we're in an endless loop
		if lastHead == head {
			s.logger.Debug("infinite loop found, aborting...\n\n")
			return s.tokens, errInfiniteLoop
		} else {
			lastHead = head
		}

		// identifiers
		if isIdentifierStart(head) {
			s.scanIdentifiers()
		} else if s.hasErrors() {
			return []Token{}, s.errors[0]
		}
		if s.idxHead >= len(s.src) {
			break
		} else if s.hasErrors() {
			return []Token{}, s.errors[0]
		}

		// literals
		if s.peek() == '"' || s.peek() == '\'' {
			s.scanStringLiteral()
		}
		if s.idxHead >= len(s.src) {
			// done!
			break
		} else if s.hasErrors() {
			return []Token{}, s.errors[0]
		}

		// template literals
		if s.peek() == '`' {
			s.scanTemplateLiteral()
		}
		if s.idxHead >= len(s.src) {
			// done!
			break
		} else if s.hasErrors() {
			return []Token{}, s.errors[0]
		}

		// numeric literals
		if isDecimalDigit(head) || head == '.' {
			s.scanDigits()
		}
		if s.idxHead >= len(s.src) {
			// done!
			break
		} else if s.hasErrors() {
			return []Token{}, s.errors[0]
		}

		// punctuators
		if isPunctuation(rune(head)) {
			s.scanPunctuators()
		}
		if s.idxHead >= len(s.src) {
			// done!
			break
		} else if s.hasErrors() {
			return []Token{}, s.errors[0]
		}

		// we should only advance the head if we didn't match anything
		// otherwise, we would've already advanced the head because of the
		// addTokenSafe call, which expects to be one index ahead of the end of the token
		//
		// example 1:
		// "fooobar boo"
		//  ^- head: 0
		//
		// "fooobar boo"
		//        ^- head: 6 (at this point, we are still in scanIdentifier, because `r` is a valid one)
		//
		// "fooobar boo"
		//         ^- head: 7 (now, whitespace is not a valid identifier, and we advance the head)
		//
		// example 2:
		// "??&&=>"
		//  ^- head: 0
		//
		// "??&&=>"
		//    ^- head: 2 (quit scanPunctuators because we found a valid `??` token with length 2)
		//               (we are now at the `&` character and we can't advance the head otherwise we'd skip it)
		if isWhitespace(head) {
			s.advanceBy(1)
		}
	}

	tokens := strings.Builder{}
	for _, t := range s.tokens {
		var lit string
		if t.Literal != nil {
			lit = fmt.Sprintf("%s", t.Literal)
		} else {
			lit = fmt.Sprintf("%s", t.Lexeme)
		}
		tokens.Write([]byte(lit))
	}
	s.logger.Debug("\nTOKENS:\n%s\n", tokens.String())

	return s.tokens, nil
}

// Template literals
//
// https://262.ecma-international.org/#prod-TemplateLiteral
func (s *Scanner) scanTemplateLiteral() {

}
