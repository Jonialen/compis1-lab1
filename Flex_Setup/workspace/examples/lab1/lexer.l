%{
#include <stdio.h>
%}

DIGIT       [0-9]
HEXDIGIT    [0-9a-fA-F]
LETTER      [a-zA-Z_]
EXP         [eE][+-]?{DIGIT}+

%x BLOCK_COMMENT

%%

"//"[^\n]*          { printf("Token(COMMENT_LINE,  '%s')\n", yytext); }

"/*"                { printf("Token(COMMENT_BLOCK, '"); BEGIN(BLOCK_COMMENT); }
<BLOCK_COMMENT>"*/" { printf("')\n");                  BEGIN(INITIAL); }
<BLOCK_COMMENT>\n   { printf("\n"); }
<BLOCK_COMMENT>.    { printf("%s", yytext); }

\"(\\.|[^\"\\])*\"  { printf("Token(STRING,        '%s')\n", yytext); }

0[xX]{HEXDIGIT}+                { printf("Token(HEX,        '%s')\n", yytext); }
{DIGIT}+"."{DIGIT}*{EXP}        { printf("Token(SCIENTIFIC, '%s')\n", yytext); }
"."{DIGIT}+{EXP}                { printf("Token(SCIENTIFIC, '%s')\n", yytext); }
{DIGIT}+{EXP}                   { printf("Token(SCIENTIFIC, '%s')\n", yytext); }
{DIGIT}+"."{DIGIT}*             { printf("Token(FLOAT,      '%s')\n", yytext); }
"."{DIGIT}+                     { printf("Token(FLOAT,      '%s')\n", yytext); }
{DIGIT}+                        { printf("Token(INT,        '%s')\n", yytext); }

"&&"    { printf("Token(AND,   '%s')\n", yytext); }
"||"    { printf("Token(OR,    '%s')\n", yytext); }
"!"     { printf("Token(NOT,   '%s')\n", yytext); }
"=="    { printf("Token(EQ,    '%s')\n", yytext); }
"!="    { printf("Token(NEQ,   '%s')\n", yytext); }
"<="    { printf("Token(LEQ,   '%s')\n", yytext); }
">="    { printf("Token(GEQ,   '%s')\n", yytext); }
"<"     { printf("Token(LT,    '%s')\n", yytext); }
">"     { printf("Token(GT,    '%s')\n", yytext); }
"+"     { printf("Token(PLUS,  '%s')\n", yytext); }
"-"     { printf("Token(MINUS, '%s')\n", yytext); }
"*"     { printf("Token(MUL,   '%s')\n", yytext); }
"/"     { printf("Token(DIV,   '%s')\n", yytext); }

{LETTER}({LETTER}|{DIGIT})*  { printf("Token(IDENT,      '%s')\n", yytext); }

[ \t\n\r]+  { /* whitespace -- ignorar */ }
.           { printf("Token(UNKNOWN,    '%s')\n", yytext); }

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
