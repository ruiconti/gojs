package lex

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
	case ',':
		s.addToken(TComma, nil)
	case ':':
		s.addToken(TColon, nil)
	case '~':
		s.addToken(TTilde, nil)
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
