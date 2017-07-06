---
layout: default
title: Common Table - Reference Manual - csvq
category: reference
---

# Common Table

A Common Table is a temporary view that can be referenced in a single query.
You can use common tables in a [Select Query]({{ '/reference/select-query.html' | relative_url }}), [Insert Query]({{ '/reference/insert-query.html' | relative_url }}), [Update Query]({{ '/reference/update-query.html' | relative_url }}), or [Delete Query]({{ '/reference/delete-query.html' | relative_url }}).

## Syntax

```sql
common_table_clause
  : WITH common_table [, common_table ...]

common_table
  : [RECURSIVE] table_name [(column_name [, column_name ...])] AS (select_query)
```

_table_name_
: [identifier]({{ '/reference/statement.html#parsing' | relative_url }})

_column_name_
: [identifier]({{ '/reference/statement.html#parsing' | relative_url }})

_select_query_
: [select_query]({{ '/reference/select-query.html' | relative_url }})

### Recursion

If you specified a RECURSIVE keyword, the _select_query_ in the _common_table_clause_ can retrieve the result recursively.
A RECURSIVE keyword is usually used with a [UNION]({{ '/reference/set-operators.html#union' | relative_url }}) operator.

```sql
WITH
  RECURSIVE table_name [(column_name [, column_name ...])]
  AS (
    base_select_query
    UNION [ALL]
    recursive_select_query
  )
```

At first, the result set of the _base_select_query_ is stored in the _temporary view_ for recursion.
Next, the _recursive_select_query_ that reference the _temporary view_ is excuted and the _temporary view_ is replaced by the result set of the _recursive_select_query_.
The execution of the _recursive_select_query_ is iterated until the result set is empty.
All the result sets are combined by the [UNION]({{ '/reference/set-operators.html#union' | relative_url }}) operator.

Example:
```sql
WITH RECURSIVE ct (n)
  AS (
    SELECT 1
    UNION ALL
    SELECT n + 1
      FROM ct
     WHERE n < 5
  )
SELECT n FROM ct;


/* Result Set
+---+
| n |
+---+
| 1 |
| 2 |
| 3 |
| 4 |
| 5 |
+---+
*/
```