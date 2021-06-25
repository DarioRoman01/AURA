package evaluator

import (
	obj "aura/src/object"
	"unicode/utf8"
)

func makeStringList(str string) []obj.Object {
	list := make([]obj.Object, 0, utf8.RuneCountInString(str))
	for _, char := range str {
		list = append(list, &obj.String{Value: string(char)})
	}

	return list
}
