# Sistema de Tipado en Aura

Aura ahora incluye soporte para anotaciones de tipo opcionales en funciones, proporcionando verificación de tipos en tiempo de ejecución.

## Tipos Disponibles

- `Entero` - Números enteros (integers)
- `Texto` - Cadenas de texto (strings)
- `Booleano` - Valores booleanos (true/false)
- `Flotante` - Números de punto flotante (floats)
- `Lista` - Listas/arrays
- `Mapa` - Diccionarios/mapas
- `Nulo` - Valor nulo

## Sintaxis

### Parámetros Tipados

```dart
funcion suma(a: Entero, b: Entero) {
    regresa a + b;
}
```

### Tipo de Retorno

```dart
funcion obtenerSaludo(nombre: Texto) -> Texto {
    regresa "Hola, " + nombre;
}
```

### Funciones Flecha con Tipos

```dart
multiplicar := |x: Entero, y: Entero| -> Entero => x * y;
```

### Métodos de Clase con Tipos

```dart
clase Calculadora(valor: Entero) {
    sumar(x: Entero) -> Entero {
        regresa valor + x;
    }
}
```

## Verificación de Tipos

El sistema realiza verificación de tipos en tiempo de ejecución:

### Verificación de Parámetros

```dart
funcion suma(a: Entero, b: Entero) {
    regresa a + b;
}

suma(5, "texto");  // Error: el parametro 'b' espera tipo Entero pero recibio texto
```

### Verificación de Retorno

```dart
funcion obtenerNumero() -> Entero {
    regresa "no es un numero";  // Error: la funcion debe regresar tipo Entero pero regreso texto
}
```

## Ejemplos

Ver los archivos en la carpeta `examples/`:
- `examples/typing.aura` - Ejemplos completos de uso del sistema de tipado
- `examples/typing_errors.aura` - Ejemplos de detección de errores de tipo

## Características

- ✅ Anotaciones de tipo opcionales (no son obligatorias)
- ✅ Verificación de tipos en parámetros de función
- ✅ Verificación de tipos en valores de retorno
- ✅ Soporte para funciones normales y funciones flecha
- ✅ Soporte para métodos de clase
- ✅ Mensajes de error descriptivos en español
- ✅ Compatible con código existente sin tipos

## Notas

- Las anotaciones de tipo son **opcionales** - el código sin tipos sigue funcionando
- La verificación se realiza en **tiempo de ejecución**, no en tiempo de compilación
- Los tipos se escriben con mayúscula inicial (Entero, Texto, etc.)
- El operador `->` se usa para especificar el tipo de retorno
- El operador `:` se usa para especificar el tipo de parámetros
