package lpp

import "fmt"

type ObjectType int

const (
	ObjTypeHead = iota
	BOOLEAN
	INTEGERS
	NULL
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

// int object type
type Number struct{ Value int }

func NewNumber(value int) *Number  { return &Number{Value: value} }
func (i *Number) Type() ObjectType { return INTEGERS }
func (i *Number) Inspect() string  { return fmt.Sprint(i.Value) }

// Bool object type
type Bool struct{ value bool }

func NewBool(value bool) *Bool   { return &Bool{value: value} }
func (b *Bool) Type() ObjectType { return BOOLEAN }
func (b *Bool) Inspect() string {
	if b.value {
		return "verdero"
	}

	return "falso"
}

// null object type
type Null struct{}

func (n *Null) Type() ObjectType { return NULL }
func (n *Null) Inspect() string  { return "nulo" }
