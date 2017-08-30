---
layout: default
title: Temporary Table - Reference Manual - csvq
category: reference
---

# Temporary Table

A Temporary Table is a temporary view that can be used in a procedure.
You can refer, insert, update, or delete records in temporary tables.

When current transaction is rolled back, all of the changes in temporary tables are discarded, and the records that are created in the declarations are set to the tables.

## Declare Temporary Table
{: #declare}

### Declare Empty Table

```sql
DECLARE table_name TABLE (column_name [, column_name ...]);
```

_table_name_
: [identifier]({{ '/reference/statement.html#parsing' | relative_url }})

_column_name_
: [identifier]({{ '/reference/statement.html#parsing' | relative_url }})


### Declare from the Result-Set of a Select Query

```sql
DECLARE table_name TABLE [(column_name [, column_name ...])] AS select_query;
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
DISPOSE TABLE table_name;
```

_table_name_
: [identifier]({{ '/reference/statement.html#parsing' | relative_url }})
