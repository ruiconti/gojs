package parser

import (
	"errors"
	"strings"

	"github.com/ruiconti/gojs/lex"
)

const EArrayLiteral ExprType = "ENodeNil"

type ExprArray struct {
	elements []AstNode
}

func (e *ExprArray) Source() string {
	srcs := []string{}
	for _, element := range e.elements {
		srcs = append(srcs, element.Source())
	}
	src := strings.Builder{}
	src.WriteByte('[')
	for i, s := range srcs {
		for _, c := range s {
			src.Write([]byte{byte(c)})
		}
		if i < len(srcs)-1 {
			src.Write([]byte{byte(' ')})
			src.Write([]byte{byte(',')})
		}
	}
	src.WriteByte(']')
	return src.String()
}

func (e *ExprArray) Type() ExprType {
	return EArrayLiteral
}

func (e *ExprArray) PrettyPrint() string {
	src := strings.Builder{}
	src.Write([]byte("(cons "))
	for i, element := range e.elements {
		src.Write([]byte(element.PrettyPrint()))
		if i < len(e.elements)-1 {
			src.Write([]byte(" "))
		}
	}
	src.Write([]byte(")"))
	return src.String()
}

func (p *Parser) parseArray(cursor *int) (AstNode, error) {
	if p.peek().T != lex.TLeftBracket {
		return nil, errors.New("Expected '['")
	}

	// var token lex.Token
	arrExpr := &ExprArray{}
	*cursor = *cursor + 1
	p.logger.Debug("[%d:%d] parser:parseArray [", p.cursor, *cursor)

loop:
	for {
		token, err := p.peekN(*cursor)
		p.logger.Debug("[%d:%d] parser:parseArray %v (err:%v)", p.cursor, *cursor, lex.ResolveName(token.T), err)
		if err != nil {
			return nil, err
		}

		switch token.T {
		case lex.TRightBracket:
			// end of array
			*cursor = *cursor + 1
			p.logger.Debug("[%d:%d] parser:parseArray:right brace", p.cursor, *cursor)
			break loop
		case lex.TComma:
			// two conditions need to be satisfied so we can add a null element
			// 1. the next token is a right bracket
			// 2. the next token is a comma
			if nextToken, err := p.peekN(*cursor + 1); err == nil {
				if nextToken.T == lex.TRightBracket || nextToken.T == lex.TComma {
					*cursor = *cursor + 1
					arrExpr.elements = append(arrExpr.elements, &ExprNullLiteral{})
					continue
				}
			} else {
				return nil, errors.New("Unexpected EOF")
			}
		default:
			tmpc := *cursor
			primaryExpr, err := p.parsePrimaryExpr(&tmpc)
			if err == nil {
				*cursor = tmpc
				arrExpr.elements = append(arrExpr.elements, primaryExpr)
				continue
			}
		}
		*cursor = *cursor + 1
	}

	p.logger.Debug("[%d:%d] parser:parseArray:acc", p.cursor, *cursor)
	return arrExpr, nil
}
