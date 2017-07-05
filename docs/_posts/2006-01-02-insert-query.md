---
layout: default
title: Insert Query - Reference Manual - csvq
category: reference
---

# Insert Query

Insert query is used to insert records to a csv file.

## Insert Values

```sql
[common_table_clause] INSERT INTO table_name
  [(column_name [, column_name ...])]
  VALUES row_value [, row_value ...]
```

_common_table_clause_
: [Common Table Clause]({{ '/reference/common-table.html' | relative_url }})

_table_name_
: [identifier]({{ '/reference/statement.html#parsing' | relative_url }})

_column_name_
: [identifier]({{ '/reference/statement.html#parsing' | relative_url }})

_row_value_
: [Row Value]({{ '/reference/row-value.html' | relative_url }})

## Insert From Select Query

```sql
[common_table_clause] INSERT INTO table_name
  [(column_name [, column_name ...])]
  select_query
```

_common_table_clause_
: [Common Table Clause]({{ '/reference/common-table.html' | relative_url }})

_table_name_
: [identifier]({{ '/reference/statement.html#parsing' | relative_url }})

_column_name_
: [identifier]({{ '/reference/statement.html#parsing' | relative_url }})

_select_query_
: [Select Query]({{ '/reference/select-query.html' | relative_url }})
