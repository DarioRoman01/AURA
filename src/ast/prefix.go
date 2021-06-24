package ast

import (
	l "aura/src/lexer"
	"fmt"
)

// Represents a prefix expression like:
//		return 5;
type Prefix struct {
	Token    l.Token    // represents the token of the expression
	Operator string     // represents the operator of the expression
	Rigth    Expression // represents the obj of the expression
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

func (p Prefix) expressNode() {}

func (p Prefix) Str() string {
	return fmt.Sprintf("(%s %s)", p.Operator, p.Rigth.Str())
}

// represents a variable or function name
type Identifier struct {
	Token l.Token // represents the token of the expression
	Value string  // represents the name of the identifier
}

// generates a new identifier instance
func NewIdentifier(token l.Token, value string) *Identifier {
	return &Identifier{Token: token, Value: value}
}

func (i Identifier) expressNode() {}

func (i Identifier) TokenLiteral() string {
	return i.Token.Literal
}

func (i Identifier) Str() string {
	return i.Value
}

// Represents a Integer expression
type Integer struct {
	Token l.Token // represents the token of the expression
	Value *int    // represents the value of the integer
}

// geneerate a new instance of integer
func NewInteger(token l.Token, value *int) *Integer {
	return &Integer{Token: token, Value: value}
}

func (i Integer) TokenLiteral() string {
	return i.Token.Literal
}

func (i Integer) expressNode() {}

func (i Integer) Str() string {
	return fmt.Sprintf("%d", *i.Value)
}

// Represents a boolean expression
type Boolean struct {
	Token l.Token // represents the token of the expression
	Value *bool   // represents the value of the expression
}

// return a new boolean instance
func NewBoolean(token l.Token, value *bool) *Boolean {
	return &Boolean{Token: token, Value: value}
}

func (b Boolean) TokenLiteral() string {
	return b.Token.Literal
}

func (b Boolean) expressNode() {}

func (b Boolean) Str() string {
	return b.TokenLiteral()
}

// represents a string literal expression
type StringLiteral struct {
	Token l.Token // represents the token of the expression
	Value string  // represents the value of the expression
}

// return a new string literal instance
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

// represents a null expression
type NullExpression struct {
	Token l.Token // represents the token of the expression
}

// generates a new null instance
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
