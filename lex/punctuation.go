package lex

import "fmt"

// Punctuators
//
// Scan for punctuators is straightforward:
// we group tokens by their first character, and always try to match
// the longest possible token, iteratively until we find a match.
func isPunctuation(r rune) bool {
	return r == '!' || r == '.' || r == ',' || r == '>' || r == '<' || r == '=' || r == '+' || r == '-' || r == '*' || r == '/' || r == '%' || r == '&' || r == '|' || r == '^' || r == '(' || r == ')' || r == '[' || r == ']' || r == '{' || r == '}' || r == ';' || r == ':' || r == '?' || r == '~'
}

func (s *Scanner) scanPunctuators() (bool, []error) {
	token, c, err := s.innerScanPunctuators(s.peek())
	if err != nil {
		return false, []error{err}
	}

	s.advanceBy(c)
	s.addTokenSafe(token)
	return true, []error{}
}

func (s *Scanner) innerScanPunctuators(char rune) (TokenType, int, error) {
	switch s.peek() {
	// Simple punctuators
	case ' ':
		return TUnknown, -1, fmt.Errorf("not valid punctuation")
	case '}':
		return TRightBrace, 1, nil
	case '{':
		return TLeftBrace, 1, nil
	case '(':
		return TLeftParen, 1, nil
	case ')':
		return TRightParen, 1, nil
	case '[':
		return TLeftBracket, 1, nil
	case ']':
		return TRightBracket, 1, nil
	case ';':
		return TSemicolon, 1, nil
	case ',':
		return TComma, 1, nil
	case ':':
		return TColon, 1, nil
	case '^':
		return TXor, 1, nil
	case '~':
		return TTilde, 1, nil
	case '%':
		return TPercent, 1, nil
	case '>':
		if s.seekMatchSequence([]rune{'>', '>', '='}) {
			// >>>= is unsigned right shift assign
			return TUnsignedRightShiftAssign, 4, nil
		} else if s.seekMatchSequence([]rune{'>', '>'}) {
			// >>> is unsigned right shift
			return TUnsignedRightShift, 3, nil
		} else if s.seekMatchSequence([]rune{'>', '='}) {
			// >>= is right shift assign
			return TRightShiftAssign, 3, nil
		} else if s.seekMatchSequence([]rune{'='}) {
			// >= is greater equal
			return TGreaterThanEqual, 2, nil
		} else if s.seekMatchSequence([]rune{'>'}) {
			// >> is right shift
			return TRightShift, 2, nil
		} else {
			// > is greater than
			return TGreaterThan, 1, nil
		}
	case '<':
		if s.seekMatchSequence([]rune{'='}) {
			// <= is greater equal
			return TLessThanEqual, 2, nil
		} else if s.seekMatchSequence([]rune{'<', '='}) {
			// <<= is left shift assign
			return TLeftShiftAssign, 3, nil
		} else if s.seekMatchSequence([]rune{'<'}) {
			// << is left shift
			return TLeftShift, 2, nil
		} else {
			// < is greater than
			return TLessThan, 1, nil
		}
	case '.':
		if s.seekMatchSequence([]rune{'.', '.'}) {
			return TEllipsis, 3, nil
			// ... is ellipsis
		} else {
			// . is period
			return TPeriod, 1, nil
		}
	case '?':
		if s.seekMatchSequence([]rune{'?'}) {
			// ?? is nullish coalescing operator
			return TDoubleQuestionMark, 2, nil
		} else {
			// ? is question mark
			return TQuestionMark, 1, nil
		}
	case '!':
		if s.seekMatchSequence([]rune{'=', '='}) {
			// !== is strict not equal
			return TStrictNotEqual, 3, nil
		} else if s.seekMatchSequence([]rune{'='}) {
			// != is not equal
			return TNotEqual, 2, nil
		} else {
			// ! is logical not
			return TBang, 1, nil
		}
	case '=':
		if s.seekMatchSequence([]rune{'=', '='}) {
			// === is strict equal
			return TStrictEqual, 3, nil
		} else if s.seekMatchSequence([]rune{'>'}) {
			// => is arrow
			return TArrow, 2, nil
		} else if s.seekMatchSequence([]rune{'='}) {
			// == is equal
			return TEqual, 2, nil
		} else {
			// = is assign
			return TAssign, 1, nil
		}
	case '&':
		if s.seekMatchSequence([]rune{'&', '='}) {
			// &&= is logical and operator and assign
			return TLogicalAndAssign, 3, nil
		} else if s.seekMatchSequence([]rune{'&'}) {
			// && is logical and operator
			return TLogicalAnd, 2, nil
		} else if s.seekMatchSequence([]rune{'='}) {
			// &= is bitwise and operator and assign
			return TAndAssign, 2, nil
		} else {
			// & is bitwise and operator
			return TAnd, 1, nil
		}
	case '|':
		if s.seekMatchSequence([]rune{'|', '='}) {
			// ||= is logical or operator and assign
			return TLogicalOrAssign, 3, nil
		} else if s.seekMatchSequence([]rune{'|'}) {
			// || is logical or operator
			return TLogicalOr, 2, nil
		} else if s.seekMatchSequence([]rune{'='}) {
			// |= is logical bitwise or operator and assign
			return TOrAssign, 2, nil
		} else {
			// | is logical bitwise or operator
			return TOr, 1, nil
		}
	case '+':
		if s.seekMatchSequence([]rune{'+'}) {
			// ++ is increment unary operator
			return TPlusPlus, 2, nil
		} else if s.seekMatchSequence([]rune{'='}) {
			// += is increment and assign
			return TPlusAssign, 2, nil
		} else {
			// + is plus binary operator
			return TPlus, 1, nil
		}
	case '-':
		if s.seekMatchSequence([]rune{'-'}) {
			// -- is decrement unary operator
			return TMinusMinus, 2, nil
		} else if s.seekMatchSequence([]rune{'='}) {
			// -= is decrement and assign
			return TMinusAssign, 2, nil
		} else {
			// - is minus binary operator
			return TMinus, 1, nil
		}
	case '*':
		if s.seekMatchSequence([]rune{'*'}) {
			// ** is exponential operator
			return TStarStar, 2, nil
		} else if s.seekMatchSequence([]rune{'='}) {
			// *= is multiply and assign
			return TStarAssign, 2, nil
		} else {
			// * is multiply binary operator
			return TStar, 1, nil
		}
	case '/':
		if s.seekMatchSequence([]rune{'='}) {
			// /= is divide and assign
			return TSlashAssign, 2, nil
		} else {
			// / is divide binary operator
			return TSlash, 1, nil
		}
	default:
		return TUnknown, -1, fmt.Errorf("unknown punctuation: %s", string(s.peek()))
	}
}
