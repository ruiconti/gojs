package main

import (
	"fmt"
	"log"
)

func isPunctuation(r rune) bool {
	return r == '!' || r == '.' || r == ',' || r == '>' || r == '<' || r == '=' || r == '+' || r == '-' || r == '*' || r == '/' || r == '%' || r == '&' || r == '|' || r == '^' || r == '~' || r == '(' || r == ')' || r == '[' || r == ']' || r == '{' || r == '}' || r == ';' || r == ':' || r == '?' || r == ' '
}

func Scan(src string) []Token {
	headIdx := 0
	tokens := []Token{}

	advance := func() {
		headIdx++
	}
	// Check whether the sequence is found next, if it is, advance headIdx e.g.
	// >> src := "()===!";
	// >> headIdx := 2;
	// >> head := '=';
	// >> candidates := []rune{'=', '='}
	// >> seekMatchSequence(candidates)
	// true
	// >> fmt.Println(headIdx)
	// 4
	log.Printf("new")
	seekMatchSequence := func(candidates []rune) bool {
		if len(candidates) == 0 {
			panic("candidates must not be empty")
		}

		i, j := headIdx+1, 0
		for j < len(candidates) {
			if i > len(src)-1 || j > len(candidates)-1 {
				// out of bounds
				return false
			}
			log.Printf("(i:%d;j:%d) head: %c seekSeq: (eq %c %c)", i, j, rune(src[headIdx]), rune(src[i]), candidates[j])
			if rune(src[i]) != candidates[j] {
				log.Printf("(i:%d;j:%d) head: %c seekSeq: false", i, j, rune(src[headIdx]))
				return false
			}
			i++
			j++
		}
		// we got here, everything is the same :)
		// runesstr := ``
		// for _, r := range candidates {
		// 	runesstr += fmt.Sprintf("%c", r)
		// }
		// log.Printf("%d: matchSequence true! %v", headIdx, runesstr)
		headIdx += len(candidates)
		log.Printf("(i:%d;j:%d) head: %c seekSeq: true", i, j, rune(src[headIdx]))
		return true
	}

	add := func(t Token) {
		log.Printf("%d: added: %v", headIdx, TokenMap[t])
		tokens = append(tokens, t)
	}

	for headIdx < len(src) {
		head := src[headIdx]

		if isPunctuation(rune(head)) {
			switch head {
			// Simple punctuators
			case ' ':
				log.Printf("%d: whitespace", headIdx)
				// nothing
			case '}':
				add(TRightBrace)
			case '{':
				add(TLeftBrace)
			case '(':
				add(TLeftParen)
			case ')':
				add(TRightParen)
			case '[':
				add(TLeftBracket)
			case ']':
				add(TRightBracket)
			case ';':
				add(TSemicolon)
			case ':':
				add(TColon)
			case '>':
				if seekMatchSequence([]rune{'='}) {
					// >= is greater equal
					add(TGreaterThanEqual)
				} else if seekMatchSequence([]rune{'>', '>', '='}) {
					// >>>= is unsigned right shift assign
					add(TUnsignedRightShiftAssign)
				} else if seekMatchSequence([]rune{'>', '>'}) {
					// >>> is unsigned right shift
					add(TUnsignedRightShift)
				} else if seekMatchSequence([]rune{'>', '='}) {
					// >>= is right shift assign
					add(TRightShiftAssign)
				} else if seekMatchSequence([]rune{'>'}) {
					// >> is right shift
					add(TRightShift)
				} else {
					// > is greater than
					add(TGreaterThan)
				}
			case '<':
				if seekMatchSequence([]rune{'='}) {
					// <= is greater equal
					add(TLessThanEqual)
				} else if seekMatchSequence([]rune{'<', '='}) {
					// <<= is left shift assign
					add(TLeftShiftAssign)
				} else if seekMatchSequence([]rune{'<'}) {
					// << is left shift
					add(TLeftShift)
				} else {
					// < is greater than
					add(TLessThan)
				}
			case '.':
				if seekMatchSequence([]rune{'.', '.'}) {
					add(TEllipsis)
				} else {
					add(TPeriod)
				}
			case '?':
				if seekMatchSequence([]rune{'?'}) {
					add(TDoubleQuestionMark)
				} else {
					add(TQuestionMark)
				}
			case '!':
				if seekMatchSequence([]rune{'=', '='}) {
					add(TStrictNotEqual)
				} else if seekMatchSequence([]rune{'='}) {
					add(TNotEqual)
				} else {
					add(TBang)
				}
			case '=':
				if seekMatchSequence([]rune{'=', '='}) {
					add(TStrictEqual)
				} else if seekMatchSequence([]rune{'>'}) {
					add(TArrow)
				} else if seekMatchSequence([]rune{'='}) {
					add(TEqual)
				} else {
					add(TAssign)
				}
			case '&':
				if seekMatchSequence([]rune{'&', '='}) {
					add(TLogicalAndAssign)
				} else if seekMatchSequence([]rune{'&'}) {
					add(TLogicalAnd)
				} else if seekMatchSequence([]rune{'='}) {
					add(TAndAssign)
				} else {
					add(TAnd)
				}
			case '|':
				if seekMatchSequence([]rune{'|', '='}) {
					add(TLogicalOrAssign)
				} else if seekMatchSequence([]rune{'|'}) {
					add(TLogicalOr)
				} else if seekMatchSequence([]rune{'='}) {
					add(TOrAssign)
				} else {
					add(TOr)
				}
			case '+':
				if seekMatchSequence([]rune{'+'}) {
					add(TPlusPlus)
				} else if seekMatchSequence([]rune{'='}) {
					add(TPlusAssign)
				} else {
					add(TPlus)
				}
			case '-':
				if seekMatchSequence([]rune{'-'}) {
					add(TMinusMinus)
				} else if seekMatchSequence([]rune{'='}) {
					add(TMinusAssign)
				} else {
					add(TMinus)
				}
			case '*':
				if seekMatchSequence([]rune{'='}) {
					add(TStarAssign)
				} else {
					add(TStar)
				}
			case '/':
				if seekMatchSequence([]rune{'='}) {
					add(TSlashAssign)
				} else {
					add(TSlash)
				}
			}
		}
		advance()
	}
	return tokens
}

func main() {
	runes := []rune{'!', '.', ',', '>', '<', '=', '+', '-', '*', '/', '%', '&', '|', '^', '~', '(', ')', '[', ']', '{', '}', ';', ':', '?', ' ', '\t', '\r', '\n'}
	for _, r := range runes {
		fmt.Printf("%s - %d\n", string(r), r)
	}
}
