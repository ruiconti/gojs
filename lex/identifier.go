package lex

import "github.com/ruiconti/gojs/internal"

func isIdentifierStart(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || r == '$' || r == '_'
}

func isIdentifierPart(r rune) bool {
	return isIdentifierStart(r) || isDecimalDigit(r) || r == '\\'
}

var keywords = internal.MapInvert(ReservedWordNames)

func (s *Scanner) acceptEscapedUnicode(cursor *int) bool {
	errs := []error{}
	tmpCursor := cursor
	if char, err := s.peekN(*cursor); err != nil || char != '\\' {
		s.logger.Debug("[%d:%d] scanIdentifier:escape: invalid first character, backtracking: ", s.idxHead, *cursor)
		cursor = tmpCursor // backtracks
		return false
	}
	// 0: current token is \\
	s.logger.Debug("[%d:%d] scanIdentifier:escape: \\", s.idxHead, *cursor)
	// consume the escape
	*cursor++

	char, err := s.peekN(*cursor)
	if err != nil {
		// can't EOF here, because we already consumed the escape
		s.logger.Debug("[%d:%d] scanIdentifier:escape: EOF: ", s.idxHead, *cursor)
		cursor = tmpCursor // backtracks
		errs = append(errs, err)
		return false
	}

	switch char {
	case 'u':
		// 1: current token is u
		s.logger.Debug("[%d:%d] scanIdentifier:unicode: u", s.idxHead, *cursor)
		// consume u
		*cursor++

		// 2-5: consume 4 hex digits
		for i := 0; i < 4; i++ {
			char, err := s.peekN(*cursor)
			s.logger.Debug("[%d:%d] scanIdentifier:unicode: %c", s.idxHead, *cursor, char)
			if err != nil || !isHexDigit(char) {
				s.logger.Debug("[%d:%d] scanIdentifier:invalid unicode escape sequence: %c", s.idxHead, *cursor, char)
				cursor = tmpCursor // backtracks
				errs = append(errs, err)
				return false
			}
			*cursor++
		}
	case 'x':
		// 1: current token is x
		s.logger.Debug("[%d:%d] scanIdentifier:hex: x", s.idxHead, *cursor)
		// consume x
		*cursor++

		char, err := s.peekN(*cursor)
		// 2: consume the next hexadecimal digit
		if err != nil || !isHexDigit(char) {
			s.logger.Debug("[%d:%d] scanIdentifier:hex:invalid hex escape sequence: %c", s.idxHead, *cursor, char)
			cursor = tmpCursor // backtracks
			errs = append(errs, err)
			return false
		}
	default:
		// 1: not a valid numerical escape sequence
		s.logger.Debug("[%d:%d] scanIdentifier:escape:invalid numerical escape sequence: %c", s.idxHead, *cursor, char)
		cursor = tmpCursor // backtracks
		return false
	}

	// if we got here, we need to back track the cursor by 1
	// because it would have gone past the end of the unicode in the last iteration,
	// before failing the loop condition.
	*cursor--
	s.logger.Debug("[%d:%d] scanIdentifier:unicode:accepted", s.idxHead, *cursor)
	return true
}

// Identifiers
//
// https://262.ecma-international.org/#sec-names-and-keywords
func (s *Scanner) scanIdentifiers() (bool, []error) {
	errs := []error{}
	if !isIdentifierStart(s.peek()) {
		return false, errs
	}

	// consume the start, as identifierStart != identifierPart
	s.advanceBy(1)
	cursor := 0
	// consume each valid identifier character up until the first invalid one
	for {
		char, err := s.peekN(cursor)
		if err != nil || !isIdentifierPart(char) {
			break
		}
		if char == '\\' && !s.acceptEscapedUnicode(&cursor) {
			break
		}
		s.logger.Debug("[%d:%d] scanIdentifier: %c", s.idxHead, cursor, char)
		cursor++
	}

	// we may have left the loop because we hit EOF
	lower, upper := s.idxHeadStart, s.idxHead+cursor
	s.idxHead = upper

	// Try to parse it as a reserved word
	candidate := s.src[lower:upper]
	if token, ok := keywords[candidate]; ok {
		s.addTokenSafe(token)
		return true, errs
	}

	s.logger.Info("scanIdentifier: lowerBound:%d upperBound:%d literal:%s", lower, upper, candidate)
	s.addTokenSafe(TIdentifier)
	return true, errs
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
