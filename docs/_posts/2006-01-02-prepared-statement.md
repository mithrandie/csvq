---
layout: default
title: Prepared Statement - Reference Manual - csvq
category: reference
---

# Prepared Statement

A prepared statement is used to execute the same statement repeatedly with embedded values.

* [Usage Flow](#usage_flow)
* [Prepare Statement](#prepare)
* [Execute Prepared Statement](#execute)
* [Dispose Prepared Statement](#dispose)
* [Placeholder](#placeholder)

## Usage Flow
{: #usage_flow}

1. Prepare a statement.
2. Execute the statement using embedded values.
3. Dispose the statement as necessary.

## Prepare Statement
{: #prepare}

```sql
PREPARE statement_name FROM statement;
```

_statement_name_
: [identifier]({{ '/reference/statement.html#parsing' | relative_url }})

_statement_
: [string]({{ '/reference/value.html#string' | relative_url }})


## Execute Prepared Statement
{: #execute}

```sql
execute_statement
  : EXECUTE statement_name;
  | EXECUTE statement_name USING replace_value [, replace_value ...];
  
replace_value
  : value
  | value AS placeholder_name
```

_statement_name_
: [identifier]({{ '/reference/statement.html#parsing' | relative_url }})

_value_
: [value]({{ '/reference/value.html' | relative_url }})

_placeholder_name_
: [identifier]({{ '/reference/statement.html#parsing' | relative_url }})

  Placeholder names are used for named placeholders.


## Dispose Prepared Statement
{: #dispose}

```sql
DISPOSE PREPARE statement_name;
```

_statement_name_
: [identifier]({{ '/reference/statement.html#parsing' | relative_url }})

## Placeholder
{: #placeholder}

Positional Placeholder
: Question Mark(U+003F `?`)

Named Placeholder
: Colon(U+003A `:`) and followd by [identifier]({{ '/reference/statement.html#parsing' | relative_url }})

### Example

```sql
-- Positional Placeholder
PREPARE stmt1 FROM 'SELECT ?, ?, ?;';
EXECUTE stmt1 USING 'a', 'b', 'c';

-- Named Placeholder
PREPARE stmt2 FROM 'SELECT :second, :third, :first;';
EXECUTE stmt2 USING 'a' AS `first`, 'b' AS `second`, 'c' AS `third`;
```