package parser

import "strings"

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

type EEllision struct{}
