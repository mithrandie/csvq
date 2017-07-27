---
layout: default
title: Alter Table Query - Reference Manual - csvq
category: reference
---

# Alter Table Query

Alter Table query is used to modify data structure on csv files.

## Add Columns

```sql
ALTER TABLE table_name
  ADD column_name [DEFAULT value]
  [FIRST|LAST|AFTER column|BEFORE column]

ALTER TABLE table_name
  ADD (column_name [DEFAULT value] [, column_name [DEFAULT value] ...])
  [FIRST|LAST|AFTER column|BEFORE column]
```

_table_name_
: [identifier]({{ '/reference/statement.html#parsing' | relative_url }})

_column_name_
: [identifier]({{ '/reference/statement.html#parsing' | relative_url }})

_value_
: [value]({{ '/reference/value.html' | relative_url }})
  
  If default value is not specified, new fields are set null.

_column_
: [field reference]({{ '/reference/value.html#field_reference' | relative_url }})


## Drop Columns

```sql
ALTER TABLE table_name DROP column
ALTER TABLE table_name DROP (column, [column, ...])
```

_table_name_
: [identifier]({{ '/reference/statement.html#parsing' | relative_url }})

_column_
: [field reference]({{ '/reference/value.html#field_reference' | relative_url }})

## Rename Column

```sql
ALTER TABLE table_name RENAME old_column TO new_column_name
```

_table_name_
: [identifier]({{ '/reference/statement.html#parsing' | relative_url }})

_old_column_
: [field reference]({{ '/reference/value.html#field_reference' | relative_url }})

_new_column_name_
: [identifier]({{ '/reference/statement.html#parsing' | relative_url }})
