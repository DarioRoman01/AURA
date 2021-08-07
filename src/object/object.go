package object

import (
	"aura/src/ast"
	"fmt"
	"strings"
)

// represents all the types in the programming lenguage
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
	FLOATING
	CLASS
	BREAK
	CONTINUE
)

// represents the methods in the standar library
type MethodsTypes int

const (
	Lhead MethodsTypes = iota
	POP
	MAP
	COUNT
	FILTER
	FOREACH
	APPEND
	REMOVE
	CONTAIS
	KEYS
	VALUES
	UPPER
	LOWER
	ISUPPER
	ISLOWER
	SPLIT
)

// string representation of the types
var Types = [...]string{
	BOOLEAN:    "booleano",
	BUILTIN:    "builtin",
	DEF:        "funcion",
	ERROR:      "error",
	INTEGERS:   "entero",
	ITER:       "iterador",
	NULL:       "nulo",
	RETURNTYPE: "regesa",
	STRINGTYPE: "texto",
	LIST:       "lista",
	METHOD:     "metodo",
	DICT:       "mapa",
	CLASS:      "clase",
}

// Object is an interface for abstract all the structs
type Object interface {
	Type() ObjectType // return the object type of the object
	Inspect() string  // return the value of the object
}

// represent the int object
type Number struct{ Value int }

func (i *Number) Type() ObjectType { return INTEGERS }
func (i *Number) Inspect() string  { return fmt.Sprint(i.Value) }

// represents the float object type
type Float struct{ Value float64 }

func NewFloat(val float64) *Float {
	return &Float{val}
}

func (f *Float) Type() ObjectType { return FLOATING }
func (f *Float) Inspect() string  { return fmt.Sprint(f.Value) }

// represent the bool object
type Bool struct{ Value bool }

func (b *Bool) Type() ObjectType { return BOOLEAN }
func (b *Bool) Inspect() string {
	if b.Value {
		return "verdadero"
	}

	return "falso"
}

// represent the null object
type Null struct {
	Val interface{}
}

func (n *Null) Type() ObjectType { return NULL }
func (n *Null) Inspect() string  { return "nulo" }

// represent the return object
type Return struct {
	Value Object // represents the value to be returned
}

func (r *Return) Type() ObjectType { return RETURNTYPE }
func (r *Return) Inspect() string  { return r.Value.Inspect() }

// represents the error object
type Error struct {
	Message string // represents the error message
}

func (e *Error) Type() ObjectType { return ERROR }
func (e *Error) Inspect() string {
	return fmt.Sprintf("Error: %s", e.Message)
}

// represents the function object
type Def struct {
	Parameters []*ast.Identifier // represents the parameters of the function
	Body       *ast.Block        // represents the body of the function
	Env        *Enviroment       // represents the scope of the function
}

// return a new function object instance
func NewDef(body *ast.Block, env *Enviroment, parameters ...*ast.Identifier) *Def {
	return &Def{Parameters: parameters, Body: body, Env: env}
}

func (d *Def) Type() ObjectType {
	return DEF
}

func (d *Def) Inspect() string {
	var buf strings.Builder
	for idx, arg := range d.Parameters {
		if idx == len(d.Parameters)-1 {
			buf.WriteString(arg.Str())
		} else {
			buf.WriteString(arg.Str() + ", ")
		}
	}

	return fmt.Sprintf("funcion(%s) {\n %s \n}", buf.String(), d.Body.Str())
}

// signature for builtin functions
type BuiltinFunction func(args ...Object) Object

// represents a builtin function
type Builtin struct {
	Fn BuiltinFunction // represents the function of the builtin
}

// return a new builtin instance
func NewBuiltin(fn BuiltinFunction) *Builtin { return &Builtin{Fn: fn} }

func (b *Builtin) Type() ObjectType { return BUILTIN }
func (b *Builtin) Inspect() string  { return "builtin function" }

// Represents a escope in the programming lengauge
type Enviroment struct {
	Store map[string]Object // repesents the store of all variables
	outer *Enviroment       // represents a posible outer scope
}

// return a new enviroment instance
func NewEnviroment(outer *Enviroment) *Enviroment {
	return &Enviroment{
		Store: make(map[string]Object),
		outer: outer,
	}
}

// return a optional object if exists in the scope
func (e *Enviroment) GetItem(key string) (Object, bool) {
	val, exists := e.Store[key]
	if !exists {
		// we check if there is an outer env and call the same method to find the object
		if e.outer != nil {
			return e.outer.GetItem(key)
		}

		return nil, false
	}
	return val, true
}

// store an object in the eviroment
func (e *Enviroment) SetItem(key string, val Object) {
	e.Store[key] = val
}

// delete an item form the enviroment
func (e *Enviroment) DelItem(key string) {
	delete(e.Store, key)
}

func (e *Enviroment) SetOuter(env *Enviroment) {
	e.outer = env
}

// repesents an iterator object
type Iterator struct {
	Current Object      // represents the current object
	List    []Object    // represents the values in the iter
	Env     *Enviroment // represents the iterator enviroment
}

// return a new iterator instance
func NewIterator(current Object, values []Object, env *Enviroment) *Iterator {
	return &Iterator{Current: current, List: values, Env: env}
}

// return the next value in the iter if there is any
// and remove the value from the iter
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
	var buf strings.Builder
	for idx, val := range i.List {
		if idx == len(i.List)-1 {
			buf.WriteString(val.Inspect())
		} else {
			buf.WriteString(val.Inspect() + ", ")
		}
	}

	return fmt.Sprintf("[%s]", buf.String())
}

// represents a method object
type Method struct {
	Value      Object       // represents the value evaluated from the arguments
	MethodType MethodsTypes // represents the method type
}

// generates a new method instance
func NewMethod(val Object, methodType MethodsTypes) *Method {
	return &Method{Value: val, MethodType: methodType}
}

func (m *Method) Type() ObjectType { return METHOD }
func (m *Method) Inspect() string {
	return fmt.Sprintf(":%d(%s)", m.MethodType, m.Value.Inspect())
}

// represets the class object
type Class struct {
	Name    *ast.Identifier       // represents the class name
	Params  []*ast.Identifier     // represents the constructor params
	Methods []*ast.ClassMethodExp // represents all the methods in the class
}

func (c *Class) Type() ObjectType { return CLASS }
func (c *Class) Inspect() string {
	var buf strings.Builder
	for idx, param := range c.Params {
		if idx == len(c.Params)-1 {
			buf.WriteString(param.Value)
		} else {
			buf.WriteString(param.Value + ", ")
		}
	}

	for _, value := range c.Methods {
		buf.WriteString("\n")
		buf.WriteString(value.Str())
	}

	return buf.String()
}

func NewClass(name *ast.Identifier, params []*ast.Identifier, methods []*ast.ClassMethodExp) *Class {
	return &Class{
		Name:    name,
		Params:  params,
		Methods: methods,
	}
}

// Represents a class instance object
type ClassInstance struct {
	Name string      // represents the class instance name
	Env  *Enviroment // represents the scope of the class
}

// generates a new class instance
func NewClassInstance(name string, env *Enviroment) *ClassInstance {
	return &ClassInstance{
		Name: name,
		Env:  env,
	}
}

func (c *ClassInstance) Type() ObjectType { return CLASS }
func (c *ClassInstance) Inspect() string {
	return fmt.Sprintf("clase %s", c.Name)
}

type BreakObj struct{}

func (b *BreakObj) Type() ObjectType { return BREAK }
func (b *BreakObj) Inspect() string  { return "romper" }

type ContinueObj struct{}

func (c *ContinueObj) Type() ObjectType { return CONTINUE }
func (c *ContinueObj) Inspect() string  { return "continuar" }

// use singleton patern with true false and null
var (
	SingletonTRUE     = &Bool{Value: true}
	SingletonFALSE    = &Bool{Value: false}
	SingletonNUll     = &Null{nil} //  this null is for functions that dont return anything
	NullVAlue         = &Null{1}   // this null is the null value
	SingleTonBreak    = &BreakObj{}
	SingletonContinue = &ContinueObj{}
)
