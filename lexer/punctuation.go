package lexer

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
func (s *Lexer) MatchSequence(chars ...rune) bool {
	if len(chars) == 0 {
		panic("incorrect usage: chars must not be empty")
	}

	for i, charExp := range chars {
		charGot := s.PeekN(uint(i + 1))
		if charGot == EOF || charGot != charExp {
			return false
		}
	}
	return true
}

func isPunctuation(r rune) bool {
	return r == '!' || r == '.' || r == ',' || r == '>' || r == '<' || r == '=' || r == '+' || r == '-' || r == '*' || r == '/' || r == '%' || r == '&' || r == '|' || r == '^' || r == '(' || r == ')' || r == '[' || r == ']' || r == '{' || r == '}' || r == ';' || r == ':' || r == '?' || r == '~' || r == '#'
}

// Punctuators
//
// Scan for punctuators is straightforward:
// we group tokens by their first character, and always try to match
// the longest possible token, iteratively until we find a match.
func (s *Lexer) scanPunctuation() Token {
	token, lexemec := s.matchPunctuation(s.Peek())
	lexeme := s.src[s.srcCursorHead : s.srcCursorHead+int(lexemec)]
	s.Jump(lexemec)
	return Token{
		Type:   token,
		Lexeme: lexeme,
	}
}

func (s *Lexer) matchPunctuation(char rune) (TokenType, uint) {
	switch char {
	case '}':
		return TRightBrace, 1
	case '{':
		return TLeftBrace, 1
	case '(':
		return TLeftParen, 1
	case ')':
		return TRightParen, 1
	case '[':
		return TLeftBracket, 1
	case ']':
		return TRightBracket, 1
	case ';':
		return TSemicolon, 1
	case ',':
		return TComma, 1
	case ':':
		return TColon, 1
	case '^':
		return TXor, 1
	case '~':
		return TTilde, 1
	case '%':
		return TPercent, 1
	case '#':
		return TNumberSign, 1
	case '>':
		if s.MatchSequence('>', '>', '=') {
			// >>>= is unsigned right shift assign
			return TUnsignedRightShiftAssign, 4
		} else if s.MatchSequence('>', '>') {
			// >>> is unsigned right shift
			return TUnsignedRightShift, 3
		} else if s.MatchSequence('>', '=') {
			// >>= is right shift assign
			return TRightShiftAssign, 3
		} else if s.MatchSequence('=') {
			// >= is greater equal
			return TGreaterThanEqual, 2
		} else if s.MatchSequence('>') {
			// >> is right shift
			return TRightShift, 2
		} else {
			// > is greater than
			return TGreaterThan, 1
		}
	case '<':
		if s.MatchSequence('=') {
			// <= is greater equal
			return TLessThanEqual, 2
		} else if s.MatchSequence('<', '=') {
			// <<= is left shift assign
			return TLeftShiftAssign, 3
		} else if s.MatchSequence('<') {
			// << is left shift
			return TLeftShift, 2
		} else {
			// < is greater than
			return TLessThan, 1
		}
	case '.':
		if s.MatchSequence('.', '.') {
			return TEllipsis, 3
			// ... is ellipsis
		} else {
			// . is period
			return TPeriod, 1
		}
	case '?':
		if s.MatchSequence('?') {
			// ?? is nullish coalescing operator
			return TDoubleQuestionMark, 2
		} else {
			// ? is question mark
			return TQuestionMark, 1
		}
	case '!':
		if s.MatchSequence('=', '=') {
			// !== is strict not equal
			return TStrictNotEqual, 3
		} else if s.MatchSequence('=') {
			// != is not equal
			return TNotEqual, 2
		} else {
			// ! is logical not
			return TBang, 1
		}
	case '=':
		if s.MatchSequence('=', '=') {
			// === is strict equal
			return TStrictEqual, 3
		} else if s.MatchSequence('>') {
			// => is arrow
			return TArrow, 2
		} else if s.MatchSequence('=') {
			// == is equal
			return TEqual, 2
		} else {
			// = is assign
			return TAssign, 1
		}
	case '&':
		if s.MatchSequence('&', '=') {
			// &&= is logical and operator and assign
			return TLogicalAndAssign, 3
		} else if s.MatchSequence('&') {
			// && is logical and operator
			return TLogicalAnd, 2
		} else if s.MatchSequence('=') {
			// &= is bitwise and operator and assign
			return TAndAssign, 2
		} else {
			// & is bitwise and operator
			return TAnd, 1
		}
	case '|':
		if s.MatchSequence('|', '=') {
			// ||= is logical or operator and assign
			return TLogicalOrAssign, 3
		} else if s.MatchSequence('|') {
			// || is logical or operator
			return TLogicalOr, 2
		} else if s.MatchSequence('=') {
			// |= is logical bitwise or operator and assign
			return TOrAssign, 2
		} else {
			// | is logical bitwise or operator
			return TOr, 1
		}
	case '+':
		if s.MatchSequence('+') {
			// ++ is increment unary operator
			return TPlusPlus, 2
		} else if s.MatchSequence('=') {
			// += is increment and assign
			return TPlusAssign, 2
		} else {
			// + is plus binary operator
			return TPlus, 1
		}
	case '-':
		if s.MatchSequence('-') {
			// -- is decrement unary operator
			return TMinusMinus, 2
		} else if s.MatchSequence('=') {
			// -= is decrement and assign
			return TMinusAssign, 2
		} else {
			// - is minus binary operator
			return TMinus, 1
		}
	case '*':
		if s.MatchSequence('*') {
			// ** is exponential operator
			return TStarStar, 2
		} else if s.MatchSequence('=') {
			// *= is multiply and assign
			return TStarAssign, 2
		} else {
			// * is multiply binary operator
			return TStar, 1
		}
	case '/':
		if s.MatchSequence('=') {
			// /= is divide and assign
			return TSlashAssign, 2
		} else {
			// / is divide binary operator
			return TSlash, 1
		}
	default:
		s.Errorf("Unknown punctuator: %s", string(char))
		return TUnknown, 0
	}
}
