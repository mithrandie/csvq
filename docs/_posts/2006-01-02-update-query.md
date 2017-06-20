---
layout: default
title: Update Query - Reference Manual - csvq
category: reference
---

# Update Query

Update query is used to update records on csv files.

## Update on a single file

```sql
UPDATE table_name
  SET column_name = value [, column_name = value ...]
  [where_clause]
```

_table_name_
: [identifier]({{ '/reference/statement.html#parsing' | relative_url }})

_column_name_
: [field reference]({{ '/reference/value.html#field_reference' | relative_url }})

_value_
: [value]({{ '/reference/value.html' | relative_url }})

_where_clause_
: [Where Clause]({{ '/reference/select-query.html#where_clause' | relative_url }})

## Update on multiple files

```sql
UPDATE table_name [, table_name ...]
  SET column_name = value [, column_name = value ...]
  from_clause
  [where_clause]
```

_table_name_
: [identifier]({{ '/reference/statement.html#parsing' | relative_url }})
  
  _table_name_ is not a file path, it is any one of table name aliases specified in _from_clause_. 

_column_name_
: [field reference]({{ '/reference/value.html#field_reference' | relative_url }})

_value_
: [value]({{ '/reference/value.html' | relative_url }})

_from_clause_
: [From Clause]({{ '/reference/select-query.html#from_clause' | relative_url }})

_where_clause_
: [Where Clause]({{ '/reference/select-query.html#where_clause' | relative_url }})
