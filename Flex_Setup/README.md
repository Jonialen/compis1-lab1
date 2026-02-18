# Flex Docker Environment

Docker environment for [flex](https://github.com/westes/flex) (Fast Lexical Analyzer) built from source.

## Getting Started

### 1. Build and start the container

```bash
docker compose up -d --build
```

### 2. Enter the container

```bash
docker compose exec flex bash
```

### 3. Run the example

```bash
cd /workspace/examples
flex example.l
gcc lex.yy.c -o scanner
echo "Hello World from flex!" | ./scanner
```

Output:

```
Lines: 1, Words: 4, Characters: 23
```

### 4. Stop the container

```bash
docker compose down
```

Inside the container, navigate to your folder:

```bash
cd /workspace/examples/tu-carpeta
flex tu-archivo.l
gcc lex.yy.c -o scanner
./scanner < input.txt
```

## Useful flex Commands

| Command                   | Description                                 |
| ------------------------- | ------------------------------------------- |
| `flex file.l`             | Generate lex.yy.c                           |
| `flex -d file.l`          | Enable debug mode (prints each token match) |
| `flex -v file.l`          | Verbose output (shows DFA stats)            |
| `flex -o output.c file.l` | Custom output filename                      |

## File Structure

```
flex-docker/
├── Dockerfile
├── docker-compose.yml
├── README.md
├── .gitignore
└── workspace/
    └── examples/
        └── example.l
```

# Python Function Scanner

A Flex-based lexical analyzer for scanning Python function syntax.

### python.l

The Flex specification file consists of three sections:

**Section 1: Declarations** (between `%{` and `%}`)

-   C includes and global variables
-   Example: `int codeLine = 1;`

**Section 2: Rules** (between `%%` markers)

-   Pattern-action pairs
-   Example: `"def" { print_token("DEF", yytext); }`

**Section 3: User Code** (after second `%%`)

-   Helper functions
-   `yywrap()` function (required)
-   `main()` function

## How to Build

```bash
# Generate C code from Flex specification
flex python.l

# Compile the generated code
gcc lex.yy.c -o scanner -lfl

# Run the scanner
./scanner test.py
```

## Input

Text file containing Python code:

```python
def add(x, y):
    return x + y
```

## Output

```
Token(DEF, 'def', line=1, pos=1)
Token(IDENTIFIER, 'add', line=1, pos=5)
Token(LPAREN, '(', line=1, pos=8)
...
```

## Key Variables

-   `yytext` - Matched string
-   `yyleng` - Length of matched string
-   `codeLine` - Current line number
-   `position` - Current character position

## Notes

-   Generated files (`lex.yy.c`, `scanner`) are excluded from git
-   Only commit your `.l` source files
-   The container includes `gcc` for compiling the generated C code
