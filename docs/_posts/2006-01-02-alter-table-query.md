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
  [FIRST|LAST|AFTER column_name|BEFORE column_name]

ALTER TABLE table_name
  ADD (column_name [DEFAULT value] [, column_name [DEFAULT value] ...])
  [FIRST|LAST|AFTER column_name|BEFORE column_name]
```

_table_name_
: [identifier]({{ '/reference/statement.html#parsing' | relative_url }})

_column_name_
: [identifier]({{ '/reference/statement.html#parsing' | relative_url }})

_value_
: [value]({{ '/reference/value.html' | relative_url }})
  
  If default value is not specified, new fields are set null.


## Drop Columns

```sql
ALTER TABLE table_name DROP column_name
ALTER TABLE table_name DROP (column_name, [column_name, ...])
```

_table_name_
: [identifier]({{ '/reference/statement.html#parsing' | relative_url }})

_column_name_
: [identifier]({{ '/reference/statement.html#parsing' | relative_url }})

## Rename Column

```sql
ALTER TABLE table_name RENAME old_column_name TO new_column_name
```

_table_name_
: [identifier]({{ '/reference/statement.html#parsing' | relative_url }})

_old_column_name_
: [identifier]({{ '/reference/statement.html#parsing' | relative_url }})

_new_column_name_
: [identifier]({{ '/reference/statement.html#parsing' | relative_url }})
