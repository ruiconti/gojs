package parser

import (
	"fmt"
	"strings"

	"github.com/ruiconti/gojs/internal"
	l "github.com/ruiconti/gojs/lexer"
)

var TokenEOF = l.Token{l.TEOF, "EOF", "EOF", 0, 0}
var TokenBOF = l.Token{l.TEOF, "BOF", "BOF", 0, 0}

// TODO: define clear boundary between Expression, Statement and Declaration
// through a clear type model

type Node interface {
	S() string
}

type NodeRoot struct {
	children []Node
}

func (n *NodeRoot) S() string {
	var src strings.Builder
	src.WriteString("(js ")
	for i, child := range n.children {
		src.WriteString(child.S())
		if i < len(n.children)-1 {
			src.WriteString(" ")
		}
	}
	src.WriteString(")")
	return src.String()
}

type Expr interface {
	S() string
}

type ExprStmt interface {
	S() string
}

type Parser struct {
	tokens      []l.Token // token slice
	checkpoints []uint32  // checkpoints for backtracking
	cursor      uint32    // current index of the token slice
	cursorOOB   bool      // whether cursor is out of bounds
	seqEnd      uint32    // last index of the token slice

	logger *internal.SimpleLogger
}

func NewParser(tokens []l.Token, logger *internal.SimpleLogger) *Parser {
	return &Parser{
		tokens:      tokens,
		cursor:      0,
		seqEnd:      uint32(len(tokens) - 1),
		checkpoints: make([]uint32, 0),
		cursorOOB:   false,
		logger:      logger,
	}
}

func (p *Parser) Peek() l.Token {
	return p.PeekN(0)
}

// look-ahead and look-behind
func (p *Parser) PeekN(n int32) l.Token {
	idx := int32(p.cursor) + n
	if idx > int32(p.seqEnd) {
		return TokenEOF
	} else if idx < 0 {
		return TokenBOF
	}

	return p.tokens[idx]
}

func (p *Parser) Next() {
	p.consume(1)
}

func (p *Parser) Backtrack() {
	p.consume(-1)
}

func (p *Parser) consume(offset int32) {
	width := int32(p.cursor) + offset
	if width == int32(p.seqEnd)+1 {
		// the last consume
		p.cursorOOB = true
	} else if width > int32(p.seqEnd)+1 || width < 0 {
		p.cursorOOB = true
		return
	}

	var consumed strings.Builder
	for i := int32(p.cursor); i < width; i++ {
		consumed.WriteString(p.tokens[i].String())
		if i < width-1 {
			consumed.WriteString(", ")
		}
	}
	p.Log("consuming %v", consumed.String())
	p.cursor = uint32(width)
}

func (p *Parser) Log(msg string, format ...interface{}) {
	fmsg := fmt.Sprintf(msg, format...)
	var logmsg string
	if p.cursor > p.seqEnd {
		logmsg = fmt.Sprintf("%d: (EOF) %s", p.cursor, fmsg)
	} else {
		// current := p.Peek()
		logmsg = fmt.Sprintf("%d: %s", p.cursor, fmsg)
	}

	p.logger.Debug(logmsg)
}

func (p *Parser) guardInfiniteLoop(lastCursor *uint32) {
	if p.cursor > 0 && p.cursor == *lastCursor {
		panic("infinite loop detected")
	} else {
		*lastCursor = p.cursor
	}
}

func (p *Parser) saveCheckpoint() uint32 {
	return p.cursor
}

func (p *Parser) restoreCheckpoint(cursor uint32) {
	if cursor > p.seqEnd {
		panic("invalid checkpoint: out-of-bounds")
	}
	p.cursor = cursor
}

func Parse(logger *internal.SimpleLogger, src string) Node {
	var (
		ast *NodeRoot
		err error
	)
	lexer := l.NewLexer(src, logger)
	tokens, errs := lexer.ScanAll()
	if len(errs) > 0 {
		for _, e := range errs {
			logger.Error(e.Error())
		}
		panic(1)
	}

	parser := NewParser(tokens, logger)
	defer func() {
		stack := recover()
		if stack != nil {
			parser.logger.Debug("AST :: %v", ast.S())
			parser.logger.DumpLogs()
			panic(stack)
		}
	}()

	parser.logger.Debug("PARSER ::")
	for i, token := range parser.tokens {
		parser.logger.Debug("T(%d) :: %v", i, token.String())
	}
	parser.logger.Debug("\n")
	ast, err = parser.parseProgram()
	if err != nil {
		parser.logger.Error(err.Error())
		panic(err)
	}
	parser.logger.Debug("AST :: %v", ast.S())
	return ast
}

func (p *Parser) parseProgram() (*NodeRoot, error) {
	var (
		statements    []Node
		lastCursorPos uint32 = 0
	)
	defer func() {
		stack := recover()
		if stack != nil {
			p.logger.DumpLogs()
			panic(stack)
		}
	}()

	for !p.cursorOOB {
		token := p.Peek()
		p.Log("loop %v", token.String())

		if token.Type == l.TSemicolon {
			p.Log("skipping ;") // TODO: Ignoring semicolons for now
			p.Next()
			continue
		}

		stmt, err := p.parseStatement()
		if err == nil {
			statements = append(statements, stmt)
		} else if stmt == nil {
			panic("boo")
		} else {
			return &NodeRoot{children: statements}, err
		}
		p.guardInfiniteLoop(&lastCursorPos)
	}
	return &NodeRoot{children: statements}, nil
}
