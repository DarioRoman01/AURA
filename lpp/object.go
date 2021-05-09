package lpp

import (
	"fmt"
	"strings"
)

type ObjectType int

const (
	ObjTypeHead = iota
	BOOLEAN
	DEF
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

type Def struct {
	Parameters []*Identifier
	Body       *Block
	Env        *Enviroment
}

func NewDef(body *Block, env *Enviroment, parameters ...*Identifier) *Def {
	return &Def{Parameters: parameters, Body: body, Env: env}
}

func (d *Def) Type() ObjectType {
	return DEF
}

func (d *Def) Inspect() string {
	var argsList []string

	for _, arg := range d.Parameters {
		argsList = append(argsList, arg.Str())
	}

	return fmt.Sprintf("funcion(%s) {\n %s \n}", strings.Join(argsList, ", "), d.Body.Str())
}

// enviroment handles stores the variables of the given program
type Enviroment struct {
	store map[interface{}]Object
	outer *Enviroment
}

func NewEnviroment(outer *Enviroment) *Enviroment {
	return &Enviroment{
		store: make(map[interface{}]Object),
		outer: outer,
	}
}

func (e *Enviroment) GetItem(key interface{}) (Object, bool) {
	val, exists := e.store[key]
	if !exists {
		if e.outer != nil {
			return e.outer.GetItem(key)
		}

		return nil, false
	}
	return val, true
}

func (e *Enviroment) SetItem(key interface{}, val Object) {
	e.store[key] = val
}

func (e *Enviroment) DelItem(key interface{}) {
	delete(e.store, key)
}
