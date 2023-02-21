package lex

import (
	"fmt"
	"strings"

	gojs "github.com/ruiconti/gojs/internal"
)

var errEOF = fmt.Errorf("EOF")
var errOOB = fmt.Errorf("out of bounds")
var errNoLiteralAfterNumber = fmt.Errorf("no literal after number")
var errDigitExpected = fmt.Errorf("digit expected")

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

func (s *Scanner) advance() {
	s.idxHead += 1
}

func (s *Scanner) advanceBy(n int) {
	s.idxHead += n
}

func (s *Scanner) peek() rune {
	return rune(s.src[s.idxHead])
}

func (s *Scanner) isEOF(idx int) bool {
	return idx >= len(s.src)
}

func (s *Scanner) isPeekAheadEOF() bool {
	return s.isEOF(s.idxHead + 1)
}

func (s *Scanner) peekBehind() rune {
	return rune(s.src[s.idxHead-1])
}

func (s *Scanner) peekBehindSafe(n int) (rune, error) {
	if s.idxHead-n < 0 {
		return 0, errOOB
	}
	return rune(s.src[s.idxHead-n]), nil
}

func (s *Scanner) peekAhead(n int) rune {
	if s.idxHead+n >= len(s.src) {
		panic(errEOF)
	}
	return rune(s.src[s.idxHead+n])
}

func (s *Scanner) peekAheadSafe(n int) (rune, error) {
	if s.idxHead+n >= len(s.src) {
		return 0, errEOF
	}
	return rune(s.src[s.idxHead+n]), nil
}

func (s *Scanner) addToken(t TokenType, literal interface{}) {
	var lexeme string
	// we only advance the head before calling `addToken` in `seekMatchSequence`
	// for single matches e.g. `>`, we need to check it here so we don't overly complicate
	// the scan() function
	if s.idxHeadStart == s.idxHead {
		s.logger.Debug("%d: addToken (idxHeadStart == idxHead)", s.idxHead)
		lexeme = string(s.peek())
	} else {
		s.logger.Debug("%d: addToken (idxHead+1 > len(src)) %d %d", s.idxHead, s.idxHeadStart, s.idxHead)
		lexeme = s.src[s.idxHeadStart : s.idxHead+1]
	}
	s.logger.Debug("%d: addToken: %v %v", s.idxHead, t, lexeme)
	s.tokens = append(s.tokens, Token{
		T:       t,
		Lexeme:  lexeme,
		Literal: literal,
		// TODO: implement positioning
		Line:   0,
		Column: 0,
	})
}

func (s *Scanner) addTokenSafe(t TokenType, literal interface{}) {
	var lexeme string
	// we only advance the head before calling `addToken` in `seekMatchSequence`
	// for single matches e.g. `>`, we need to check it here so we don't overly complicate
	// the scan() function
	if s.idxHeadStart == s.idxHead {
		s.logger.Debug("%d: addToken (idxHeadStart == idxHead)", s.idxHead)
		lexeme = string(s.peek())
	} else {
		s.logger.Debug("%d: addToken (idxHead+1 > len(src)) %d %d", s.idxHead, s.idxHeadStart, s.idxHead)
		lexeme = s.src[s.idxHeadStart:s.idxHead]
	}
	s.logger.Debug("%d: addToken: %v %v", s.idxHead, t, lexeme)
	s.tokens = append(s.tokens, Token{
		T:       t,
		Lexeme:  lexeme,
		Literal: literal,
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
		// s.logger.Debug("(i:%d;j:%d) head: %c seekSeq: (eq %c %c)", i, j, s.peek(), cursorGot, cursorExpected)

		if cursorGot != cursorExpected {
			// s.logger.Debug("(i:%d;j:%d) head: %c seekSeq (leaving)", i, j, s.peek())
			return false
		}
		i++
		j++
	}
	// s.logger.Debug("%d: matchSequence true! %v", headIdx, runesstr)
	s.advanceBy(len(sequence))
	// s.logger.Debug("(i:%d;j:%d) head: %c seekSeq (advanced %d)", i, j, rune(s.peek()), len(sequence))
	return true
}

// TODO: RegExp
// TODO: Template literals
func (s *Scanner) Scan() ([]Token, error) {
	s.logger.Debug("%d: parsing: %s", s.idxHead, s.src)
	for s.idxHead < len(s.src) {
		s.idxHeadStart = s.idxHead

		head := s.peek()

		// identifiers
		if isIdentifierStart(head) {
			s.scanIdentifiers()
		}
		if s.idxHead == len(s.src) {
			// done!
			break
		} else if s.hasErrors() {
			return []Token{}, s.errors[0]
		}

		// literals
		if s.peek() == '"' || s.peek() == '\'' {
			s.scanStringLiteral()
		}
		if s.idxHead == len(s.src) {
			// done!
			break
		} else if s.hasErrors() {
			return []Token{}, s.errors[0]
		}

		// template literals
		if s.peek() == '`' {
			s.scanTemplateLiteral()
		}
		if s.idxHead == len(s.src) {
			// done!
			break
		} else if s.hasErrors() {
			return []Token{}, s.errors[0]
		}

		// numeric literals
		if isDecimalDigit(head) || head == '.' {
			s.scanDigits()
		}
		if s.idxHead == len(s.src) {
			// done!
			break
		} else if s.hasErrors() {
			return []Token{}, s.errors[0]
		}

		// punctuators
		if isPunctuation(rune(head)) {
			s.scanPunctuators()
		}
		if s.idxHead == len(s.src) {
			// done!
			break
		} else if s.hasErrors() {
			return []Token{}, s.errors[0]
		}

		s.logger.Debug("%d: (%c) next iter..", s.idxHead, s.peek())
		s.advance()
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

	return s.tokens, nil
}

// Template literals
//
// https://262.ecma-international.org/#prod-TemplateLiteral
func (s *Scanner) scanTemplateLiteral() {

}
