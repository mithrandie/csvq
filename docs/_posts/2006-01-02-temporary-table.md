---
layout: default
title: Temporary Table - Reference Manual - csvq
category: reference
---

# Temporary Table

A Temporary Table that is treated as "VIEW" can be used in a procedure.
You can refer, insert, update, or delete records in temporary tables.

Temporary tables are affected by transactions.
When current transaction is rolled back, the records that saved at the previous commit are restored. 

## Declare Temporary Table
{: #declare}

### Declare Empty Table

```sql
DECLARE table_name VIEW (column_name [, column_name ...]);
```

_table_name_
: [identifier]({{ '/reference/statement.html#parsing' | relative_url }})

_column_name_
: [identifier]({{ '/reference/statement.html#parsing' | relative_url }})


### Declare from the Result-Set of a Select Query

```sql
DECLARE table_name VIEW [(column_name [, column_name ...])] AS select_query;
```

_table_name_
: [identifier]({{ '/reference/statement.html#parsing' | relative_url }})

_column_name_
: [identifier]({{ '/reference/statement.html#parsing' | relative_url }})

_select_query_
: [Select Query]({{ '/reference/select-query.html' | relative_url }})


## Dispose Temporary Table
{: #dispose}

```sql
DISPOSE VIEW table_name;
```

_table_name_
: [identifier]({{ '/reference/statement.html#parsing' | relative_url }})
