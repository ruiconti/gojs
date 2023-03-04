package lex

// String literals
//
// https://262.ecma-international.org/#sec-literals-string-literals

func isValidEscapeSequence(r rune) bool {
	return r == '0' || r == 'b' || r == 'f' || r == 'n' || r == 'r' || r == 't' || r == 'v' || r == '\\' || r == '\'' || r == '"'
}

func (s *Scanner) rejectEscape(c *int) error {
	char, errCur := s.peekN(*c)
	if char != '\\' || errCur != nil {
		// not in an escape sequence
		return nil
	}

	// consume escape
	*c = *c + 1
	char, errCur = s.peekN(*c)
	if errCur != nil {
		// can't EOF in the middle of an escape sequence
		return errBadEscapeSequence
	} else if char == 'u' || char == 'x' {
		// try to consume a numeric escape sequence

		// go back to the escape, as the entry point to that production
		// is the escape itself
		*c = *c - 1
		if !s.acceptEscapedUnicode(c) {
			return errBadEscapeSequence
		}
	} else if isDecimalDigit(char) && char != '0' {
		// not allowing octal escape sequences for now
		return errInvalidNumericEscape
	}

	s.logger.Debug("[%d:%d] scanStringLiteral:escape: %c", s.idxHead, *c, char)
	// consume whatever is after the escape
	// *c = *c + 1
	return nil
}

func (s *Scanner) scanStringLiteral() {
	var strType StringLiteralType
	char := s.peek()

	if char == '"' {
		strType = DoubleQuote
	} else if char == '\'' {
		strType = SingleQuote
	} else {
		return
	}

	// consume the quote
	s.advanceBy(1)
	// we are capturing only the contents, literal part of the string, so we advance
	// the start index too
	s.idxHeadStart = s.idxHead
	cursor := 0
	for {
		char, err := s.peekN(cursor)
		if err != nil {
			// we found an EOF in mid-string, so this is an error
			s.errors = append(s.errors, errUnterminatedStringLiteral)
			return
		}

		s.logger.Debug("[%d:%d] scanStringLiteral: %c", s.idxHead, cursor, char)
		// strings can accept virtually all characters, so we just check for the exceptions
		// which are the escape sequence and the end quote
		if char == '\\' {
			// first try to find escape sequences
			escapeErr := s.rejectEscape(&cursor)
			if escapeErr != nil {
				s.errors = append(s.errors, escapeErr)
				return
			}
		}

		if strType == DoubleQuote && char == '"' {
			// end quote; consume it and be done with it :)
			s.logger.Debug("[%d:%d] scanStringLiteral:end %c", s.idxHead, cursor, char)
			break
		} else if strType == SingleQuote && char == '\'' {
			// end quote; consume it and be done with it :)
			s.logger.Debug("[%d:%d] scanStringLiteral:end %c", s.idxHead, cursor, char)
			break
		} else {
			// carry on :)
			cursor++
		}
	}

	if s.idxHead+cursor == s.idxHead {
		// we found an empty string, "" or '', which doesn't increment the cursor
		// so we do that here
		// THIS IS BEING HANDLED AT addTokenSafe
	}

	// update the head so we can capture the literal
	s.idxHead += cursor

	if strType == SingleQuote {
		s.addTokenSafe(TStringLiteral_SingleQuote)
	} else if strType == DoubleQuote {
		s.addTokenSafe(TStringLiteral_DoubleQuote)
	}

	// walk the cursor past the end quote
	s.advanceBy(1)
}
