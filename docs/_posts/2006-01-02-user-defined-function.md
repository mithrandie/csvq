---
layout: default
title: User Defined Function - Reference Manual - csvq
category: reference
---

# User Defined Function

A User Defined Function is a routine that can be called just like built-in functions.
A function has some input parameters, and [returns](#return) a single value.

Functions create local scopes.
[Variables]({{ '/reference/variable.html' | relative_url }}), [cursors]({{ '/reference/cursor.html' | relative_url }}), [temporary tables]({{ '/reference/temporary-table.html' | relative_url }}), and [functions]({{ '/reference/user-defined-function.html' | relative_url }}) declared in user defined functions can be refered only within the functions. 

* [Scala Function](#scala)
* [Aggregate Function](#aggregate)
* [Return Statement](#return)

## Scala Function
{: #scala}

### Declaration
{: #scala_declaration}

```sql
DECLARE function_name FUNCTION ([parameter [, parameter ...]])
AS
BEGIN
  statements
END;
```

_function_name_
: [identifier]({{ '/reference/statement.html#parsing' | relative_url }})

_parameter_
: [Variable]({{ '/reference/variable.html' | relative_url }})

_statements_
: [Statements]({{ '/reference/statement.html' | relative_url }})

A scala function takes some arguments, and return a value.

#### Usage

```sql
function_name([argument, [, argument ...]])
```

_function_name_
: [identifier]({{ '/reference/statement.html#parsing' | relative_url }})

_argument_
: [value]({{ '/reference/value.html' | relative_url }})

In the statements, arguments are set to variables specified in the declaration as _parameters_.


## Aggregate Function
{: #aggregate}

### Declaration
{: #aggregate_declaration}

```sql
DECLARE function_name AGGREGATE (cursor_name)
AS
BEGIN
  statements
END;
```

_function_name_
: [identifier]({{ '/reference/statement.html#parsing' | relative_url }})

_cursor_name_
: [identifier]({{ '/reference/statement.html#parsing' | relative_url }})

_statements_
: [Statements]({{ '/reference/statement.html' | relative_url }})

A aggregate function takes exactly one argument, and return a value.

#### Usage

You can use a user defined aggregate function as a [Aggregate Function]({{ '/reference/aggregate-functions.html' | relative_url }}) or a [Analytic Function]({{ '/reference/analytic-functions.html' | relative_url }}).

##### As a Aggregate Function

```sql
function_name([DISTINCT] expr)
```

_function_name_
: [identifier]({{ '/reference/statement.html#parsing' | relative_url }})

_expr_
: [value]({{ '/reference/value.html' | relative_url }})

##### As a Analytic Function

```sql
function_name([DISTINCT] expr) OVER ([partition_clause] [order_by_clause])
```

_function_name_
: [identifier]({{ '/reference/statement.html#parsing' | relative_url }})

_expr_
: [value]({{ '/reference/value.html' | relative_url }})

In the statements, grouped values represented by _expr_ can be referred with a pseudo cursor named as _cursor_name_ in the declaration.
You can use the [Fetch Statement]({{ '/reference/cursor.html#fetch' | relative_url }}), [While In Loop Statement]({{ '/reference/control-flow.html#while_in_loop' | relative_url }}) or the [Cursor Status Expressions]({{ '/reference/cursor.html#status' | relative_url }}) against the pseudo cursor. 

Example:

```sql
DECLARE multiply AGGREGATE (list)
AS
BEGIN
    VAR @value, @fetch;

    WHILE @fetch IN list
    DO
        IF @fetch IS NULL THEN
            CONTINUE;
        END IF;

        IF @value IS NULL THEN
            @value := @fetch
            CONTINUE;
        END IF;

        @value := @value * @fetch;
    END WHILE;

    RETURN @value;
END;

SELECT multiply(i) FROM numbers;

SELECT i, multiply(i) OVER () FROM numbers;
```


## RETURN Statement
{: #return}

A Return statement terminates executing function, then returns a value.
If the return value is not specified, then returns a null.

When there is no return statement, the function executes all of the statements and returns a null.

```sql
return_statement
  : RETURN;
  | RETURN value;
```

_value_
: [value]({{ '/reference/value.html' | relative_url }})
