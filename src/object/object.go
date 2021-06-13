package object

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"katan/src/ast"
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
	METHOD
	DICT
)

type MethodsTypes int

const (
	Lhead MethodsTypes = iota
	POP
	APPEND
	REMOVE
	CONTAIS
	KEYS
	VALUES
)

var Types = [...]string{
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
	METHOD:     "METODO",
	DICT:       "MAPA",
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

func NewBool(value bool) *Bool   { return &Bool{Value: value} }
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
	Parameters []*ast.Identifier
	Body       *ast.Block
	Env        *Enviroment
}

func NewDef(body *ast.Block, env *Enviroment, parameters ...*ast.Identifier) *Def {
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
	Store map[string]Object
	outer *Enviroment
}

func NewEnviroment(outer *Enviroment) *Enviroment {
	return &Enviroment{
		Store: make(map[string]Object),
		outer: outer,
	}
}

func (e *Enviroment) GetItem(key string) (Object, bool) {
	val, exists := e.Store[key]
	if !exists {
		if e.outer != nil {
			return e.outer.GetItem(key)
		}

		return nil, false
	}
	return val, true
}

func (e *Enviroment) SetItem(key string, val Object) {
	e.Store[key] = val
}

func (e *Enviroment) DelItem(key string) {
	delete(e.Store, key)
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

func (l *List) Add(obj Object) {
	l.Values = append(l.Values, obj)
}

func (l *List) Pop() Object {
	if len(l.Values) == 0 {
		return &Error{"La lista esta vacia"}
	}

	obj := l.Values[len(l.Values)-1]
	l.Values = l.Values[:len(l.Values)-1]
	return obj
}

func (l *List) RemoveAt(index int) Object {
	if index >= len(l.Values) || len(l.Values) == 0 {
		return &Error{"Indice fuera de rango"}
	}

	val := l.Values[index]
	l.Values = append(l.Values[:index], l.Values[index+1:]...)
	return val
}

type Iterator struct {
	Current Object
	List    []Object
}

func NewIterator(current Object, values []Object) *Iterator {
	return &Iterator{Current: current, List: values}
}

func (i *Iterator) Next() Object {
	if len(i.List) == 0 {
		return nil
	}

	val := i.Current
	i.List = i.List[1:]

	if len(i.List) != 0 {
		i.Current = i.List[0]
	}

	return val
}

func (i *Iterator) Type() ObjectType { return ITER }
func (i *Iterator) Inspect() string {
	var buff []string

	for _, val := range i.List {
		buff = append(buff, val.Inspect())
	}

	return fmt.Sprintf("[%s]", strings.Join(buff, ", "))
}

type Method struct {
	Value      Object
	MethodType MethodsTypes
}

func NewMethod(val Object, methodType MethodsTypes) *Method {
	return &Method{Value: val, MethodType: methodType}
}

func (m *Method) Type() ObjectType { return METHOD }
func (m *Method) Inspect() string {
	return fmt.Sprintf(":%d(%s)", m.MethodType, m.Value.Inspect())
}

type Map struct {
	Store map[string]Object
}

func (m *Map) Type() ObjectType { return DICT }
func (m *Map) Inspect() string {
	var buff []string

	for key, val := range m.Store {
		str := fmt.Sprintf("%s => %s", m.Deserialize([]byte(key)), val.Inspect())
		buff = append(buff, str)
	}

	return fmt.Sprintf("[%s]", strings.Join(buff, ", "))
}

func (m *Map) Get(key string) Object {
	obj, exists := m.Store[key]
	if !exists {
		return NullVAlue
	}

	return obj
}

func (m *Map) UpdateKey(key, newVal Object) {
	hashedKey := m.Serialize(key)
	m.Store[string(hashedKey)] = newVal
}

func (m *Map) Serialize(obj Object) []byte {
	var buff bytes.Buffer
	encoder := gob.NewEncoder(&buff)
	encoder.Encode(obj.Inspect())
	return buff.Bytes()
}

func (m *Map) Deserialize(data []byte) string {
	var str string
	decoder := gob.NewDecoder(bytes.NewReader(data))
	decoder.Decode(&str)
	return str
}

func (m *Map) SetValues(key Object, value Object) error {
	hashedKey := m.Serialize(key)
	if _, exists := m.Store[string(hashedKey)]; exists {
		return errors.New("la llave ya existe en el mapa")
	}

	m.Store[string(hashedKey)] = value
	return nil
}

// use singleton patern with true false and null
var (
	SingletonTRUE  = &Bool{Value: true}
	SingletonFALSE = &Bool{Value: false}
	SingletonNUll  = &Null{} //  this null is for functions that dont return anything
	NullVAlue      = &Null{} // this null is the null value
)
