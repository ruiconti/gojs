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

func (p *Parser) LookAhead(n int) (lex.Token, error) {
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

func Parse(src string) *ExprRootNode {
	logger := internal.NewSimpleLogger(internal.ModeDebug)
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
	p.logger.Debug("parser:tokens::")
	for _, token := range p.seq {
		p.logger.Debug("parser:token: %v", token)
	}
	for {
		token, err := p.LookAhead(0)
		if err != nil {
			p.logger.Error("parser:loop: %s", err.Error())
			break
		}
		// EOF
		if token.T == lex.TEOF {
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
			// TODO: [ could also be a property access
			// consumes '['
			p.logger.Debug("parsing [: %d: %v", p.cursor, token)
			tmpCursor := p.cursor
			tmpCursor++

			next, err := p.LookAhead(tmpCursor)
			if err != nil {
				break
			}
			arrNode := &ExprArray{}
			arrNode.elements = []AstNode{}
			// simple case: []
			// look for N consecutive elisions
			// if we have 0 elisions, we don't even start the loop
			for next.T == lex.TComma {
				node := &ExprElision{}
				arrNode.elements = append(arrNode.elements, node)
				tmpCursor++
				next, _ = p.LookAhead(tmpCursor)
			}
			if next.T == lex.TRightBracket {
				tmpCursor++
				p.cursor = tmpCursor
				rootNode.children = append(rootNode.children, arrNode)
			} else {
				p.logger.Debug("TODO: parse array")
				// bail
			}
		}
		// unary operator
		if isUnaryOperator(token) {
			tmpCursor := p.cursor // backtracks if we don't find a unary operator
			unaryOp, rej := p.parseUnaryOperator(nil)
			if rej != nil {
				rootNode.children = append(rootNode.children, unaryOp)
			} else {
				p.cursor = tmpCursor
			}
		}
	}
	p.logger.Debug("parsed %s", rootNode)
	return rootNode
}
