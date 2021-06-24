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
	stmtNode() // method to distinguish statements and expression
}

// represents a expression
type Expression interface {
	ASTNode       // ensure all statements implements the ASTnode Interface
	expressNode() // method to distinguish statements and expression
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
	var out = make([]string, 0, len(p.Staments))
	for _, v := range p.Staments {
		out = append(out, v.Str())
	}

	return strings.Join(out, " ")
}

// Represents a variable or function declaration
type LetStatement struct {
	Token l.Token     // represent the token of the statement
	Name  *Identifier // represents the name of the variable
	Value Expression  // represents the values assing to the variable
}

// generate a new let stament instance
func NewLetStatement(token l.Token, name *Identifier, value Expression) *LetStatement {
	return &LetStatement{
		Token: token,
		Name:  name,
		Value: value,
	}
}

func (l LetStatement) TokenLiteral() string {
	return l.Token.Literal
}

func (l LetStatement) stmtNode() {}

func (l LetStatement) Str() string {
	return fmt.Sprintf("%s %s = %s;", l.TokenLiteral(), l.Name.Str(), l.Value.Str())
}

// Represents a return statement
type ReturnStament struct {
	Token       l.Token    // represents the token
	ReturnValue Expression // represents the value to be returned
}

// generates a new return statement instance
func NewReturnStatement(token l.Token, returnValue Expression) *ReturnStament {
	return &ReturnStament{Token: token, ReturnValue: returnValue}
}

func (r ReturnStament) TokenLiteral() string {
	return r.Token.Literal
}

func (r ReturnStament) stmtNode() {}

func (r ReturnStament) Str() string {
	return fmt.Sprintf("%s %s;", r.TokenLiteral(), r.ReturnValue.Str())
}

// handle expressions statements
type ExpressionStament struct {
	token      l.Token
	Expression Expression
}

// generates a new expression statement instance
func NewExpressionStament(token l.Token, expression Expression) *ExpressionStament {
	return &ExpressionStament{token: token, Expression: expression}
}

func (e ExpressionStament) TokenLiteral() string {
	return e.token.Literal
}

func (e ExpressionStament) stmtNode() {}
func (e ExpressionStament) Str() string {
	return e.Expression.Str()
}

// Suffix representrs a suffix expression
type Suffix struct {
	Token    l.Token    // represents the token of the expression
	Left     Expression // represents the object that will be apply the suffix expression
	Operator string     // represents the operator to be apply to the object
}

// generates a new suffix instance
func NewSuffix(token l.Token, left Expression, operator string) *Suffix {
	return &Suffix{
		Token:    token,
		Left:     left,
		Operator: operator,
	}
}

func (s *Suffix) TokenLiteral() string {
	return s.Token.Literal
}

func (s *Suffix) expressNode() {}

func (s *Suffix) Str() string {
	return fmt.Sprintf("%s%s", s.Left.Str(), s.Operator)
}

// Represents a block of code delimited by curly braces
type Block struct {
	Token    l.Token // represents the token of the expression
	Staments []Stmt  // represents all the statements inside the block
}

// generates a new block instance
func NewBlock(token l.Token, staments ...Stmt) *Block {
	return &Block{
		Token:    token,
		Staments: staments,
	}
}

func (b Block) TokenLiteral() string {
	return b.Token.Literal
}

func (b Block) stmtNode() {}

func (b Block) Str() string {
	var out = make([]string, 0, len(b.Staments))
	for _, stament := range b.Staments {
		out = append(out, stament.Str())
	}

	return strings.Join(out, " ")
}

// IF represents an If expression
type If struct {
	Token       l.Token    // represents the token of the expression
	Condition   Expression // represents the condition of the expression
	Consequence *Block     // represents the consequence if the condition is trythy
	Alternative *Block     // represents the alternative if the condition is not trythy
}

// generates a new if instance
func NewIf(token l.Token, condition Expression, consequence, alternative *Block) *If {
	return &If{
		Token:       token,
		Condition:   condition,
		Consequence: consequence,
		Alternative: alternative,
	}
}

func (i If) TokenLiteral() string {
	return i.Token.Literal
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
	Token      l.Token       // represents the token of the expression
	Parameters []*Identifier // represents the parameters of the function
	Body       *Block        // represents the function body
}

// create a new function instance
func NewFunction(token l.Token, body *Block, parameters ...*Identifier) *Function {
	return &Function{
		Token:      token,
		Parameters: parameters,
		Body:       body,
	}
}

func (f Function) TokenLiteral() string {
	return f.Token.Literal
}

func (f Function) expressNode() {}

func (f Function) Str() string {
	var paramList = make([]string, 0, len(f.Parameters))
	for _, parameter := range f.Parameters {
		paramList = append(paramList, parameter.Str())
	}

	params := strings.Join(paramList, ", ")
	return fmt.Sprintf("%s(%s) %s", f.TokenLiteral(), params, f.Body.Str())
}

// represents a function call
type Call struct {
	Token     l.Token      // represents the token of the expresion
	Function  Expression   // represents the function to be call
	Arguments []Expression // represents the arguments given to call the function
}

// generates a new Call instance
func NewCall(token l.Token, function Expression, arguments ...Expression) *Call {
	return &Call{
		Token:     token,
		Function:  function,
		Arguments: arguments,
	}
}

func (c Call) TokenLiteral() string {
	return c.Token.Literal
}

func (C Call) expressNode() {}

func (c Call) Str() string {
	var argsList = make([]string, 0, len(c.Arguments))
	for _, arg := range c.Arguments {
		argsList = append(argsList, arg.Str())
	}

	args := strings.Join(argsList, ", ")
	return fmt.Sprintf("%s(%s)", c.Function.Str(), args)
}

// Represents a for expression
type For struct {
	Token     l.Token    // represents the token of the expression
	Condition Expression // represents the iterable expression
	Body      *Block     // represents the body of the forloop
}

// generates a new For instance
func NewFor(token l.Token, condition Expression, body *Block) *For {
	return &For{Token: token, Condition: condition, Body: body}
}

func (f *For) TokenLiteral() string {
	return f.Token.Literal
}

func (f *For) expressNode() {}

func (f *For) Str() string {
	return fmt.Sprintf("%s %s { %s }", f.TokenLiteral(), f.Condition.Str(), f.Body.Str())
}

// Represents a WhileLoop expression
type While struct {
	Token     l.Token    // represents the token of the expression
	Condition Expression // represents the condition of the while loop
	Body      *Block     // represents the body of the while loop
}

// generates a new whileloop instance
func NewWhile(token l.Token, cond Expression, body *Block) *While {
	return &While{Token: token, Condition: cond, Body: body}
}

func (w *While) TokenLiteral() string {
	return w.Token.Literal
}

func (w *While) expressNode() {}

func (w *While) Str() string {
	return fmt.Sprintf("%s %s { %s }", w.TokenLiteral(), w.Condition.Str(), w.Body.Str())
}

// represents an Array Expression
type Array struct {
	Token  l.Token      // represents the token of the expression
	Values []Expression // represents the values inside the array
}

// generates a new array instance
func NewArray(token l.Token, values ...Expression) *Array {
	return &Array{Token: token, Values: values}
}

func (a *Array) TokenLiteral() string {
	return a.Token.Literal
}

func (a *Array) expressNode() {}

func (a *Array) Str() string {
	var out = make([]string, 0, len(a.Values))
	for _, val := range a.Values {
		out = append(out, val.Str())
	}

	return strings.Join(out, ", ")
}

// represents a call to a data structure like maps, arrays or strings
type CallList struct {
	Token     l.Token    // represents the token of the expression
	ListIdent Expression // represents the data structure to be call
	Index     Expression // represents where is the values in the data structure
}

// generates a new CallList instance
func NewCallList(token l.Token, listIdent Expression, index Expression) *CallList {
	return &CallList{
		Token:     token,
		ListIdent: listIdent,
		Index:     index,
	}
}

func (c *CallList) TokenLiteral() string {
	return c.Token.Literal
}

func (c *CallList) expressNode() {}
func (c *CallList) Str() string {
	return fmt.Sprintf("%s[%s]", c.ListIdent.Str(), c.Index.Str())
}

// Represents a HashMap expression
type MapExpression struct {
	Token l.Token     // represents the token of the expression
	Body  []*KeyValue // represents all the key values pairs in the HashMap
}

// generates a new MapExpression instance
func NewMapExpression(token l.Token, body []*KeyValue) *MapExpression {
	return &MapExpression{token, body}
}

func (m *MapExpression) TokenLiteral() string {
	return m.Token.Literal
}

func (m *MapExpression) expressNode() {}

func (m *MapExpression) Str() string {
	var buff = make([]string, 0, len(m.Body))
	for _, keyVal := range m.Body {
		buff = append(buff, keyVal.Str())
	}

	return fmt.Sprintf("mapa{%s}", strings.Join(buff, ", "))
}
