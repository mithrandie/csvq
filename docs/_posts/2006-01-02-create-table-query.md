---
layout: default
title: Create Table Query - Reference Manual - csvq
category: reference
---

# Create Table Query

Create Table query is used to create new csv files.

## Create Empty Table

```sql
CREATE TABLE file_path (column_name [, column_name ...])
```

_file_path_
: [identifier]({{ '/reference/statement.html#parsing' | relative_url }})

_column_name_
: [identifier]({{ '/reference/statement.html#parsing' | relative_url }})


## Create from the Result-Set of a Select Query

```sql
CREATE TABLE file_path [(column_name [, column_name ...])] [AS] select_query
```

_file_path_
: [identifier]({{ '/reference/statement.html#parsing' | relative_url }})

_column_name_
: [identifier]({{ '/reference/statement.html#parsing' | relative_url }})

_select_query_
: [Select Query]({{ '/reference/select-query.html' | relative_url }})
