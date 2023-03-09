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

func (p *Parser) advanceBy(n int) {
	p.cursor += n
}

func (p *Parser) peekN(n int) (lex.Token, error) {
	if p.cursor+n >= len(p.seq) {
		return lex.Token{}, errors.New("EOF")
	}
	return p.seq[p.cursor+n], nil
}

func (p *Parser) matchAny(types ...lex.TokenType) bool {
	for _, t := range types {
		if p.peek().T == t {
			p.advanceBy(1)
			return true
		}
	}
	return false
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
	p.logger.Debug("PARSER ::")
	for i, token := range p.seq {
		p.logger.Debug("T(%d) :: %v", i, token)
	}
	p.logger.Debug("\n")

	defer func() {
		stack := recover()
		if stack != nil {
			p.logger.EmitStdout()
			panic(stack)
		}
	}()

	lastToken := lex.Token{}
	for {
		cursor := 0
		token, err := p.peekN(cursor)
		if lastToken == token {
			p.logger.Debug("[%d:0] LOOP: bailing to prevent infinite loop", p.cursor)
			p.logger.EmitStdout()
			break
		} else {
			lastToken = token
		}

		p.logger.Debug("[%d:0] LOOP: %v", p.cursor, token)
		if err != nil {
			p.logger.Error("[%d:0] LOOP: %s", p.cursor, err.Error())
			break
		}
		if token.T == lex.TSemicolon {
			// TODO: Ignoring semicolons for now
			p.logger.Debug("[%d:0] LOOP:skipping (;)", p.cursor)
		}

		// [
		if token.T == lex.TLeftBracket {
			var c int
			p.logger.Debug("c:%d (%d)", c, &c)
			node, err := p.parseArray(&c)
			p.logger.Debug("c:%d (%d)", c, &c)
			p.logger.Debug("[%d:%d] parser:root:endArray: %s", p.cursor, c, node.PrettyPrint())
			if err == nil {
				rootNode.children = append(rootNode.children, node)
				p.cursor += c // accept the cursor
				p.logger.Debug("[%d:%d] parser:root:pushToken: %s", p.cursor, c, node.PrettyPrint())
				continue
			} else {
				p.logger.Debug("[%d:%d] parser:array ERR: %s", p.cursor, c, err)
			}
		}

		// primary expressions
		// id | await | yield | literals
		if isLiteralToken(token) || token.T == lex.TIdentifier || token.T == lex.TAwait || token.T == lex.TYield || token.T == lex.TUndefined {
			node, err := p.parsePrimaryExpr(&cursor)
			if err == nil {
				rootNode.children = append(rootNode.children, node)
				p.cursor += cursor // accept the cursor
				p.logger.Debug("[%d:%d] parser:root:pushToken: %s", p.cursor, cursor, node.PrettyPrint())
				continue
			}
		}

		// unary operator expression
		if isUnaryOperator(token) {
			node, err := p.parseUnaryOperator(&cursor)
			if err == nil {
				rootNode.children = append(rootNode.children, node)
				p.logger.Debug("[%d:%d] parser:root:pushToken: %s", p.cursor, cursor, node.PrettyPrint())
				p.cursor += cursor // accept the cursor
				continue
			}
		}
		// update expression
		if isUpdateExpression(token) {
			node, err := p.parseUpdateExpr(&cursor)
			if err == nil {
				rootNode.children = append(rootNode.children, node)
				p.logger.Debug("[%d:%d] parser:root:pushToken: %s", p.cursor, cursor, node.PrettyPrint())
				p.cursor += cursor // accept the cursor
				continue
			}
		}

		p.cursor++
	}
	p.logger.Debug("parsed %s", rootNode)
	return rootNode
}
