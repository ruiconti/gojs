package parser

import (
	"errors"

	"github.com/ruiconti/gojs/internal"
	"github.com/ruiconti/gojs/lex"
)

type ExprType string

const (
	ENodeRoot ExprType = "ENodeRoot"
	EElision  ExprType = "EElision"
)

type AstNode interface {
	Source() string
	Type() ExprType
	PrettyPrint() string
}

type Parser struct {
	seq    []lex.Token
	cursor int
	logger *internal.SimpleLogger
}

type ExprElision struct{}

func (e *ExprElision) Type() ExprType {
	return EElision
}

func (e *ExprElision) Source() string {
	return ""
}

func (p *Parser) peek() lex.Token {
	return p.seq[p.cursor]
}

func (p *Parser) peekN(n int) (lex.Token, error) {
	if p.cursor+n >= len(p.seq) {
		return lex.Token{}, errors.New("EOF")
	}
	return p.seq[p.cursor+n], nil
}

func NewParser(seq []lex.Token, logger *internal.SimpleLogger) *Parser {
	return &Parser{
		seq:    seq,
		cursor: 0,
		logger: logger,
	}
}

func Parse(logger *internal.SimpleLogger, src string) *ExprRootNode {
	scanner := lex.NewScanner(src, logger)
	tokens, err := scanner.Scan()
	if err != nil {
		logger.Error("Error scanning source: %s", err.Error())
		panic(1)
	}
	parser := NewParser(tokens, logger)
	ast := parser.parseTokens()
	return ast
}

func (p *Parser) parseTokens() *ExprRootNode {
	rootNode := &ExprRootNode{
		children: []AstNode{},
	}
	p.logger.Debug("PARSER::")
	defer func() {
		stack := recover()
		if stack != nil {
			p.logger.EmitStdout()
			panic(stack)
		}
	}()

	for {
		token, err := p.peekN(p.cursor)
		p.logger.Debug("[%d] parser:loop: %v", p.cursor, token)
		if err != nil {
			p.logger.Error("[%d] parser:loop: %s", p.cursor, err.Error())
			break
		}
		// id | await | yield
		if token.T == lex.TIdentifier || token.T == lex.TAwait || token.T == lex.TYield {
			p.cursor++
			node := &ExprIdentifierReference{reference: token.Lexeme}
			rootNode.children = append(rootNode.children, node)
		}
		// ;
		if token.T == lex.TSemicolon {
			// TODO: Ignoring semicolons for now
			p.cursor++
		}
		// [
		if token.T == lex.TLeftBracket {
			// TODO: implement
			p.cursor++
		}
		// unary operator expression
		if isUnaryOperator(token) {
			cursor := 0
			unaryOp, err := p.parseUnaryOperator(&cursor)
			if err == nil {
				rootNode.children = append(rootNode.children, unaryOp)
				p.logger.Debug("[%d:%d] parser:root:pushToken: %s", p.cursor, cursor, unaryOp.PrettyPrint())
				p.cursor = p.cursor + cursor
			}
		}
		// update expression
		if isUpdateExpression(token) {
			cursor := 0
			updateExpr, err := p.parseUpdateExpr(&cursor)
			if err == nil {
				rootNode.children = append(rootNode.children, updateExpr)
				p.logger.Debug("[%d:%d] parser:root:pushToken: %s", p.cursor, cursor, updateExpr.PrettyPrint())
				p.cursor = p.cursor + cursor
			}
		}

		p.cursor++
	}
	p.logger.Debug("parsed %s", rootNode)
	return rootNode
}
