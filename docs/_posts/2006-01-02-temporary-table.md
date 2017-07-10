---
layout: default
title: Temporary Table - Reference Manual - csvq
category: reference
---

# Temporary Table

A Temporary Table is a temporary view that can be used in a procedure.
You can refer, insert, update, or delete temporary tables.

## Declare Temporary Table
{: #declare}

```sql
temporary_table_declaration
  : DECLARE table_name TABLE (column_name [, column_name ...]);
  | DECLARE table_name TABLE [(column_name [, column_name ...])] FOR select_query;
```

_table_name_
: [identifier]({{ '/reference/statement.html#parsing' | relative_url }})

_column_name_
: [identifier]({{ '/reference/statement.html#parsing' | relative_url }})

_select_query_
: [Select Query]({{ '/reference/select-query.html' | relative_url }})
