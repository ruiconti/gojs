package parser

import (
	"bytes"
	"fmt"
	"strings"

	l "github.com/ruiconti/gojs/lexer"
)

type StmtType string

type Stmt interface {
	S() string
}

// Statement[Yield, Await, Return] :
// | BlockStatement[?Yield, ?Await, ?Return]
// | VariableStatement[?Yield, ?Await]
// | EmptyStatement
// | ExpressionStatement[?Yield, ?Await]
// | IfStatement[?Yield, ?Await, ?Return]
// | BreakableStatement[?Yield, ?Await, ?Return]
// | ContinueStatement[?Yield, ?Await]
// | BreakStatement[?Yield, ?Await]
// | [+Return] ReturnStatement[?Yield, ?Await]
// | WithStatement[?Yield, ?Await, ?Return]
// | LabelledStatement[?Yield, ?Await, ?Return]
// | ThrowStatement[?Yield, ?Await]
// | TryStatement[?Yield, ?Await, ?Return]
// | DebuggerStatement
func (p *Parser) parseStatement() (Stmt, error) {
	token := p.Peek()
	var stmt Stmt
	var err error

	cp := p.saveCheckpoint()
	switch token.Type {
	// todo: resolve conflict with object initialization
	case l.TLeftBrace:
		stmt, err = p.parseBlockStatement()
	case l.TVar, l.TConst, l.TLet:
		stmt, err = p.parseVariableStatement()
	case l.TSemicolon:
		stmt, err = p.parseEmptyStatement()
	case l.TIf:
		stmt, err = p.parseIfStatement()
	}

	if err != nil || stmt == nil {
		p.restoreCheckpoint(cp)
		return p.parseExpressionStatement()
	}
	return stmt, err
}

// EmptyStatement : ';'
type EmptyStatement struct{}

const SStmt StmtType = "SStmt"

func (s *EmptyStatement) Type() StmtType { return SStmt }
func (s *EmptyStatement) S() string      { return ";" }

func (p *Parser) parseEmptyStatement() (*EmptyStatement, error) {
	if p.Peek().Type != l.TSemicolon {
		return nil, fmt.Errorf("expected ';', got %v", p.Peek().Type)
	}
	p.Next() // Consume the ';' token
	return &EmptyStatement{}, nil
}

// IfStatement[Yield, Await, Return] :
// 'if' '(' Expression[+In, ?Yield, ?Await] ')' Statement[?Yield, ?Await, ?Return] 'else' Statement[?Yield, ?Await, ?Return]
// 'if' '(' Expression[+In, ?Yield, ?Await] ')' Statement[?Yield, ?Await, ?Return] [lookahead ≠ else]
type IfStatement struct {
	Condition Expr
	ThenStmt  Stmt
	ElseStmt  Stmt
}

func (s *IfStatement) Type() StmtType {
	return SStmt
}

func (s *IfStatement) S() string {
	if s.ElseStmt == nil {
		return fmt.Sprintf("(if %s %s)", s.Condition.S(), s.ThenStmt.S())
	} else {
		return fmt.Sprintf("(if %s %s %s)", s.Condition.S(), s.ThenStmt.S(), s.ElseStmt.S())
	}
}

func (p *Parser) parseIfStatement() (*IfStatement, error) {
	if p.Peek().Type != l.TIf {
		return nil, fmt.Errorf("expected 'if' keyword, got %v", p.Peek().Type)
	}
	p.Next() // consume 'if'
	if p.Peek().Type != l.TLeftParen {
		return nil, fmt.Errorf("expected '(' after 'if' keyword, got %v", p.Peek().Type)
	}
	p.Next() // consume '('
	condition, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	if p.Peek().Type != l.TRightParen {
		return nil, fmt.Errorf("expected ')' after expression in 'if' statement, got %v", p.Peek().Type)
	}
	p.Next() // consume ')'
	thenStmt, err := p.parseStatement()
	if err != nil {
		return nil, err
	}
	var elseStmt Stmt
	if p.Peek().Type == l.TElse {
		p.Next() // Consume the 'else' token
		elseStmt, err = p.parseStatement()
		if err != nil {
			return nil, err
		}
	}
	return &IfStatement{Condition: condition, ThenStmt: thenStmt, ElseStmt: elseStmt}, nil
}

// BlockStatement[Yield, Await, Return] :
// | Block[?Yield, ?Await, ?Return]
//
// Block[Yield, Await, Return] :
// | '{' StatementList[?Yield, ?Await, ?Return]? '}'
//
// StatementList[Yield, Await, Return] :
// | StatementListItem[?Yield, ?Await, ?Return]
// | StatementList[?Yield, ?Await, ?Return] StatementListItem[?Yield, ?Await, ?Return]
//
// StatementListItem[Yield, Await, Return] :
// | Statement[?Yield, ?Await, ?Return]
// | Declaration[?Yield, ?Await]
type BlockStatement struct {
	Stmts []Stmt
}

func (s *BlockStatement) S() string {
	src := strings.Builder{}
	src.Write([]byte("(block "))
	for i, stmt := range s.Stmts {
		src.Write([]byte(stmt.S()))
		if i < len(s.Stmts)-1 {
			src.Write([]byte("\n"))
		}
	}
	src.Write([]byte(")"))
	return src.String()
}
func (p *Parser) parseBlockStatement() (Stmt, error) {
	p.Next() // Consume the '{' token
	var stmtList []Stmt
	for p.Peek().Type != l.TRightBrace {
		stmt, err := p.parseStatement()
		if err != nil {
			return nil, err
		}
		stmtList = append(stmtList, stmt)
	}
	p.Next() // Consume the '}' token
	return &BlockStatement{Stmts: stmtList}, nil
}

const StmtExpression = "StmtExpression"

func (s *ExpressionStatement) S() string {
	return s.expression.S()
}

type ExpressionStatement struct {
	expression Expr
}

// ExpressionStatement[Yield, Await] :
// | Expression[+In, ?Yield, ?Await] ';'
func (p *Parser) parseExpressionStatement() (*ExpressionStatement, error) {
	expr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}

	if p.Peek().Type == l.TSemicolon {
		p.Next() // consume ';'
	}

	return &ExpressionStatement{expression: expr}, nil
}

// BreakableStatement[Yield, Await, Return] :
// | IterationStatement[?Yield, ?Await, ?Return]
// | SwitchStatement[?Yield, ?Await, ?Return]

// IterationStatement[Yield, Await, Return] :
// | 'do' Statement 'while' '(' Expression ')' ';'
// | 'while' '(' Expression[+In, ?Yield, ?Await] ')' Statement[?Yield, ?Await, ?Return]
// | ForStatement[?Yield, ?Await, ?Return]
// | ForInOfStatement[?Yield, ?Await, ?Return]

// ForStatement[Yield, Await, Return] :
// | 'for' '(' [lookahead ≠ let [] Expression[~In, ?Yield, ?Await]? ';' Expression[+In, ?Yield, ?Await]? ';' Expression[+In, ?Yield, ?Await]? ')' Statement[?Yield, ?Await, ?Return]
// | 'for' '(' var VariableDeclarationList[~In, ?Yield, ?Await] ';' Expression[+In, ?Yield, ?Await]? ';' Expression[+In, ?Yield, ?Await]? ')' Statement[?Yield, ?Await, ?Return]
// | 'for' '(' LexicalDeclaration[~In, ?Yield, ?Await] Expression[+In, ?Yield, ?Await]? ';' Expression[+In, ?Yield, ?Await]? ')' Statement[?Yield, ?Await, ?Return]

// ForInOfStatement[Yield, Await, Return] :
// | 'for' '(' [lookahead ≠ 'let' [] LeftHandSideExpression[?Yield, ?Await] 'in' Expression[+In, ?Yield, ?Await] ')' Statement[?Yield, ?Await, ?Return]
// | 'for' '(' 'var' ForBinding[?Yield, ?Await] in Expression[+In, ?Yield, ?Await] ')' Statement[?Yield, ?Await, ?Return]
// | 'for' '(' ForDeclaration[?Yield, ?Await] in Expression[+In, ?Yield, ?Await] ')' Statement[?Yield, ?Await, ?Return]
// | 'for' '(' [lookahead ∉ { 'let', 'async' 'of' }] LeftHandSideExpression[?Yield, ?Await] 'of' AssignmentExpression[+In, ?Yield, ?Await] ')' Statement[?Yield, ?Await, ?Return]
// | 'for' '(' var ForBinding[?Yield, ?Await] 'of' AssignmentExpression[+In, ?Yield, ?Await] ')' Statement[?Yield, ?Await, ?Return]
// | 'for' '(' ForDeclaration[?Yield, ?Await] 'of' AssignmentExpression[+In, ?Yield, ?Await] ')' Statement[?Yield, ?Await, ?Return]
// | [+Await] for await ( [lookahead ≠ let] LeftHandSideExpression[?Yield, ?Await] 'of' AssignmentExpression[+In, ?Yield, ?Await] ')' Statement[?Yield, ?Await, ?Return]
// | [+Await] for await ( var ForBinding[?Yield, ?Await] 'of' AssignmentExpression[+In, ?Yield, ?Await] ')' Statement[?Yield, ?Await, ?Return]
// | [+Await] for await ( ForDeclaration[?Yield, ?Await] 'of' AssignmentExpression[+In, ?Yield, ?Await] ')' Statement[?Yield, ?Await, ?Return]
//
// ForDeclaration[Yield, Await] :
// | LetOrConst ForBinding[?Yield, ?Await]
// ForBinding[Yield, Await] :
// | BindingIdentifier[?Yield, ?Await]
// | BindingPattern[?Yield, ?Await]

// SwitchStatement[Yield, Await, Return] :
// 'switch' '(' Expression[+In, ?Yield, ?Await] ) CaseBlock[?Yield, ?Await, ?Return]
//
// CaseBlock[Yield, Await, Return] :
// | '{' CaseClauses[?Yield, ?Await, ?Return]? '}'
// | '{' CaseClauses[?Yield, ?Await, ?Return]? DefaultClause[?Yield, ?Await, ?Return] CaseClauses[?Yield, ?Await, ?Return]? '}'
//
// CaseClauses[Yield, Await, Return] :
// | CaseClause[?Yield, ?Await, ?Return]
// | CaseClauses[?Yield, ?Await, ?Return] CaseClause[?Yield, ?Await, ?Return]
//
// CaseClause[Yield, Await, Return] :
// | 'case' Expression[+In, ?Yield, ?Await] ':' StatementList[?Yield, ?Await, ?Return]?
//
// DefaultClause[Yield, Await, Return] :
// | 'default' ':' StatementList[?Yield, ?Await, ?Return]?

// Declaration[Yield, Await] :
// | HoistableDeclaration[?Yield, ?Await, ~Default]
// | ClassDeclaration[?Yield, ?Await, ~Default]
// | LexicalDeclaration[+In, ?Yield, ?Await]
//
// HoistableDeclaration[Yield, Await, Default] :
// | FunctionDeclaration[?Yield, ?Await, ?Default]
// | GeneratorDeclaration[?Yield, ?Await, ?Default]
// | AsyncFunctionDeclaration[?Yield, ?Await, ?Default]
// | AsyncGeneratorDeclaration[?Yield, ?Await, ?Default]

// LexicalDeclaration[In, Yield, Await] :
// | ('let' | 'const') BindingList ';'
// BindingList[In, Yield, Await] :
// | LexicalBinding
// | BindingList ',' LexicalBinding
// LexicalBinding[In, Yield, Await] :
// | BindingIdentifier Initializer?
// | BindingPattern Initializer

//
// BreakableStatement[Yield, Await, Return] :
// | IterationStatement[?Yield, ?Await, ?Return]
// | SwitchStatement[?Yield, ?Await, ?Return]

// VariableStatement[Yield, Await] :
// | 'var' VariableDeclarationList[+In, ?Yield, ?Await] ';'
//
// VariableDeclarationList[In, Yield, Await] :
// | VariableDeclaration[?In, ?Yield, ?Await]
// | VariableDeclarationList[?In, ?Yield, ?Await] ',' VariableDeclaration[?In, ?Yield, ?Await]
//
// VariableDeclaration[In, Yield, Await] :
// | BindingIdentifier[?Yield, ?Await] Initializer[?In, ?Yield, ?Await]opt
// | BindingPattern[?Yield, ?Await] Initializer[?In, ?Yield, ?Await]
type VariableDeclaration struct {
	identifier *ExprIdentifier // todo: more accurate naming
	init       Expr
	pattern    Expr
}

func (s *VariableDeclaration) S() string {
	var lhs, rhs string
	if s.identifier != nil {
		lhs = s.identifier.name
	} else {
		lhs = s.pattern.S()
	}

	if s.init != nil {
		rhs = fmt.Sprintf(" <- %s", s.init.S())
	} else {
		rhs = ""
	}
	return fmt.Sprintf("(%s%s)", lhs, rhs)
}

type VariableStatement struct {
	kind         l.Token
	declarations []*VariableDeclaration
}

func (s *VariableStatement) S() string {
	var buf bytes.Buffer

	buf.WriteString(fmt.Sprintf("(%s ", s.kind.Type.S()))
	for i, decl := range s.declarations {
		buf.WriteString(decl.S())
		if i < len(s.declarations)-1 {
			buf.WriteString(" ")
		}
	}
	buf.WriteString(")")

	return buf.String()
}

// Two productions:
//
// VariableStatement[Yield, Await] :
// | 'var' VariableDeclarationList[+In, ?Yield, ?Await] ';'
//
// LexicalDeclaration[In, Yield, Await] :
// | ('let' | 'const') BindingList ';'
// BindingList[In, Yield, Await] :
// | LexicalBinding
// | BindingList ',' LexicalBinding
// LexicalBinding[In, Yield, Await] :
// | BindingIdentifier Initializer?
// | BindingPattern Initializer

func (p *Parser) parseVariableStatement() (*VariableStatement, error) {
	kind := p.Peek()
	if kind.Type != l.TVar && kind.Type != l.TConst && kind.Type != l.TLet {
		return nil, fmt.Errorf("unexpected token %s, expected 'var' | 'const' | 'let'", kind.String())
	}

	p.Next() // consume 'var' | 'const' | 'let'
	varDeclList, err := p.parseVariableDeclarationList()
	if err != nil {
		return nil, err
	}

	if p.Peek().Type == l.TSemicolon {
		p.Next() // consume ';'
	}

	return &VariableStatement{declarations: varDeclList, kind: kind}, nil
}

// VariableDeclarationList[In, Yield, Await] :
// | VariableDeclaration[?In, ?Yield, ?Await]
// | VariableDeclarationList[?In, ?Yield, ?Await] ',' VariableDeclaration[?In, ?Yield, ?Await]
func (p *Parser) parseVariableDeclarationList() ([]*VariableDeclaration, error) {
	var declarations []*VariableDeclaration

	for {
		decl, err := p.parseVariableDeclaration()
		if err != nil {
			return nil, err
		}
		declarations = append(declarations, decl)

		if p.Peek().Type != l.TComma {
			break
		}
		p.Next() // consume ','
	}

	return declarations, nil
}

// VariableDeclaration[In, Yield, Await] :
// | BindingIdentifier[?Yield, ?Await] Initializer[?In, ?Yield, ?Await]?
// | BindingPattern[?Yield, ?Await] Initializer[?In, ?Yield, ?Await]
func (p *Parser) parseVariableDeclaration() (*VariableDeclaration, error) {
	var (
		identifier *ExprIdentifier
		pattern    Expr
		init       Expr
		err        error
	)

	token := p.Peek()

	if token.Type == l.TIdentifier {
		identifier = &ExprIdentifier{name: token.Lexeme}
		p.Next() // consume identifier
	} else {
		pattern, err = p.parseBindingPattern()
		if err != nil {
			return nil, err
		}
	}

	// Initializer[In, Yield, Await] : '=' AssignmentExpression[?In, ?Yield, ?Await]
	if p.Peek().Type == l.TAssign {
		p.Next() // consume '='
		initExpr, err := p.parseAssignExpr()
		if err != nil {
			return nil, err
		}
		init = initExpr
	}

	if identifier == nil && pattern == nil {
		return nil, fmt.Errorf("expected identifier or pattern, got %s", token.String())
	}

	if pattern != nil && init == nil {
		return nil, fmt.Errorf("expected initializer for pattern, got %s", token.String())
	}

	return &VariableDeclaration{identifier: identifier, pattern: pattern, init: init}, nil
}

// BindingIdentifier[Yield, Await] :
// | Identifier
// | yield
// | await
//
// BindingPattern[Yield, Await] :
// | ObjectBindingPattern[?Yield, ?Await]
// | ArrayBindingPattern[?Yield, ?Await]
//
// ObjectBindingPattern[Yield, Await] :
// | '{' '}'
// | '{' BindingRestProperty[?Yield, ?Await] '}'
// | '{' BindingPropertyList[?Yield, ?Await] '}'
// | '{' BindingPropertyList[?Yield, ?Await] , BindingRestProperty[?Yield, ?Await]? '}'
//
// ArrayBindingPattern[Yield, Await] :
// | '[' Elision? BindingRestElement[?Yield, ?Await]? ']'
// | '[' BindingElementList[?Yield, ?Await] ']'
// | '[' BindingElementList[?Yield, ?Await] ',' Elision? BindingRestElement[?Yield, ?Await]? ']'
//
// BindingRestProperty[Yield, Await] :
// | '...' BindingIdentifier[?Yield, ?Await]
//
// BindingPropertyList[Yield, Await] :
// | BindingProperty[?Yield, ?Await]
// | BindingPropertyList[?Yield, ?Await] ',' BindingProperty[?Yield, ?Await]
//
// BindingElementList[Yield, Await] :
// | BindingElisionElement[?Yield, ?Await]
// | BindingElementList[?Yield, ?Await] ',' BindingElisionElement[?Yield, ?Await]
//
// BindingElisionElement[Yield, Await] :
// | Elision? BindingElement[?Yield, ?Await]
//
// BindingProperty[Yield, Await] :
// | SingleNameBinding[?Yield, ?Await]
//
// PropertyName[?Yield, ?Await] : BindingElement[?Yield, ?Await]
//
// BindingElement[Yield, Await] :
// | SingleNameBinding[?Yield, ?Await]
// | BindingPattern[?Yield, ?Await] Initializer[+In, ?Yield, ?Await?]
//
// SingleNameBinding[Yield, Await] :
// | BindingIdentifier[?Yield, ?Await] Initializer[+In, ?Yield, ?Await?]
//
// BindingRestElement[Yield, Await] :
// | '...' BindingIdentifier[?Yield, ?Await]
// | '...' BindingPattern[?Yield, ?Await]
//
// Initializer[In, Yield, Await] : '=' AssignmentExpression[?In, ?Yield, ?Await]
//
// BindingPattern[Yield, Await] :
// | ObjectBindingPattern[?Yield, ?Await]
// | ArrayBindingPattern[?Yield, ?Await]
func (p *Parser) parseBindingPattern() (Expr, error) {
	switch p.Peek().Type {
	// TODO: there are legal syntax in object init that are not legal in binding pattern
	// e.g {a: 1, b: 2}, { a(), async b() } is legal in object init but not in binding pattern
	case l.TLeftBrace:
		return p.parseObjectInitializer()
	case l.TLeftBracket:
		return p.parseArrayInitializer()
	default:
		return nil, fmt.Errorf("expected an object or array binding pattern")
	}
}
