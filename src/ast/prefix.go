package ast

import (
	l "aura/src/lexer"
	"fmt"
)

// Represents a prefix expression like:
//		return 5;
type Prefix struct {
	BaseNode            // Extends base node struct
	Operator string     // represents the operator of the expression
	Rigth    Expression // represents the obj of the expression
}

// generates a new prefix instance
func NewPrefix(token l.Token, operator string, rigth Expression) *Prefix {
	return &Prefix{
		BaseNode: BaseNode{token},
		Operator: operator,
		Rigth:    rigth,
	}
}
func (p Prefix) expressNode() {}

func (p Prefix) Str() string {
	return fmt.Sprintf("(%s %s)", p.Operator, p.Rigth.Str())
}

// represents a variable or function name
type Identifier struct {
	BaseNode        // Extends base node struct
	Value    string // represents the name of the identifier
}

// generates a new identifier instance
func NewIdentifier(token l.Token, value string) *Identifier {
	return &Identifier{BaseNode: BaseNode{token}, Value: value}
}

func (i Identifier) expressNode() {}

func (i Identifier) Str() string {
	return i.Value
}

// Represents a Integer expression
type Integer struct {
	BaseNode      // Extends base node struct
	Value    *int // represents the value of the integer
}

// geneerate a new instance of integer
func NewInteger(token l.Token, value *int) *Integer {
	return &Integer{BaseNode: BaseNode{token}, Value: value}
}

func (i Integer) expressNode() {}

func (i Integer) Str() string {
	return fmt.Sprintf("%d", *i.Value)
}

// Represents a float expression
type FloatExp struct {
	BaseNode         // Extends base node struct
	Value    float64 // represents the value of the expression
}

func NewFloatExp(token l.Token, value float64) *FloatExp {
	return &FloatExp{BaseNode: BaseNode{token}, Value: value}
}

func (f *FloatExp) expressNode() {}

func (f *FloatExp) Str() string {
	return fmt.Sprintf("%f", f.Value)
}

// Represents a boolean expression
type Boolean struct {
	BaseNode       // Extends base node struct
	Value    *bool // represents the value of the expression
}

// return a new boolean instance
func NewBoolean(token l.Token, value *bool) *Boolean {
	return &Boolean{BaseNode: BaseNode{token}, Value: value}
}
func (b Boolean) expressNode() {}

func (b Boolean) Str() string {
	return b.TokenLiteral()
}

// represents a string literal expression
type StringLiteral struct {
	BaseNode        // Extends base node struct
	Value    string // represents the value of the expression
}

// return a new string literal instance
func NewStringLiteral(token l.Token, value string) *StringLiteral {
	return &StringLiteral{BaseNode: BaseNode{token}, Value: value}
}

func (s StringLiteral) expressNode() {}

func (s StringLiteral) Str() string {
	return s.Value
}

// represents a null expression
type NullExpression struct {
	BaseNode // Extends base node struct
}

// generates a new null instance
func NewNull(token l.Token) *NullExpression {
	return &NullExpression{BaseNode: BaseNode{token}}
}

func (n *NullExpression) expressNode() {}

func (n *NullExpression) Str() string {
	return "nulo"
}
