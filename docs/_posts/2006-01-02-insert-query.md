---
layout: default
title: Insert Query - Reference Manual - csvq
category: reference
---

# Insert Query

Insert query is used to insert records to a csv file.

## Insert Values

```sql
INSERT INTO table_name
  [(column_name [, column_name ...])]
  VALUES (value [, value ...]) [, (value [, value ...]) ...]
```

_table_name_
: [identifier]({{ '/reference/statement.html#parsing' | relative_url }})

_column_name_
: [identifier]({{ '/reference/statement.html#parsing' | relative_url }})

_value_
: [value]({{ '/reference/value.html' | relative_url }})

## Insert From Select Query

```sql
INSERT INTO table_name
  [(column_name [, column_name ...])]
  select_query
```

_table_name_
: [identifier]({{ '/reference/statement.html#parsing' | relative_url }})

_column_name_
: [identifier]({{ '/reference/statement.html#parsing' | relative_url }})

_select_query_
: [Select Query]({{ '/reference/select-query.html' | relative_url }})
