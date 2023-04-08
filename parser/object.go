package parser

import (
	"fmt"
	"strings"

	l "github.com/ruiconti/gojs/lexer"
)

// /////////////////////
// PropertyDefinition //
// /////////////////////
const EPropertyDefinition ExprType = "EPropertyDefinition"

type PropertyDefinition struct {
	key       Expr
	value     Expr
	computed  bool // { [foo]: 1 }
	method    bool // { foo() {} }
	shorthand bool // { foo }
}

func (p *PropertyDefinition) Type() ExprType { return EPropertyDefinition }
func (p *PropertyDefinition) S() string {
	return fmt.Sprintf("(k:%s v:%s (%v %v %v))", p.key.S(), p.value.S(), p.computed, p.method, p.shorthand)
}

// /////////////
// ExprObject //
// /////////////
const EObjectInitialization ExprType = "EObjectInitialization"

type ExprObject struct {
	properties []*PropertyDefinition
}

func (e *ExprObject) Type() ExprType {
	return EObjectInitialization
}

func (e *ExprObject) S() string {
	src := strings.Builder{}
	src.Write([]byte("(dict "))
	for i, prop := range e.properties {
		src.Write([]byte(prop.S()))
		if i < len(e.properties)-1 {
			src.Write([]byte(" "))
		}
	}
	src.Write([]byte(")"))
	return src.String()
}

// ObjectLiteral :
// '{' '}'
// '{' PropertyDefinitionList ','? '}'

// PropertyDefinitionList :
// PropertyDefinition (',' PropertyDefinition)* ','?
//
// PropertyDefinition :
// | IdentifierReference
// | IdentifierReference '=' AssignmentExpression
// | PropertyName ':' AssignmentExpression
// | MethodDefinition (TODO)
// | '...' AssignmentExpression
// | CoverInitializedName (TODO)
//
// PropertyName :
// | LiteralPropertyName
// | ComputedPropertyName
//
// LiteralPropertyName :
// | IdentifierName
// | StringLiteral
// | NumericLiteral
//
// ComputedPropertyName :
// '[' AssignmentExpression ']'
//
// CoverInitializedName : IdentifierReference '=' AssignmentExpression
func (p *Parser) parseObjectInitializer() (Expr, error) {
	var exprObject ExprObject
	if p.Peek().Type == l.TLeftBrace {
		p.Next() // consume '{'

	loop:
		for {
			switch token := p.Peek(); token.Type {
			case l.TEOF:
				break loop
			case l.TRightBrace:
				p.Next() // consume '}'
				return &exprObject, nil
			case l.TComma:
				p.Next() // consume ','
			default:
				propDef, err := p.parsePropertyDefinition()
				if err != nil {
					return nil, err
				}
				exprObject.properties = append(exprObject.properties, propDef)
			}
		}
		return &exprObject, nil
	}
	return nil, fmt.Errorf("rejected on parseObjectInitializer")
}

// PropertyDefinition :
// | (('yield' | 'await')? Identifier) ('=' AssignmentExpression)?
// | (Identifier | StringLiteral| NumericLiteral | ComputedPropertyName) ':' AssignmentExpression
// | MethodDefinition
// | '...' AssignmentExpression
func (p *Parser) parsePropertyDefinition() (*PropertyDefinition, error) {
	var err error
	propName, computed, err := p.parsePropertyName()
	if err != nil {
		// PropertyDefinition : '...' AssignmentExpression
		if p.Peek().Type == l.TEllipsis {
			p.Next() // consume '...'
			expr, err := p.parseAssignExpr()
			if err != nil {
				return nil, err
			}
			// TODO: reduce this to a single Expr
			return &PropertyDefinition{key: expr, value: &SpreadElement{argument: expr}}, nil
		}
		return nil, err
	}

	// continuation of Identifier
	switch token := p.Peek(); token.Type {
	case l.TEqual:
		// PropertyDefinition : (('yield' | 'await')? Identifier) ('=' AssignmentExpression)?
		p.Next() // consume '='
		expr, err := p.parseAssignExpr()
		if err != nil {
			return nil, err
		}
		return &PropertyDefinition{key: &ExprIdentifier{name: token.Lexeme}, value: expr, shorthand: true}, nil
	case l.TColon:
		// PropertyDefinition : (Identifier | StringLiteral| NumericLiteral | ComputedPropertyName) ':' AssignmentExpression
		p.Next() // consume ':'
		expr, err := p.parseAssignExpr()
		if err != nil {
			return nil, err
		}
		return &PropertyDefinition{key: propName, value: expr, computed: computed}, nil
	case l.TRightBrace, l.TComma:
		// PropertyDefinition : Identifier ('}' | ', )
		if computed {
			// invalid syntax
			return nil, fmt.Errorf("can't use a computed property name in a shorthand fashion")
		}
		// we do not consume the token here, because it will be consumed by the caller
		return &PropertyDefinition{key: propName, value: propName, computed: false, shorthand: true}, nil
	}

	// TODO: Implement MethodDefinition

	return nil, fmt.Errorf("rejected on parsePropertyDefinition")
}

func (p *Parser) parsePropertyName() (Expr, bool /* computed */, error) {
	token := p.Peek()
	switch token.Type {
	case l.TIdentifier:
		p.Next() // consume identifier
		return &ExprIdentifier{name: token.Lexeme}, false, nil
	case l.TStringLiteral_DoubleQuote, l.TStringLiteral_SingleQuote:
		p.Next() // consume string
		return &ExprLiteral[string]{tok: token}, false, nil
	case l.TNumericLiteral:
		p.Next() // consume numeric
		return &ExprLiteral[float64]{tok: token}, false, nil
	case l.TLeftBracket:
		p.Next() // consume '['
		expr, err := p.parseAssignExpr()
		if err != nil {
			return nil, false, err
		}
		if p.Peek().Type == l.TRightBracket {
			p.Next() // consume ']'
			return expr, true, nil
		}
		return nil, false, fmt.Errorf("expected ']' after ComputedPropertyName")
	}
	return nil, false, fmt.Errorf("rejected on parsePropertyName")
}

// func (p *Parser) parseMethodDefinition() (*ExprFunction, error) {
// 	// Check for async keyword
// 	isAsync := false
// 	if p.Peek().Type == l.TAsync {
// 		p.Next() // consume 'async'
// 		isAsync = true
// 	}

// 	// Check for generator
// 	isGenerator := false
// 	if p.Peek().Type == l.TStar {
// 		p.Next() // consume '*'
// 		isGenerator = true
// 	}

// 	// Parse property name
// 	propName, _, err := p.parsePropertyName()
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Parse function parameters and body
// 	params, body, err := p.parseFunctionParametersAndBody()
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &MethodDefinition{
// 		key:       propName,
// 		params:    params,
// 		body:      body,
// 		async:     isAsync,
// 		generator: isGenerator,
// 	}, nil
// }

// MethodDefinition[Yield, Await] :
// | ClassElementName[?Yield, ?Await] ( UniqueFormalParameters[~Yield, ~Await] ) { FunctionBody[~Yield, ~Await] }
// | GeneratorMethod[?Yield, ?Await]
// | AsyncMethod[?Yield, ?Await]
// | AsyncGeneratorMethod[?Yield, ?Await]
// | get ClassElementName[?Yield, ?Await] ( ) { FunctionBody[~Yield, ~Await] }
// | set ClassElementName[?Yield, ?Await] ( PropertySetParameterList ) { FunctionBody[~Yield, ~Await] }
// | PropertySetParameterList :
// | FormalParameter[~Yield, ~Await]

// ClassElement[Yield, Await] :
// | MethodDefinition[?Yield, ?Await]
// | static MethodDefinition[?Yield, ?Await]
// | FieldDefinition[?Yield, ?Await] ;
// | static FieldDefinition[?Yield, ?Await] ;
// | ClassStaticBlock
// | ';'
//
// FieldDefinition[Yield, Await] :
// | ClassElementName[?Yield, ?Await] Initializer[+In, ?Yield, ?Await]opt
//
// ClassElementName[Yield, Await] :
// | PropertyName[?Yield, ?Await]
// | PrivateIdentifier
//
// ClassStaticBlock :
// | static { ClassStaticBlockBody }
//
// ClassStaticBlockBody :
// | ClassStaticBlockStatementList
//
// ClassStaticBlockStatementList :
// | StatementList[~Yield, +Await, ~Return]op

// UniqueFormalParameters[Yield, Await] :
// | FormalParameters[?Yield, ?Await]
//
// FormalParameters[Yield, Await] :
// | [empty]
// | FunctionRestParameter[?Yield, ?Await]
// | FormalParameterList[?Yield, ?Await]
// | FormalParameterList[?Yield, ?Await] ,
// | FormalParameterList[?Yield, ?Await] , FunctionRestParameter[?Yield, ?Await]
//
// FormalParameterList[Yield, Await] :
// | FormalParameter[?Yield, ?Await]
// | FormalParameterList[?Yield, ?Await] , FormalParameter[?Yield, ?Await]
//
// FunctionRestParameter[Yield, Await] :
// | BindingRestElement[?Yield, ?Await]
//
// FormalParameter[Yield, Await] :
// | BindingElement[?Yield, ?Await]
//
// BindingRestProperty[Yield, Await] :
// | ... BindingIdentifier[?Yield, ?Await]
//
// BindingPropertyList[Yield, Await] :
// | BindingProperty[?Yield, ?Await]
// | BindingPropertyList[?Yield, ?Await] , BindingProperty[?Yield, ?Await]
//
// BindingElementList[Yield, Await] :
// | BindingElisionElement[?Yield, ?Await]
// | BindingElementList[?Yield, ?Await] , BindingElisionElement[?Yield, ?Await]
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
// | BindingPattern[?Yield, ?Await] Initializer[+In, ?Yield, ?Await]opt
//
// SingleNameBinding[Yield, Await] :
// | BindingIdentifier[?Yield, ?Await] Initializer[+In, ?Yield, ?Await]opt
//
// BindingRestElement[Yield, Await] :
// | ... BindingIdentifier[?Yield, ?Await]
// | ... BindingPattern[?Yield, ?Await]
