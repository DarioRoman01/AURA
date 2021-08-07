package object

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type applyFunc func(Object, ...Object) Object
type isTruthyFunc func(Object) bool

// Represents an Array
type List struct {
	Values []Object // represents all the values in the array
}

func (l *List) Type() ObjectType { return LIST }
func (l *List) Inspect() string {
	var buf strings.Builder
	for idx, val := range l.Values {
		if idx == len(l.Values)-1 {
			buf.WriteString(val.Inspect())
		} else {
			buf.WriteString(val.Inspect() + ", ")
		}
	}

	return fmt.Sprintf("[%s]", buf.String())
}

// add a object to the values of the array
func (l *List) Add(obj Object) {
	l.Values = append(l.Values, obj)
}

// pop the last item in the array
func (l *List) Pop() Object {
	if len(l.Values) == 0 {
		return &Error{"La lista esta vacia"}
	}

	obj := l.Values[len(l.Values)-1]
	l.Values = l.Values[:len(l.Values)-1]
	return obj
}

// remove elements by index
func (l *List) RemoveAt(index int) Object {
	if index >= len(l.Values) || len(l.Values) == 0 {
		return &Error{"Indice fuera de rango"}
	}

	val := l.Values[index]
	l.Values = append(l.Values[:index], l.Values[index+1:]...)
	return val
}

func (l *List) Contains(obj Object) Object {
	for _, val := range l.Values {
		if reflect.DeepEqual(val, obj) {
			return SingletonTRUE
		}
	}

	return SingletonFALSE
}

func (l *List) Map(fn *Def, applyFunction applyFunc) *List {
	newList := new(List)
	for _, val := range l.Values {
		newList.Values = append(newList.Values, applyFunction(fn, val))
	}

	return newList
}

func (l *List) ForEach(fn *Def, applyFunction applyFunc) Object {
	for _, val := range l.Values {
		applyFunction(fn, val)
	}

	return SingletonNUll
}

func (l *List) Filter(fn *Def, applyFunction applyFunc, isTruthy isTruthyFunc) *List {
	newList := new(List)
	for _, val := range l.Values {
		eval := applyFunction(fn, val)
		if isTruthy(eval) {
			newList.Values = append(newList.Values, val)
		}
	}

	return newList
}

func (l *List) Count(fn *Def, applyFunction applyFunc, isTruthy isTruthyFunc) *Number {
	count := new(Number)
	for _, val := range l.Values {
		eval := applyFunction(fn, val)
		if isTruthy(eval) {
			count.Value++
		}
	}

	return count
}

// represents a HashMap
type Map struct {
	Store map[string]Object // represents the hashmap it self
}

func (m *Map) Type() ObjectType { return DICT }
func (m *Map) Inspect() string {
	var buff = make([]string, 0, len(m.Store))
	for key, val := range m.Store {
		str := fmt.Sprintf("%s => %s", key, val.Inspect())
		buff = append(buff, str)
	}

	return fmt.Sprintf("{%s}", strings.Join(buff, ", "))
}

// get the value associeted with the given key if exists
func (m *Map) Get(key string) Object {
	obj, exists := m.Store[key]
	if !exists {
		return NullVAlue
	}

	return obj
}

// update the value associeted with the given key if exists
// if not exists is just added to the map
func (m *Map) UpdateKey(key, newVal Object) {
	m.Store[key.Inspect()] = newVal
}

// Set the key value pair in the map and ckeck if the key already exists
func (m *Map) SetValues(key Object, value Object) error {
	if _, exists := m.Store[key.Inspect()]; exists {
		return errors.New("la llave ya existe en el mapa")
	}
	m.Store[key.Inspect()] = value
	return nil
}

// represents the strings object
type String struct {
	Value string // represents the value of the string
}

func (s *String) Type() ObjectType { return STRINGTYPE }
func (s *String) Inspect() string  { return s.Value }

// check if the string is upper case
func (s *String) IsUpper() Object {
	if strings.ToUpper(s.Value) == s.Value {
		return SingletonTRUE
	}

	return SingletonFALSE
}

// check if the string is lower case
func (s *String) IsLower() Object {
	if strings.ToLower(s.Value) == s.Value {
		return SingletonTRUE
	}

	return SingletonFALSE
}

// check if the string contains another string
func (s *String) Contains(val string) Object {
	if strings.Contains(s.Value, val) {
		return SingletonTRUE
	}

	return SingletonFALSE
}

func (s *String) Split(sep string) *List {
	list := strings.Split(s.Value, sep)
	splited := new(List)
	for _, val := range list {
		splited.Values = append(splited.Values, &String{val})
	}
	return splited
}
