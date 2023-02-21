package lex

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
