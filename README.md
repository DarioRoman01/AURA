# The Katan programing lenguage

## What is this?
This is a very basic programming lenguage inspired in javascript, the goal is to create a programming lenguage in spanish
that helps people that are begining in software development or computer science, probiding a high level programming lenguage
in spanish that is very simple to use.

It is an interpreted language fully built using Go.

## Syntax example

## Variables
```ts
var a = 5;
var b = 5;
var c = a + b;

escribir(c); // prints 10
```

### Operators

These are operators:
| Operator             | Symbol |
|----------------------|--------|
| Plus                 |    +   |
| Increment            |   ++   |
| Add assigment        |   +=   |
| Minus                |    -   |
| Decrement            |    --  |
| Subtract assigment   |   -=   |
| Multiplication       |    *   |
| Exponential          |   **   |
| Multiply assigment   |   *=   |
| Division             |    /   |
| Division assigment   |   /=   |
| Modulus              |    %   |
| Equal                |   ==   |
| Not Equal            |   !=   |
| Not                  |    !   |
| Less than            |    <   |
| Greater than         |    >   |
| Less or equal than   |   <=   |
| Greater or equal than|   >=   |
| And                  |   &&   |
| Or                   |  \|\|  |
<br/>

## Functions
For declaring a function, you need to use the next syntax:
```ts
var example = funcion(<Argmuent name>, <Argmuent name>) {
    regresa <return value>;
};
```

binary search exapmle:
```dart
var binary_search = funcion(elements, val) {
    var left = 0;
    var rigth = largo(elements) - 1;
    var mid = 0;

	mientras(left <= rigth) {
		mid = (left + rigth) / 2;
		var mid_number = elements[mid];

		si (mid_number == val) {
			regresa mid;
		}

		si (mid_number < val) {
			left = mid + 1;
		} si_no {
			rigth = mid - 1;
		}
	}

	regresa -1
}

var numbers = lista[1,4,6,9,10,12,26];
var index = binary_search(numbers, 1);
escribir("numero encontrado en el indice ", index);
// output: 0;
```

## Lists
List allows you to group a list of data, 
lists are escential in any programming lengauge
```ts
var mi_lista = lista[2, 3, 4, "hello", "world"];
mi_lista[0] // output: 2
```

Also list have methods:
```js
mi_lista:agregar(5);   // add 5 to the list
mi_lista:pop();        // pop the last item and return it
mi_lista:popIndice(0); // remove by index and return it
```


## HashMaps
HashMaps are datastructures that help you store data by key => value
representation

For declaring a HashMap, you need to use the next syntax:
```ts
var example = mapa{key => value, key => value, key => value};

// get the value of the given key
example[key]
```

for example:
```dart
var mi_mapa = mapa{
    "a" => 1,
    "b" => 2,
    "c" => 3,
}

mi_mapa["a"] // output: 1
```

## Loops
WhileLoop syntax:
```ts
mientras(<condition>) {
    // code to be execute
}
```

for example:
```ts
var i = 0;
mientras(i <= 5) {
    escribir("hola mundo");
    i++;
}
```

For loop syntax:
```ts
por(i en <iterable>) {
    // code to be execute
}
```

for example:
```ts
por(i en rango(5)) {
    escribir("hola mundo");
}
```

you can also can iterate lists:
```ts
var mi_lista = lista[2,3,4];
por(i en mi_lista) {
    escribir(i);
}
```


you can see more examples in the examples folder.

## Usage
first copy the repository and change to the directory created:
```shell
$ git clone https://github.com/Haizza1/Katan && cd Katan
```
download the dependencies:
```shell
$ go mod download
```
check that tests pass:
```shell
$ go test -v ./...
```
compile the package:
```shell
$ go build -o katan
```

then you can create a file or play with the repl to play with the repl just run:
```shell
$ ./katan rpl
```

to use a file you can create a file with the .lpp extension and run:
```shell
$ ./katan file -path <path to your file>
```


## Contributions
Should you like to provide any feedback, please open up an Issue, I appreciate feedback and comments, although please keep in 
mind the project is incomplete, and I'm doing my best to keep it up to date.
