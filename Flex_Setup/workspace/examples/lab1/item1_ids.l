%{
/*
 * item1_ids.l
 * Reconoce identificadores válidos de Java:
 *   - Comienzan con letra o guion bajo
 *   - Seguidos de letras, dígitos o guiones bajos
 *
 * Compilar:  flex item1_ids.l && gcc lex.yy.c -o item1 -lfl
 * Ejecutar:  ./item1 < test_ids.txt
 */
#include <stdio.h>
%}

LETTER   [a-zA-Z_]
DIGIT    [0-9]
ID       {LETTER}({LETTER}|{DIGIT})*

%%

{ID}          { printf("Token(IDENTIFIER, '%s')\n", yytext); }
[ \t\n\r]+    { /* whitespace — ignorar */ }
.             { printf("Token(INVALID,    '%s')\n", yytext); }

%%

int yywrap() { return 1; }

int main(int argc, char **argv) {
    if (argc > 1) {
        yyin = fopen(argv[1], "r");
        if (!yyin) { perror(argv[1]); return 1; }
    }
    yylex();
    if (argc > 1) fclose(yyin);
    return 0;
}
