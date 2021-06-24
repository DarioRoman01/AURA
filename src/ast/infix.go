package ast

import (
	l "aura/src/lexer"
	"fmt"
)

// Represents a infix expression like 5 + 5:
type Infix struct {
	Token    l.Token    // represents the token of the expression
	Rigth    Expression // represents the rigth object of the expresion
	Operator string     // represents the operator between the objects
	Left     Expression // represents the left object of the expresion
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

func (i Infix) expressNode() {}

func (i Infix) Str() string {
	return fmt.Sprintf("(%s %s %s)", i.Left.Str(), i.Operator, i.Rigth.Str())
}

// represents an in expression like:
//		for i in range(10):
//			do something
//
type RangeExpression struct {
	Token    l.Token    // represents the token of the expression
	Variable Expression // represents the variable in the expression
	Range    Expression // represents the iterable in the expression
}

// Generates a new Range instance
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

// represents a key value expression:
//		key => value
type KeyValue struct {
	Token l.Token    // represents the token of the expression
	Key   Expression // represents the key of the expression
	Value Expression // represents the value of the expression
}

// generates a new key value instance
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

// represents a method expression like:
//		obj.do_somthing()
type MethodExpression struct {
	Token  l.Token    // represents the token of the expression
	Obj    Expression // represents the object to be apply the method
	Method Expression // represents the method it self
}

// generates a method expression instance
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

// represents a reassigment expression
type Reassignment struct {
	Token      l.Token    // represents the token of the expression
	Identifier Expression // represents the variable to be reassing
	NewVal     Expression // represents the new value to the variable
}

// generates a new reassigmment instacne
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
