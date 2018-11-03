%{
package json

import "strconv"
%}

%union{
    structure      Structure
    structures     []Structure
    object_member  ObjectMember
    object_members []ObjectMember
    token          Token
}

%type<structure>      structure
%type<structure>      container
%type<object_member>  object_member
%type<object_members> object_members
%type<structures>     array_items
%type<structure>      value

%token<token> NUMBER STRING BOOLEAN NULL

%%

structure
    :
    {
        $$ = nil
        yylex.(*Lexer).structure = $$
    }
    | container
    {
        $$ = $1
        yylex.(*Lexer).structure = $$
    }

container
    : '{' object_members '}'
    {
        $$ = Object{Members: $2}
    }
    | '[' array_items ']'
    {
        $$ = Array($2)
    }

object_member
    : STRING ':' value
    {
        $$ = ObjectMember{Key: $1.Literal, Value: $3}
    }

object_members
    :
    {
        $$ = nil
    }
    | object_member
    {
        $$ = []ObjectMember{$1}
    }
    | object_member ',' object_members
    {
        $$ = append([]ObjectMember{$1}, $3...)
    }

array_items
    :
    {
        $$ = []Structure{}
    }
    | value
    {
        $$ = []Structure{$1}
    }
    | value ',' array_items
    {
        $$ = append([]Structure{$1}, $3...)
    }

value
    : container
    {
        $$ = $1
    }
    | STRING
    {
        $$ = String($1.Literal)
    }
    | NUMBER
    {
        f, _ := strconv.ParseFloat($1.Literal, 64)
        $$ = Number(f)
    }
    | BOOLEAN
    {
        b, _ := strconv.ParseBool($1.Literal)
        $$ = Boolean(b)
    }
    | NULL
    {
        $$ = Null{}
    }

%%

func ParseJson(src string) (Structure, EscapeType, error) {
	l := new(Lexer)
	l.Init(src)
	yyParse(l)
	return l.structure, l.EscapeType(), l.err
}