package src

import (
	"fmt"
	"strings"
)

type ObjectType int

const (
	ObjTypeHead = iota
	BOOLEAN
	BUILTIN
	DEF
	ERROR
	INTEGERS
	ITER
	NULL
	RETURNTYPE
	STRINGTYPE
	LIST
)

var types = [...]string{
	BOOLEAN:    "BOOLEANO",
	BUILTIN:    "BUILTIN",
	DEF:        "FUNCION",
	ERROR:      "ERROR",
	INTEGERS:   "ENTERO",
	ITER:       "ITERADOR",
	NULL:       "NULO",
	RETURNTYPE: "REGRESA",
	STRINGTYPE: "TEXTO",
	LIST:       "LISTA",
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
func (n *Null) Inspect() string  { return "" }

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

type String struct {
	Value string
}

func (s *String) Type() ObjectType { return STRINGTYPE }
func (s *String) Inspect() string  { return s.Value }

type BuiltinFunction func(args ...Object) Object

type Builtin struct {
	Fn BuiltinFunction
}

func NewBuiltin(fn BuiltinFunction) *Builtin { return &Builtin{Fn: fn} }

func (b *Builtin) Type() ObjectType { return BUILTIN }
func (b *Builtin) Inspect() string  { return "builtin function" }

// enviroment handles stores the variables of the given program
type Enviroment struct {
	store map[string]Object
	outer *Enviroment
}

func NewEnviroment(outer *Enviroment) *Enviroment {
	return &Enviroment{
		store: make(map[string]Object),
		outer: outer,
	}
}

func (e *Enviroment) GetItem(key string) (Object, bool) {
	val, exists := e.store[key]
	if !exists {
		if e.outer != nil {
			return e.outer.GetItem(key)
		}

		return nil, false
	}
	return val, true
}

func (e *Enviroment) SetItem(key string, val Object) {
	e.store[key] = val
}

func (e *Enviroment) DelItem(key string) {
	delete(e.store, key)
}

type List struct {
	Values []Object
}

func (l *List) Type() ObjectType { return LIST }
func (l *List) Inspect() string {
	var buff []string

	for _, val := range l.Values {
		buff = append(buff, val.Inspect())
	}

	return fmt.Sprintf("[%s]", strings.Join(buff, ", "))
}

type Iterator struct {
	current Object
	list    []Object
}

func (l *List) Iter() *Iterator {
	if len(l.Values) == 0 {
		return nil
	}

	return &Iterator{current: l.Values[0], list: l.Values}
}

func (i *Iterator) Next() Object {
	if len(i.list) == 0 {
		return nil
	}

	val := i.current
	i.list = i.list[1:]
	i.current = i.list[0]
	return val
}

func (i *Iterator) Type() ObjectType { return ITER }
func (i *Iterator) Inspect() string {
	var buff []string

	for _, val := range i.list {
		buff = append(buff, val.Inspect())
	}

	return fmt.Sprintf("[%s]", strings.Join(buff, ", "))
}
