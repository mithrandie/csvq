---
layout: default
title: Analytic Functions - Reference Manual - csvq
category: reference
---

# Analytic Functions

Analytic functions calculate values of groups.
Analytic Functions can be used only in [Select Clause]({{ '/reference/select-query.html#select_clause' | relative_url }}) and [Order By Clause]({{ '/reference/select-query.html#order_by_clause' | relative_url }})

| name | description |
| :- | :- |
| [ROW_NUMBER](#row_number) | Return sequential numbers |
| [RANK](#rank)             | Return ranks |
| [DENSE_RANK](#dense_rank) | Return ranks without any gaps in the ranking |

## Syntax
{: #syntax}

```sql
analytic_function
  : function_name([args]) OVER ([partition_clause] [order_by_clause])

args
  : value [, value ...]

partition_clause
  : PARTITION BY value [, value ...]
```

_value_
: [value]({{ '/reference/value.html' | relative_url }})

_order_by_clause_
: [Order By Clause]({{ '/reference/select_query.html#order_by_clause' | relative_url }})

Analytic Functions sort the result set by _order_by_clause_ and calculate values within each of groups partitioned by _partition_clause_.
If there is no _partition_clause_, then all records of the result set are dealt with as one group. 

## Definitions

### ROW_NUMBER
{: #row_number}

```
ROW_NUMBER() OVER ([partition_clause] [order_by_clause])
```

_partition_clause_
: [Partition Clause](#syntax)

_order_by_clause_
: [Order By Clause]({{ '/reference/select_query.html#order_by_clause' | relative_url }})

Return sequential numbers of records in a group.


### RANK
{: #rank}

```
RANK() OVER ([partition_clause] [order_by_clause])
```

_partition_clause_
: [Partition Clause](#syntax)

_order_by_clause_
: [Order By Clause]({{ '/reference/select_query.html#order_by_clause' | relative_url }})

Return ranks of records in a group.


### DENSE_RANK
{: #dense_rank}

```
DENSE_RANK() OVER ([partition_clause] [order_by_clause])
```

_partition_clause_
: [Partition Clause](#syntax)

_order_by_clause_
: [Order By Clause]({{ '/reference/select_query.html#order_by_clause' | relative_url }})

Return ranks of records without any gaps in the ranking in a group.

