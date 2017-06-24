%{
package parser
%}

%union{
    program     []Statement
    statement   Statement
    expression  Expression
    expressions []Expression
    procexpr    ProcExpr
    procexprs   []ProcExpr
    primary     Primary
    identifier  Identifier
    text        String
    integer     Integer
    float       Float
    ternary     Ternary
    datetime    Datetime
    null        Null
    variable    Variable
    variables   []Variable
    token       Token
}

%type<program>     program
%type<program>     in_loop_program
%type<statement>   statement
%type<statement>   in_loop_statement
%type<statement>   variable_statement
%type<statement>   transaction_statement
%type<statement>   cursor_statement
%type<statement>   flow_control_statement
%type<statement>   in_loop_flow_control_statement
%type<statement>   command_statement
%type<expression>  select_query
%type<expression>  select_entity
%type<expression>  select_clause
%type<expression>  from_clause
%type<expression>  where_clause
%type<expression>  group_by_clause
%type<expression>  having_clause
%type<expression>  order_by_clause
%type<expression>  limit_clause
%type<primary>     primary
%type<expression>  field_reference
%type<expression>  value
%type<expression>  order_item
%type<expression>  subquery
%type<expression>  string_operation
%type<expression>  comparison
%type<expression>  arithmetic
%type<expression>  logic
%type<expression>  function
%type<expression>  option
%type<expression>  group_concat
%type<expression>  identified_table
%type<expression>  virtual_table
%type<expression>  table
%type<expression>  join
%type<expression>  join_condition
%type<expression>  field_object
%type<expression>  field
%type<expression>  case
%type<expression>  case_value
%type<expression>  case_else
%type<expressions> field_references
%type<expressions> values
%type<expressions> order_items
%type<expressions> tables
%type<expressions> identified_tables
%type<expressions> using_fields
%type<expressions> fields
%type<expressions> case_when
%type<expression>  insert_query
%type<expression>  insert_values
%type<expressions> insert_values_list
%type<expression>  update_query
%type<expression>  update_set
%type<expressions> update_set_list
%type<expression>  delete_query
%type<expression>  create_table
%type<expression>  add_columns
%type<expression>  column_default
%type<expressions> column_defaults
%type<expression>  column_position
%type<expression>  drop_columns
%type<expression>  rename_column
%type<procexprs>   elseif
%type<procexpr>    else
%type<procexprs>   in_loop_elseif
%type<procexpr>    in_loop_else
%type<identifier>  identifier
%type<text>        text
%type<integer>     integer
%type<float>       float
%type<ternary>     ternary
%type<datetime>    datetime
%type<null>        null
%type<variable>    variable
%type<variables>   variables
%type<expression>  variable_substitution
%type<expression>  variable_assignment
%type<expressions> variable_assignments
%type<token>       distinct
%type<token>       negation
%type<token>       order_direction
%type<token>       join_inner
%type<token>       join_outer
%type<token>       join_direction
%type<token>       all
%type<token>       comparison_operator
%type<token>       statement_terminal

%token<token> IDENTIFIER STRING INTEGER FLOAT BOOLEAN TERNARY DATETIME VARIABLE FLAG
%token<token> SELECT FROM UPDATE SET DELETE WHERE INSERT INTO VALUES AS DUAL STDIN
%token<token> CREATE ADD DROP ALTER TABLE FIRST LAST AFTER BEFORE DEFAULT RENAME TO
%token<token> ORDER GROUP HAVING BY ASC DESC LIMIT
%token<token> JOIN INNER OUTER LEFT RIGHT FULL CROSS ON USING NATURAL
%token<token> UNION INTERSECT EXCEPT
%token<token> ALL ANY EXISTS IN
%token<token> AND OR NOT BETWEEN LIKE IS NULL
%token<token> DISTINCT WITH
%token<token> CASE IF ELSEIF WHILE WHEN THEN ELSE DO END
%token<token> DECLARE CURSOR FOR FETCH OPEN CLOSE DISPOSE
%token<token> GROUP_CONCAT SEPARATOR
%token<token> COMMIT ROLLBACK
%token<token> CONTINUE BREAK EXIT
%token<token> PRINT
%token<token> VAR
%token<token> COMPARISON_OP STRING_OP SUBSTITUTION_OP

%left UNION EXCEPT
%left INTERSECT
%left OR
%left AND
%right NOT
%nonassoc '=' COMPARISON_OP IS BETWEEN IN LIKE
%left STRING_OP
%left '+' '-'
%left '*' '/' '%'

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

in_loop_program
    :
    {
        $$ = nil
        yylex.(*Lexer).program = $$
    }
    | in_loop_statement in_loop_program
    {
        $$ = append([]Statement{$1}, $2...)
        yylex.(*Lexer).program = $$
    }

statement
    : select_query statement_terminal
    {
        $$ = $1
    }
    | insert_query statement_terminal
    {
        $$ = $1
    }
    | update_query statement_terminal
    {
        $$ = $1
    }
    | delete_query statement_terminal
    {
        $$ = $1
    }
    | create_table statement_terminal
    {
        $$ = $1
    }
    | add_columns statement_terminal
    {
        $$ = $1
    }
    | drop_columns statement_terminal
    {
        $$ = $1
    }
    | rename_column statement_terminal
    {
        $$ = $1
    }
    | variable_statement
    {
        $$ = $1
    }
    | transaction_statement
    {
        $$ = $1
    }
    | cursor_statement
    {
        $$ = $1
    }
    | flow_control_statement
    {
        $$ = $1
    }
    | command_statement
    {
        $$ = $1
    }

in_loop_statement
    : statement
    {
        $$ = $1
    }
    | in_loop_flow_control_statement
    {
        $$ = $1
    }

variable_statement
    : VAR variable_assignments statement_terminal
    {
        $$ = VariableDeclaration{Assignments:$2}
    }
    | variable_substitution statement_terminal
    {
        $$ = $1
    }

transaction_statement
    : COMMIT statement_terminal
    {
        $$ = TransactionControl{Token: $1.Token}
    }
    | ROLLBACK statement_terminal
    {
        $$ = TransactionControl{Token: $1.Token}
    }

cursor_statement
    : DECLARE identifier CURSOR FOR select_query statement_terminal
    {
        $$ = CursorDeclaration{Cursor:$2, Query: $5.(SelectQuery)}
    }
    | OPEN identifier statement_terminal
    {
        $$ = OpenCursor{Cursor: $2}
    }
    | CLOSE identifier statement_terminal
    {
        $$ = CloseCursor{Cursor: $2}
    }
    | DISPOSE identifier statement_terminal
    {
        $$ = DisposeCursor{Cursor: $2}
    }
    | FETCH identifier INTO variables statement_terminal
    {
        $$ = FetchCursor{Cursor: $2, Variables: $4}
    }

flow_control_statement
    : IF value THEN program else END IF statement_terminal
    {
        $$ = If{Condition: $2, Statements: $4, Else: $5}
    }
    | IF value THEN program elseif else END IF statement_terminal
    {
        $$ = If{Condition: $2, Statements: $4, ElseIf: $5, Else: $6}
    }
    | WHILE value DO in_loop_program END WHILE statement_terminal
    {
        $$ = While{Condition: $2, Statements: $4}
    }
    | WHILE variables IN identifier DO in_loop_program END WHILE statement_terminal
    {
        $$ = WhileInCursor{Variables: $2, Cursor: $4, Statements: $6}
    }
    | EXIT statement_terminal
    {
        $$ = FlowControl{Token: $1.Token}
    }

in_loop_flow_control_statement
    : IF value THEN in_loop_program in_loop_else END IF statement_terminal
    {
        $$ = If{Condition: $2, Statements: $4, Else: $5}
    }
    | IF value THEN in_loop_program in_loop_elseif in_loop_else END IF statement_terminal
    {
        $$ = If{Condition: $2, Statements: $4, ElseIf: $5, Else: $6}
    }
    | CONTINUE statement_terminal
    {
        $$ = FlowControl{Token: $1.Token}
    }
    | BREAK statement_terminal
    {
        $$ = FlowControl{Token: $1.Token}
    }

command_statement
    : SET FLAG '=' primary statement_terminal
    {
        $$ = SetFlag{Name: $2.Literal, Value: $4}
    }
    | PRINT value statement_terminal
    {
        $$ = Print{Value: $2}
    }

select_query
    : select_entity order_by_clause limit_clause
    {
        $$ = SelectQuery{
            SelectEntity:  $1,
            OrderByClause: $2,
            LimitClause:   $3,
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
    | select_entity UNION all select_entity
    {
        $$ = SelectSet{
            LHS:      $1,
            Operator: $2,
            All:      $3,
            RHS:      $4,
        }
    }
    | select_entity INTERSECT all select_entity
    {
        $$ = SelectSet{
            LHS:      $1,
            Operator: $2,
            All:      $3,
            RHS:      $4,
        }
    }
    | select_entity EXCEPT all select_entity
    {
        $$ = SelectSet{
            LHS:      $1,
            Operator: $2,
            All:      $3,
            RHS:      $4,
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
    | LIMIT integer
    {
        $$ = LimitClause{Limit: $1.Literal, Number: $2.Value()}
    }

primary
    : text
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
    | datetime
    {
        $$ = $1
    }
    | null
    {
        $$ = $1
    }

field_reference
    : identifier
    {
        $$ = FieldReference{Column: $1}
    }
    | identifier '.' identifier
    {
        $$ = FieldReference{View: $1, Column: $3}
    }

value
    : field_reference
    {
        $$ = $1
    }
    | primary
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
    | value '=' value
    {
        $$ = Comparison{LHS: $1, Operator: Token{Token: COMPARISON_OP, Literal: "="}, RHS: $3}
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
    | value comparison_operator ANY subquery
    {
        $$ = Any{Any: $3.Literal, LHS: $1, Operator: $2, Query: $4.(Subquery)}
    }
    | value comparison_operator ALL subquery
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
        $$ = Logic{LHS: nil, Operator: $1, RHS: $2}
    }

function
    : identifier '(' option ')'
    {
        $$ = Function{Name: $1.Literal, Option: $3.(Option)}
    }
    | group_concat
    {
        $$ = $1
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
    | distinct values
    {
        $$ = Option{Distinct: $1, Args: $2}
    }

group_concat
    : GROUP_CONCAT '(' option order_by_clause ')'
    {
        $$ = GroupConcat{GroupConcat: $1.Literal, Option: $3.(Option), OrderBy: $4}
    }
    | GROUP_CONCAT '(' option order_by_clause SEPARATOR STRING ')'
    {
        $$ = GroupConcat{GroupConcat: $1.Literal, Option: $3.(Option), OrderBy: $4, SeparatorLit: $5.Literal, Separator: $6.Literal}
    }

identified_table
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

virtual_table
    : subquery
    {
        $$ = $1
    }
    | STDIN
    {
        $$ = Stdin{Stdin: $1.Literal}
    }

table
    : identified_table
    {
        $$ = $1
    }
    | virtual_table
    {
        $$ = Table{Object: $1}
    }
    | virtual_table identifier
    {
        $$ = Table{Object: $1, Alias: $2}
    }
    | virtual_table AS identifier
    {
        $$ = Table{Object: $1, As: $2, Alias: $3}
    }
    | join
    {
        $$ = Table{Object: $1}
    }
    | DUAL
    {
        $$ = Table{Object: Dual{Dual: $1.Literal}}
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
    | ON value
    {
        $$ = JoinCondition{Literal:$1.Literal, On: $2}
    }
    | USING '(' using_fields ')'
    {
        $$ = JoinCondition{Literal:$1.Literal, Using: $3}
    }

field_object
    : value
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

field_references
    : field_reference
    {
        $$ = []Expression{$1}
    }
    | field_reference ',' field_references
    {
        $$ = append([]Expression{$1}, $3...)
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

identified_tables
    : identified_table
    {
        $$ = []Expression{$1}
    }
    | identified_table ',' identified_tables
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
    : WHEN value THEN value
    {
        $$ = []Expression{CaseWhen{When: $1.Literal, Then: $3.Literal, Condition: $2, Result: $4}}
    }
    | case_when case_when
    {
        $$ = append($1, $2...)
    }

insert_query
    : INSERT INTO identifier VALUES insert_values_list
    {
        $$ = InsertQuery{Insert: $1.Literal, Into: $2.Literal, Table: $3, Values: $4.Literal, ValuesList: $5}
    }
    | INSERT INTO identifier '(' field_references ')' VALUES insert_values_list
    {
        $$ = InsertQuery{Insert: $1.Literal, Into: $2.Literal, Table: $3, Fields: $5, Values: $7.Literal, ValuesList: $8}
    }
    | INSERT INTO identifier select_query
    {
        $$ = InsertQuery{Insert: $1.Literal, Into: $2.Literal, Table: $3, Query: $4.(SelectQuery)}
    }
    | INSERT INTO identifier '(' field_references ')' select_query
    {
        $$ = InsertQuery{Insert: $1.Literal, Into: $2.Literal, Table: $3, Fields: $5, Query: $7.(SelectQuery)}
    }

insert_values
    : '(' values ')'
    {
        $$ = InsertValues{Values: $2}
    }

insert_values_list
    : insert_values
    {
        $$ = []Expression{$1}
    }
    | insert_values ',' insert_values_list
    {
        $$ = append([]Expression{$1}, $3...)
    }

update_query
    : UPDATE identified_tables SET update_set_list from_clause where_clause
    {
        $$ = UpdateQuery{Update: $1.Literal, Tables: $2, Set: $3.Literal, SetList: $4, FromClause: $5, WhereClause: $6}
    }

update_set
    : field_reference '=' value
    {
        $$ = UpdateSet{Field: $1.(FieldReference), Value: $3}
    }

update_set_list
    : update_set
    {
        $$ = []Expression{$1}
    }
    | update_set ',' update_set_list
    {
        $$ = append([]Expression{$1}, $3...)
    }

delete_query
    : DELETE FROM tables where_clause
    {
        from := FromClause{From: $2.Literal, Tables: $3}
        $$ = DeleteQuery{Delete: $1.Literal, FromClause: from, WhereClause: $4}
    }
    | DELETE identified_tables FROM tables where_clause
    {
        from := FromClause{From: $3.Literal, Tables: $4}
        $$ = DeleteQuery{Delete: $1.Literal, Tables: $2, FromClause: from, WhereClause: $5}
    }

create_table
    : CREATE TABLE identifier '(' using_fields ')'
    {
        $$ = CreateTable{CreateTable: $1.Literal + " " + $2.Literal, Table: $3, Fields: $5}
    }

add_columns
    : ALTER TABLE identifier ADD column_default column_position
    {
        $$ = AddColumns{AlterTable: $1.Literal + " " + $2.Literal, Table: $3, Add: $4.Literal, Columns: []Expression{$5}, Position: $6}
    }
    | ALTER TABLE identifier ADD '(' column_defaults ')' column_position
    {
        $$ = AddColumns{AlterTable: $1.Literal + " " + $2.Literal, Table: $3, Add: $4.Literal, Columns: $6, Position: $8}
    }

column_default
    : identifier
    {
        $$ = ColumnDefault{Column: $1}
    }
    | identifier DEFAULT value
    {
        $$ = ColumnDefault{Column: $1, Default: $2.Literal, Value: $3}
    }

column_defaults
    : column_default
    {
        $$ = []Expression{$1}
    }
    | column_default ',' column_defaults
    {
        $$ = append([]Expression{$1}, $3...)
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

drop_columns
    : ALTER TABLE identifier DROP field_reference
    {
        $$ = DropColumns{AlterTable: $1.Literal + " " + $2.Literal, Table: $3, Drop: $4.Literal, Columns: []Expression{$5}}
    }
    | ALTER TABLE identifier DROP '(' field_references ')'
    {
        $$ = DropColumns{AlterTable: $1.Literal + " " + $2.Literal, Table: $3, Drop: $4.Literal, Columns: $6}
    }

rename_column
    : ALTER TABLE identifier RENAME field_reference TO identifier
    {
        $$ = RenameColumn{AlterTable: $1.Literal + " " + $2.Literal, Table: $3, Rename: $4.Literal, Old: $5.(FieldReference), To: $6.Literal, New: $7}
    }

elseif
    : ELSEIF value THEN program
    {
        $$ = []ProcExpr{ElseIf{Condition: $2, Statements: $4}}
    }
    | elseif elseif
    {
        $$ = append($1, $2...)
    }

else
    :
    {
        $$ = nil
    }
    | ELSE program
    {
        $$ = Else{Statements: $2}
    }

in_loop_elseif
    : ELSEIF value THEN in_loop_program
    {
        $$ = []ProcExpr{ElseIf{Condition: $2, Statements: $4}}
    }
    | in_loop_elseif in_loop_elseif
    {
        $$ = append($1, $2...)
    }

in_loop_else
    :
    {
        $$ = nil
    }
    | ELSE in_loop_program
    {
        $$ = Else{Statements: $2}
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

datetime
    : DATETIME
    {
        $$ = NewDatetimeFromString($1.Literal)
    }

null
    : NULL
    {
        $$ = NewNullFromString($1.Literal)
    }

variable
    : VARIABLE
    {
        $$ = Variable{Name:$1.Literal}
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
    : VARIABLE
    {
        $$ = VariableAssignment{Name:$1.Literal}
    }
    | VARIABLE SUBSTITUTION_OP value
    {
        $$ = VariableAssignment{Name: $1.Literal, Value: $3}
    }

variable_assignments
    : variable_assignment
    {
        $$ = []Expression{$1}
    }
    | variable_assignment ',' variable_assignments
    {
        $$ = append([]Expression{$1}, $3...)
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

all
    :
    {
        $$ = Token{}
    }
    | ALL
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
        $$ = Token{Token:COMPARISON_OP, Literal:string('=')}
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