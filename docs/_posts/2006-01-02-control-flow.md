---
layout: default
title: Control Flow - Reference Manual - csvq
category: reference
---

# Control Flow

* [IF](#if)
* [CASE](#case)
* [WHILE](#while_loop)
* [WHILE IN](#while_in_loop)
* [CONTINUE](#continue)
* [BREAK](#break)
* [EXIT](#exit)
* [TRIGGER ERROR](#trigger_error)

_IF_ statements and _WHILE_ statements create local scopes.
[Variables]({{ '/reference/variable.html' | relative_url }}), [cursors]({{ '/reference/cursor.html' | relative_url }}), [temporary tables]({{ '/reference/temporary-table.html' | relative_url }}), and [functions]({{ '/reference/user-defined-function.html' | relative_url }}) declared in statement blocks can be refered only within the blocks. 

## IF
{: #if}

```sql
IF condition THEN statements
  [ELSEIF condition THEN statements ...]
  [ELSE statements]
END IF;
```

_condition_
: [value]({{ '/reference/value.html' | relative_url }})

_statements_
: [Statements]({{ '/reference/statement.html' | relative_url }})

IF statement executes the first _statements_ that _condition_ is TRUE.
If no condition is TRUE, the _statements_ of the ELSE expression are executed.

## CASE
{: #case}

### Case with condition

```sql
CASE
  WHEN condition THEN statements
  [WHEN condition THEN statements]
  [ELSE statements]
END CASE;
```

_condition_
: [value]({{ '/reference/value.html' | relative_url }})

_statements_
: [Statements]({{ '/reference/statement.html' | relative_url }})

Execute _statements_ of the first WHEN expression that _condition_ is TRUE.
If no condition is TRUE, then execute _statements_ of the ELSE expression.


### Case with comparison

```sql
CASE value
  WHEN comparison_value THEN statements
  [WHEN comparison_value THEN statements]
  [ELSE statements]
END CASE;
```

_value_
: [value]({{ '/reference/value.html' | relative_url }})

_comparison_value_
: [value]({{ '/reference/value.html' | relative_url }})

_statements_
: [Statements]({{ '/reference/statement.html' | relative_url }})

Execute _statements_ of the first WHEN expression that _comparison_value_ is equal to _value_.
If no _comparison_value_ is match, then execute _statements_ of the ELSE expression.


## WHILE
{: #while_loop}

```sql
WHILE condition
DO
  statements
END WHILE;
```

_condition_
: [value]({{ '/reference/value.html' | relative_url }})

_statements_
: [Statements]({{ '/reference/statement.html' | relative_url }})

A While statement evaluate _condition_, then if condition is TRUE, executes _statements_. 
The While statement iterates it while _condition_ is TRUE.

## WHILE IN
{: #while_in_loop}
```sql
WHILE [DECLARE|VAR] variable [, variable ...] IN cursor_name
DO
  statements
END WHILE;
```

_variable_
: [Variable]({{ '/reference/variable.html' | relative_url }})

_cursor_name_
: [identifier]({{ '/reference/statement.html#parsing' | relative_url }})

_statements_
: [Statements]({{ '/reference/statement.html' | relative_url }})

A While In statement fetch the data from the [cursor]({{ '/reference/cursor.html' | relative_url }}) into variables, then execute _statements_.
The While In statement iterates it until the _cursor_name_ pointer reaches the last record in the referring view.

If DECLARE or VAR keyword is specified, then variables are declared in the child scope. 
Otherwise, variables in the current scope is used to fetch.

## CONTINUE
{: #continue}

```sql
CONTINUE;
```

A Continue statement stops statements execution in loop, then jumps to the next iteration.

## BREAK
{: #break}

```sql
BREAK;
```

A Break statement stops statements execution in loop, then exit from current loop.

## EXIT
{: #exit}

```sql
EXIT [exit_code];
```

_exit_code_
: [integer]({{ '/reference/value.html#integer' | relative_url }})

  0 is the default.

Exit statement stops statements execution, then terminates the executing procedure without commit.

## TRIGGER ERROR
{: #trigger_error}

```sql
TRIGGER ERROR [exit_code] [error_message];
```

_exit_code_
: [integer]({{ '/reference/value.html#integer' | relative_url }})

  64 is the default.

_error_message_
: [string]({{ '/reference/value.html#string' | relative_url }})

A trigger error statement stops statements execution, then terminates the executing procedure with an error.