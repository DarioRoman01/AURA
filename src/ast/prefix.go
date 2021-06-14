package ast

import (
	l "aura/src/lexer"
	"fmt"
)

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
