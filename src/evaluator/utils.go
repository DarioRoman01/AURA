package evaluator

import (
	obj "katan/src/object"
)

func makeStringList(str string) []obj.Object {
	list := []obj.Object{}
	for _, char := range str {
		list = append(list, &obj.String{Value: string(char)})
	}

	return list
}