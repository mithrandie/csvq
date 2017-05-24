%{
package parser
%}

%union{
    program     []Statement
    statement   Statement
    expression  Expression
    expressions []Expression
    identifier  Identifier
    text        String
    integer     Integer
    float       Float
    ternary     Ternary
    null        Null
    token       Token
}

%type<program>     program
%type<statement>   statement
%type<expression>  expression
%type<expression>  select_query
%type<expression>  select_clause
%type<expression>  from_clause
%type<expression>  where_clause
%type<expression>  group_by_clause
%type<expression>  having_clause
%type<expression>  order_by_clause
%type<expression>  limit_clause
%type<expression>  primary
%type<expression>  value
%type<expression>  order_item
%type<expression>  subquery
%type<expression>  string_operation
%type<expression>  comparison
%type<expression>  arithmetic
%type<expression>  logic
%type<expression>  filter
%type<expression>  function
%type<expression>  option
%type<expression>  table
%type<expression>  join
%type<expression>  join_condition
%type<expression>  field_object
%type<expression>  field
%type<expression>  case
%type<expression>  case_value
%type<expression>  case_else
%type<expressions> values
%type<expressions> filters
%type<expressions> order_items
%type<expressions> tables
%type<expressions> using_fields
%type<expressions> fields
%type<expressions> case_when
%type<identifier>  identifier
%type<text>        text
%type<integer>     integer
%type<float>       float
%type<ternary>     ternary
%type<null>        null
%type<token>       distinct
%type<token>       negation
%type<token>       order_direction
%type<token>       join_inner
%type<token>       join_outer
%type<token>       join_direction
%type<token>       statement_terminal
%token<token> IDENTIFIER STRING INTEGER FLOAT BOOLEAN TERNARY DATETIME
%token<token> SELECT FROM UPDATE SET DELETE WHERE INSERT INTO VALUES AS DUAL
%token<token> CREATE DROP ALTER TABLE COLUMN
%token<token> ORDER GROUP HAVING BY ASC DESC LIMIT
%token<token> JOIN INNER OUTER LEFT RIGHT FULL CROSS ON USING NATURAL UNION
%token<token> ALL ANY EXISTS IN
%token<token> AND OR NOT BETWEEN LIKE IS NULL
%token<token> DISTINCT WITH
%token<token> TRUE FALSE UNKNOWN
%token<token> CASE WHEN THEN ELSE END
%token<token> COMPARISON_OP STRING_OP

%left OR
%left AND
%left NOT
%left COMPARISON_OP STRING_OP
%left '+' '-'
%left '*' '/'

%%

program
    :
    {
        $$ = nil
        yylex.(*Lexer).program = $$
    }
    | statement program
    {
        $$ = append([]Statement{$1}, $2...)
        yylex.(*Lexer).program = $$
    }

statement
    : expression statement_terminal
    {
        $$ = $1
    }

expression
    : select_query
    {
        $$ = $1
    }

select_query
    : select_clause from_clause where_clause group_by_clause having_clause order_by_clause limit_clause
    {
        $$ = SelectQuery{
            SelectClause:  $1,
            FromClause:    $2,
            WhereClause:   $3,
            GroupByClause: $4,
            HavingClause:  $5,
            OrderByClause: $6,
            LimitClause:   $7,
        }
    }

select_clause
    : SELECT distinct fields
    {
        $$ = SelectClause{Select: $1.Literal, Distinct: $2, Fields: $3}
    }

from_clause
    :
    {
        $$ = nil
    }
    | FROM DUAL
    {
        $$ = FromClause{From: $1.Literal, Tables: []Expression{Dual{Dual: $2.Literal}}}
    }
    | FROM tables
    {
        $$ = FromClause{From: $1.Literal, Tables: $2}
    }

where_clause
    :
    {
        $$ = nil
    }
    | WHERE filter
    {
        $$ = WhereClause{Where: $1.Literal, Filter: $2}
    }

group_by_clause
    :
    {
        $$ = nil
    }
    | GROUP BY values
    {
        $$ = GroupByClause{GroupBy: $1.Literal + " " + $2.Literal, Items: $3}
    }

having_clause
    :
    {
        $$ = nil
    }
    | HAVING filter
    {
        $$ = HavingClause{Having: $1.Literal, Filter: $2}
    }

order_by_clause
    :
    {
        $$ = nil
    }
    | ORDER BY order_items
    {
        $$ = OrderByClause{OrderBy: $1.Literal + " " + $2.Literal, Items: $3}
    }

limit_clause
    :
    {
        $$ = nil
    }
    | LIMIT integer
    {
        $$ = LimitClause{Limit: $1.Literal, Number: $2.Value()}
    }

primary
    : identifier
    {
        $$ = $1
    }
    | text
    {
        $$ = $1
    }
    | integer
    {
        $$ = $1
    }
    | float
    {
        $$ = $1
    }
    | ternary
    {
        $$ = $1
    }
    | DATETIME
    {
        $$ = NewDatetimeFromString($1.Literal)
    }
    | null
    {
        $$ = $1
    }

value
    : primary
    {
        $$ = $1
    }
    | arithmetic
    {
        $$ = $1
    }
    | string_operation
    {
        $$ = $1
    }
    | subquery
    {
        $$ = $1
    }
    | function
    {
        $$ = $1
    }
    | case
    {
        $$ = $1
    }
    | '(' value ')'
    {
        $$ = Parentheses{Expr: $2}
    }

order_item
    : value order_direction
    {
        $$ = OrderItem{Item: $1, Direction: $2}
    }

order_direction
    :
    {
        $$ = Token{}
    }
    | ASC
    {
        $$ = $1
    }
    | DESC
    {
        $$ = $1
    }

subquery
    : '(' select_query ')'
    {
        $$ = Subquery{Query: $2.(SelectQuery)}
    }

string_operation
    : value STRING_OP value
    {
        var item1 []Expression
        var item2 []Expression

        c1, ok := $1.(Concat)
        if ok {
            item1 = c1.Items
        } else {
            item1 = []Expression{$1}
        }

        c2, ok := $3.(Concat)
        if ok {
            item2 = c2.Items
        } else {
            item2 = []Expression{$3}
        }

        $$ = Concat{Items: append(item1, item2...)}
    }

comparison
    : value COMPARISON_OP value
    {
        $$ = Comparison{LHS: $1, Operator: $2, RHS: $3}
    }
    | value IS negation ternary
    {
        $$ = Is{Is: $2.Literal, LHS: $1, RHS: $4, Negation: $3}
    }
    | value IS negation null
    {
        $$ = Is{Is: $2.Literal, LHS: $1, RHS: $4, Negation: $3}
    }
    | value negation BETWEEN value AND value
    {
        $$ = Between{Between: $3.Literal, And: $5.Literal, LHS: $1, Low: $4, High: $6, Negation: $2}
    }
    | value negation IN '(' values ')'
    {
        $$ = In{In: $3.Literal, LHS: $1, List: $5, Negation: $2}
    }
    | value negation IN subquery
    {
        $$ = In{In: $3.Literal, LHS: $1, Query: $4.(Subquery), Negation: $2}
    }
    | value negation LIKE value
    {
        $$ = Like{Like: $3.Literal, LHS: $1, Pattern: $4, Negation: $2}
    }
    | value COMPARISON_OP ANY subquery
    {
        $$ = Any{Any: $3.Literal, LHS: $1, Operator: $2, Query: $4.(Subquery)}
    }
    | value COMPARISON_OP ALL subquery
    {
        $$ = All{All: $3.Literal, LHS: $1, Operator: $2, Query: $4.(Subquery)}
    }
    | EXISTS subquery
    {
        $$ = Exists{Exists: $1.Literal, Query: $2.(Subquery)}
    }

arithmetic
    : value '+' value
    {
        $$ = Arithmetic{LHS: $1, Operator: int('+'), RHS: $3}
    }
    | value '-' value
    {
        $$ = Arithmetic{LHS: $1, Operator: int('-'), RHS: $3}
    }
    | value '*' value
    {
        $$ = Arithmetic{LHS: $1, Operator: int('*'), RHS: $3}
    }
    | value '/' value
    {
        $$ = Arithmetic{LHS: $1, Operator: int('/'), RHS: $3}
    }
    | value '%' value
    {
        $$ = Arithmetic{LHS: $1, Operator: int('%'), RHS: $3}
    }

logic
    : filter OR filter
    {
        $$ = Logic{LHS: $1, Operator: $2, RHS: $3}
    }
    | filter AND filter
    {
        $$ = Logic{LHS: $1, Operator: $2, RHS: $3}
    }
    | NOT filter
    {
        $$ = Logic{LHS: nil, Operator: $1, RHS: $2}
    }

filter
    : value
    {
        $$ = $1
    }
    | comparison
    {
        $$ = $1
    }
    | logic
    {
        $$ = $1
    }
    | '(' filter ')'
    {
        $$ = Parentheses{Expr: $2}
    }

function
    : identifier '(' option ')'
    {
        $$ = Function{Name: $1.Literal, Option: $3.(Option)}
    }

option
    :
    {
        $$ = Option{}
    }
    | distinct '*'
    {
        $$ = Option{Distinct: $1, Args: []Expression{AllColumns{}}}
    }
    | distinct filters
    {
        $$ = Option{Distinct: $1, Args: $2}
    }

table
    : identifier
    {
        $$ = Table{Object: $1}
    }
    | identifier identifier
    {
        $$ = Table{Object: $1, Alias: $2}
    }
    | identifier AS identifier
    {
        $$ = Table{Object: $1, As: $2, Alias: $3}
    }
    | subquery
    {
        $$ = Table{Object: $1}
    }
    | subquery identifier
    {
        $$ = Table{Object: $1, Alias: $2}
    }
    | subquery AS identifier
    {
        $$ = Table{Object: $1, As: $2, Alias: $3}
    }
    | join
    {
        $$ = Table{Object: $1}
    }

join
    : table join_inner JOIN table join_condition
    {
        $$ = Join{Join: $3.Literal, Table: $1.(Table), JoinTable: $4.(Table), Natural: Token{}, JoinType: $2, Condition: $5}
	}
    | table NATURAL join_inner JOIN table
    {
        $$ = Join{Join: $4.Literal, Table: $1.(Table), JoinTable: $5.(Table), Natural: $2, JoinType: $3, Condition: nil}
	}
    | table join_direction join_outer JOIN table join_condition
    {
        $$ = Join{Join: $4.Literal, Table: $1.(Table), JoinTable: $5.(Table), Natural: Token{}, JoinType: $3, Direction: $2, Condition: $6}
    }
    | table NATURAL join_direction join_outer JOIN table
    {
        $$ = Join{Join: $5.Literal, Table: $1.(Table), JoinTable: $6.(Table), Natural: $2, JoinType: $4, Direction: $3, Condition: nil}
    }
    | table CROSS JOIN table
    {
        $$ = Join{Join: $3.Literal, Table: $1.(Table), JoinTable: $4.(Table), Natural: Token{}, JoinType: $2, Condition: nil}
    }

join_condition
    :
    {
        $$ = nil
    }
    | ON filter
    {
        $$ = JoinCondition{Literal:$1.Literal, On: $2}
    }
    | USING '(' using_fields ')'
    {
        $$ = JoinCondition{Literal:$1.Literal, Using: $3}
    }

field_object
    : filter
    {
        $$ = $1
    }
    | '*'
    {
        $$ = AllColumns{}
    }

field
    : field_object
    {
        $$ = Field{Object: $1}
    }
    | field_object AS identifier
    {
        $$ = Field{Object: $1, As: $2, Alias: $3}
    }

case
    : CASE case_value case_when case_else END
    {
        $$ = Case{Case: $1.Literal, End: $5.Literal, Value: $2, When: $3, Else: $4}
    }

case_value
    :
    {
        $$ = nil
    }
    | value
    {
        $$ = $1
    }

case_else
    :
    {
        $$ = nil
    }
    | ELSE value
    {
        $$ = CaseElse{Else: $1.Literal, Result: $2}
    }

values
    : value
    {
        $$ = []Expression{$1}
    }
    | value ',' values
    {
        $$ = append([]Expression{$1}, $3...)
    }

filters
    : filter
    {
        $$ = []Expression{$1}
    }
    | filter ',' filters
    {
        $$ = append([]Expression{$1}, $3...)
    }

order_items
    : order_item
    {
        $$ = []Expression{$1}
    }
    | order_item ',' order_items
    {
        $$ = append([]Expression{$1}, $3...)
    }

tables
    : table
    {
        $$ = []Expression{$1}
    }
    | table ',' tables
    {
        $$ = append([]Expression{$1}, $3...)
    }

using_fields
    : identifier
    {
        $$ = []Expression{$1}
    }
    | identifier ',' using_fields
    {
        $$ = append([]Expression{$1}, $3...)
    }

fields
    : field
    {
        $$ = []Expression{$1}
    }
    | field ',' fields
    {
        $$ = append([]Expression{$1}, $3...)
    }

case_when
    : WHEN filter THEN value
    {
        $$ = []Expression{CaseWhen{When: $1.Literal, Then: $3.Literal, Condition: $2, Result: $4}}
    }
    | case_when case_when
    {
        $$ = append($1, $2...)
    }

identifier
    : IDENTIFIER
    {
        $$ = Identifier{Literal: $1.Literal, Quoted: $1.Quoted}
    }

text
    : STRING
    {
        $$ = NewString($1.Literal)
    }

integer
    : INTEGER
    {
        $$ = NewIntegerFromString($1.Literal)
    }
    | '-' integer
    {
        i := $2.Value() * -1
        $$ = NewInteger(i)
    }

float
    : FLOAT
    {
        $$ = NewFloatFromString($1.Literal)
    }
    | '-' float
    {
        f := $2.Value() * -1
        $$ = NewFloat(f)
    }

ternary
    : TERNARY
    {
        $$ = NewTernaryFromString($1.Literal)
    }

null
    : NULL
    {
        $$ = NewNullFromString($1.Literal)
    }

distinct
    :
    {
        $$ = Token{}
    }
    | DISTINCT
    {
        $$ = $1
    }

negation
    :
    {
        $$ = Token{}
    }
    | NOT
    {
        $$ = $1
    }

join_inner
    :
    {
        $$ = Token{}
    }
    | INNER
    {
        $$ = $1
    }

join_outer
    :
    {
        $$ = Token{}
    }
    | OUTER
    {
        $$ = $1
    }

join_direction
    :
    {
        $$ = Token{}
    }
    | LEFT
    {
        $$ = $1
    }
    | RIGHT
    {
        $$ = $1
    }
    | FULL
    {
        $$ = $1
    }

statement_terminal
    :
    {
        $$ = Token{}
    }
    | ';'
    {
        $$ = Token{Token: ';', Literal: string(';')}
    }

%%

func SetDebugLevel(level int, verbose bool) {
	yyDebug        = level
	yyErrorVerbose = verbose
}

func Parse(s string) ([]Statement, error) {
    l := new(Lexer)
    l.Init(s)
    yyParse(l)
    return l.program, l.err
}