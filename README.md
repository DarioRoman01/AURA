# lpp programing lenguage

## What is this?
This is a very basic programming lenguage inspired in javascript, the goal is to create a programming lenguage in spanish
that helps people that are begining in software development or computer science, probiding a high lever programming lenguage
in spanish

## Syntax example

```
var a = 5;
var b = 5;
var c = a + b;

escribir(c); // prints 10
```

```
var suma = funcion(x, y) {
    regresa x + y;
};

// this will print 9
var resultado = suma(5, 4);
escribir(resultado); 
```

```
var mayor_de_edad = funcion(edad) {
    si (edad >= 18) {
        regresa verdadero;
    } si_no {
        regresa falso;
    }
}

var resultado = mayor_de_edad(18);
escribir(resultado); // prints verdadero

var resultado = mayor_de_edad(10);   
escribir(resultado); // prints falso

var edad = 30;
var resultado = mayor_de_edad(edad); 
escribir(resultado); // prints verdadero
```

you can see more examples in the examples folder.

## Usage
first copy the repository and change to the directory created:
```
$ git clone https://github.com/Haizza1/lpp && cd lpp
```
download the dependencies:
```
$ go mod download
```
check that tests pass:
```
$ go test -v ./...
```
compile the package:
```
$ go build -o ppl
```

then you can create a file or play with the repl to play with the repl just run:
```
$ ./ppl rpl
```

to use a file you can create a file with the .lpp extension and run:
```
$ ./ppl file -path <path to your file>
```


## Contributions
Should you like to provide any feedback, please open up an Issue, I appreciate feedback and comments, although please keep in 
mind the project is incomplete, and I'm doing my best to keep it up to date.
