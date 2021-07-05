package ast

import (
	l "aura/src/lexer"
	"fmt"
	"strings"
)

// Represents and AST node
type ASTNode interface {
	TokenLiteral() string // return the token literal of the node
	Str() string          // return  a string representation of the node
}

// represents a statement
type Stmt interface {
	ASTNode    // ensure all statements implements the ASTnode Interface
	stmtNode() // method to distinguish statements and expressions
}

// represents a expression
type Expression interface {
	ASTNode       // ensure all statements implements the ASTnode Interface
	expressNode() // method to distinguish statements and expressions
}

// BaseNode is a struct wich all the expressions and statements
// whill inherit providing the default TokenLiteral implementation
type BaseNode struct {
	Token l.Token // represents the token of the node
}

// return the token literal of the node
func (b BaseNode) TokenLiteral() string {
	return b.Token.Literal
}

// Program represents all the program
type Program struct {
	Staments []Stmt // represents all the statements in the program
}

// generates a new program instance
func NewProgram(statements []Stmt) *Program {
	return &Program{Staments: statements}
}

func (p Program) TokenLiteral() string {
	if len(p.Staments) > 0 {
		return p.Staments[0].TokenLiteral()
	}

	return ""
}

func (p Program) Str() string {
	var buf strings.Builder
	for idx, v := range p.Staments {
		if idx >= len(p.Staments)-1 {
			buf.WriteString(v.Str())
		} else {
			buf.WriteString(v.Str() + " ")
		}
	}

	return buf.String()
}

// Represents a variable or function declaration
type LetStatement struct {
	BaseNode             // represent the token of the statement
	Name     *Identifier // represents the name of the variable
	Value    Expression  // represents the values assing to the variable
}

// generate a new let stament instance
func NewLetStatement(token l.Token, name *Identifier, value Expression) *LetStatement {
	return &LetStatement{
		BaseNode: BaseNode{token},
		Name:     name,
		Value:    value,
	}
}

func (l LetStatement) stmtNode() {}

func (l LetStatement) Str() string {
	return fmt.Sprintf("%s %s = %s;", l.TokenLiteral(), l.Name.Str(), l.Value.Str())
}

// Represents a return statement
type ReturnStament struct {
	BaseNode               // represents the token
	ReturnValue Expression // represents the value to be returned
}

// generates a new return statement instance
func NewReturnStatement(token l.Token, returnValue Expression) *ReturnStament {
	return &ReturnStament{BaseNode: BaseNode{token}, ReturnValue: returnValue}
}

func (r ReturnStament) stmtNode() {}

func (r ReturnStament) Str() string {
	return fmt.Sprintf("%s %s;", r.TokenLiteral(), r.ReturnValue.Str())
}

// handle expressions statements
type ExpressionStament struct {
	BaseNode
	Expression Expression
}

// generates a new expression statement instance
func NewExpressionStament(token l.Token, expression Expression) *ExpressionStament {
	return &ExpressionStament{BaseNode: BaseNode{token}, Expression: expression}
}

func (e ExpressionStament) stmtNode() {}
func (e ExpressionStament) Str() string {
	return e.Expression.Str()
}

// Suffix representrs a suffix expression
type Suffix struct {
	BaseNode            // Extends base node struct
	Left     Expression // represents the object that will be apply the suffix expression
	Operator string     // represents the operator to be apply to the object
}

// generates a new suffix instance
func NewSuffix(token l.Token, left Expression, operator string) *Suffix {
	return &Suffix{BaseNode: BaseNode{token}, Left: left, Operator: operator}
}

func (s *Suffix) expressNode() {}

func (s *Suffix) Str() string {
	return fmt.Sprintf("%s%s", s.Left.Str(), s.Operator)
}

// Represents a block of code delimited by curly braces
type Block struct {
	BaseNode        // Extends base node struct
	Staments []Stmt // represents all the statements inside the block
}

// generates a new block instance
func NewBlock(token l.Token, staments ...Stmt) *Block {
	return &Block{BaseNode: BaseNode{token}, Staments: staments}
}

func (b Block) stmtNode() {}

func (b Block) Str() string {
	var buf strings.Builder
	for idx, stament := range b.Staments {
		if idx >= len(b.Staments)-1 {
			buf.WriteString(stament.Str())
		} else {
			buf.WriteString(stament.Str() + " ")
		}
	}

	return buf.String()
}

// IF represents an If expression
type If struct {
	BaseNode               // Extends base node struct
	Condition   Expression // represents the condition of the expression
	Consequence *Block     // represents the consequence if the condition is trythy
	Alternative *Block     // represents the alternative if the condition is not trythy
}

// generates a new if instance
func NewIf(token l.Token, condition Expression, consequence, alternative *Block) *If {
	return &If{
		BaseNode:    BaseNode{token},
		Condition:   condition,
		Consequence: consequence,
		Alternative: alternative,
	}
}

func (i If) expressNode() {}

func (i If) Str() string {
	var out strings.Builder
	out.WriteString(fmt.Sprintf("si %s %s ", i.Condition.Str(), i.Consequence.Str()))
	if i.Alternative != nil {
		out.WriteString(fmt.Sprintf("si_no %s", i.Alternative.Str()))
	}

	return out.String()
}

// Represents a function declaration
type Function struct {
	BaseNode                 // Extends base node struct
	Parameters []*Identifier // represents the parameters of the function
	Body       *Block        // represents the function body
}

// create a new function instance
func NewFunction(token l.Token, body *Block, parameters ...*Identifier) *Function {
	return &Function{
		BaseNode:   BaseNode{token},
		Parameters: parameters,
		Body:       body,
	}
}

func (f Function) expressNode() {}

func (f Function) Str() string {
	var buf strings.Builder
	for idx, parameter := range f.Parameters {
		if idx == len(f.Parameters)-1 {
			buf.WriteString(parameter.Str())
		} else {
			buf.WriteString(parameter.Str() + " ")
		}
	}

	return fmt.Sprintf("%s(%s) %s", f.TokenLiteral(), buf.String(), f.Body.Str())
}

// represents a function call
type Call struct {
	BaseNode               // extends base node struct
	Function  Expression   // represents the function to be call
	Arguments []Expression // represents the arguments given to call the function
}

// generates a new Call instance
func NewCall(token l.Token, function Expression, arguments ...Expression) *Call {
	return &Call{
		BaseNode:  BaseNode{token},
		Function:  function,
		Arguments: arguments,
	}
}

func (C Call) expressNode() {}

func (c Call) Str() string {
	var buf strings.Builder
	for idx, arg := range c.Arguments {
		if idx == len(c.Arguments)-1 {
			buf.WriteString(arg.Str())
		} else {
			buf.WriteString(arg.Str() + ", ")
		}
	}

	return fmt.Sprintf("%s(%s)", c.Function.Str(), buf.String())
}

// Represents a for expression
type For struct {
	BaseNode             // Extends base node struct
	Condition Expression // represents the iterable expression
	Body      *Block     // represents the body of the forloop
}

// generates a new For instance
func NewFor(token l.Token, condition Expression, body *Block) *For {
	return &For{BaseNode: BaseNode{token}, Condition: condition, Body: body}
}

func (f *For) expressNode() {}

func (f *For) Str() string {
	return fmt.Sprintf("%s %s { %s }", f.TokenLiteral(), f.Condition.Str(), f.Body.Str())
}

// Represents a WhileLoop expression
type While struct {
	BaseNode             // Extends base node struct
	Condition Expression // represents the condition of the while loop
	Body      *Block     // represents the body of the while loop
}

// generates a new whileloop instance
func NewWhile(token l.Token, cond Expression, body *Block) *While {
	return &While{BaseNode: BaseNode{token}, Condition: cond, Body: body}
}

func (w *While) expressNode() {}

func (w *While) Str() string {
	return fmt.Sprintf("%s %s { %s }", w.TokenLiteral(), w.Condition.Str(), w.Body.Str())
}

// represents an Array Expression
type Array struct {
	BaseNode              // Extends base node struct
	Values   []Expression // represents the values inside the array
}

// generates a new array instance
func NewArray(token l.Token, values ...Expression) *Array {
	return &Array{BaseNode: BaseNode{token}, Values: values}
}

func (a *Array) expressNode() {}

func (a *Array) Str() string {
	var buf strings.Builder
	for idx, val := range a.Values {
		if idx == len(a.Values)-1 {
			buf.WriteString(val.Str())
		} else {
			buf.WriteString(val.Str() + ", ")
		}
	}

	return buf.String()
}

// represents a call to a data structure like maps, arrays or strings
type CallList struct {
	BaseNode             // Extends base node struct
	ListIdent Expression // represents the data structure to be call
	Index     Expression // represents where is the values in the data structure
}

// generates a new CallList instance
func NewCallList(token l.Token, listIdent Expression, index Expression) *CallList {
	return &CallList{
		BaseNode:  BaseNode{token},
		ListIdent: listIdent,
		Index:     index,
	}
}

func (c *CallList) expressNode() {}
func (c *CallList) Str() string {
	return fmt.Sprintf("%s[%s]", c.ListIdent.Str(), c.Index.Str())
}

// Represents a HashMap expression
type MapExpression struct {
	BaseNode             // Extends base node struct
	Body     []*KeyValue // represents all the key values pairs in the HashMap
}

// generates a new MapExpression instance
func NewMapExpression(token l.Token, body []*KeyValue) *MapExpression {
	return &MapExpression{BaseNode{token}, body}
}
func (m *MapExpression) expressNode() {}

func (m *MapExpression) Str() string {
	var buf strings.Builder
	for idx, keyVal := range m.Body {
		if idx >= len(m.Body)-1 {
			buf.WriteString(keyVal.Str())
		} else {
			buf.WriteString(keyVal.Str() + ", ")
		}
	}

	return fmt.Sprintf("mapa{%s}", buf.String())
}

// Represents a class statement
type ClassStatement struct {
	BaseNode                   // extends base node
	Name     *Identifier       // represents the class name
	Params   []*Identifier     // represents the constructor params
	Methods  []*ClassMethodExp // represents all the methods in the class
}

// generates a new class statement instance
func NewClassStatement(token l.Token, name *Identifier, params []*Identifier, methods []*ClassMethodExp) *ClassStatement {
	return &ClassStatement{
		BaseNode: BaseNode{token},
		Params:   params,
		Methods:  methods,
	}
}

func (c *ClassStatement) stmtNode() {}

func (c *ClassStatement) Str() string {
	var paramsBuf strings.Builder
	for idx, param := range c.Params {
		if idx == len(c.Params)-1 {
			paramsBuf.WriteString(param.Str())
		} else {
			paramsBuf.WriteString(param.Str() + ", ")
		}
	}
	var buf strings.Builder
	for idx, method := range c.Methods {
		if idx == len(c.Methods)-1 {
			buf.WriteString(method.Str())
		} else {
			buf.WriteString(method.Str() + " ")
		}
	}

	return fmt.Sprintf(
		"clase %s(%s) {\n %s \n }",
		c.Name.Value,
		paramsBuf.String(),
		buf.String(),
	)
}
