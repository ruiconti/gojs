package main

import (
	"fmt"
	"log"
)

func isPunctuation(r rune) bool {
	return r == '!' || r == '.' || r == ',' || r == '>' || r == '<' || r == '=' || r == '+' || r == '-' || r == '*' || r == '/' || r == '%' || r == '&' || r == '|' || r == '^' || r == '~' || r == '(' || r == ')' || r == '[' || r == ']' || r == '{' || r == '}' || r == ';' || r == ':' || r == '?' || r == ' '
}

func isDecimalDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

func isValidDecDigitIntermediate(r rune) bool {
	return isDecimalDigit(r) || r == '_' || r == '.' || r == 'e' || r == 'E' || r == '+' || r == '-'
}

func isHexDigit(r rune) bool {
	return isDecimalDigit(r) || (r >= 'a' && r <= 'f') || (r >= 'A' && r <= 'F')
}

var errEOF = fmt.Errorf("EOF")

type Scanner struct {
	// source string being scanned
	src string
	// offset of the current index of the current lexeme
	idxHead int
	// offset of the beginning of the current lexeme
	idxHeadStart int
	// token slice
	tokens []Token
}

func NewScanner(src string) *Scanner {
	return &Scanner{
		src:          src,
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

func (s *Scanner) peekAhead(n int) (rune, error) {
	if s.idxHead+n >= len(s.src) {
		return rune(0), errEOF
	}
	return rune(s.src[s.idxHead+n]), nil
}

func (s *Scanner) addToken(t TokenType, literal interface{}) {
	var lexeme string
	// we only advance the head before calling `addToken` in `seekMatchSequence`
	// for single matches e.g. `>`, we need to check it here so we don't overly complicate
	// the scan() function
	if s.idxHeadStart == s.idxHead {
		log.Printf("%d: addToken (idxHeadStart == idxHead)", s.idxHead)
		lexeme = string(s.peek())
	} else {
		lexeme = s.src[s.idxHeadStart : s.idxHead+1]
	}
	log.Printf("%d: addToken: %v %v", s.idxHead, t, lexeme)
	s.tokens = append(s.tokens, Token{
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
		// log.Printf("(i:%d;j:%d) head: %c seekSeq: (eq %c %c)", i, j, s.peek(), cursorGot, cursorExpected)

		if cursorGot != cursorExpected {
			// log.Printf("(i:%d;j:%d) head: %c seekSeq (leaving)", i, j, s.peek())
			return false
		}
		i++
		j++
	}
	// we got here, everything is the same :)
	// runesstr := ``
	// for _, r := range candidates {
	// 	runesstr += fmt.Sprintf("%c", r)
	// }
	// log.Printf("%d: matchSequence true! %v", headIdx, runesstr)
	s.advanceBy(len(sequence))
	// log.Printf("(i:%d;j:%d) head: %c seekSeq (advanced %d)", i, j, rune(s.peek()), len(sequence))
	return true

}

func (s *Scanner) Scan() []Token {
	for s.idxHead < len(s.src) {
		s.idxHeadStart = s.idxHead
		head := s.peek()

		if isDecimalDigit(head) || head == '.' {
			s.scanDigits()
		}
		if isPunctuation(rune(head)) {
			s.scanPunctuators()
		}
		log.Printf("%d: (%c) next iter..", s.idxHead, s.peek())
		s.advance()
	}
	return s.tokens
}

// Numeric literals are defined exclusively by the _following_ production:
// "NumericLiteral":
// | DecimalLiteral
// | HexIntegerLiteral
//
// "DecimalLiteral":
// | DecimalIntegerLiteral . DecimalDigits? ExponentPart?
// | . DecimalDigits ExponentPart?
// | DecimalIntegerLiteral ExponentPart?
//
// "DecimalDigits":
// | DecimalDigit
// | DecimalDigits DecimalDigit
//
// "DecimalIntegerLiteral"
// | 0
// | NonZeroDigit DecimalDigit?
//
// "DecimalDigit": 0 1 2 3 4 5 6 7 8 9
// "NonZeroDigit": 1 2 3 4 5 6 7 8 9
//
// "SignedInteger":
// | + DecimalDigit
// | - DecimalDigit
//
// "HexIntegerLiteral":
// | 0x HexDigit
// | 0X HexDigit
// | HexIntegerLiteral HexDigit
//
// "HexDigit": 0 1 2 3 4 5 6 7 8 9 a b c d e f A B C D E F
// "ExponentPart": ExponentIndicator SignedInteger
// "ExponentIndicator": e E
type NumberType string

const (
	NumberTypeDecimal NumberType = "decimal"
	NumberTypeHex     NumberType = "hex"
)

func (s *Scanner) scanDigits() {
	headNext, _ := s.peekAhead(1)

	log.Printf("%d: (scanDigits:entry) head:%c headNext:%c", s.idxHead, s.peek(), headNext)
	// first bytes will tell us what type of number, if any, we have to parse
	var numberType NumberType
	if s.peek() == '0' /* head is '0' */ &&
		(headNext == 'x' || headNext == 'X' /* allowed next char */) {
		numberType = NumberTypeHex
	} else if isDecimalDigit(s.peek()) /* head is a digit OR */ ||
		s.peek() == '.' && isDecimalDigit(headNext) /* head is '.' and headNext is a digit */ {
		numberType = NumberTypeDecimal
	} else {
		// invalid number
		return
	}

	digitLength := 1
	switch numberType {
	case NumberTypeDecimal:
		for isValidDecDigitIntermediate(s.peek()) {
			headNext, err := s.peekAhead(1)
			log.Printf("%d: (scanDecimalDigits) head:%c headNext:%c digitLength:%d err:%v isDecimalDigit(headNext):%v", s.idxHead, s.peek(), headNext, digitLength, err, isDecimalDigit(headNext))
			// if next digit is a valid intermediate representation, then we can store it:)
			// this is overly naive and will allow tons of shit
			if err == nil && isValidDecDigitIntermediate(headNext) {
				log.Printf("%d: (scanDecimalDigits) advancing decimal number...", s.idxHead)
				s.advance()
				digitLength++
			} else {
				break
			}
		}
	case NumberTypeHex:
		for isHexDigit(s.peek()) || s.peek() == 'x' || s.peek() == 'X' {
			headNext, err := s.peekAhead(1)
			log.Printf("%d: (scanHexDigits) head:%c headNext:%c digitLength:%d err:%v isDecimalDigit(headNext):%v", s.idxHead, s.peek(), headNext, digitLength, err, isDecimalDigit(headNext))
			// if next digit is a valid intermediate representation, then we can store it:)
			// this is overly naive and will allow tons of shit
			if err == nil && (isHexDigit(headNext) || headNext == 'x' || headNext == 'X') {
				log.Printf("%d: (scanHexDigits) advancing decimal number...", s.idxHead)
				s.advance()
				digitLength++
			} else {
				break
			}
		}
	}

	// reached here: we found something that **is not** a digit
	// so we can add the literal token

	var literalUpperBound int
	if s.idxHead+digitLength >= len(s.src) {
		literalUpperBound = len(s.src) - 1
	} else {
		literalUpperBound = s.idxHead + digitLength
	}

	s.addToken(TNumericLiteral, s.src[s.idxHeadStart:literalUpperBound])
}

// Scan for punctuators is straightforward:
// we group tokens by their first character, and always try to match
// the longest possible token, iteratively until we find a match.
func (s *Scanner) scanPunctuators() {
	switch s.peek() {
	// Simple punctuators
	case ' ':
		// log.Printf("%d: <whitespace>", s.idxHead)
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

func main() {
	runes := []rune{'!', '.', ',', '>', '<', '=', '+', '-', '*', '/', '%', '&', '|', '^', '~', '(', ')', '[', ']', '{', '}', ';', ':', '?', ' ', '\t', '\r', '\n'}
	for _, r := range runes {
		fmt.Printf("%s - %d\n", string(r), r)
	}
}
