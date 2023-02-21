package lex

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
