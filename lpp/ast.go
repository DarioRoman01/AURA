package lpp

import (
	"fmt"
	"strings"
)

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

type Expression interface {
	ASTNode
	expressNode()
}

type Program struct {
	Staments []Stmt
}

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
	var out []string

	for _, v := range p.Staments {
		out = append(out, v.Str())
	}

	return strings.Join(out, " ")
}

type Identifier struct {
	token Token
	value string
}

func NewIdentifier(token Token, value string) *Identifier {
	return &Identifier{token: token, value: value}
}

func (i Identifier) expressNode() {}

func (i Identifier) TokenLiteral() string {
	return i.token.Literal
}

func (i Identifier) Str() string {
	return i.value
}

type LetStatement struct {
	Token Token
	Name  *Identifier
	Value Expression
}

func NewLetStatement(token Token, name *Identifier, value Expression) *LetStatement {
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
	return fmt.Sprintf("%s %s = %s;", l.TokenLiteral(), l.Name.Str(), l.Value.TokenLiteral())
}

type ReturnStament struct {
	Token       Token
	ReturnValue Expression
}

func NewReturnStatement(token Token, returnValue Expression) *ReturnStament {
	return &ReturnStament{Token: token, ReturnValue: returnValue}
}

func (r ReturnStament) TokenLiteral() string {
	return r.Token.Literal
}

func (r ReturnStament) stmtNode() {}

func (r ReturnStament) Str() string {
	return fmt.Sprintf("%s %s;", r.TokenLiteral(), r.ReturnValue.TokenLiteral())
}

type ExpressionStament struct {
	token      Token
	Expression Expression
}

func NewExpressionStament(token Token, expression Expression) *ExpressionStament {
	return &ExpressionStament{token: token, Expression: expression}
}

func (e ExpressionStament) TokenLiteral() string {
	return e.token.Literal
}

func (e ExpressionStament) stmtNode() {}

func (e ExpressionStament) Str() string {
	return e.Expression.Str()
}

type Integer struct {
	Token Token
	Value *int
}

func NewInteger(token Token, value *int) *Integer {
	return &Integer{Token: token, Value: value}
}

func (i Integer) TokenLiteral() string {
	return i.Token.Literal
}

func (i Integer) expressNode() {}
func (i Integer) Str() string {
	return fmt.Sprint(*i.Value)
}

type Prefix struct {
	Token    Token
	Operator string
	Rigth    Expression
}

func NewPrefix(token Token, operator string, rigth Expression) *Prefix {
	return &Prefix{
		Token:    token,
		Operator: operator,
		Rigth:    rigth,
	}
}

func (p *Prefix) TokenLiteral() string {
	return p.Token.Literal
}

func (p *Prefix) expressNode() {}
func (p *Prefix) Str() string {
	return fmt.Sprintf("(%s %s)", p.Operator, p.Rigth.Str())
}

type Infix struct {
	Token    Token
	Rigth    Expression
	Operator string
	Left     Expression
}

func Newinfix(token Token, r Expression, operator string, l Expression) *Infix {
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

func (i Infix) expressNode() {}
func (i Infix) Str() string {
	return fmt.Sprintf("(%s %s %s)", i.Left.Str(), i.Operator, i.Rigth.Str())
}

type Boolean struct {
	Token Token
	Value *bool
}

func NewBoolean(token Token, value *bool) *Boolean {
	return &Boolean{Token: token, Value: value}
}

func (b Boolean) TokenLiteral() string {
	return b.Token.Literal
}

func (b Boolean) expressNode() {}
func (b Boolean) Str() string {
	return b.TokenLiteral()
}

type Block struct {
	Token    Token
	Staments []Stmt
}

func NewBlock(token Token, staments ...Stmt) *Block {
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
	var out []string
	for _, stament := range b.Staments {
		out = append(out, stament.Str())
	}

	return strings.Join(out, " ")
}

type If struct {
	Token       Token
	Condition   Expression
	Consequence *Block
	Alternative *Block
}

func NewIf(token Token, condition Expression, consequence, alternative *Block) *If {
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

type Function struct {
	Token      Token
	Parameters []*Identifier
	Body       *Block
}

func NewFunction(token Token, body *Block, parameters ...*Identifier) *Function {
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
	var paramList []string
	for _, parameter := range f.Parameters {
		paramList = append(paramList, parameter.Str())
	}

	params := strings.Join(paramList, ", ")
	return fmt.Sprintf("%s(%s) %s", f.TokenLiteral(), params, f.Body.Str())
}

type Call struct {
	Token     Token
	Function  Expression
	Arguments []Expression
}

func NewCall(token Token, function Expression, arguments ...Expression) *Call {
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
	var argsList []string
	for _, arg := range c.Arguments {
		argsList = append(argsList, arg.Str())
	}

	args := strings.Join(argsList, ", ")
	return fmt.Sprintf("%s(%s)", c.Function.Str(), args)
}
