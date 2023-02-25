package lex

import "github.com/ruiconti/gojs/internal"

func isIdentifierStart(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || r == '$' || r == '_'
}

func isIdentifierPart(r rune) bool {
	return isIdentifierStart(r) || isDecimalDigit(r) || r == '\\'
}

var keywords = internal.MapInvert(ReservedWordNames)

// Identifiers
//
// https://262.ecma-international.org/#sec-names-and-keywords
func (s *Scanner) scanIdentifiers() {
	if !isIdentifierStart(s.peek()) {
		return
	}

	s.advance()
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

	// Try to parse it as a reserved word
	candidate := s.src[lower:upper]
	if token, ok := keywords[candidate]; ok {
		s.addToken(token, candidate)
		return
	}

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
