package parser

import (
	"fmt"
	"strings"

	l "github.com/ruiconti/gojs/lexer"
)

const EArrayLiteral ExprType = "EExprArray"

type ExprArray struct {
	elements []Expr
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

// ArrayLiteral :
// | '[' Elision? ']'
// | '[' ElementList ']'
// | '[' ElementList ',' Elision? ']'
//
// simplifying
// ArrayLiteral :
// '[' (ElementList? ','?)* ']'
//
// ElementList :
// | (Elision? (AssignmentExpression | SpreadElement))*
func (p *Parser) parseArrayInitializer() (Expr, error) {
	var exprArray ExprArray
	if p.Peek().Type == l.TLeftBracket {
		p.Next() // consume '['

	loop:
		for {
			switch token := p.Peek(); token.Type {
			case l.TEOF:
				break loop
			case l.TRightBracket:
				p.Next() // consume ']'
				break loop
			case l.TComma:
				if p.PeekN(-1).Type == l.TComma || p.PeekN(-1).Type == l.TLeftBracket {
					// need to add only one element because we move to the second comma
					// and don't consume it
					// [ ,null, ]
					//   ^    ^ leave cursor at this point
					//   |
					//   | consumed on this iteration
					//
					exprArray.elements = append(exprArray.elements, ExprLitNull)
				}
				p.Next() // consume ','
			case l.TEllipsis:
				p.Next() // consume '...'
				arg, err := p.parseAssignExpr()
				if err != nil {
					return nil, err
				}
				exprArray.elements = append(exprArray.elements, &SpreadElement{argument: arg})

			default:
				exprAssign, err := p.parseAssignExpr()
				if err != nil {
					return nil, err
				}
				exprArray.elements = append(exprArray.elements, exprAssign)
			}
		}
		return &exprArray, nil
	}
	return nil, fmt.Errorf("rejected on parseArrayInitializer")
}
