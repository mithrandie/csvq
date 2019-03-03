---
layout: default
title: Alter Table Query - Reference Manual - csvq
category: reference
---

# Alter Table Query

Alter Table query is used to modify data structure.

* [ADD COLUMNS](#add-columns)
* [DROP COLUMNS](#drop-columns)
* [RENAME COLUMN](#rename-column)
* [SET ATTRIBUTE](#set-attribute)

## Add Columns
{: #add-columns}

```sql
ALTER TABLE table_name
  ADD column_name [DEFAULT value]
  [FIRST|LAST|AFTER column|BEFORE column]

ALTER TABLE table_name
  ADD (column_name [DEFAULT value] [, column_name [DEFAULT value] ...])
  [FIRST|LAST|AFTER column|BEFORE column]
```

_table_name_
: [identifier]({{ '/reference/statement.html#parsing' | relative_url }}) or [Table Object]({{ '/reference/select-query.html#from_clause' | relative_url }})

_column_name_
: [identifier]({{ '/reference/statement.html#parsing' | relative_url }})

_value_
: [value]({{ '/reference/value.html' | relative_url }})
  
  If default value is not specified, new fields are set null.

_column_
: [field reference]({{ '/reference/value.html#field_reference' | relative_url }})

_LAST_ is the default position.


## Drop Columns
{: #drop-columns}

```sql
ALTER TABLE table_name DROP column
ALTER TABLE table_name DROP (column [, column ...])
```

_table_name_
: [identifier]({{ '/reference/statement.html#parsing' | relative_url }}) or [Table Object]({{ '/reference/select-query.html#from_clause' | relative_url }})

_column_
: [field reference]({{ '/reference/value.html#field_reference' | relative_url }})

## Rename Column
{: #rename-column}

```sql
ALTER TABLE table_name RENAME old_column_name TO new_column_name
```

_table_name_
: [identifier]({{ '/reference/statement.html#parsing' | relative_url }}) or [Table Object]({{ '/reference/select-query.html#from_clause' | relative_url }})

_old_column_name_
: [field reference]({{ '/reference/value.html#field_reference' | relative_url }})

_new_column_name_
: [identifier]({{ '/reference/statement.html#parsing' | relative_url }})

## Set Attribute
{: #set-attribute}

Set file attributes. 
File attributes is used to create or update files by the results of queries.

Changes to the file attributes are retained until the end of the transaction.


```sql
ALTER TABLE table_name SET attribute TO value
```

_table_name_
: [identifier]({{ '/reference/statement.html#parsing' | relative_url }}) or [Table Object]({{ '/reference/select-query.html#from_clause' | relative_url }})

_attribute_
: [identifier]({{ '/reference/statement.html#parsing' | relative_url }})

  | name | type | description |
  | :- | :- | :- |
  | FORMAT              | string  | Format |
  | DELIMITER           | string  | Field delimiter for CSV |
  | DELIMITER_POSITIONS | string  | Delimiter positions for Fixed-Length Format |
  | JSON_ESCAPE         | string  | JSON escape type |
  | ENCODING            | string  | File Encoding |
  | LINE_BREAK          | string  | Line Break |
  | HEADER              | boolean | Write header line in the file |
  | ENCLOSE_ALL         | boolean | Enclose all string values in CSV |
  | PRETTY_PRINT        | boolean | Make JSON output easier to read |

_value_
: [value]({{ '/reference/value.html' | relative_url }}) or [identifier]({{ '/reference/statement.html#parsing' | relative_url }})
