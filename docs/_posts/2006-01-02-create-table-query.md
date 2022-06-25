---
layout: default
title: Create Table Query - Reference Manual - csvq
category: reference
---

# Create Table Query

Create Table query is used to create new csv files.

## Create Empty Table

```sql
CREATE TABLE [IF NOT EXISTS] file_path (column_name [, column_name ...])
```

_file_path_
: [identifier]({{ '/reference/statement.html#parsing' | relative_url }})

_column_name_
: [identifier]({{ '/reference/statement.html#parsing' | relative_url }})

If the _IF NOT EXISTS_ clause is specified and the file already exists, no operation is performed.
In this case, an error is raised if the specified columns are different from the existing file.

## Create from the Result-Set of a Select Query

```sql
CREATE TABLE [IF NOT EXISTS] file_path [(column_name [, column_name ...])] [AS] select_query
```

_file_path_
: [identifier]({{ '/reference/statement.html#parsing' | relative_url }})

_column_name_
: [identifier]({{ '/reference/statement.html#parsing' | relative_url }})

_select_query_
: [Select Query]({{ '/reference/select-query.html' | relative_url }})

If the _IF NOT EXISTS_ clause is specified and the file already exists, no operation is performed.
In this case, an error is raised if the specified columns are different from the existing file.
