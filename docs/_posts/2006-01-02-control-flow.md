---
layout: default
title: Control Flow - Reference Manual - csvq
category: reference
---

# Control Flow

* [If](#if)
* [While Loop](#while_loop)
* [While In Loop](#while_in_loop)
* [Continue](#continue)
* [Break](#break)
* [Exit](#exit)

## If
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

A If statement executes the first _statements_ that _condition_ is TRUE.
If no condition is TRUE, it executes the _statements_ of the ELSE expression.

## While Loop
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

## While In Loop
{: #while_in_loop}
```sql
WHILE variable [, variable ...] IN cursor
DO
  statements
END WHILE
```

_variable_
: [Variable]({{ '/reference/variable.html' | relative_url }})

_cursor_
: [Cursor]({{ '/reference/cursor.html' | relative_url }})

_statements_
: [Statements]({{ '/reference/statement.html' | relative_url }})

A While In statement fetch the data from _cursor_ into variables, then execute _statements_.
The While In statement iterates it until the _cursor_ pointer reaches the last record in the referring view.

## Continue
{: #continue}

```sql
CONTINUE;
```

A Continue statement stops statements execution in loop, then jumps to the next iteration.

## Break
{: #break}

```sql
BREAK;
```

A Break statement stops statements execution in loop, then exit from current loop.

## Exit
{: #exit}

```sql
EXIT;
```

A Exit statement stops statements execution, then terminates current program without commit.