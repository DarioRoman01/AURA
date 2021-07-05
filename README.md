# The Aura programing lenguage

## What is this?
This is a very basic programming lenguage inspired in javascript, the goal is to create a programming lenguage in spanish
that helps people that are begining in software development or computer science, probiding a high level programming lenguage
in spanish that is very simple to use. Also is an expressions based lenguage so is more easy to work with it

It is an interpreted language fully built using Go standar libreary.

## Syntax example

## Variables
```go
var a = 5;
var b = 5;
var c = a + b;

or 

a := 5;
b := 5;
c := a + b;
```

## Types
```ts
var number = 5; // Integer
var float = 2.5; // float
var str = "string"; // string
var bool = verdadero; // boolean
var list = lista[1,2,3] // list
var map = mapa{1 => "a", 2 => "b"}; // map
var null = nulo; // null
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

## Functions
For declaring a function, you need to use the next syntax:
```ts
funcion example(<Argmuent name>, <Argmuent name>) {
    regresa <return value>;
};
```

simple function example:
```ts
funcion sum(a, b) {
    regresa a + b;
}

escribir(add(5,8)); // output: 13
```

## Lists
List allows you to group a list of data, 
lists are escential in any programming lengauge
```ts
var mi_lista = lista[2, 3, 4, "hello", "world"];
mi_lista[0]; // output: 2
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
```go
example := mapa{key => value, key => value, key => value};

// get the value of the given key
example[key];
```

for example:
```go
mi_mapa := mapa{
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
```go
i := 0;
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

you can also can iterate lists and strings:
```ts
var mi_lista = lista[2,3,4];
por(i en mi_lista) {
    escribir(i);
}

por(i en "Hello world") {
    escribir(i);
}
```

## Clases
you can create clases with following syntax:
```dart
clase ClassName(<constructor params>) {
    method() {
        // method body
    }
}
```

for example:
```ts
clase Persona(name, age) {
    saludar() {
        escribir("hi im ", name, " i have ", age, " years old")
    }
}

p := nuevo Persona("eddy", 24);
p.saludar(); // output: hi im eddy i have 24 years old
```

with all this lets look a real world example with bynary search:
```ts
funcion binary_search(elements, val) {
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

numbers := lista[1,4,6,9,10,12,26];
index := binary_search(numbers, 1);
escribir("numero encontrado en el indice ", index);
// output: numero encontrado en el indice 0
```

## <h1>Installation</h1>
for using it you need to have Go install check https://golang.org/ for install Go

<h3>first copy the repository and change to the directory created:</h3>

```shell
$ git clone https://github.com/DarioRoman01/AURA.git && cd AURA
```

<h3>download the dependencies:</h3>

```shell
$ go mod download
```

<h3>check that tests pass:</h3>

```shell
$ go test -v ./...
```

<h3>if you use mac or linux just run the install script.
this will install aura in your system</h3>

```shell
$ chmod a+x install.sh
``` 
```shell
$ ./install.sh
```
<br>

<h2>if you are using windows or you want to install aura in other folder follow the next steps:</h2> 


<h3>compile aura:</h3>

```
go build -o aura
```

<h3>You can discover the install path by running the go list command, as in the following example</h3>

```shell
$ go list -f '{{.Target}}'
```
example output: /home/user/Go/bin/aura <br>

<h3>you can change the install target by setting the GOBIN variable using the go env command:</h3>

* on linux:
```shell
$ go env -w GOBIN=/path/to/your/bin
```
* on windows:
```powershell
$ go env -w GOBIN=C:\path\to\your\bin
```
<br> 

<h3>Add the Go install directory to your system's shell path</h3>

* on linux:
```shell
$ export PATH=$PATH:/path/to/your/install/directory
```

* on windows:
```powershell
$ set PATH=%PATH%;C:\path\to\your\install\directory
```

<br>

<h3>Once you've updated the shell path, run the go install command to compile and install the package.
then you need to go where the package was install and rename to binary aura to aura.exe</h3>

```shell
$ go install
```

<h3>then you can create a file or play with the repl. 
to play with the repl just run:</h3>

```shell
$ aura
```

<h3>to use a file you can create a file with the .aura extension and run:</h3> 

**Is important to have the .aura extension otherwise the lenguage wont read the file**
```shell
$ aura some/file.aura
```


## Contributions
Should you like to provide any feedback, please open up an Issue, I appreciate feedback and comments, although please keep in 
mind the project is incomplete, and I'm doing my best to keep it up to date.
