package lpp

import "fmt"

type ObjectType int

const (
	ObjTypeHead = iota
	BOOLEAN
	ERROR
	INTEGERS
	NULL
	RETURNTYPE
)

var types = [...]string{
	BOOLEAN:    "BOOLEAN",
	ERROR:      "ERROR",
	INTEGERS:   "INTEGERS",
	NULL:       "NULL",
	RETURNTYPE: "RETURN",
}

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

// return object type
type Return struct {
	Value Object
}

func (r *Return) Type() ObjectType { return RETURNTYPE }
func (r *Return) Inspect() string  { return r.Value.Inspect() }

// error object type
type Error struct {
	Message string
}

func (e *Error) Type() ObjectType { return ERROR }
func (e *Error) Inspect() string {
	return fmt.Sprintf("Error: %s", e.Message)
}

type Enviroment struct {
	store map[interface{}]Object
}

func NewEnviroment() *Enviroment {
	return &Enviroment{
		store: make(map[interface{}]Object),
	}
}

func (e *Enviroment) GetItem(key interface{}) interface{} {
	return e.store[key]
}

func (e *Enviroment) SetItem(key interface{}, val Object) {
	e.store[key] = val
}

func (e *Enviroment) DelItem(key interface{}) {
	delete(e.store, key)
}
