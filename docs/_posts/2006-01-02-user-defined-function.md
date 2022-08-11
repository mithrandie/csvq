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

* [Scalar Function](#scalar)
* [Aggregate Function](#aggregate)
* [DISPOSE FUNCTION Statement](#dispose)
* [RETURN Statement](#return)

## Scalar Function
{: #scalar}

### Declaration
{: #scalar_declaration}

```sql
scalar_function_declaration
  : DECLARE function_name FUNCTION ([parameter [, parameter ...] [, optional_parameter ...]])
    AS
    BEGIN
      statements
    END;

optional_parameter
  : parameter DEFAULT value
```

_function_name_
: [identifier]({{ '/reference/statement.html#parsing' | relative_url }})

_statements_
: [Statements]({{ '/reference/statement.html' | relative_url }})

_parameter_
: [Variable]({{ '/reference/variable.html' | relative_url }})

_value_
: [value]({{ '/reference/statement.html' | relative_url }})

A scalar function takes some arguments, and returns a value.
In the statements, arguments are set to variables specified in the declaration as _parameters_.


#### Usage

```sql
function_name([argument, [, argument ...]])
```

_function_name_
: [identifier]({{ '/reference/statement.html#parsing' | relative_url }})

_argument_
: [value]({{ '/reference/value.html' | relative_url }})


## Aggregate Function
{: #aggregate}

### Declaration
{: #aggregate_declaration}

```sql
aggregate_function_declaration
  : DECLARE function_name AGGREGATE (cursor_name [, parameter ...] [, optional_parameter ...])
    AS
    BEGIN
      statements
    END;

optional_parameter
  : parameter DEFAULT value
```

_function_name_
: [identifier]({{ '/reference/statement.html#parsing' | relative_url }})

_cursor_name_
: [identifier]({{ '/reference/statement.html#parsing' | relative_url }})

_statements_
: [Statements]({{ '/reference/statement.html' | relative_url }})

_parameter_
: [Variable]({{ '/reference/variable.html' | relative_url }})
  
_value_
: [value]({{ '/reference/statement.html' | relative_url }})

An aggregate function takes at least one argument, and returns a value.
The first argument is a representation of grouped values, and the following arguments are parameters.

In the statements, grouped values represented by the first argument can be referred with a pseudo cursor named as _cursor_name_, 
and the second argument and the followings are set to variables specified in the declaration as _parameters_.
You can use the [Fetch Statement]({{ '/reference/cursor.html#fetch' | relative_url }}), [While In Statement]({{ '/reference/control-flow.html#while_in_loop' | relative_url }}) or the [Cursor Status Expressions]({{ '/reference/cursor.html#status' | relative_url }}) against the pseudo cursor. 


#### Usage

You can use a user defined aggregate function as an [Aggregate Function]({{ '/reference/aggregate-functions.html' | relative_url }}) or an [Analytic Function]({{ '/reference/analytic-functions.html' | relative_url }}).

##### As an Aggregate Function

```sql
function_name([DISTINCT] expr [, argument ...])
```

_function_name_
: [identifier]({{ '/reference/statement.html#parsing' | relative_url }})

_expr_
: [value]({{ '/reference/value.html' | relative_url }})

_argument_
: [value]({{ '/reference/value.html' | relative_url }})

##### As an Analytic Function

```sql
function_name([DISTINCT] expr [, argument ...]) OVER ([partition_clause] [order_by_clause [windowing_clause]])
```

_function_name_
: [identifier]({{ '/reference/statement.html#parsing' | relative_url }})

_expr_
: [value]({{ '/reference/value.html' | relative_url }})

_argument_
: [value]({{ '/reference/value.html' | relative_url }})

_partition_clause_
: [Partition Clause]({{ '/reference/analytic-functions.html#syntax' | relative_url }})

_order_by_clause_
: [Order By Clause]({{ '/reference/select-query.html#order_by_clause' | relative_url }})

_windowing_clause_
: [Windowing Clause]({{ '/reference/analytic-functions.html#syntax' | relative_url }})


Example:

```sql
DECLARE product AGGREGATE (list, @default DEFAULT 0)
AS
BEGIN
    VAR @value, @fetch;

    WHILE @fetch IN list
    DO
        VAR @floatVal := FLOAT(@fetch);
        
        IF @floatVal IS NULL THEN
            CONTINUE;
        END IF;

        IF @value IS NULL THEN
            @value := @floatVal;
            CONTINUE;
        END IF;

        @value := @value * @floatVal;
    END WHILE;
    
    IF @value IS NULL THEN
        @value := @default;
    END IF;

    RETURN @value;
END;

SELECT product(i) FROM numbers;

SELECT product(i, NULL) FROM numbers;

SELECT i, product(i) OVER (order by i) FROM numbers;
```

## DISPOSE FUNCTION Statement
{: #dispose}

A DISPOSE FUNCTION statement disposes user defined function named as _function_name_.

```sql
DISPOSE FUNCTION function_name; 
```

_function_name_
: [identifier]({{ '/reference/statement.html#parsing' | relative_url }})


## RETURN Statement
{: #return}

A RETURN statement terminates executing function, then returns a value.
If the return value is not specified, then returns a null.

When there is no return statement, the function executes all the statements and returns a null.

```sql
RETURN [value];
```

_value_
: [value]({{ '/reference/value.html' | relative_url }})
