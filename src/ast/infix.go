package ast

import (
	l "aura/src/lexer"
	"fmt"
)

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
