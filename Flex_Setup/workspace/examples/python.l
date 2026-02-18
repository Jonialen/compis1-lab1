%{
#include <stdio.h>
#include <string.h>

int codeLine = 1;      /* Current line number */
int position = 1;      /* Current column position */

void print_token(const char* type, const char* value) {
    printf("Token(%s, '%s', line=%d, pos=%d)\n", 
           type, value, codeLine, position);
}
%}

DIGIT       [0-9]
LETTER      [a-zA-Z_]
ID          {LETTER}({LETTER}|{DIGIT})*
NUMBER      {DIGIT}+
WHITESPACE  [ \t]

%%

"def"       { print_token("DEF", yytext); position += yyleng; }
"return"    { print_token("RETURN", yytext); position += yyleng; }
"if"        { print_token("IF", yytext); position += yyleng; }
"else"      { print_token("ELSE", yytext); position += yyleng; }

{ID}        { print_token("IDENTIFIER", yytext); position += yyleng; }
{NUMBER}    { print_token("NUMBER", yytext); position += yyleng; }

":"         { print_token("COLON", yytext); position += yyleng; }
"("         { print_token("LPAREN", yytext); position += yyleng; }
")"         { print_token("RPAREN", yytext); position += yyleng; }
","         { print_token("COMMA", yytext); position += yyleng; }
"+"         { print_token("PLUS", yytext); position += yyleng; }
"-"         { print_token("MINUS", yytext); position += yyleng; }
"*"         { print_token("TIMES", yytext); position += yyleng; }
"/"         { print_token("DIVIDE", yytext); position += yyleng; }
"="         { print_token("ASSIGN", yytext); position += yyleng; }

{WHITESPACE}+ { position += yyleng; /* Skip whitespace */ }

\n          { codeLine++; position = 1; /* New line */ }

.           { 
    printf("ERROR: Unknown character '%s' at line %d, position %d\n", 
           yytext, codeLine, position); 
    position++; 
}

%%

int yywrap() {
    return 1;  /* Stop after one file */
}

int main(int argc, char **argv) {
    /* If filename provided, read from file */
    if (argc > 1) {
        yyin = fopen(argv[1], "r");
        if (!yyin) {
            fprintf(stderr, "Cannot open file: %s\n", argv[1]);
            return 1;
        }
    }
    /* Otherwise read from stdin (keyboard) */
    
    printf("Scanning Python code...\n\n");
    yylex();  /* Start scanning */
    
    /* Close file if we opened one */
    if (argc > 1) {
        fclose(yyin);
    }
    
    return 0;
}