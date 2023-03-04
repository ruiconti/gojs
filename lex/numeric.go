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
	return isDecimalDigit(r) || r == '_' || r == '.' || r == 'e' || r == 'E'
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

// rejectExponential tries to parse `e | E` production
func (s *Scanner) rejectExponential(c *int) error {
	char, errCur := s.peekN(*c)
	if errCur != nil || (char != 'e' && char != 'E') {
		return nil
	}

	charNext, errNext := s.peekN(*c + 1)
	if errNext != nil {
		return errUnexpectedToken // ExponentIndicator cannot be reduced to a terminal
	}

	// prod: ExponentIndicator SignedInteger
	// prod: SignedInteger -> DecimalDigits
	if isDecimalDigit(charNext) {
		*c = *c + 1
		return nil
	}
	// prod: SignedInteger -> (+|-) DecimalDigits
	if charNext == '+' || charNext == '-' {
		*c = *c + 1
		charNext, errNext = s.peekN(*c + 1)

		if errNext != nil {
			return errUnexpectedToken // +|- are not terminals in this production
		}
		if isDecimalDigit(charNext) {
			*c = *c + 1
			return nil
		}
	}
	return errUnexpectedToken
}

// rejectNumberSeparator tries to parse `+ | - | _` production
func (s *Scanner) rejectNumberSeparator(c *int) error {
	char, errCur := s.peekN(*c)
	// if errCur != nil || (char != '+' && char != '-' && char != '_') {
	if errCur != nil || char != '_' {
		return nil
	}
	charNext, errNext := s.peekN(*c + 1)

	// check possible transitions from _
	shouldReject := errNext == nil && !isDecimalDigit(charNext) || errNext != nil
	if shouldReject {
		return errUnexpectedToken
	}
	*c = *c + 1
	return nil
}

func (s *Scanner) rejectDot(c *int) error {
	char, errCur := s.peekN(*c)
	if errCur != nil || char != '.' {
		return nil
	}

	charNext, errNext := s.peekN(*c + 1)
	// prod: DecimalDigit . DecimalDigits ExponentPart?
	if errNext == nil && (charNext == 'e' || charNext == 'E') {
		*c = *c + 1
		if err := s.rejectExponential(c); err != nil {
			return err
		} else {
			return nil
		}
	}

	// prod: DecimalDigit . DecimalDigits
	if errNext == nil && isDecimalDigit(charNext) {
		*c = *c + 1
		return nil
	}

	return errUnexpectedToken
}

func (s *Scanner) rejectBigInt(c *int, zeroStart bool) error {
	char, errCur := s.peekN(*c)
	if errCur != nil || (char != 'n' && char != 'N') {
		return nil
	}

	charNext, errNext := s.peekN(*c + 1)
	// check possible transitions from n | N and check whether N was created in a valid way
	shouldReject := zeroStart && *c > 1 || errNext == nil && isAlphaNumeric(charNext)
	if shouldReject {
		return errUnexpectedToken
	}

	*c = *c + 1
	return nil
}

func (s *Scanner) prettyPrintScan() {
	s.logger.Debug("%v", s.src)
	cursor := []byte{}
	for i := 0; i < s.idxHead; i++ {
		cursor = append(cursor, ' ')
	}
	cursor = append(cursor, '^')
	s.logger.Debug("%s", cursor)
}

func (s *Scanner) scanDigits() {
	char := s.peek()
	charNext, err := s.peekN(1)

	zeroStart := false     // it's possible to define 0n 000n so we need to track whether a number starts with 0
	tmpCursor := s.idxHead // allow us backtrack

	// 0: current char is digit | .
	// parse initial state; 0x | 0X | 0b | 0B | 0o | 0O | 0 | 1-9 | .
	// and advances the cursor to the next digit chars
	var numberType NumericLiteralType
	if char == '0' && err == nil {
		if charNext == 'x' || charNext == 'X' {
			numberType = LiteralHex
			s.advanceBy(2)
		} else if charNext == 'b' || charNext == 'B' {
			numberType = LiteralBinary
			s.advanceBy(2)
		} else if charNext == 'o' || charNext == 'O' {
			numberType = LiteralOctal
			s.advanceBy(2)
			zeroStart = true
		} else {
			numberType = LiteralDecimal
		}
	} else if isDecimalDigit(char) {
		// we don't need to advance the cursor here
		// because we are already at a valid number
		numberType = LiteralDecimal
	} else if char == '.' {
		dotCursor := 0
		err := s.rejectDot(&dotCursor)
		if err != nil && (charNext == '_' || charNext == 'e' || charNext == 'E' || charNext == '+' || charNext == '-') {
			// not a valid number and illegal syntax
			// maybe there's a better place to capture it, but we're doing it here for now
			s.errors = append(s.errors, err)
			return
		} else if err != nil {
			// not a valid number starting with . but can be valid syntax
			// e.g member access
			return
		}
		numberType = LiteralDecimal
		// can be a valid number
	} else {
		return
	}

	cursor := 0
	errors := []error{}

	// 1: parse digits
	switch numberType {
	case LiteralDecimal:
		for {
			char, err := s.peekN(cursor)
			if err != nil {
				break
			}
			s.logger.Debug("[%d:%d] scan: %s", s.idxHead, cursor, string(char))
			if isLegalDecDigitIntermediate(char) {
				if char == '.' {
					err = s.rejectDot(&cursor)
					if err != nil {
						errors = append(errors, err)
						s.logger.Debug("[%d:%d] foundError:%v", s.idxHead, cursor, err)
						break
					}
				} else if char == 'e' || char == 'E' {
					err = s.rejectExponential(&cursor)
					if err != nil {
						errors = append(errors, err)
						s.logger.Debug("[%d:%d] foundError:%v", s.idxHead, cursor, err)
						break

					}
					// } else if char == '+' || char == '-' || char == '_' {
				} else if char == '_' {
					// current + | - | _
					err = s.rejectNumberSeparator(&cursor)
					if err != nil {
						errors = append(errors, err)
						s.logger.Debug("[%d:%d] foundError:%v", s.idxHead, cursor, err)
						break
					}
				} else {
					cursor++
				}
			} else {
				if char == 'n' || char == 'N' {
					err = s.rejectBigInt(&cursor, zeroStart)
					if err != nil {
						errors = append(errors, err)
						s.logger.Debug("[%d:%d] foundError:%v", s.idxHead, cursor, err)
					}
				} else if isAlphaNumeric(char) {
					errors = append(errors, errNoLiteralAfterNumber)
					s.logger.Debug("[%d:%d] foundError:%v", s.idxHead, cursor, errUnexpectedToken)
				}
				break
			}
		}
	case LiteralHex:
		for isLegalHexDigitIntermediate(s.peek()) {
			char, err := s.peekN(cursor)
			if err != nil {
				break
			}
			// if next digit is a valid intermediate representation, then we can keep parsing
			if isLegalHexDigitIntermediate(char) {
				cursor++
			} else {
				break
			}
		}
	case LiteralBinary:
		for isLegalBinaryDigitIntermediate(s.peek()) {
			char, err := s.peekN(cursor)
			if err != nil {
				break
			}
			// if next digit is a valid intermediate representation, then we can keep parsing
			if isLegalBinaryDigitIntermediate(char) {
				cursor++
			} else {
				break
			}
		}
	case LiteralOctal:
		for isLegalOctalDigitIntermediate(s.peek()) {
			char, err := s.peekN(cursor)
			if err != nil {
				break
			}
			// if next digit is a valid intermediate representation, then we can keep parsing
			if isLegalOctalDigitIntermediate(char) {
				cursor++
			} else {
				break
			}
		}
	}

	if len(errors) > 0 {
		s.idxHead = tmpCursor // backtracks
		s.errors = append(s.errors, errors...)
		s.logger.Debug("%d: (scanDigits) errors found:%v", s.idxHead, errors)
		return
	}

	s.idxHead = s.idxHead + cursor
	s.addTokenSafe(TNumericLiteral)
}
