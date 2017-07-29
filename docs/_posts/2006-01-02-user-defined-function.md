---
layout: default
title: User Defined Function - Reference Manual - csvq
category: reference
---

# User Defined Function

A User Defined Function is a routine that can be called just like built-in functions.
A function has some input parameters, and returns a single value.

Functions create local scopes.
[Variables]({{ '/reference/variable.html' | relative_url }}), [cursors]({{ '/reference/cursor.html' | relative_url }}), [temporary tables]({{ '/reference/temporary-table.html' | relative_url }}), and [functions]({{ '/reference/user-defined-function.html' | relative_url }}) declared in user defined functions can be refered only within the functions. 

## Declare Function
{: #declare}

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

### RETURN Statement

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
