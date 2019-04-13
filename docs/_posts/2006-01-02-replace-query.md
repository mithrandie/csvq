---
layout: default
title: Replace Query - Reference Manual - csvq
category: reference
---

# Replace Query

Replace query is used to insert or update records to a csv file.

If records matching _key_columns_ exist, the records are updated. Otherwise insert records.

## Insert or Update Values

```sql
[WITH common_table_expression [, common_table_expression ...]]
  REPLACE INTO table_name
  [(column [, column ...])]
  USING (key_column [, key_column ...]))
  VALUES row_value [, row_value ...]
```

_common_table_expression_
: [Common Table Expression]({{ '/reference/common-table-expression.html' | relative_url }})

_table_name_
: [identifier]({{ '/reference/statement.html#parsing' | relative_url }}) or [Table Object]({{ '/reference/select-query.html#from_clause' | relative_url }})

_column_
: [field reference]({{ '/reference/value.html#field_reference' | relative_url }})

_key_column_
: [field reference]({{ '/reference/value.html#field_reference' | relative_url }})

_row_value_
: [Row Value]({{ '/reference/row-value.html' | relative_url }})

## Insert or Update From Select Query

```sql
[WITH common_table_expression [, common_table_expression ...]]
  REPLACE INTO table_name
  [(column [, column ...])]
  USING (key_column [, key_column ...]))
  select_query
```

_common_table_expression_
: [Common Table Expression]({{ '/reference/common-table-expression.html' | relative_url }})

_table_name_
: [identifier]({{ '/reference/statement.html#parsing' | relative_url }}) or [Table Object]({{ '/reference/select-query.html#from_clause' | relative_url }})

_column_
: [field reference]({{ '/reference/value.html#field_reference' | relative_url }})

_key_column_
: [field reference]({{ '/reference/value.html#field_reference' | relative_url }})

_select_query_
: [Select Query]({{ '/reference/select-query.html' | relative_url }})
