# Como correr los lexers

## Flex (items 1-5) — requiere Docker

```bash
# Levantar el contenedor
docker compose up -d --build

# Entrar al contenedor
docker compose exec flex bash

# Dentro del contenedor, ir a la carpeta
cd /workspace/examples/lab1
```

### Item 1 — Identificadores Java
```bash
flex item1_ids.l && gcc lex.yy.c -o item1 -lfl
./item1 < test_ids.txt
```

### Items 2-5 — Lexer completo
```bash
flex lexer.l && gcc lex.yy.c -o lexer -lfl
./lexer < test_all.txt
```

> Para probar con otro archivo: `./lexer < mi_archivo.txt`

---

## Go (item 6) — requiere Go instalado

```bash
cd Flex_Setup/workspace/examples/go_lexer

# Correr con el ejemplo embebido
go run main.go

# Correr con un archivo propio
go run main.go ruta/al/archivo.go
```
