package lpp

import "fmt"

type ObjectType int

const (
	ObjTypeHead = iota
	BOOLEAN
	INTEGERS
	NULL
	RETURNTYPE
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

// int object type
type Number struct{ Value int }

func (i *Number) Type() ObjectType { return INTEGERS }
func (i *Number) Inspect() string  { return fmt.Sprint(i.Value) }

// Bool object type
type Bool struct{ Value bool }

func (b *Bool) Type() ObjectType { return BOOLEAN }
func (b *Bool) Inspect() string {
	if b.Value {
		return "verdero"
	}

	return "falso"
}

// null object type
type Null struct{}

func (n *Null) Type() ObjectType { return NULL }
func (n *Null) Inspect() string  { return "nulo" }

type Return struct {
	Value Object
}

func (r *Return) Type() ObjectType { return RETURNTYPE }
func (r *Return) Inspect() string  { return r.Value.Inspect() }
