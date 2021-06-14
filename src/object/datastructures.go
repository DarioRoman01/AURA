package object

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"strings"
)

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
