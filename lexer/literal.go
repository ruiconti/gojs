package lexer

// String literals
type StringLiteralType string

const (
	SingleQuote StringLiteralType = "singleq"
	DoubleQuote StringLiteralType = "doubleq"
)

func isValidEscapeSequence(r rune) bool {
	return r == '0' || r == 'b' || r == 'f' || r == 'n' || r == 'r' || r == 't' || r == 'v' || r == '\\' || r == '\'' || r == '"'
}

// String literals
//
// https://262.ecma-international.org/#sec-literals-string-literals
func (s *Lexer) scanStringLiteral() Token {
	var (
		strType StringLiteralType
		err     error
		char    rune
	)

	char = s.Peek()
	switch char {
	case '"':
		strType = DoubleQuote
	case '\'':
		strType = SingleQuote
	default:
		s.Errorf("scanStringLiteral: unexpected char: %c", char)
		return TokenUnknown
	}

	s.Next() // consume the quote
	// we are capturing only the contents, literal part of the string, so we advance
	// the start index too
	// literalStart := s.srcCursorHead
	s.PeekLoop(func(ch rune) bool {
		if ch == EOF {
			// we found an EOF in mid-string, so this is an error
			s.Errorf(errUnterminatedStringLiteral.Error())
			return false
		}

		// strings can accept virtually all characters, so we just check for the exceptions
		// which are the escape sequence and the end quote
		switch ch {
		case '\\':
			if err = s.rejectEscapedSequence(); err != nil {
				s.Errorf(errUnterminatedStringLiteral.Error())
				return false
			}
		case '"':
			s.Next()
			if strType == DoubleQuote {
				// end of string
				return false
			}
		case '\'':
			s.Next()
			if strType == SingleQuote {
				// end of string
				return false
			}
		default:
			// carry on
			s.Next()
		}
		return true
	})

	var typ TokenType
	if strType == SingleQuote {
		typ = TStringLiteral_SingleQuote
	} else if strType == DoubleQuote {
		typ = TStringLiteral_DoubleQuote
	}
	tok := s.CreateLiteralToken(typ)
	return tok
}

// Numeric literal
type NumericLiteralType string

const (
	LiteralDecimal NumericLiteralType = "decimal"
	LiteralHex     NumericLiteralType = "hex"
	LiteralBinary  NumericLiteralType = "binary"
	LiteralOctal   NumericLiteralType = "octal"
)

func isDec(r rune) bool      { return r >= '0' && r <= '9' }
func isDecInter(r rune) bool { return isDec(r) || r == '_' || r == '.' || r == 'e' || r == 'E' }

func isHex(r rune) bool      { return isDec(r) || (r >= 'a' && r <= 'f') || (r >= 'A' && r <= 'F') }
func isHexInter(r rune) bool { return isHex(r) || r == '_' }

func isBinary(r rune) bool      { return r == '0' || r == '1' }
func isBinaryInter(r rune) bool { return isBinary(r) || r == '_' }

func isOctal(r rune) bool      { return r >= '0' && r <= '7' }
func isOctalInter(r rune) bool { return isOctal(r) || r == '_' }

// rejectExponentialPart tries to reject a cursor at ExponentPart production
//
// DecimalLiteral ::
// DecimalIntegerLiteral . DecimalDigits? ExponentPart?
// | "." DecimalDigits ExponentPart?
//
// ExponentPart[Sep] ::
// ExponentIndicator SignedInteger
func (s *Lexer) rejectExponentialPart() error {
	char := s.Peek()
	if char == EOF || (char != 'e' && char != 'E') {
		return errUnexpectedToken
	}

	char = s.PeekN(1)
	if char == EOF {
		return errUnexpectedToken // ExponentIndicator cannot be reduced to a terminal
	}

	// ExponentIndicator SignedInteger
	switch {
	case isDec(char):
		// SignedInteger -> DecimalDigits
		s.Next() // consume 'e' | 'E'
		return nil
	case char == '+' || char == '-':
		// SignedInteger -> (+|-) DecimalDigits
		s.Next() // consume 'e' | 'E'
		char = s.PeekN(1)
		if char == EOF || !isDec(char) {
			return errUnexpectedToken
		}
		s.Next() // consume '+' | '-'
		return nil
	}

	return errUnexpectedToken
}

// rejectNumericLiteralSep tries to reject a cursor at NumericLiteralSeparator
//
// Productions:
// DecimalBigIntegerLiteral -> DecimalDigits NumericLiteralSeparator DecimalDigit
// DecimalIntegerLiteral -> NonZeroDigit NumericLiteralSeparator? DecimalDigits
// DecimalDigits -> DecimalDigits NumericLiteralSeparator DecimalDigit
// BinaryDigits -> BinaryDigits NumericLiteralSeparator BinaryDigit
func (s *Lexer) rejectNumericLiteralSep() error {
	char := s.Peek()
	if char == EOF || char != '_' {
		return errUnexpectedToken
	}

	char = s.PeekN(1)
	// prod: Separator DecimalDigit
	if char == EOF || char != EOF && !isDec(char) {
		return errUnexpectedToken
	}

	s.Next() // consume '_'
	return nil
}

// rejectDot tries to reject a cursor at "."
//
// Productions:
// DecimalLiteral -> . DecimalDigits ExponentPart?
func (s *Lexer) rejectDot() error {
	char := s.Peek()
	if char == EOF || char != '.' {
		return nil
	}

	char = s.PeekN(1)
	if char == EOF {
		return errUnexpectedToken
	}
	switch {
	case char == 'e' || char == 'E':
		// DecimalDigit . DecimalDigits ExponentPart?
		s.Next() // consume '.'
		if err := s.rejectExponentialPart(); err != nil {
			return err
		} else {
			return nil
		}
	case isDec(char):
		// DecimalDigit . DecimalDigits
		s.Next() // consume '.'
		return nil
	}

	return errUnexpectedToken
}

// rejectBigInt tries to reject a cursor at DecimalBigIntegerLiteral
//
// Productions:
// NumericLiteral -> DecimalBigIntegerLiteral
func (s *Lexer) rejectBigInt(zeroStart, hasDot bool) error {
	char := s.Peek()
	if char == EOF || (char != 'n' && char != 'N') || hasDot {
		return errUnexpectedToken
	}

	char = s.PeekN(1)
	if char == EOF {
		s.Next() // consume 'n' | 'N'
		return nil
	}

	// try to validate without looking back
	// DecimalBigIntegerLiteral ::
	// | 0 BigIntLiteralSuffix
	// | NonZeroDigit DecimalDigits[+Sep]opt BigIntLiteralSuffix
	// | NonZeroDigit NumericLiteralSeparator DecimalDigits[+Sep] BigIntLiteralSuffix
	offset := s.srcCursorHead - s.srcCursor
	shouldReject := zeroStart && offset > 1 || isAlphaNumeric(char)
	if shouldReject {
		return errUnexpectedToken
	}

	s.Next() // consume 'n' | 'N'
	return nil
}

// Numeric literals
//
// https://262.ecma-international.org/#sec-literals-numeric-literals
func (s *Lexer) scanNumericLiteral() Token {
	char, charNext := s.Peek(), s.PeekN(1)
	var zeroStart, hasDot bool

	// 0: current char is digit | .
	// parse initial state; 0x | 0X | 0b | 0B | 0o | 0O | 0 | 1-9 | .
	// and advances the cursor to the next digit chars
	var numberType NumericLiteralType
	if char == '0' && charNext != EOF {
		switch charNext {
		case 'x', 'X':
			numberType = LiteralHex
			s.Jump(2)
		case 'b', 'B':
			numberType = LiteralBinary
			s.Jump(2)
		case 'o', 'O':
			numberType = LiteralOctal
			s.Jump(2)
			zeroStart = true
		default:
			numberType = LiteralDecimal
			zeroStart = true
		}
	} else if isDec(char) {
		// we don't need to advance the cursor here
		// because we are already at a valid number
		numberType = LiteralDecimal
	} else if char == '.' {
		hasDot = true
		err := s.rejectDot()
		if err != nil {
			// not a valid number starting with . but can be valid punctuation
			return s.scanPunctuation()
		}
		numberType = LiteralDecimal // can be a valid number
	}

	// 1: parse digits
	switch numberType {
	case LiteralDecimal:
		s.PeekLoop(func(ch rune) bool {
			var err error
			if char == EOF {
				return false
			}

			if isDecInter(ch) {
				// regular decimal number
				switch ch {
				case '.':
					hasDot = true
					err = s.rejectDot()
				case 'e', 'E':
					err = s.rejectExponentialPart()
				case '_':
					err = s.rejectNumericLiteralSep()
				default:
					s.Next()
				}
				if err != nil {
					s.Errorf(err.Error())
					return false
				}
			} else {
				// potential bigint
				switch ch {
				case 'n', 'N':
					err = s.rejectBigInt(zeroStart, hasDot)
					if err != nil {
						s.Errorf(err.Error())
						return false
					}
				default:
					if isAlphaNumeric(ch) {
						s.Errorf(errNoLiteralAfterNumber.Error())
					}
					return false
				}
			}
			return true
		})

	case LiteralHex:
		s.PeekLoop(func(ch rune) bool {
			if isHexInter(ch) {
				s.Next()
				return true
			} else {
				return false
			}
		})

	case LiteralBinary:
		s.PeekLoop(func(ch rune) bool {
			if isBinaryInter(ch) {
				s.Next()
				return true
			} else {
				return false
			}
		})

	case LiteralOctal:
		s.PeekLoop(func(ch rune) bool {
			if isOctalInter(ch) {
				s.Next()
				return true
			} else {
				return false
			}
		})
	}

	if len(s.errors) > 0 {
		s.srcCursorHead = s.srcCursor
		return TokenUnknown
	}

	return s.CreateLiteralToken(TNumericLiteral)
}

func isId(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || r == '$' || r == '_' || r == '\\'
}
func isIdInter(r rune) bool { return isId(r) || isDec(r) || r == '\\' }

// rejectEscapedUnicode tries to reject a cursor at "\"
//
// EscapeSequence -> UnicodeEscapeSequence ::
// | u Hex4Digits
// | u{ CodePoint }
func (s *Lexer) rejectEscapedUnicode() error {
	var char rune

	if char = s.Peek(); char == EOF || char != '\\' {
		return errInvalidEscapedSequence
	}

	// 0: current token is \\
	char = s.PeekN(1) // lookAhead past '\\'

	switch char {
	case EOF:
		s.Errorf(errEOF.Error())
		return errEOF
	case 'u':
		s.Next() // 0: consume '\\'
		s.Next() // 1: consume 'u'

		// 2-5: consume 4 hex digits
		for i := 0; i < 4; i++ {
			char = s.Peek()
			if isHex(char) {
				s.Next() // 2-5: consume '0'-'9' | 'a'-'f' | 'A'-'F'
			} else {
				s.Errorf(errInvalidEscapedSequence.Error())
				return errInvalidEscapedSequence
			}
		}
	default:
		return errInvalidEscapedSequence
	}

	return nil
}

// rejectEscapeInString tries to reject a cursor at "\"
//
// LineContinuation ::
// \ EscapeSequence
//
// EscapeSequence ::
// | CharacterEscapeSequence
// | 0 [lookahead ∉ DecimalDigit]
// | LegacyOctalEscapeSequence(!)
// | NonOctalDecimalEscapeSequence(!)
// | HexEscapeSequence
// | UnicodeEscapeSequence
//
// !: TODO
func (s *Lexer) rejectEscapedSequence() error {
	var char rune

	if char = s.Peek(); char == EOF || char != '\\' {
		return errInvalidEscapedSequence
	}

	// 0: current token is \\
	char = s.PeekN(1) // lookAhead past '\\'

	switch {
	case char == EOF:
		s.Errorf(errEOF.Error())
		return errEOF
	case char == 'u':
		err := s.rejectEscapedUnicode()
		if err != nil {
			return err
		}
	case char == 'x':
		// HexEscapeSequence -> x HexDigit HexDigit
		s.Next() // 0: consume '\\'
		s.Next() // 1: consume 'x'

		char = s.Peek()
		if isHex(char) {
			s.Next() // 2: consume '0'-'9' | 'a'-'f' | 'A'-'F'
		} else {
			s.Errorf(errInvalidEscapedSequence.Error())
			return errInvalidEscapedSequence
		}
	case isDec(char) && char != '0':
		// EscapeSequence -> 0 [lookahead ∉ DecimalDigit]
		return errInvalidEscapedSequence
	default:
		// EscapeSequence ->
		// CharacterEscapeSequence :: SingleEscapeCharacter | NonEscapeCharacter
		s.Jump(2) // consume '\\' and the escaped char
	}

	return nil
}

// Identifiers
//
// https://262.ecma-international.org/#sec-names-and-keywords
func (s *Lexer) scanIdentifier() Token {
	var char rune
	char = s.Peek()
	if !isId(char) {
		return TokenUnknown
	}

	s.Next() // consume the start, as identifierStart != identifierPart

	s.PeekLoop(func(char rune) bool {
		switch {
		case isIdInter(char):
			s.Next()
		case char == '\\':
			if err := s.rejectEscapedUnicode(); err != nil {
				return false
			}
			s.Next()
		default:
			return false
		}

		return true
	})

	return s.CreateLiteralToken(TIdentifier)
}

func isLegalStringLiteralIntermediate(r rune) bool {
	return r != '"'
}
func isAlphaNumeric(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9')
}

func (s *Lexer) isValidEscapeSequence(r rune) bool {
	switch r {
	case 'b', 'f', 'n', 'r', 't', 'v', '\\', '"', '\'', '0':
		return true
	}
	return false
}
