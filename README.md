# The Katan programing lenguage

## What is this?
This is a very basic programming lenguage inspired in javascript, the goal is to create a programming lenguage in spanish
that helps people that are begining in software development or computer science, probiding a high level programming lenguage
in spanish that is very simple to use

## Syntax example

## Variables
```ts
var a = 5;
var b = 5;
var c = a + b;

escribir(c); // prints 10
```

## Functions

```ts
var suma = funcion(x, y) {
    regresa x + y;
};

// this will print 9
var resultado = suma(5, 4);
escribir(resultado); 
```

```ts 
var mayor_de_edad = funcion(edad) {
    si (edad >= 18) {
        regresa verdadero;
    } si_no {
        regresa falso;
    }
}

var resultado = mayor_de_edad(18);
escribir(resultado); // prints verdadero

resultado = mayor_de_edad(10);   
escribir(resultado); // prints falso

var edad = recibir_entero();
resultado = mayor_de_edad(edad); 
escribir(resultado);
```

## Lists
```ts
var mi_lista = lista[2,3,4];
mi_lista:agregar(4); // add 4 to the list
mi_lista:pop(); // pop the last item
mi_lista:popIndice(0); // remove by index
```

## HashMaps
```dart
var mi_mapa = mapa{
    "a" => 1,
    "b" => 2,
    "c" => 3,
}

mi_mapa["a"] // return 1
```

## Loops
```ts
var i = 0;
mientras(i <= 5) {
    escribir("hola mundo");
    i += 1;
}
```

```ts
por(i en rango(5)) {
    escribir("hola mundo");
}

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
