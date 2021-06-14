package ast

import (
	"fmt"
	l "katan/src/lexer"
	"strings"
)

// ast node interface
type ASTNode interface {
	TokenLiteral() string
	Str() string
}

// statement interface
// ensure all statements implements ast node
type Stmt interface {
	ASTNode
	stmtNode()
}

// expression interfaces
// ensure all expressions implements ast node
type Expression interface {
	ASTNode
	expressNode()
}

// program is the base node of the ast
type Program struct {
	Staments []Stmt
}

// generates a new program instance
func NewProgram(statements []Stmt) *Program {
	return &Program{Staments: statements}
}

// return the token literal of the first stament
func (p Program) TokenLiteral() string {
	if len(p.Staments) > 0 {
		return p.Staments[0].TokenLiteral()
	}

	return ""
}

// return a string representation of the program
func (p Program) Str() string {
	var out []string

	for _, v := range p.Staments {
		out = append(out, v.Str())
	}

	return strings.Join(out, " ")
}

// identifier handles variables names and function names
type Identifier struct {
	Token l.Token
	Value string
}

// generates a new identifier instance
func NewIdentifier(token l.Token, value string) *Identifier {
	return &Identifier{Token: token, Value: value}
}

// ensure the identifier is a expression
func (i Identifier) expressNode() {}

// return literal of the identifier token
func (i Identifier) TokenLiteral() string {
	return i.Token.Literal
}

// return the value of the identifier wich is the name
func (i Identifier) Str() string {
	return i.Value
}

// let stament handles assing operations
type LetStatement struct {
	Token l.Token
	Name  *Identifier
	Value Expression
}

// generate a new let stament instance
func NewLetStatement(token l.Token, name *Identifier, value Expression) *LetStatement {
	return &LetStatement{
		Token: token,
		Name:  name,
		Value: value,
	}
}

// return the literal of the stament token
func (l LetStatement) TokenLiteral() string {
	return l.Token.Literal
}

// ensure let statement is a expression node
func (l LetStatement) stmtNode() {}

// return a string representation of the stament
func (l LetStatement) Str() string {
	return fmt.Sprintf("%s %s = %s;", l.TokenLiteral(), l.Name.Str(), l.Value.Str())
}

// handles return staments
type ReturnStament struct {
	Token       l.Token
	ReturnValue Expression
}

// generates a new return statement instance
func NewReturnStatement(token l.Token, returnValue Expression) *ReturnStament {
	return &ReturnStament{Token: token, ReturnValue: returnValue}
}

// return the literal of the return stament token
func (r ReturnStament) TokenLiteral() string {
	return r.Token.Literal
}

// ensure return statement implements expression
func (r ReturnStament) stmtNode() {}

// return a string representation of the return stament node
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

// ensure expression statement implements expression
func (e ExpressionStament) stmtNode() {}

// return a string representation of the string
func (e ExpressionStament) Str() string {
	return e.Expression.Str()
}

// integer type
type Integer struct {
	Token l.Token
	Value *int
}

// geneerate a new instance of integer
func NewInteger(token l.Token, value *int) *Integer {
	return &Integer{Token: token, Value: value}
}

func (i Integer) TokenLiteral() string {
	return i.Token.Literal
}

// ensure integer implements expression
func (i Integer) expressNode() {}

// return the integer value as a string
func (i Integer) Str() string {
	return fmt.Sprintf("%d", *i.Value)
}

type Suffix struct {
	Token    l.Token
	Left     Expression
	Operator string
}

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

// prefix handles prefix staments like regresa x;
type Prefix struct {
	Token    l.Token
	Operator string
	Rigth    Expression
}

// generates a new prefix instance
func NewPrefix(token l.Token, operator string, rigth Expression) *Prefix {
	return &Prefix{
		Token:    token,
		Operator: operator,
		Rigth:    rigth,
	}
}

func (p Prefix) TokenLiteral() string {
	return p.Token.Literal
}

// ensure that prefix implements expression
func (p Prefix) expressNode() {}

// return a string representation of the prefix stament
func (p Prefix) Str() string {
	return fmt.Sprintf("(%s %s)", p.Operator, p.Rigth.Str())
}

// infix handles expressions like 5 + 5; where the operator is in the middle of two values
type Infix struct {
	Token    l.Token
	Rigth    Expression
	Operator string
	Left     Expression
}

// generates a new infix instance
func Newinfix(token l.Token, r Expression, operator string, l Expression) *Infix {
	return &Infix{
		Token:    token,
		Rigth:    r,
		Operator: operator,
		Left:     l,
	}
}

func (i Infix) TokenLiteral() string {
	return i.Token.Literal
}

// ensure that infix implements expression node
func (i Infix) expressNode() {}

// return a string representation of the infix stament
func (i Infix) Str() string {
	return fmt.Sprintf("(%s %s %s)", i.Left.Str(), i.Operator, i.Rigth.Str())
}

// boolean type
type Boolean struct {
	Token l.Token
	Value *bool
}

// return a new boolean instance
func NewBoolean(token l.Token, value *bool) *Boolean {
	return &Boolean{Token: token, Value: value}
}

func (b Boolean) TokenLiteral() string {
	return b.Token.Literal
}

// ensure boolean implements expression
func (b Boolean) expressNode() {}

// return a string representation of the boolean
func (b Boolean) Str() string {
	return b.TokenLiteral()
}

// Block group a chunk of staments
type Block struct {
	Token    l.Token
	Staments []Stmt
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

// ensure that block implements stament
func (b Block) stmtNode() {}

// return a string representation of the block
func (b Block) Str() string {
	var out []string
	for _, stament := range b.Staments {
		out = append(out, stament.Str())
	}

	return strings.Join(out, " ")
}

// if handles all the logic of if staments
type If struct {
	Token       l.Token
	Condition   Expression
	Consequence *Block
	Alternative *Block
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

// ensure If implements expression
func (i If) expressNode() {}

// return a string representation of the if else statement
func (i If) Str() string {
	var out strings.Builder
	out.WriteString(fmt.Sprintf("si %s %s ", i.Condition.Str(), i.Consequence.Str()))
	if i.Alternative != nil {
		out.WriteString(fmt.Sprintf("si_no %s", i.Alternative.Str()))
	}

	return out.String()
}

// function type handles functions declarations
type Function struct {
	Token      l.Token
	Parameters []*Identifier
	Body       *Block
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

// ensure that function implements expression
func (f Function) expressNode() {}

// return a string representation of the function
func (f Function) Str() string {
	var paramList []string
	for _, parameter := range f.Parameters {
		paramList = append(paramList, parameter.Str())
	}

	params := strings.Join(paramList, ", ")
	return fmt.Sprintf("%s(%s) %s", f.TokenLiteral(), params, f.Body.Str())
}

// Call type handles function calls and its arguments
type Call struct {
	Token     l.Token
	Function  Expression
	Arguments []Expression
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

// ensure that call implements expression
func (C Call) expressNode() {}

// return a string representation of the call
func (c Call) Str() string {
	var argsList []string
	for _, arg := range c.Arguments {
		argsList = append(argsList, arg.Str())
	}

	args := strings.Join(argsList, ", ")
	return fmt.Sprintf("%s(%s)", c.Function.Str(), args)
}

type StringLiteral struct {
	Token l.Token
	Value string
}

func NewStringLiteral(token l.Token, value string) *StringLiteral {
	return &StringLiteral{Token: token, Value: value}
}

func (s StringLiteral) TokenLiteral() string {
	return s.Token.Literal
}

func (s StringLiteral) expressNode() {}

func (s StringLiteral) Str() string {
	return s.Value
}

type For struct {
	Token     l.Token
	Condition Expression
	Body      *Block
}

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

type While struct {
	Token     l.Token
	Condition Expression
	Body      *Block
}

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

type Array struct {
	Token  l.Token
	Values []Expression
}

func NewArray(token l.Token, values ...Expression) *Array {
	return &Array{Token: token, Values: values}
}

func (a *Array) TokenLiteral() string {
	return a.Token.Literal
}

func (a *Array) expressNode() {}

func (a *Array) Str() string {
	var out []string
	for _, val := range a.Values {
		out = append(out, val.Str())
	}

	return strings.Join(out, ", ")
}

type CallList struct {
	Token     l.Token
	ListIdent Expression
	Index     Expression
}

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

type Reassignment struct {
	Token      l.Token
	Identifier Expression
	NewVal     Expression
}

func NewReassignment(token l.Token, ident Expression, newVal Expression) *Reassignment {
	return &Reassignment{
		Token:      token,
		Identifier: ident,
		NewVal:     newVal,
	}
}

func (r *Reassignment) TokenLiteral() string {
	return r.Token.Literal
}

func (r *Reassignment) expressNode() {}

func (r *Reassignment) Str() string {
	return fmt.Sprintf("%s = %s", r.Identifier.Str(), r.NewVal.Str())
}

type RangeExpression struct {
	Token    l.Token
	Variable Expression
	Range    Expression
}

func NewRange(token l.Token, variable Expression, Range Expression) *RangeExpression {
	return &RangeExpression{
		Token:    token,
		Variable: variable,
		Range:    Range,
	}
}

func (r *RangeExpression) TokenLiteral() string {
	return r.Token.Literal
}

func (r *RangeExpression) expressNode() {}

func (r *RangeExpression) Str() string {
	return fmt.Sprintf("%s en %s", r.Variable.Str(), r.Range.Str())
}

type MethodExpression struct {
	Token  l.Token
	Obj    Expression
	Method Expression
}

func NewMethodExpression(token l.Token, obj Expression, method Expression) *MethodExpression {
	return &MethodExpression{
		Token:  token,
		Obj:    obj,
		Method: method,
	}
}

func (m *MethodExpression) TokenLiteral() string {
	return m.Token.Literal
}

func (m *MethodExpression) expressNode() {}

func (m *MethodExpression) Str() string {
	return fmt.Sprintf("%s:%s", m.Obj.Str(), m.Method.Str())
}

type NullExpression struct {
	Token l.Token
}

func NewNull(token l.Token) *NullExpression {
	return &NullExpression{Token: token}
}

func (n *NullExpression) TokenLiteral() string {
	return n.Token.Literal
}

func (n *NullExpression) expressNode() {}

func (n *NullExpression) Str() string {
	return "nulo"
}

type MapExpression struct {
	Token l.Token
	Body  []*KeyValue
}

func NewMapExpression(token l.Token, body []*KeyValue) *MapExpression {
	return &MapExpression{token, body}
}

func (m *MapExpression) TokenLiteral() string {
	return m.Token.Literal
}

func (m *MapExpression) expressNode() {}

func (m *MapExpression) Str() string {
	var buff []string
	for _, keyVal := range m.Body {
		buff = append(buff, keyVal.Str())
	}

	return fmt.Sprintf("mapa{%s}", strings.Join(buff, ", "))
}

type KeyValue struct {
	Token l.Token
	Key   Expression
	Value Expression
}

func NewKeyVal(token l.Token, key, value Expression) *KeyValue {
	return &KeyValue{
		Token: token,
		Key:   key,
		Value: value,
	}
}

func (k *KeyValue) TokenLiteral() string {
	return k.Token.Literal
}

func (k *KeyValue) expressNode() {}

func (k *KeyValue) Str() string {
	return fmt.Sprintf("%s => %s", k.Key.Str(), k.Value.Str())
}