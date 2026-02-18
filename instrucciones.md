# Diseño de lenguajes de programación

**Febrero, 2026**

## Laboratorio 1. Diseño de un generador de análisis léxico

### Generadores de Analizadores Léxicos

Un generador de analizadores léxicos es una herramienta que automatiza la construcción de lexers a partir de especificaciones de alto nivel.

En lugar de programar manualmente cada token, el desarrollador escribe una especificación declarativa usando expresiones regulares (regex) para definir los patrones, y el generador produce automáticamente el código del analizador léxico.

El proceso de generación convierte las regex en un NFA, luego en un DFA mediante la construcción de subconjuntos, minimiza el DFA para optimizar estados, y finalmente genera código que implementa el autómata minimizado.

Los generadores modernos incluyen:

* Precedencia de tokens (resolución de ambigüedades).
* Acciones semánticas asociadas a cada patrón.
* Estados léxicos (por ejemplo, dentro de comentarios o strings).
* Reportes de errores informativos.

### Instrucciones

Lea cuidadosamente cada inciso y realizar lo solicitado.

* Para la parte de generación de código crear un repositorio en github.
* Incluya en su repositorio un video breve explicando sus ejemplos y lo solicitado en cada inciso.

---

### Problema 1: 50%

**Repaso conceptos:**

1. **Defina qué hace la función retractar dentro del proceso de reconocimiento de tokens. Explique:**
* En qué momento se invoca esta función.
* Qué sucede con el puntero del buffer al ejecutarse.
* Por qué es necesaria para aceptar correctamente un token.
* Cómo influye en el comportamiento del autómata finito.


2. **Defina que hace la función fallo en el contexto del análisis léxico. Explique:**
* Qué significa que un autómata "falle".
* Qué ocurre cuando no existe una transición válida para el carácter leído.
* Cómo permite esta función intentar el reconocimiento de otro patrón.
* Su relación con la estrategia de máxima coincidencia.


3. **Defina la función error dentro del proceso de análisis léxico. Su explicación debe incluir:**
* En qué situación se produce un error léxico.
* Qué ocurre cuando ningún autómata reconoce la secuencia de entrada.
* Cómo se reporta el error.
* Qué mecanismos de recuperación pueden aplicarse.


4. **Provea un ejemplo como el que se vió en clase.**

---

### Problema 2: 50%

**Actividades de Experimentación con Flex:**
Para este ejercicio clone el siguiente repositorio y realice lo solicitado: **Flex Setup Repository**

1. Cree un archivo de especificación Flex (`.l`) que reconozca identificadores válidos en Java (que comiencen con letra o guion bajo, seguidos de letras, dígitos o guiones bajos). Compile y pruebe con diferentes casos de prueba.
2. Implemente un lexer que reconozca y clasifique los diferentes tipos de literales numéricos: enteros, flotantes, notación científica, y números hexadecimales. Para cada tipo, el lexer debe imprimir el token y su valor.
3. Extienda su lexer para reconocer operadores aritméticos `(+, -, *, /)`, relacionales `(==, !=, <, >, <=, >=)`, y lógicos `(&&, !)`.
4. Agregue soporte para reconocer comentarios de una línea (`//`) y multilínea (`/**/`). Los comentarios deben ser descartados por el lexer y no generar tokens.
5. Implemente el reconocimiento de cadenas literales entre comillas dobles, manejando correctamente secuencias de escape como `\n`, `\t`, `\"`, y `\\`.
6. Realice pruebas en otro lenguaje de su elección, que no sea python ni java.
7. Explique que tendría que hacer si usted quisiera cambiar el color de algunas palabras (al estilo visual studio code).
8. Investigue otra herramienta similar a Flex y de una breve descripción de como funciona. Puede que más adelante tenga que realizar actividades con dicha herramienta.
