---
layout: default
title: Delete Query - Reference Manual - csvq
category: reference
---

# Delete Query

Delete query is used to delete records on csv files.

## Delete on a sigle file.

```sql
[with_clause] DELETE
  FROM table_name
  [where_clause]
```

_with_clause_
: [Common Table Expression]({{ '/reference/common-table-expression.html' | relative_url }})

_table_name_
: [identifier]({{ '/reference/statement.html#parsing' | relative_url }})

_where_clause_
: [Where Clause]({{ '/reference/select-query.html#where_clause' | relative_url }})

## Delete on multiple files

```sql
[with_clause] DELETE table_name [, table_name ...]
  from_clause
  [where_clause]
```

_with_clause_
: [Common Table Expression]({{ '/reference/common-table-expression.html' | relative_url }})

_table_name_
: [identifier]({{ '/reference/statement.html#parsing' | relative_url }})
  
  _table_name_ is not a file path, it is any one of table name aliases specified in _from_clause_. 

_from_clause_
: [From Clause]({{ '/reference/select-query.html#from_clause' | relative_url }})

_where_clause_
: [Where Clause]({{ '/reference/select-query.html#where_clause' | relative_url }})
