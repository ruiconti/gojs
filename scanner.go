package main

import (
	"fmt"
	"strings"
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
	logger *SimpleLogger
	// errors
	errors []error
}

func (s *Scanner) hasErrors() bool {
	return len(s.errors) > 0
}

func NewScanner(src string, logger *SimpleLogger) *Scanner {
	if logger == nil {
		logger = NewSimpleLogger(ModeDebug)
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

// func isIdentifierPart(r rune) bool {
// 	return isIdentifierStart(r) || isDecimalDigit(r)
// }

func isIdentifierStart(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || r == '$' || r == '_'
}

func isIdentifierPart(r rune) bool {
	return isIdentifierStart(r) || isDecimalDigit(r) || r == '\\'
}

// Identifiers
//
// https://262.ecma-international.org/#sec-names-and-keywords
func (s *Scanner) scanIdentifiers() {
	if !isIdentifierStart(s.peek()) {
		return
	}

	for isIdentifierPart(s.peek()) {
		if s.idxHead == len(s.src)-1 /* if EOF */ {
			break
		}
		headNext := s.peekAhead(1)
		if s.peek() == '\\' && headNext == 'u' {
			if s.idxHead+4 > len(s.src)-1 {
				// invalid unicode escape
				break
			} else if !isHexDigit(s.peekAhead(2)) || !isHexDigit(s.peekAhead(3)) || !isHexDigit(s.peekAhead(4)) || !isHexDigit(s.peekAhead(5)) {
				// invalid unicode escape
				break
			}
		}

		s.logger.Debug("%d: (scanIdentifier): %c", s.idxHead, s.peek())
		s.advance()
	}

	// we reach here when s.idxHead is **not** a valid identifier part
	// meaning that s.peek() is not part of the identifier
	s.logger.Debug("%d: (scanIdentifier): finished parsing %c", s.idxHead, s.peek())
	if !isIdentifierPart(s.peek()) {
		s.idxHead--
	} else {
		s.logger.Debug("%d: (scanIdentifier): finished parsing but idxHead was still valid: %c", s.idxHead, s.peek())
	}

	var upper, lower int
	lower = s.idxHeadStart
	upper = s.idxHeadStart + s.idxHead // collect final quote
	if upper >= len(s.src) {
		upper = len(s.src)
	}
	// s.idxHead = s.idxHead - 1 // collect final quote

	s.logger.Info("%d: (scanIdentifier) lowerBound:%d upperBound:%d literal:%s", s.idxHead, lower, upper, s.src[lower:upper])
	// s.logger.Debug("%d: (scanIdentifier) lowerBound:%d upperBound:%d literal:%s", s.idxHead, lower, upper, s.src[lower:upper])
	s.addToken(TIdentifier, s.src[s.idxHeadStart:upper])
}

func isLegalStringLiteralIntermediate(r rune) bool {
	return r != '"'
}
func isAlphaNumeric(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9')
}

func (s *Scanner) isValidEscapeSequence(r rune) bool {
	switch r {
	case 'b', 'f', 'n', 'r', 't', 'v', '\\', '"', '\'', '0':
		return true
	}
	return false
}

type StringLiteralType string

const (
	SingleQuote StringLiteralType = "singleq"
	DoubleQuote StringLiteralType = "doubleq"
)

// String literals
//
// https://262.ecma-international.org/#sec-literals-string-literals
func (s *Scanner) scanStringLiteral() {
	// If we got here, it means that we spotted a double quote
	var strType StringLiteralType
	if s.peek() == '"' {
		strType = DoubleQuote
	} else if s.peek() == '\'' {
		strType = SingleQuote
	} else {
		return
	}

	s.logger.Debug("%d: (scanStringLiteral) entry!", s.idxHead)

	strLength := 1 // advance first quote
	cursor := 0
	invalidString := false
	switch strType {
	case DoubleQuote:
		for {
			if s.idxHead == len(s.src)-1 /* if EOF */ {
				break
			}
			headNext := s.peekAhead(1)
			if s.idxHead > len(s.src)-1 /* if EOF */ ||
				(s.peek() == '"' && cursor > 0 && s.peekBehind() != '\\') /* it's an end quote " */ ||
				(s.peek() == '"' && cursor > 0 && s.peekBehind() == '\\' && headNext == ' ') {
				// TODO: when dealing with escapes, we should probably be better off with unicode code points instead of comparing runes
				break
			}
			s.logger.Debug("%d: (scanStringLiteral): %c", s.idxHead, s.peek())
			strLength++
			s.advance()
			cursor++
		}
	case SingleQuote:
		for {
			if s.idxHead == len(s.src)-1 /* if EOF */ {
				break
			}
			headNext := s.peekAhead(1)
			if s.idxHead > len(s.src)-1 /* if EOF */ ||
				(s.peek() == '\'' && cursor > 0 && s.peekBehind() != '\\') /* it's an end quote " */ ||
				(s.peek() == '\'' && cursor > 0 && s.peekBehind() == '\\' && headNext == ' ') {
				// TODO: when dealing with escapes, we should probably be better off with unicode code points instead of comparing runes
				break
			}
			s.logger.Debug("%d: (scanStringLiteral): %c", s.idxHead, s.peek())
			strLength++
			s.advance()
			cursor++
		}
	}

	if invalidString {
		return
	}

	var upper, lower int
	lower = s.idxHeadStart
	upper = s.idxHeadStart + strLength + 1 // collect final quote
	if upper >= len(s.src)-1 {
		upper = len(s.src) - 1
	}

	s.logger.Debug("%d: (scanStringLiteral) lowerBound:%d upperBound:%d literal:%s", s.idxHead, lower, upper, s.src[lower:upper])
	s.addToken(TStringLiteral, s.src[s.idxHeadStart:upper])
}

// Template literals
//
// https://262.ecma-international.org/#prod-TemplateLiteral
func (s *Scanner) scanTemplateLiteral() {

}

// Numeric literals
//
// https://262.ecma-international.org/#sec-literals-numeric-literals
type NumericLiteralType string

const (
	// TODO: LegacyOctal and BigInt are being handled in "decimal"
	LiteralDecimal NumericLiteralType = "decimal"
	LiteralHex     NumericLiteralType = "hex"
	LiteralBinary  NumericLiteralType = "binary"
	LiteralOctal   NumericLiteralType = "octal"
)

// Dec
func isDecimalDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

func isLegalDecDigitIntermediate(r rune) bool {
	return isDecimalDigit(r) || r == '_' || r == '.' || r == 'e' || r == 'E' || r == '+' || r == '-'
}

// BigInt
func isLegalBigIntDigitIntermediate(r rune) bool {
	return isDecimalDigit(r) || r == 'n'
}

// Hex
func isHexDigit(r rune) bool {
	return isDecimalDigit(r) || (r >= 'a' && r <= 'f') || (r >= 'A' && r <= 'F')
}

func isLegalHexDigitIntermediate(r rune) bool {
	return isHexDigit(r) || r == '_'
}

// Binary
func isBinaryDigit(r rune) bool {
	return r == '0' || r == '1'
}

func isLegalBinaryDigitIntermediate(r rune) bool {
	return isBinaryDigit(r) || r == '_'
}

// Octal
func isOctalDigit(r rune) bool {
	return r >= '0' && r <= '7'
}

func isLegalOctalDigitIntermediate(r rune) bool {
	return isOctalDigit(r) || r == '_'
}

func (s *Scanner) scanDigits() {
	headNext, _ := s.peekAheadSafe(1)
	zeroStart := false

	// find out what type of number we're dealing with
	var numberType NumericLiteralType
	if s.peek() == '0' {
		if headNext == 'x' || headNext == 'X' {
			numberType = LiteralHex
			s.advance()
			s.advance()
		} else if headNext == 'b' || headNext == 'B' {
			numberType = LiteralBinary
			s.advance()
			s.advance()
		} else if headNext == 'o' || headNext == 'O' {
			numberType = LiteralOctal
			s.advance()
			s.advance()
			zeroStart = true
		} else {
			numberType = LiteralDecimal
		}
	} else if isDecimalDigit(s.peek()) || s.peek() == '.' {
		// it can be a dec literal
		numberType = LiteralDecimal
	} else {
		s.logger.Debug("scanDigits:Invalid number!")
		// invalid number
		return
	}

	s.logger.Debug("%d: (scanDigits:entry) head:%c headNext:%c type:%s", s.idxHead, s.peek(), headNext, numberType)
	digitLength := 1
	errors := []error{}
	safe := false
	// limitedE := 0
	switch numberType {
	case LiteralDecimal:
		for {
			if s.idxHead == len(s.src) {
				break
			}
			s.logger.Debug("%d: (scanDecimalDigits) head:%c headNext:%c digitLength:%d isDecimalDigit(headNext):%v", s.idxHead, s.peek(), headNext, digitLength, isDecimalDigit(headNext))
			headNext, errAhead := s.peekAheadSafe(1)
			headBefore, errBehind := s.peekBehindSafe(1)
			if isLegalDecDigitIntermediate(s.peek()) {
				// current: e | E
				if (errAhead == nil && (s.peek() == 'e' || s.peek() == 'E') && !isDecimalDigit(headNext) && headNext != '+' && headNext != '-') ||
					(errAhead != nil && (s.peek() == 'e' || s.peek() == 'E')) {
					// only digits allowed after e: 10e, 10e-, 10e, 10e.
					errors = append(errors, errDigitExpected)
					break
				}
				// current + | - | _
				if (errAhead == nil && (s.peek() == '+' || s.peek() == '-' || s.peek() == '_') && !isDecimalDigit(headNext)) ||
					(errAhead != nil && (s.peek() == '+' || s.peek() == '-' || s.peek() == '_')) {
					// 10e+, 10_, 10e-
					errors = append(errors, errDigitExpected)
					break
				}
				if errBehind == nil && (s.peek() == '+' || s.peek() == '-') && headBefore != 'e' && headBefore != 'E' {
					// 1+1, 1-1
					errors = append(errors, errDigitExpected)
					break
				}
				// current: .
				if (errAhead == nil /* not last char */ && s.peek() == '.' && headNext != 'e' && headNext != 'E' && (headNext == '+' || headNext == '-' || headNext == '_')) ||
					// only non-digit allowed after . is e or E
					(errAhead != nil /* last char */ && s.peek() == '.' && headBefore != 'e' && headBefore != 'E') {
					// it's ok to end at the end as long as it's digit behind
					errors = append(errors, errDigitExpected)
					break
				}

				s.advance()
				digitLength++
			} else {
				if (zeroStart && digitLength == 1 && (s.peek() == 'n' || s.peek() == 'N')) /* 0n */ ||
					(!zeroStart && s.peek() == 'n' || s.peek() == 'N') /* 122930n */ {
					s.advance()
					// early drop; we find an N we are done! (BigInt)
				} else if isAlphaNumeric(s.peek()) {
					errors = append(errors, errNoLiteralAfterNumber)
				} else if (errAhead == nil /* not last char */ && s.peek() == '.' && !isDecimalDigit(headNext) && headNext != 'e' && headNext != 'E') ||
					// only non-digit allowed after . is e or E
					(errAhead != nil /* last char */ && s.peek() == '.' && headBefore != 'e' && headBefore != 'E') {
					// it's ok to end at the end as long as it's digit behind
					errors = append(errors, errDigitExpected)
					break
				}
				break
			}
		}
		safe = true
	case LiteralHex:
		for isLegalHexDigitIntermediate(s.peek()) {
			if s.idxHead == len(s.src)-1 {
				break
			}
			headNext := s.peekAhead(1)
			s.logger.Debug("%d: (scanHexDigits) head:%c headNext:%c digitLength:%d isDecimalDigit(headNext):%v", s.idxHead, s.peek(), headNext, digitLength, isHexDigit(headNext))
			// if next digit is a valid intermediate representation, then we can keep parsing
			// WARNING: this is naively permissive and will allow tons of illegal combinations
			if isLegalHexDigitIntermediate(headNext) {
				s.logger.Debug("%d: (scanHexDigits) advancing hex literal...", s.idxHead)
				s.advance()
				digitLength++
			} else {
				break
			}
		}
	case LiteralBinary:
		for isLegalBinaryDigitIntermediate(s.peek()) {
			if s.idxHead == len(s.src)-1 {
				break
			}
			headNext := s.peekAhead(1)
			s.logger.Debug("%d: (scanBinaryDigits) head:%c headNext:%c digitLength:%d isBinaryDigit(headNext):%v", s.idxHead, s.peek(), headNext, digitLength, isBinaryDigit(headNext))
			// if next digit is a valid intermediate representation, then we can keep parsing
			// WARNING: this is naively permissive and will allow tons of illegal combinations
			if isLegalBinaryDigitIntermediate(headNext) {
				s.logger.Debug("%d: (scanBinaryDigits) advancing binary literal...", s.idxHead)
				s.advance()
				digitLength++
			} else {
				break
			}
		}
	case LiteralOctal:
		for isLegalOctalDigitIntermediate(s.peek()) {
			if s.idxHead == len(s.src)-1 {
				break
			}
			headNext := s.peekAhead(1)
			s.logger.Debug("%d: (scanOctalDigit) head:%c headNext:%c digitLength:%d isOctalDigit(headNext):%v", s.idxHead, s.peek(), headNext, digitLength, isOctalDigit(headNext))
			// if next digit is a valid intermediate representation, then we can keep parsing
			// WARNING: this is naively permissive and will allow tons of illegal combinations
			if isLegalOctalDigitIntermediate(headNext) {
				s.logger.Debug("%d: (scanOctalDigit) advancing octal literal...", s.idxHead)
				s.advance()
				digitLength++
			} else {
				break
			}
		}
	}

	// reached here: we found something that **is not** part of a numeric literal
	// so we can add the literal token
	if len(errors) > 0 {
		s.errors = append(s.errors, errors...)
		s.logger.Debug("%d: (scanDigits) errors found:%v", s.idxHead, errors)
		return
	}

	var literalUpperBound int
	if s.idxHead+digitLength >= len(s.src) {
		literalUpperBound = len(s.src) - 1
	} else {
		literalUpperBound = s.idxHead + digitLength
	}

	if safe {
		s.addTokenSafe(TNumericLiteral, s.src[s.idxHeadStart:literalUpperBound])
	} else {
		s.addToken(TNumericLiteral, s.src[s.idxHeadStart:literalUpperBound])
	}
}

// Punctuators
//
// Scan for punctuators is straightforward:
// we group tokens by their first character, and always try to match
// the longest possible token, iteratively until we find a match.
func isPunctuation(r rune) bool {
	return r == '!' || r == '.' || r == ',' || r == '>' || r == '<' || r == '=' || r == '+' || r == '-' || r == '*' || r == '/' || r == '%' || r == '&' || r == '|' || r == '^' || r == '~' || r == '(' || r == ')' || r == '[' || r == ']' || r == '{' || r == '}' || r == ';' || r == ':' || r == '?' || r == ' '
}

func (s *Scanner) scanPunctuators() {
	switch s.peek() {
	// Simple punctuators
	case ' ':
		// s.logger.Debug("%d: <whitespace>", s.idxHead)
		return
		// nothing
	case '}':
		s.addToken(TRightBrace, nil)
	case '{':
		s.addToken(TLeftBrace, nil)
	case '(':
		s.addToken(TLeftParen, nil)
	case ')':
		s.addToken(TRightParen, nil)
	case '[':
		s.addToken(TLeftBracket, nil)
	case ']':
		s.addToken(TRightBracket, nil)
	case ';':
		s.addToken(TSemicolon, nil)
	case ':':
		s.addToken(TColon, nil)
	case '>':
		if s.seekMatchSequence([]rune{'='}) {
			// >= is greater equal
			s.addToken(TGreaterThanEqual, nil)
		} else if s.seekMatchSequence([]rune{'>', '>', '='}) {
			// >>>= is unsigned right shift assign
			s.addToken(TUnsignedRightShiftAssign, nil)
		} else if s.seekMatchSequence([]rune{'>', '>'}) {
			// >>> is unsigned right shift
			s.addToken(TUnsignedRightShift, nil)
		} else if s.seekMatchSequence([]rune{'>', '='}) {
			// >>= is right shift assign
			s.addToken(TRightShiftAssign, nil)
		} else if s.seekMatchSequence([]rune{'>'}) {
			// >> is right shift
			s.addToken(TRightShift, nil)
		} else {
			// > is greater than
			s.addToken(TGreaterThan, nil)
		}
	case '<':
		if s.seekMatchSequence([]rune{'='}) {
			// <= is greater equal
			s.addToken(TLessThanEqual, nil)
		} else if s.seekMatchSequence([]rune{'<', '='}) {
			// <<= is left shift assign
			s.addToken(TLeftShiftAssign, nil)
		} else if s.seekMatchSequence([]rune{'<'}) {
			// << is left shift
			s.addToken(TLeftShift, nil)
		} else {
			// < is greater than
			s.addToken(TLessThan, nil)
		}
	case '.':
		if s.seekMatchSequence([]rune{'.', '.'}) {
			s.addToken(TEllipsis, nil)
		} else {
			s.addToken(TPeriod, nil)
		}
	case '?':
		if s.seekMatchSequence([]rune{'?'}) {
			s.addToken(TDoubleQuestionMark, nil)
		} else {
			s.addToken(TQuestionMark, nil)
		}
	case '!':
		if s.seekMatchSequence([]rune{'=', '='}) {
			s.addToken(TStrictNotEqual, nil)
		} else if s.seekMatchSequence([]rune{'='}) {
			s.addToken(TNotEqual, nil)
		} else {
			s.addToken(TBang, nil)
		}
	case '=':
		if s.seekMatchSequence([]rune{'=', '='}) {
			s.addToken(TStrictEqual, nil)
		} else if s.seekMatchSequence([]rune{'>'}) {
			s.addToken(TArrow, nil)
		} else if s.seekMatchSequence([]rune{'='}) {
			s.addToken(TEqual, nil)
		} else {
			s.addToken(TAssign, nil)
		}
	case '&':
		if s.seekMatchSequence([]rune{'&', '='}) {
			s.addToken(TLogicalAndAssign, nil)
		} else if s.seekMatchSequence([]rune{'&'}) {
			s.addToken(TLogicalAnd, nil)
		} else if s.seekMatchSequence([]rune{'='}) {
			s.addToken(TAndAssign, nil)
		} else {
			s.addToken(TAnd, nil)
		}
	case '|':
		if s.seekMatchSequence([]rune{'|', '='}) {
			s.addToken(TLogicalOrAssign, nil)
		} else if s.seekMatchSequence([]rune{'|'}) {
			s.addToken(TLogicalOr, nil)
		} else if s.seekMatchSequence([]rune{'='}) {
			s.addToken(TOrAssign, nil)
		} else {
			s.addToken(TOr, nil)
		}
	case '+':
		if s.seekMatchSequence([]rune{'+'}) {
			s.addToken(TPlusPlus, nil)
		} else if s.seekMatchSequence([]rune{'='}) {
			s.addToken(TPlusAssign, nil)
		} else {
			s.addToken(TPlus, nil)
		}
	case '-':
		if s.seekMatchSequence([]rune{'-'}) {
			s.addToken(TMinusMinus, nil)
		} else if s.seekMatchSequence([]rune{'='}) {
			s.addToken(TMinusAssign, nil)
		} else {
			s.addToken(TMinus, nil)
		}
	case '*':
		if s.seekMatchSequence([]rune{'='}) {
			s.addToken(TStarAssign, nil)
		} else {
			s.addToken(TStar, nil)
		}
	case '/':
		if s.seekMatchSequence([]rune{'='}) {
			s.addToken(TSlashAssign, nil)
		} else {
			s.addToken(TSlash, nil)
		}
	}
}
