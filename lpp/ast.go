package lpp

import (
	"fmt"
	"strings"
)

type ASTNode interface {
	TokenLiteral() string
	Str() string
}

type Statement struct {
	Token        Token
	LetStatement *LetStatement
}

func NewStatement(token Token, letStatement *LetStatement) *Statement {
	return &Statement{
		Token:        token,
		LetStatement: letStatement,
	}
}

func (s Statement) TokenLiteral() string {
	return s.Token.Literal
}

func (s Statement) Str() string {
	return s.TokenLiteral()
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
	Staments []Statement
}

func NewProgram(statements []Statement) *Program {
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
		out = append(out, v.TokenLiteral())
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
	name  *Identifier
	value *Expression
}

func NewLetStatement(name *Identifier, value *Expression) *LetStatement {
	return &LetStatement{
		name:  name,
		value: value,
	}
}

func (l *LetStatement) TokenLiteral() string {
	return l.name.TokenLiteral()
}

func (l *LetStatement) Str() string {
	return fmt.Sprintf("%s %s = %s;", l.TokenLiteral(), l.name.Str(), l.value.TokenLiteral())
}
