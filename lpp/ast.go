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

type Expression struct {
	Token Token
}

func NewExpression(token Token) *Expression {
	return &Expression{Token: token}
}

func (e Expression) TokenLiteral() string {
	return e.Token.Literal
}

func (e Expression) Str() string {
	return e.Token.PrintToken()
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

func (i Identifier) TokenLiteral() string {
	return i.token.Literal
}

func (i Identifier) Str() string {
	return i.value
}

type LetStatement struct {
	Token Token
	Name  *Identifier
	Value *Expression
}

func NewLetStatement(token Token, name *Identifier, value *Expression) *LetStatement {
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
	ReturnValue *Expression
}

func NewReturnStatement(token Token, returnValue *Expression) *ReturnStament {
	return &ReturnStament{Token: token, ReturnValue: returnValue}
}

func (r ReturnStament) TokenLiteral() string {
	return r.Token.Literal
}

func (r ReturnStament) stmtNode() {}

func (r ReturnStament) Str() string {
	return fmt.Sprintf("%s %s;", r.TokenLiteral(), r.ReturnValue.TokenLiteral())
}
