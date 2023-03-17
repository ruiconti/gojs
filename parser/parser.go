package parser

import (
	"errors"
	"fmt"

	"github.com/ruiconti/gojs/internal"
	"github.com/ruiconti/gojs/lex"
)

type ExprType string

const (
	ENodeRoot ExprType = "ENodeRoot"
	EElision  ExprType = "EElision"
)

type Node interface {
	S() string
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

func (p *Parser) isEOF() bool {
	return p.cursor >= len(p.seq)
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

func (p *Parser) log(cursor *int, msg string, format ...interface{}) {
	if cursor == nil || *cursor < 0 {
		*cursor = -1
	}
	fmsg := fmt.Sprintf(msg, format...)
	if p.isEOF() {
		p.logger.Debug("[%d:%d] (EOF) %s", p.cursor, *cursor, fmsg)
	} else {
		current := p.peek()
		tname := lex.ResolveName(current.T)
		p.logger.Debug("[%d:%d] [%v:%v] %s", p.cursor, *cursor, tname, current.Lexeme, fmsg)
	}
}

func (p *Parser) matchAny(types ...lex.TokenType) bool {
	if p.cursor >= len(p.seq) {
		return false
	}

	for _, t := range types {
		if p.peek().Type == t {

			p.logger.Debug("[%d] matched %s, stepping 1..", p.cursor, lex.ResolveName(t))
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
	lexer := lex.NewLexer(src, logger)
	tokens, err := lexer.ScanAll()
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
		children: []Node{},
	}
	p.logger.Debug("PARSER ::")
	for i, token := range p.seq {
		p.logger.Debug("T(%d) :: %v", i, token)
	}
	p.logger.Debug("\n")

	defer func() {
		stack := recover()
		if stack != nil {
			p.logger.DumpLogs()
			panic(stack)
		}
	}()

	lastToken := lex.Token{}
	for {
		cursor := 0
		token, err := p.peekN(cursor)
		if lastToken == token {
			p.logger.Debug("[%d:0] LOOP: bailing to prevent infinite loop", p.cursor)
			p.logger.DumpLogs()
			break
		} else {
			lastToken = token
		}

		p.logger.Debug("[%d:0] LOOP: %v", p.cursor, token)
		if err != nil {
			p.logger.Error("[%d:0] LOOP: %s", p.cursor, err.Error())
			break
		}
		if token.Type == lex.TSemicolon {
			// TODO: Ignoring semicolons for now
			p.logger.Debug("[%d:0] LOOP:skipping (;)", p.cursor)
		}

		// [
		if token.Type == lex.TLeftBracket {
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

		node, err := p.parseExpr(&cursor)
		if err == nil {
			rootNode.children = append(rootNode.children, node)
			p.logger.Debug("[%d:%d] parser:root ACC: %s", p.cursor, cursor, node.PrettyPrint())
		} else {
			p.logger.Debug("[%d:%d] parser:root ERR: %s", p.cursor, cursor, err)
		}

		// primary expressions
		// id | await | yield | literals
		if isLiteralToken(token) || token.Type == lex.TIdentifier || token.Type == lex.TAwait || token.Type == lex.TYield || token.Type == lex.TUndefined {
			node, err := p.parsePrimaryExpr(&cursor)
			if err == nil {
				rootNode.children = append(rootNode.children, node)
				p.cursor += cursor // accept the cursor
				p.logger.Debug("[%d:%d] parser:root:pushToken: %s", p.cursor, cursor, node.PrettyPrint())
				continue
			}
		}

		// unary operator expression
		// if isUnaryOperator(token) {
		// 	node, err := p.parseUnaryOperator(&cursor)
		// 	if err == nil {
		// 		rootNode.children = append(rootNode.children, node)
		// 		p.logger.Debug("[%d:%d] parser:root:pushToken: %s", p.cursor, cursor, node.PrettyPrint())
		// 		p.cursor += cursor // accept the cursor
		// 		continue
		// 	}
		// }
		// update expression
		// if isUpdateExpression(token) {
		// 	node, err := p.parseUpdateExpr(&cursor)
		// 	if err == nil {
		// 		rootNode.children = append(rootNode.children, node)
		// 		p.logger.Debug("[%d:%d] parser:root:pushToken: %s", p.cursor, cursor, node.PrettyPrint())
		// 		p.cursor += cursor // accept the cursor
		// 		continue
		// 	}
		// }

		p.cursor++
	}
	p.logger.Debug("parsed %s", rootNode)
	return rootNode
}
