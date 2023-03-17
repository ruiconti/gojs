package parser

import (
	"errors"
	"strings"

	"github.com/ruiconti/gojs/lex"
)

const EArrayLiteral ExprType = "ENodeNil"

type ExprArray struct {
	elements []Node
}

func (e *ExprArray) Type() ExprType {
	return EArrayLiteral
}

func (e *ExprArray) S() string {
	src := strings.Builder{}
	src.Write([]byte("(cons "))
	for i, element := range e.elements {
		src.Write([]byte(element.S()))
		if i < len(e.elements)-1 {
			src.Write([]byte(" "))
		}
	}
	src.Write([]byte(")"))
	return src.String()
}

func (p *Parser) parseArray(cursor *int) (Node, error) {
	if p.peek().Type != lex.TLeftBracket {
		return nil, errors.New("Expected '['")
	}

	// var token lex.Token
	arrExpr := &ExprArray{}
	*cursor = *cursor + 1
	p.logger.Debug("[%d:%d] parser:parseArray [", p.cursor, *cursor)

loop:
	for {
		token, err := p.peekN(*cursor)
		p.logger.Debug("[%d:%d] parser:parseArray %v (err:%v)", p.cursor, *cursor, token.Type.S(), err)
		if err != nil {
			return nil, err
		}

		switch token.Type {
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
				if nextToken.Type == lex.TRightBracket || nextToken.Type == lex.TComma {
					*cursor = *cursor + 1
					arrExpr.elements = append(arrExpr.elements, &ExprLiteral[string]{value: "null", typ: lex.TNull})
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
