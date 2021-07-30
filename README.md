# The Aura programing lenguage

## What is this?
This is a very basic programming lenguage inspired in javascript, python and golang, the goal is to create a programming lenguage in spanish
that helps people that are begining in software development or computer science, probiding a high level programming lenguage
in spanish that is very simple to use. Also is an expressions based lenguage so is more easy to work with it

It is an interpreted language fully built using Go standar libreary.

## Code Snippet
```dart
// insertion sort implementation in aura

funcion insertion_sort(elements) {
    por(i en rango(1, largo(elements))) {
        anchor := elements[i];
        j := i - 1;
        mientras(j >= 0 && anchor < elements[j]) {
            elements[j + 1] = elements[j];
            j--;
        }

        elements[j + 1] = anchor;
    }

    escribirF("Array ordernado {}", elements)
}

funcion main() {
    tests := lista[
        lista[11,9,29,7,2,15,28],
        lista[3, 7, 9, 11],
        lista[25, 22, 21, 10],
        lista[29, 15, 28]
    ];

    tests:porCada(|test| => insertion_sort(test));
}
```

## <h1>Installation</h1>
1. ### Go to Realeses and download the version that prefer **if you are using windows you need to download the .exe file**

2. <h3>move the binary to a folder of your preference for example:</h3>

    ```shell
    mv aura /path/to/your/install/directory
    ```

3. <h3>Then you have to set aura in your path:</h3>

    * on linux:
    ```shell
    export PATH=$PATH:/path/to/your/install/directory
    ``` 

    * on windows:
    ```powershell
    set PATH=%PATH%;C:\path\to\your\install\directory
    ```

<h3>then you can create a file or play with the repl. 
to play with the repl just run:</h3>

```shell
$ aura
```

<h3>to use a file you can create a file with the .aura extension and run:</h3> 

**Is important to have the .aura extension otherwise the lenguage wont read the file**
```shell
$ aura file.aura
```


## Contributions
Should you like to provide any feedback, please open up an Issue, I appreciate feedback and comments, although please keep in 
mind the project is incomplete, and I'm doing my best to keep it up to date.
