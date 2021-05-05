package lpp

import "strings"

type ASTNode interface {
	TokenLiteral() string
	Str() string
}

type Statement struct {
	Token Token
}

func NewStatement(token Token) *Statement {
	return &Statement{Token: token}
}

func (s Statement) TokenLiteral() string {
	return s.Token.Literal
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
