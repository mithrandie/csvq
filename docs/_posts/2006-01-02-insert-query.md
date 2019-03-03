---
layout: default
title: Insert Query - Reference Manual - csvq
category: reference
---

# Insert Query

Insert query is used to insert records to a csv file.

## Insert Values

```sql
[WITH common_table_expression [, common_table_expression ...]]
  INSERT INTO table_name
  [(column [, column ...])]
  VALUES row_value [, row_value ...]
```

_common_table_expression_
: [Common Table Expression]({{ '/reference/common-table-expression.html' | relative_url }})

_table_name_
: [identifier]({{ '/reference/statement.html#parsing' | relative_url }}) or [Table Object]({{ '/reference/select-query.html#from_clause' | relative_url }})

_column_
: [field reference]({{ '/reference/value.html#field_reference' | relative_url }})

_row_value_
: [Row Value]({{ '/reference/row-value.html' | relative_url }})

## Insert From Select Query

```sql
[WITH common_table_expression [, common_table_expression ...]]
  INSERT INTO table_name
  [(column [, column ...])]
  select_query
```

_common_table_expression_
: [Common Table Expression]({{ '/reference/common-table-expression.html' | relative_url }})

_table_name_
: [identifier]({{ '/reference/statement.html#parsing' | relative_url }}) or [Table Object]({{ '/reference/select-query.html#from_clause' | relative_url }})

_column_
: [field reference]({{ '/reference/value.html#field_reference' | relative_url }})

_select_query_
: [Select Query]({{ '/reference/select-query.html' | relative_url }})
