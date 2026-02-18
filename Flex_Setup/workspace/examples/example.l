/* example.l - A simple flex example that counts words, lines, and characters */

%{
int chars = 0;
int words = 0;
int lines = 0;
%}

%%

[a-zA-Z]+   { words++; chars += yyleng; }
\n          { chars++; lines++; }
.           { chars++; }

%%

int main(int argc, char **argv) {
    yylex();
    printf("Lines: %d, Words: %d, Characters: %d\n", lines, words, chars);
    return 0;
}

int yywrap() {
    return 1;
}
