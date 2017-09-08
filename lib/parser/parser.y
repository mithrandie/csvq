%{
package parser

import "github.com/mithrandie/csvq/lib/value"
%}

%union{
    program     []Statement
    statement   Statement
    queryexpr   QueryExpression
    queryexprs  []QueryExpression
    expression  Expression
    expressions []Expression
    identifier  Identifier
    table       Table
    variable    Variable
    variables   []Variable
    varassign   VariableAssignment
    varassigns  []VariableAssignment
    updateset   UpdateSet
    updatesets  []UpdateSet
    columndef   ColumnDefault
    columndefs  []ColumnDefault
    elseif      []ElseIf
    elseexpr    Else
    casewhen    []CaseWhen
    caseelse    CaseElse
    fetchpos    FetchPosition
    token       Token
}

%type<program>     program
%type<program>     loop_program
%type<program>     function_program
%type<program>     function_loop_program
%type<statement>   common_statement
%type<statement>   common_loop_flow_control_statement
%type<statement>   procedure_statement
%type<statement>   while_statement
%type<statement>   exit_statement
%type<statement>   flow_control_statement
%type<statement>   loop_statement
%type<statement>   loop_flow_control_statement
%type<statement>   function_statement
%type<statement>   function_while_statement
%type<statement>   function_exit_statement
%type<statement>   function_loop_statement
%type<statement>   function_flow_control_statement
%type<statement>   function_loop_flow_control_statement
%type<statement>   variable_statement
%type<statement>   transaction_statement
%type<statement>   table_operation_statement
%type<columndef>   column_default
%type<columndefs>  column_defaults
%type<expression>  column_position
%type<statement>   cursor_statement
%type<statement>   temporary_table_statement
%type<varassign>   parameter
%type<varassigns>  parameters
%type<varassign>   optional_parameter
%type<varassigns>  optional_parameters
%type<varassigns>  function_parameters
%type<statement>   user_defined_function_statement
%type<fetchpos>    fetch_position
%type<queryexpr>   cursor_status
%type<statement>   command_statement
%type<statement>   trigger_statement
%type<queryexpr>   select_query
%type<queryexpr>   select_entity
%type<queryexpr>   select_set_entity
%type<queryexpr>   select_clause
%type<queryexpr>   from_clause
%type<queryexpr>   where_clause
%type<queryexpr>   group_by_clause
%type<queryexpr>   having_clause
%type<queryexpr>   order_by_clause
%type<queryexpr>   limit_clause
%type<queryexpr>   limit_with
%type<queryexpr>   offset_clause
%type<queryexpr>   with_clause
%type<queryexpr>   inline_table
%type<queryexprs>  inline_tables
%type<queryexpr>   primitive_type
%type<queryexpr>   ternary
%type<queryexpr>   null
%type<queryexpr>   field_reference
%type<queryexpr>   value
%type<queryexpr>   wildcard
%type<queryexpr>   row_value
%type<queryexprs>  row_values
%type<queryexprs>  order_items
%type<queryexpr>   order_item
%type<queryexpr>   order_value
%type<token>       order_direction
%type<token>       order_null_position
%type<queryexpr>   subquery
%type<queryexpr>   string_operation
%type<queryexpr>   comparison
%type<queryexpr>   arithmetic
%type<queryexpr>   logic
%type<queryexprs>  arguments
%type<queryexpr>   function
%type<queryexpr>   aggregate_function
%type<queryexpr>   listagg
%type<queryexpr>   analytic_function
%type<queryexpr>   analytic_clause
%type<queryexpr>   partition
%type<queryexpr>   table_identifier
%type<table>       identified_table
%type<queryexprs>  operate_tables
%type<queryexpr>   virtual_table_object
%type<queryexpr>   table
%type<queryexpr>   join
%type<queryexpr>   join_condition
%type<queryexpr>   field_object
%type<queryexpr>   field
%type<queryexpr>   case_expr
%type<queryexpr>   case_value
%type<queryexprs>  case_expr_when
%type<queryexpr>   case_expr_else
%type<queryexprs>  field_references
%type<queryexprs>  values
%type<queryexprs>  tables
%type<queryexprs>  identifiers
%type<queryexprs>  fields
%type<expression>  insert_query
%type<expression>  update_query
%type<updateset>   update_set
%type<updatesets>  update_set_list
%type<expression>  delete_query
%type<elseif>      elseif
%type<elseexpr>    else
%type<elseif>      in_loop_elseif
%type<elseexpr>    in_loop_else
%type<elseif>      in_function_elseif
%type<elseexpr>    in_function_else
%type<elseif>      in_function_in_loop_elseif
%type<elseexpr>    in_function_in_loop_else
%type<casewhen>    case_when
%type<caseelse>    case_else
%type<casewhen>    in_loop_case_when
%type<caseelse>    in_loop_case_else
%type<casewhen>    in_function_case_when
%type<caseelse>    in_function_case_else
%type<casewhen>    in_function_in_loop_case_when
%type<caseelse>    in_function_in_loop_case_else
%type<identifier>  identifier
%type<variable>    variable
%type<variables>   variables
%type<queryexpr>   variable_substitution
%type<varassign>   variable_assignment
%type<varassigns>  variable_assignments
%type<token>       distinct
%type<token>       negation
%type<token>       join_type_inner
%type<token>       join_type_outer
%type<token>       join_outer_direction
%type<token>       all
%type<token>       recursive
%type<token>       as
%type<token>       comparison_operator

%token<token> IDENTIFIER STRING INTEGER FLOAT BOOLEAN TERNARY DATETIME VARIABLE FLAG
%token<token> SELECT FROM UPDATE SET DELETE WHERE INSERT INTO VALUES AS DUAL STDIN
%token<token> RECURSIVE
%token<token> CREATE ADD DROP ALTER TABLE FIRST LAST AFTER BEFORE DEFAULT RENAME TO
%token<token> ORDER GROUP HAVING BY ASC DESC LIMIT OFFSET PERCENT
%token<token> JOIN INNER OUTER LEFT RIGHT FULL CROSS ON USING NATURAL
%token<token> UNION INTERSECT EXCEPT
%token<token> ALL ANY EXISTS IN
%token<token> AND OR NOT BETWEEN LIKE IS NULL
%token<token> DISTINCT WITH
%token<token> CASE IF ELSEIF WHILE WHEN THEN ELSE DO END
%token<token> DECLARE CURSOR FOR FETCH OPEN CLOSE DISPOSE
%token<token> NEXT PRIOR ABSOLUTE RELATIVE RANGE
%token<token> SEPARATOR PARTITION OVER
%token<token> COMMIT ROLLBACK
%token<token> CONTINUE BREAK EXIT
%token<token> PRINT PRINTF SOURCE TRIGGER
%token<token> FUNCTION AGGREGATE BEGIN RETURN
%token<token> IGNORE WITHIN
%token<token> VAR
%token<token> TIES NULLS
%token<token> ERROR
%token<token> COUNT LISTAGG
%token<token> AGGREGATE_FUNCTION FUNCTION_WITH_INS
%token<token> COMPARISON_OP STRING_OP SUBSTITUTION_OP
%token<token> UMINUS UPLUS
%token<token> ';' '*' '=' '-' '+' '!' '(' ')'

%right SUBSTITUTION_OP
%left UNION EXCEPT
%left INTERSECT
%left CROSS FULL NATURAL JOIN
%left OR
%left AND
%right NOT
%nonassoc '=' COMPARISON_OP IS BETWEEN IN LIKE
%left STRING_OP
%left '+' '-'
%left '*' '/' '%'
%right UMINUS UPLUS '!'

%%

program
    :
    {
        $$ = nil
        yylex.(*Lexer).program = $$
    }
    | procedure_statement
    {
        $$ = []Statement{$1}
        yylex.(*Lexer).program = $$
    }
    | procedure_statement ';' program
    {
        $$ = append([]Statement{$1}, $3...)
        yylex.(*Lexer).program = $$
    }

loop_program
    :
    {
        $$ = nil
    }
    | loop_statement ';' loop_program
    {
        $$ = append([]Statement{$1}, $3...)
    }

function_program
    :
    {
        $$ = nil
    }
    | function_statement ';' function_program
    {
        $$ = append([]Statement{$1}, $3...)
    }

function_loop_program
    :
    {
        $$ = nil
    }
    | function_loop_statement ';' function_loop_program
    {
        $$ = append([]Statement{$1}, $3...)
    }

common_statement
    : select_query
    {
        $$ = $1
    }
    | insert_query
    {
        $$ = $1
    }
    | update_query
    {
        $$ = $1
    }
    | delete_query
    {
        $$ = $1
    }
    | function
    {
        $$ = $1
    }
    | table_operation_statement
    {
        $$ = $1
    }
    | variable_statement
    {
        $$ = $1
    }
    | cursor_statement
    {
        $$ = $1
    }
    | temporary_table_statement
    {
        $$ = $1
    }
    | user_defined_function_statement
    {
        $$ = $1
    }
    | transaction_statement
    {
        $$ = $1
    }
    | command_statement
    {
        $$ = $1
    }
    | trigger_statement
    {
        $$ = $1
    }

common_loop_flow_control_statement
    : CONTINUE
    {
        $$ = FlowControl{Token: $1.Token}
    }
    | BREAK
    {
        $$ = FlowControl{Token: $1.Token}
    }

procedure_statement
    : common_statement
    {
        $$ = $1
    }
    | flow_control_statement
    {
        $$ = $1
    }

while_statement
    : WHILE value DO loop_program END WHILE
    {
        $$ = While{Condition: $2, Statements: $4}
    }
    | WHILE variable IN identifier DO loop_program END WHILE
    {
        $$ = WhileInCursor{Variables: []Variable{$2}, Cursor: $4, Statements: $6}
    }
    | WHILE variables IN identifier DO loop_program END WHILE
    {
        $$ = WhileInCursor{Variables: $2, Cursor: $4, Statements: $6}
    }

exit_statement
    : EXIT
    {
        $$ = Exit{}
    }
    | EXIT INTEGER
    {
        $$ = Exit{Code: value.NewIntegerFromString($2.Literal)}
    }

loop_statement
    : common_statement
    {
        $$ = $1
    }
    | loop_flow_control_statement
    {
        $$ = $1
    }

flow_control_statement
    : IF value THEN program else END IF
    {
        $$ = If{Condition: $2, Statements: $4, Else: $5}
    }
    | IF value THEN program elseif else END IF
    {
        $$ = If{Condition: $2, Statements: $4, ElseIf: $5, Else: $6}
    }
    | CASE case_value case_when case_else END CASE
    {
        $$ = Case{Value: $2, When: $3, Else: $4}
    }
    | while_statement
    {
        $$ = $1
    }
    | exit_statement
    {
        $$ = $1
    }

loop_flow_control_statement
    : IF value THEN loop_program in_loop_else END IF
    {
        $$ = If{Condition: $2, Statements: $4, Else: $5}
    }
    | IF value THEN loop_program in_loop_elseif in_loop_else END IF
    {
        $$ = If{Condition: $2, Statements: $4, ElseIf: $5, Else: $6}
    }
    | CASE case_value in_loop_case_when in_loop_case_else END CASE
    {
        $$ = Case{Value: $2, When: $3, Else: $4}
    }
    | while_statement
    {
        $$ = $1
    }
    | exit_statement
    {
        $$ = $1
    }
    | common_loop_flow_control_statement
    {
        $$ = $1
    }

function_statement
    : common_statement
    {
        $$ = $1
    }
    | function_flow_control_statement
    {
        $$ = $1
    }

function_while_statement
    : WHILE value DO function_loop_program END WHILE
    {
        $$ = While{Condition: $2, Statements: $4}
    }
    | WHILE variable IN identifier DO function_loop_program END WHILE
    {
        $$ = WhileInCursor{Variables: []Variable{$2}, Cursor: $4, Statements: $6}
    }
    | WHILE variables IN identifier DO function_loop_program END WHILE
    {
        $$ = WhileInCursor{Variables: $2, Cursor: $4, Statements: $6}
    }

function_exit_statement
    : RETURN
    {
        $$ = Return{Value: NewNullValue()}
    }
    | RETURN value
    {
        $$ = Return{Value: $2}
    }

function_loop_statement
    : common_statement
    {
        $$ = $1
    }
    | function_loop_flow_control_statement
    {
        $$ = $1
    }

function_flow_control_statement
    : IF value THEN function_program in_function_else END IF
    {
        $$ = If{Condition: $2, Statements: $4, Else: $5}
    }
    | IF value THEN function_program in_function_elseif in_function_else END IF
    {
        $$ = If{Condition: $2, Statements: $4, ElseIf: $5, Else: $6}
    }
    | CASE case_value in_function_case_when in_function_case_else END CASE
    {
        $$ = Case{Value: $2, When: $3, Else: $4}
    }
    | function_while_statement
    {
        $$ = $1
    }
    | function_exit_statement
    {
        $$ = $1
    }

function_loop_flow_control_statement
    : IF value THEN function_loop_program in_function_in_loop_else END IF
    {
        $$ = If{Condition: $2, Statements: $4, Else: $5}
    }
    | IF value THEN function_loop_program in_function_in_loop_elseif in_function_in_loop_else END IF
    {
        $$ = If{Condition: $2, Statements: $4, ElseIf: $5, Else: $6}
    }
    | CASE case_value in_function_in_loop_case_when in_function_in_loop_case_else END CASE
    {
        $$ = Case{Value: $2, When: $3, Else: $4}
    }
    | function_while_statement
    {
        $$ = $1
    }
    | function_exit_statement
    {
        $$ = $1
    }
    | common_loop_flow_control_statement
    {
        $$ = $1
    }

variable_statement
    : VAR variable_assignments
    {
        $$ = VariableDeclaration{Assignments:$2}
    }
    | DECLARE variable_assignments
    {
        $$ = VariableDeclaration{Assignments:$2}
    }
    | variable_substitution
    {
        $$ = $1
    }
    | DISPOSE variable
    {
        $$ = DisposeVariable{Variable:$2}
    }

transaction_statement
    : COMMIT
    {
        $$ = TransactionControl{BaseExpr: NewBaseExpr($1), Token: $1.Token}
    }
    | ROLLBACK
    {
        $$ = TransactionControl{BaseExpr: NewBaseExpr($1), Token: $1.Token}
    }

table_operation_statement
    : CREATE TABLE identifier '(' identifiers ')'
    {
        $$ = CreateTable{Table: $3, Fields: $5}
    }
    | CREATE TABLE identifier '(' identifiers ')' as select_query
    {
        $$ = CreateTable{Table: $3, Fields: $5, Query: $8}
    }
    | CREATE TABLE identifier as select_query
    {
        $$ = CreateTable{Table: $3, Query: $5}
    }
    | ALTER TABLE table_identifier ADD column_default column_position
    {
        $$ = AddColumns{Table: $3, Columns: []ColumnDefault{$5}, Position: $6}
    }
    | ALTER TABLE table_identifier ADD '(' column_defaults ')' column_position
    {
        $$ = AddColumns{Table: $3, Columns: $6, Position: $8}
    }
    | ALTER TABLE table_identifier DROP field_reference
    {
        $$ = DropColumns{Table: $3, Columns: []QueryExpression{$5}}
    }
    | ALTER TABLE table_identifier DROP '(' field_references ')'
    {
        $$ = DropColumns{Table: $3, Columns: $6}
    }
    | ALTER TABLE table_identifier RENAME field_reference TO identifier
    {
        $$ = RenameColumn{Table: $3, Old: $5, New: $7}
    }

column_default
    : identifier
    {
        $$ = ColumnDefault{Column: $1}
    }
    | identifier DEFAULT value
    {
        $$ = ColumnDefault{Column: $1, Value: $3}
    }

column_defaults
    : column_default
    {
        $$ = []ColumnDefault{$1}
    }
    | column_default ',' column_defaults
    {
        $$ = append([]ColumnDefault{$1}, $3...)
    }

column_position
    :
    {
        $$ = nil
    }
    | FIRST
    {
        $$ = ColumnPosition{Position: $1}
    }
    | LAST
    {
        $$ = ColumnPosition{Position: $1}
    }
    | AFTER field_reference
    {
        $$ = ColumnPosition{Position: $1, Column: $2}
    }
    | BEFORE field_reference
    {
        $$ = ColumnPosition{Position: $1, Column: $2}
    }

cursor_statement
    : DECLARE identifier CURSOR FOR select_query
    {
        $$ = CursorDeclaration{Cursor:$2, Query: $5.(SelectQuery)}
    }
    | OPEN identifier
    {
        $$ = OpenCursor{Cursor: $2}
    }
    | CLOSE identifier
    {
        $$ = CloseCursor{Cursor: $2}
    }
    | DISPOSE CURSOR identifier
    {
        $$ = DisposeCursor{Cursor: $3}
    }
    | FETCH fetch_position identifier INTO variables
    {
        $$ = FetchCursor{Position: $2, Cursor: $3, Variables: $5}
    }

temporary_table_statement
    : DECLARE identifier TABLE '(' identifiers ')'
    {
        $$ = TableDeclaration{Table: $2, Fields: $5}
    }
    | DECLARE identifier TABLE '(' identifiers ')' AS select_query
    {
        $$ = TableDeclaration{Table: $2, Fields: $5, Query: $8}
    }
    | DECLARE identifier TABLE AS select_query
    {
        $$ = TableDeclaration{Table: $2, Query: $5}
    }
    | DISPOSE TABLE identifier
    {
        $$ = DisposeTable{Table: $3}
    }

parameter
    : variable
    {
        $$ = VariableAssignment{Variable:$1}
    }

parameters
    : parameter
    {
        $$ = []VariableAssignment{$1}
    }
    | parameters ',' parameter
    {
        $$ = append($1, $3)
    }

optional_parameter
    : variable DEFAULT value
    {
        $$ = VariableAssignment{Variable: $1, Value: $3}
    }

optional_parameters
    : optional_parameter
    {
        $$ = []VariableAssignment{$1}
    }
    | optional_parameter ',' optional_parameters
    {
        $$ = append([]VariableAssignment{$1}, $3...)
    }

function_parameters
    : parameters
    {
        $$ = $1
    }
    | optional_parameters
    {
        $$ = $1
    }
    | parameters ',' optional_parameters
    {
        $$ = append($1, $3...)
    }

user_defined_function_statement
    : DECLARE identifier FUNCTION '(' ')' AS BEGIN function_program END
    {
        $$ = FunctionDeclaration{Name: $2, Statements: $8}
    }
    | DECLARE identifier FUNCTION '(' function_parameters ')' AS BEGIN function_program END
    {
        $$ = FunctionDeclaration{Name: $2, Parameters: $5, Statements: $9}
    }
    | DECLARE identifier AGGREGATE '(' identifier ')' AS BEGIN function_program END
    {
        $$ = AggregateDeclaration{Name: $2, Cursor: $5, Statements: $9}
    }
    | DECLARE identifier AGGREGATE '(' identifier ',' function_parameters ')' AS BEGIN function_program END
    {
        $$ = AggregateDeclaration{Name: $2, Cursor: $5, Parameters: $7, Statements: $11}
    }

fetch_position
    :
    {
        $$ = FetchPosition{}
    }
    | NEXT
    {
        $$ = FetchPosition{Position: $1}
    }
    | PRIOR
    {
        $$ = FetchPosition{Position: $1}
    }
    | FIRST
    {
        $$ = FetchPosition{Position: $1}
    }
    | LAST
    {
        $$ = FetchPosition{Position: $1}
    }
    | ABSOLUTE value
    {
        $$ = FetchPosition{BaseExpr: NewBaseExpr($1), Position: $1, Number: $2}
    }
    | RELATIVE value
    {
        $$ = FetchPosition{BaseExpr: NewBaseExpr($1), Position: $1, Number: $2}
    }

cursor_status
    : CURSOR identifier IS negation OPEN
    {
        $$ = CursorStatus{CursorLit: $1.Literal, Cursor: $2, Is: $3.Literal, Negation: $4, Type: $5.Token, TypeLit: $5.Literal}
    }
    | CURSOR identifier IS negation IN RANGE
    {
        $$ = CursorStatus{CursorLit: $1.Literal, Cursor: $2, Is: $3.Literal, Negation: $4, Type: $6.Token, TypeLit: $5.Literal + " " + $6.Literal}
    }
    | CURSOR identifier COUNT
    {
        $$ = CursorAttrebute{CursorLit: $1.Literal, Cursor: $2, Attrebute: $3}
    }

command_statement
    : SET FLAG '=' primitive_type
    {
        $$ = SetFlag{BaseExpr: NewBaseExpr($1), Name: $2.Literal, Value: $4.(PrimitiveType).Value}
    }
    | PRINT value
    {
        $$ = Print{Value: $2}
    }
    | PRINTF value
    {
        $$ = Printf{BaseExpr: NewBaseExpr($1), Format: $2}
    }
    | PRINTF value ',' values
    {
        $$ = Printf{BaseExpr: NewBaseExpr($1), Format: $2, Values: $4}
    }
    | SOURCE value
    {
        $$ = Source{BaseExpr: NewBaseExpr($1), FilePath: $2}
    }

trigger_statement
    : TRIGGER ERROR
    {
        $$ = Trigger{BaseExpr: NewBaseExpr($1), Token: $2.Token}
    }
    | TRIGGER ERROR value
    {
        $$ = Trigger{BaseExpr: NewBaseExpr($1), Token: $2.Token, Message: $3}
    }
    | TRIGGER ERROR INTEGER value
    {
        $$ = Trigger{BaseExpr: NewBaseExpr($1), Token: $2.Token, Message: $4, Code: value.NewIntegerFromString($3.Literal)}
    }

select_query
    : with_clause select_entity order_by_clause limit_clause offset_clause
    {
        $$ = SelectQuery{
            WithClause:    $1,
            SelectEntity:  $2,
            OrderByClause: $3,
            LimitClause:   $4,
            OffsetClause:  $5,
        }
    }

select_entity
    : select_clause from_clause where_clause group_by_clause having_clause
    {
        $$ = SelectEntity{
            SelectClause:  $1,
            FromClause:    $2,
            WhereClause:   $3,
            GroupByClause: $4,
            HavingClause:  $5,
        }
    }
    | select_set_entity UNION all select_set_entity
    {
        $$ = SelectSet{
            LHS:      $1,
            Operator: $2,
            All:      $3,
            RHS:      $4,
        }
    }
    | select_set_entity INTERSECT all select_set_entity
    {
        $$ = SelectSet{
            LHS:      $1,
            Operator: $2,
            All:      $3,
            RHS:      $4,
        }
    }
    | select_set_entity EXCEPT all select_set_entity
    {
        $$ = SelectSet{
            LHS:      $1,
            Operator: $2,
            All:      $3,
            RHS:      $4,
        }
    }

select_set_entity
    : select_entity
    {
        $$ = $1
    }
    | subquery
    {
        $$ = $1
    }

select_clause
    : SELECT distinct fields
    {
        $$ = SelectClause{BaseExpr: NewBaseExpr($1), Select: $1.Literal, Distinct: $2, Fields: $3}
    }

from_clause
    :
    {
        $$ = nil
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
    | WHERE value
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
    | HAVING value
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
    | LIMIT value limit_with
    {
        $$ = LimitClause{BaseExpr: NewBaseExpr($1), Limit: $1.Literal, Value: $2, With: $3}
    }
    | LIMIT value PERCENT limit_with
    {
        $$ = LimitClause{BaseExpr: NewBaseExpr($1), Limit: $1.Literal, Value: $2, Percent: $3.Literal, With: $4}
    }

limit_with
    :
    {
        $$ = nil
    }
    | WITH TIES
    {
        $$ = LimitWith{With: $1.Literal, Type: $2}
    }

offset_clause
    :
    {
        $$ = nil
    }
    | OFFSET value
    {
        $$ = OffsetClause{BaseExpr: NewBaseExpr($1), Offset: $1.Literal, Value: $2}
    }

with_clause
    :
    {
        $$ = nil
    }
    | WITH inline_tables
    {
        $$ = WithClause{With: $1.Literal, InlineTables: $2}
    }

inline_table
    : recursive identifier AS '(' select_query ')'
    {
        $$ = InlineTable{Recursive: $1, Name: $2, As: $3.Literal, Query: $5.(SelectQuery)}
    }
    | recursive identifier '(' identifiers ')' AS '(' select_query ')'
    {
        $$ = InlineTable{Recursive: $1, Name: $2, Fields: $4, As: $6.Literal, Query: $8.(SelectQuery)}
    }

inline_tables
    : inline_table
    {
        $$ = []QueryExpression{$1}
    }
    | inline_table ',' inline_tables
    {
        $$ = append([]QueryExpression{$1}, $3...)
    }

primitive_type
    : STRING
    {
        $$ = NewStringValue($1.Literal)
    }
    | INTEGER
    {
        $$ = NewIntegerValueFromString($1.Literal)
    }
    | FLOAT
    {
        $$ = NewFloatValueFromString($1.Literal)
    }
    | ternary
    {
        $$ = $1
    }
    | DATETIME
    {
        $$ = NewDatetimeValueFromString($1.Literal)
    }
    | null
    {
        $$ = $1
    }

ternary
    : TERNARY
    {
        $$ = NewTernaryValueFromString($1.Literal)
    }

null
    : NULL
    {
        $$ = NewNullValueFromString($1.Literal)
    }

field_reference
    : identifier
    {
        $$ = FieldReference{BaseExpr: $1.BaseExpr, Column: $1}
    }
    | identifier '.' identifier
    {
        $$ = FieldReference{BaseExpr: $1.BaseExpr, View: $1, Column: $3}
    }
    | STDIN '.' identifier
    {
        $$ = FieldReference{BaseExpr: NewBaseExpr($1), View: Identifier{BaseExpr: NewBaseExpr($1), Literal: $1.Literal}, Column: $3}
    }
    | identifier '.' INTEGER
    {
        $$ = ColumnNumber{BaseExpr: $1.BaseExpr, View: $1, Number: value.NewIntegerFromString($3.Literal)}
    }
    | STDIN '.' INTEGER
    {
        $$ = ColumnNumber{BaseExpr: NewBaseExpr($1), View: Identifier{BaseExpr: NewBaseExpr($1), Literal: $1.Literal}, Number: value.NewIntegerFromString($3.Literal)}
    }

value
    : field_reference
    {
        $$ = $1
    }
    | primitive_type
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
    | aggregate_function
    {
        $$ = $1
    }
    | case_expr
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
    | variable
    {
        $$ = $1
    }
    | variable_substitution
    {
        $$ = $1
    }
    | cursor_status
    {
        $$ = $1
    }
    | '(' value ')'
    {
        $$ = Parentheses{Expr: $2}
    }

wildcard
    : '*'
    {
        $$ = AllColumns{BaseExpr: NewBaseExpr($1)}
    }

row_value
    : '(' values ')'
    {
        $$ = RowValue{BaseExpr: NewBaseExpr($1), Value: ValueList{Values: $2}}
    }
    | subquery
    {
        $$ = RowValue{BaseExpr: $1.GetBaseExpr(), Value: $1}
    }

row_values
    : row_value
    {
        $$ = []QueryExpression{$1}
    }
    | row_value ',' row_values
    {
        $$ = append([]QueryExpression{$1}, $3...)
    }

order_items
    : order_item
    {
        $$ = []QueryExpression{$1}
    }
    | order_item ',' order_items
    {
        $$ = append([]QueryExpression{$1}, $3...)
    }

order_item
    : order_value order_direction
    {
        $$ = OrderItem{Value: $1, Direction: $2}
    }
    | order_value order_direction NULLS order_null_position
    {
        $$ = OrderItem{Value: $1, Direction: $2, Nulls: $3.Literal, Position: $4}
    }

order_value
    : value
    {
        $$ = $1
    }
    | analytic_function
    {
        $$ = $1
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

order_null_position
    : FIRST
    {
        $$ = $1
    }
    | LAST
    {
        $$ = $1
    }

subquery
    : '(' select_query ')'
    {
        $$ = Subquery{BaseExpr: NewBaseExpr($1), Query: $2.(SelectQuery)}
    }

string_operation
    : value STRING_OP value
    {
        var item1 []QueryExpression
        var item2 []QueryExpression

        c1, ok := $1.(Concat)
        if ok {
            item1 = c1.Items
        } else {
            item1 = []QueryExpression{$1}
        }

        c2, ok := $3.(Concat)
        if ok {
            item2 = c2.Items
        } else {
            item2 = []QueryExpression{$3}
        }

        $$ = Concat{Items: append(item1, item2...)}
    }

comparison
    : value COMPARISON_OP value
    {
        $$ = Comparison{LHS: $1, Operator: $2.Literal, RHS: $3}
    }
    | row_value COMPARISON_OP row_value
    {
        $$ = Comparison{LHS: $1, Operator: $2.Literal, RHS: $3}
    }
    | value '=' value
    {
        $$ = Comparison{LHS: $1, Operator: "=", RHS: $3}
    }
    | row_value '=' row_value
    {
        $$ = Comparison{LHS: $1, Operator: "=", RHS: $3}
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
    | row_value negation BETWEEN row_value AND row_value
    {
        $$ = Between{Between: $3.Literal, And: $5.Literal, LHS: $1, Low: $4, High: $6, Negation: $2}
    }
    | value negation IN row_value
    {
        $$ = In{In: $3.Literal, LHS: $1, Values: $4, Negation: $2}
    }
    | row_value negation IN '(' row_values ')'
    {
        $$ = In{In: $3.Literal, LHS: $1, Values: RowValueList{RowValues: $5}, Negation: $2}
    }
    | row_value negation IN subquery
    {
        $$ = In{In: $3.Literal, LHS: $1, Values: $4, Negation: $2}
    }
    | value negation LIKE value
    {
        $$ = Like{Like: $3.Literal, LHS: $1, Pattern: $4, Negation: $2}
    }
    | value comparison_operator ANY row_value
    {
        $$ = Any{Any: $3.Literal, LHS: $1, Operator: $2.Literal, Values: $4}
    }
    | row_value comparison_operator ANY '(' row_values ')'
    {
        $$ = Any{Any: $3.Literal, LHS: $1, Operator: $2.Literal, Values: RowValueList{RowValues: $5}}
    }
    | row_value comparison_operator ANY subquery
    {
        $$ = Any{Any: $3.Literal, LHS: $1, Operator: $2.Literal, Values: $4}
    }
    | value comparison_operator ALL row_value
    {
        $$ = All{All: $3.Literal, LHS: $1, Operator: $2.Literal, Values: $4}
    }
    | row_value comparison_operator ALL '(' row_values ')'
    {
        $$ = All{All: $3.Literal, LHS: $1, Operator: $2.Literal, Values: RowValueList{RowValues: $5}}
    }
    | row_value comparison_operator ALL subquery
    {
        $$ = All{All: $3.Literal, LHS: $1, Operator: $2.Literal, Values: $4}
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
    | '-' value %prec UMINUS
    {
        $$ = UnaryArithmetic{Operand: $2, Operator: $1}
    }
    | '+' value %prec UPLUS
    {
        $$ = UnaryArithmetic{Operand: $2, Operator: $1}
    }

logic
    : value OR value
    {
        $$ = Logic{LHS: $1, Operator: $2, RHS: $3}
    }
    | value AND value
    {
        $$ = Logic{LHS: $1, Operator: $2, RHS: $3}
    }
    | NOT value
    {
        $$ = UnaryLogic{Operand: $2, Operator: $1}
    }
    | '!' value
    {
        $$ = UnaryLogic{Operand: $2, Operator: $1}
    }

arguments
    :
    {
        $$ = nil
    }
    | values
    {
        $$ = $1
    }

function
    : identifier '(' arguments ')'
    {
        $$ = Function{BaseExpr: $1.BaseExpr, Name: $1.Literal, Args: $3}
    }
    | IF '(' arguments ')'
    {
        $$ = Function{BaseExpr: NewBaseExpr($1), Name: $1.Literal, Args: $3}
    }


aggregate_function
    : identifier '(' distinct arguments ')'
    {
        $$ = AggregateFunction{BaseExpr: $1.BaseExpr, Name: $1.Literal, Distinct: $3, Args: $4}
    }
    | AGGREGATE_FUNCTION '(' distinct arguments ')'
    {
        $$ = AggregateFunction{BaseExpr: NewBaseExpr($1), Name: $1.Literal, Distinct: $3, Args: $4}
    }
    | COUNT '(' distinct arguments ')'
    {
        $$ = AggregateFunction{BaseExpr: NewBaseExpr($1), Name: $1.Literal, Distinct: $3, Args: $4}
    }
    | COUNT '(' distinct wildcard ')'
    {
        $$ = AggregateFunction{BaseExpr: NewBaseExpr($1), Name: $1.Literal, Distinct: $3, Args: []QueryExpression{$4}}
    }
    | listagg
    {
        $$ = $1
    }

listagg
    : LISTAGG '(' distinct arguments ')'
    {
        $$ = ListAgg{BaseExpr: NewBaseExpr($1), ListAgg: $1.Literal, Distinct: $3, Args: $4}
    }
    | LISTAGG '(' distinct arguments ')' WITHIN GROUP '(' order_by_clause ')'
    {
        $$ = ListAgg{BaseExpr: NewBaseExpr($1), ListAgg: $1.Literal, Distinct: $3, Args: $4, WithinGroup: $6.Literal + " " + $7.Literal, OrderBy: $9}
    }

analytic_function
    : identifier '(' arguments ')' OVER '(' analytic_clause ')'
    {
        $$ = AnalyticFunction{BaseExpr: $1.BaseExpr, Name: $1.Literal, Args: $3, Over: $5.Literal, AnalyticClause: $7.(AnalyticClause)}
    }
    | identifier '(' distinct arguments ')' OVER '(' analytic_clause ')'
    {
        $$ = AnalyticFunction{BaseExpr: $1.BaseExpr, Name: $1.Literal, Distinct: $3, Args: $4, Over: $6.Literal, AnalyticClause: $8.(AnalyticClause)}
    }
    | AGGREGATE_FUNCTION '(' distinct arguments ')' OVER '(' analytic_clause ')'
    {
        $$ = AnalyticFunction{BaseExpr: NewBaseExpr($1), Name: $1.Literal, Distinct: $3, Args: $4, Over: $6.Literal, AnalyticClause: $8.(AnalyticClause)}
    }
    | COUNT '(' distinct arguments ')' OVER '(' analytic_clause ')'
    {
        $$ = AnalyticFunction{BaseExpr: NewBaseExpr($1), Name: $1.Literal, Distinct: $3, Args: $4, Over: $6.Literal, AnalyticClause: $8.(AnalyticClause)}
    }
    | COUNT '(' distinct wildcard ')' OVER '(' analytic_clause ')'
    {
        $$ = AnalyticFunction{BaseExpr: NewBaseExpr($1), Name: $1.Literal, Distinct: $3, Args: []QueryExpression{$4}, Over: $6.Literal, AnalyticClause: $8.(AnalyticClause)}
    }
    | LISTAGG '(' distinct arguments ')' OVER '(' analytic_clause ')'
    {
        $$ = AnalyticFunction{BaseExpr: NewBaseExpr($1), Name: $1.Literal, Distinct: $3, Args: $4, Over: $6.Literal, AnalyticClause: $8.(AnalyticClause)}
    }
    | FUNCTION_WITH_INS '(' arguments ')' OVER '(' analytic_clause ')'
    {
        $$ = AnalyticFunction{BaseExpr: NewBaseExpr($1), Name: $1.Literal, Args: $3, Over: $5.Literal, AnalyticClause: $7.(AnalyticClause)}
    }
    | FUNCTION_WITH_INS '(' arguments ')' IGNORE NULLS OVER '(' analytic_clause ')'
    {
        $$ = AnalyticFunction{BaseExpr: NewBaseExpr($1), Name: $1.Literal, Args: $3, IgnoreNulls: true, IgnoreNullsLit: $5.Literal + " " + $6.Literal, Over: $7.Literal, AnalyticClause: $9.(AnalyticClause)}
    }

analytic_clause
    : partition order_by_clause
    {
        $$ = AnalyticClause{Partition: $1, OrderByClause: $2}
    }

partition
    :
    {
        $$ = nil
    }
    | PARTITION BY values
    {
        $$ = Partition{PartitionBy: $1.Literal + " " + $2.Literal, Values: $3}
    }

table_identifier
    : identifier
    {
        $$ = $1
    }
    | STDIN
    {
        $$ = Stdin{BaseExpr: NewBaseExpr($1), Stdin: $1.Literal}
    }

identified_table
    : table_identifier
    {
        $$ = Table{Object: $1}
    }
    | table_identifier identifier
    {
        $$ = Table{Object: $1, Alias: $2}
    }
    | table_identifier AS identifier
    {
        $$ = Table{Object: $1, As: $2.Literal, Alias: $3}
    }

virtual_table_object
    : subquery
    {
        $$ = $1
    }

table
    : identified_table
    {
        $$ = $1
    }
    | virtual_table_object
    {
        $$ = Table{Object: $1}
    }
    | virtual_table_object identifier
    {
        $$ = Table{Object: $1, Alias: $2}
    }
    | virtual_table_object AS identifier
    {
        $$ = Table{Object: $1, As: $2.Literal, Alias: $3}
    }
    | join
    {
        $$ = Table{Object: $1}
    }
    | DUAL
    {
        $$ = Table{Object: Dual{Dual: $1.Literal}}
    }
    | '(' table ')'
    {
        $$ = Parentheses{Expr: $2}
    }

join
    : table CROSS JOIN table
    {
        $$ = Join{Join: $3.Literal, Table: $1, JoinTable: $4, JoinType: $2, Condition: nil}
    }
    | table join_type_inner JOIN table join_condition
    {
        $$ = Join{Join: $3.Literal, Table: $1, JoinTable: $4, JoinType: $2, Condition: $5}
    }
    | table join_outer_direction join_type_outer JOIN table join_condition
    {
        $$ = Join{Join: $4.Literal, Table: $1, JoinTable: $5, JoinType: $3, Direction: $2, Condition: $6}
    }
    | table FULL join_type_outer JOIN table ON value
    {
        $$ = Join{Join: $4.Literal, Table: $1, JoinTable: $5, JoinType: $3, Direction: $2, Condition: JoinCondition{Literal:$6.Literal, On: $7}}
    }
    | table NATURAL join_type_inner JOIN table
    {
        $$ = Join{Join: $4.Literal, Table: $1, JoinTable: $5, JoinType: $3, Natural: $2}
    }
    | table NATURAL join_outer_direction join_type_outer JOIN table
    {
        $$ = Join{Join: $5.Literal, Table: $1, JoinTable: $6, JoinType: $4, Direction: $3, Natural: $2}
    }

join_condition
    : ON value
    {
        $$ = JoinCondition{Literal:$1.Literal, On: $2}
    }
    | USING '(' identifiers ')'
    {
        $$ = JoinCondition{Literal:$1.Literal, Using: $3}
    }

field_object
    : value
    {
        $$ = $1
    }
    | analytic_function
    {
        $$ = $1
    }

field
    : field_object
    {
        $$ = Field{Object: $1}
    }
    | field_object AS identifier
    {
        $$ = Field{Object: $1, As: $2.Literal, Alias: $3}
    }
    | wildcard
    {
        $$ = Field{Object: $1}
    }

case_expr
    : CASE case_value case_expr_when case_expr_else END
    {
        $$ = CaseExpr{Case: $1.Literal, End: $5.Literal, Value: $2, When: $3, Else: $4}
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

case_expr_when
    : WHEN value THEN value
    {
        $$ = []QueryExpression{CaseExprWhen{When: $1.Literal, Then: $3.Literal, Condition: $2, Result: $4}}
    }
    | WHEN value THEN value case_expr_when
    {
        $$ = append([]QueryExpression{CaseExprWhen{When: $1.Literal, Then: $3.Literal, Condition: $2, Result: $4}}, $5...)
    }

case_expr_else
    :
    {
        $$ = nil
    }
    | ELSE value
    {
        $$ = CaseExprElse{Else: $1.Literal, Result: $2}
    }

field_references
    : field_reference
    {
        $$ = []QueryExpression{$1}
    }
    | field_reference ',' field_references
    {
        $$ = append([]QueryExpression{$1}, $3...)
    }

values
    : value
    {
        $$ = []QueryExpression{$1}
    }
    | value ',' values
    {
        $$ = append([]QueryExpression{$1}, $3...)
    }

tables
    : table
    {
        $$ = []QueryExpression{$1}
    }
    | table ',' tables
    {
        $$ = append([]QueryExpression{$1}, $3...)
    }

operate_tables
    : table_identifier
    {
        $$ = []QueryExpression{Table{Object: $1}}
    }
    | table_identifier ',' operate_tables
    {
        $$ = append([]QueryExpression{Table{Object: $1}}, $3...)
    }

identifiers
    : identifier
    {
        $$ = []QueryExpression{$1}
    }
    | identifier ',' identifiers
    {
        $$ = append([]QueryExpression{$1}, $3...)
    }

fields
    : field
    {
        $$ = []QueryExpression{$1}
    }
    | field ',' fields
    {
        $$ = append([]QueryExpression{$1}, $3...)
    }

insert_query
    : with_clause INSERT INTO identified_table VALUES row_values
    {
        $$ = InsertQuery{WithClause: $1, Table: $4, ValuesList: $6}
    }
    | with_clause INSERT INTO identified_table '(' field_references ')' VALUES row_values
    {
        $$ = InsertQuery{WithClause: $1, Table: $4, Fields: $6, ValuesList: $9}
    }
    | with_clause INSERT INTO identified_table select_query
    {
        $$ = InsertQuery{WithClause: $1, Table: $4, Query: $5.(SelectQuery)}
    }
    | with_clause INSERT INTO identified_table '(' field_references ')' select_query
    {
        $$ = InsertQuery{WithClause: $1, Table: $4, Fields: $6, Query: $8.(SelectQuery)}
    }

update_query
    : with_clause UPDATE operate_tables SET update_set_list from_clause where_clause
    {
        $$ = UpdateQuery{WithClause: $1, Tables: $3, SetList: $5, FromClause: $6, WhereClause: $7}
    }

update_set
    : field_reference '=' value
    {
        $$ = UpdateSet{Field: $1, Value: $3}
    }

update_set_list
    : update_set
    {
        $$ = []UpdateSet{$1}
    }
    | update_set ',' update_set_list
    {
        $$ = append([]UpdateSet{$1}, $3...)
    }

delete_query
    : with_clause DELETE FROM tables where_clause
    {
        from := FromClause{From: $3.Literal, Tables: $4}
        $$ = DeleteQuery{BaseExpr: NewBaseExpr($2), WithClause: $1, FromClause: from, WhereClause: $5}
    }
    | with_clause DELETE operate_tables FROM tables where_clause
    {
        from := FromClause{From: $4.Literal, Tables: $5}
        $$ = DeleteQuery{BaseExpr: NewBaseExpr($2), WithClause: $1, Tables: $3, FromClause: from, WhereClause: $6}
    }

elseif
    : ELSEIF value THEN program
    {
        $$ = []ElseIf{{Condition: $2, Statements: $4}}
    }
    | ELSEIF value THEN program elseif
    {
        $$ = append([]ElseIf{{Condition: $2, Statements: $4}}, $5...)
    }

else
    :
    {
        $$ = Else{}
    }
    | ELSE program
    {
        $$ = Else{Statements: $2}
    }

in_loop_elseif
    : ELSEIF value THEN loop_program
    {
        $$ = []ElseIf{{Condition: $2, Statements: $4}}
    }
    | ELSEIF value THEN loop_program in_loop_elseif
    {
        $$ = append([]ElseIf{{Condition: $2, Statements: $4}}, $5...)
    }

in_loop_else
    :
    {
        $$ = Else{}
    }
    | ELSE loop_program
    {
        $$ = Else{Statements: $2}
    }

in_function_elseif
    : ELSEIF value THEN function_program
    {
        $$ = []ElseIf{{Condition: $2, Statements: $4}}
    }
    | ELSEIF value THEN function_program in_function_elseif
    {
        $$ = append([]ElseIf{{Condition: $2, Statements: $4}}, $5...)
    }

in_function_else
    :
    {
        $$ = Else{}
    }
    | ELSE function_program
    {
        $$ = Else{Statements: $2}
    }

in_function_in_loop_elseif
    : ELSEIF value THEN function_loop_program
    {
        $$ = []ElseIf{{Condition: $2, Statements: $4}}
    }
    | ELSEIF value THEN function_loop_program in_function_in_loop_elseif
    {
        $$ = append([]ElseIf{{Condition: $2, Statements: $4}}, $5...)
    }

in_function_in_loop_else
    :
    {
        $$ = Else{}
    }
    | ELSE function_loop_program
    {
        $$ = Else{Statements: $2}
    }

case_when
    : WHEN value THEN program
    {
        $$ = []CaseWhen{{Condition: $2, Statements: $4}}
    }
    | WHEN value THEN program case_when
    {
        $$ = append([]CaseWhen{{Condition: $2, Statements: $4}}, $5...)
    }

case_else
    :
    {
        $$ = CaseElse{}
    }
    | ELSE program
    {
        $$ = CaseElse{Statements: $2}
    }

in_loop_case_when
    : WHEN value THEN loop_program
    {
        $$ = []CaseWhen{{Condition: $2, Statements: $4}}
    }
    | WHEN value THEN loop_program in_loop_case_when
    {
        $$ = append([]CaseWhen{{Condition: $2, Statements: $4}}, $5...)
    }

in_loop_case_else
    :
    {
        $$ = CaseElse{}
    }
    | ELSE loop_program
    {
        $$ = CaseElse{Statements: $2}
    }

in_function_case_when
    : WHEN value THEN function_program
    {
        $$ = []CaseWhen{{Condition: $2, Statements: $4}}
    }
    | WHEN value THEN function_program in_function_case_when
    {
        $$ = append([]CaseWhen{{Condition: $2, Statements: $4}}, $5...)
    }

in_function_case_else
    :
    {
        $$ = CaseElse{}
    }
    | ELSE function_program
    {
        $$ = CaseElse{Statements: $2}
    }

in_function_in_loop_case_when
    : WHEN value THEN function_loop_program
    {
        $$ = []CaseWhen{{Condition: $2, Statements: $4}}
    }
    | WHEN value THEN function_loop_program in_function_in_loop_case_when
    {
        $$ = append([]CaseWhen{{Condition: $2, Statements: $4}}, $5...)
    }

in_function_in_loop_case_else
    :
    {
        $$ = CaseElse{}
    }
    | ELSE function_loop_program
    {
        $$ = CaseElse{Statements: $2}
    }

identifier
    : IDENTIFIER
    {
        $$ = Identifier{BaseExpr: NewBaseExpr($1), Literal: $1.Literal, Quoted: $1.Quoted}
    }
    | TIES
    {
        $$ = Identifier{BaseExpr: NewBaseExpr($1), Literal: $1.Literal, Quoted: $1.Quoted}
    }
    | NULLS
    {
        $$ = Identifier{BaseExpr: NewBaseExpr($1), Literal: $1.Literal, Quoted: $1.Quoted}
    }
    | COUNT
    {
        $$ = Identifier{BaseExpr: NewBaseExpr($1), Literal: $1.Literal, Quoted: $1.Quoted}
    }
    | LISTAGG
    {
        $$ = Identifier{BaseExpr: NewBaseExpr($1), Literal: $1.Literal, Quoted: $1.Quoted}
    }
    | AGGREGATE_FUNCTION
    {
        $$ = Identifier{BaseExpr: NewBaseExpr($1), Literal: $1.Literal, Quoted: $1.Quoted}
    }
    | FUNCTION_WITH_INS
    {
        $$ = Identifier{BaseExpr: NewBaseExpr($1), Literal: $1.Literal, Quoted: $1.Quoted}
    }
    | ERROR
    {
        $$ = Identifier{BaseExpr: NewBaseExpr($1), Literal: $1.Literal, Quoted: $1.Quoted}
    }

variable
    : VARIABLE
    {
        $$ = Variable{BaseExpr: NewBaseExpr($1), Name:$1.Literal}
    }

variables
    : variable
    {
        $$ = []Variable{$1}
    }
    | variable ',' variables
    {
        $$ = append([]Variable{$1}, $3...)
    }

variable_substitution
    : variable SUBSTITUTION_OP value
    {
        $$ = VariableSubstitution{Variable:$1, Value:$3}
    }

variable_assignment
    : variable
    {
        $$ = VariableAssignment{Variable:$1}
    }
    | variable SUBSTITUTION_OP value
    {
        $$ = VariableAssignment{Variable: $1, Value: $3}
    }

variable_assignments
    : variable_assignment
    {
        $$ = []VariableAssignment{$1}
    }
    | variable_assignment ',' variable_assignments
    {
        $$ = append([]VariableAssignment{$1}, $3...)
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

join_type_inner
    :
    {
        $$ = Token{}
    }
    | INNER
    {
        $$ = $1
    }

join_type_outer
    :
    {
        $$ = Token{}
    }
    | OUTER
    {
        $$ = $1
    }

join_outer_direction
    : LEFT
    {
        $$ = $1
    }
    | RIGHT
    {
        $$ = $1
    }

all
    :
    {
        $$ = Token{}
    }
    | ALL
    {
        $$ = $1
    }

recursive
    :
    {
        $$ = Token{}
    }
    | RECURSIVE
    {
        $$ = $1
    }

as
    :
    {
        $$ = Token{}
    }
    | AS
    {
        $$ = $1
    }

comparison_operator
    : COMPARISON_OP
    {
        $$ = $1
    }
    | '='
    {
        $1.Token = COMPARISON_OP
        $$ = $1
    }

%%

func SetDebugLevel(level int, verbose bool) {
	yyDebug        = level
	yyErrorVerbose = verbose
}

func Parse(s string, sourceFile string) ([]Statement, error) {
    l := new(Lexer)
    l.Init(s, sourceFile)
    yyParse(l)
    return l.program, l.err
}