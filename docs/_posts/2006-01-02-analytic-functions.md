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
| [ROW_NUMBER](#row_number)   | Return sequential numbers |
| [RANK](#rank)               | Return ranks |
| [DENSE_RANK](#dense_rank)   | Return ranks without any gaps in the ranking |
| [FIRST_VALUE](#first_value) | Return the first value |
| [LAST_VALUE](#last_value)   | Return the last value |
| [LAG](#lag)                 | Return the difference of current row and previous row |
| [LEAD](#lead)               | Return the difference of current row and following row |
| [COUNT](#count)             | Return the number of values |
| [MIN](#min)                 | Return the minimum value |
| [MAX](#max)                 | Return the maximum value |
| [SUM](#sum)                 | Return the sum of values |
| [AVG](#avg)                 | Return the average of values |
| [LISTAGG](#listagg)         | Return the concatenated string of values |

## Basic Syntax
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
: [Order By Clause]({{ '/reference/select-query.html#order_by_clause' | relative_url }})

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
: [Order By Clause]({{ '/reference/select-query.html#order_by_clause' | relative_url }})

_return_
: [integer]({{ '/reference/value.html#integer' | relative_url }})

Return sequential numbers of records in a group.


### RANK
{: #rank}

```
RANK() OVER ([partition_clause] [order_by_clause])
```

_partition_clause_
: [Partition Clause](#syntax)

_order_by_clause_
: [Order By Clause]({{ '/reference/select-query.html#order_by_clause' | relative_url }})

_return_
: [integer]({{ '/reference/value.html#integer' | relative_url }})

Return ranks of records in a group.


### DENSE_RANK
{: #dense_rank}

```
DENSE_RANK() OVER ([partition_clause] [order_by_clause])
```

_partition_clause_
: [Partition Clause](#syntax)

_order_by_clause_
: [Order By Clause]({{ '/reference/select-query.html#order_by_clause' | relative_url }})

_return_
: [integer]({{ '/reference/value.html#integer' | relative_url }})

Return ranks of records without any gaps in the ranking in a group.


### FIRST_VALUE
{: #first_value}

```
FIRST_VALUE(expr [IGNORE NULLS]) OVER ([partition_clause] [order_by_clause])
```

_expr_
: [value]({{ '/reference/value.html' | relative_url }})

_partition_clause_
: [Partition Clause](#syntax)

_order_by_clause_
: [Order By Clause]({{ '/reference/select-query.html#order_by_clause' | relative_url }})

_return_
: [primitive value]({{ '/reference/value.html#primitive_types' | relative_url }})

Return the first value in a group.
If _IGNORE NULLS_ keywords are specified, then return the first value that is not a null.


### LAST_VALUE
{: #last_value}

```
LAST_VALUE(expr [IGNORE NULLS]) OVER ([partition_clause] [order by clause])
```

_expr_
: [value]({{ '/reference/value.html' | relative_url }})

_partition_clause_
: [Partition Clause](#syntax)

_order_by_clause_
: [Order By Clause]({{ '/reference/select-query.html#order_by_clause' | relative_url }})

_return_
: [primitive value]({{ '/reference/value.html#primitive_types' | relative_url }})

Return the last value in a group.
If _IGNORE NULLS_ keywords are specified, then return the last value that is not a null.


### LAG
{: #lag}

```
LAG(expr [, offset [, default]] [IGNORE NULLS]) OVER ([partition_clause] [order by clause])
```

_expr_
: [value]({{ '/reference/value.html' | relative_url }})

_offset_
: [integer]({{ '/reference/value.html#integer' | relative_url }})
  
  The number of rows from current row. The default is 1.

_default_
: [value]({{ '/reference/value.html' | relative_url }})

  The value to set when the offset row does not exist or the difference cannot be calculated.
  The default is NULL.

_partition_clause_
: [Partition Clause](#syntax)

_order_by_clause_
: [Order By Clause]({{ '/reference/select-query.html#order_by_clause' | relative_url }})

_return_
: [float]({{ '/reference/value.html#float' | relative_url }}) or [integer]({{ '/reference/value.html#integer' | relative_url }})

Return the difference of current row and previous row.
If _IGNORE NULLS_ keywords are specified, then rows that _expr_ values are nulls are skipped. 


### LEAD
{: #lead}

```
LEAD(expr [, offset [, default]] [IGNORE NULLS]) OVER ([partition_clause] [order by clause])
```

_expr_
: [value]({{ '/reference/value.html' | relative_url }})

_offset_
: [integer]({{ '/reference/value.html#integer' | relative_url }})
  
  The number of rows from current row. The default is 1.

_default_
: [value]({{ '/reference/value.html' | relative_url }})

  The value to set when the offset row does not exist or the difference cannot be calculated.
  The default is NULL.

_partition_clause_
: [Partition Clause](#syntax)

_order_by_clause_
: [Order By Clause]({{ '/reference/select-query.html#order_by_clause' | relative_url }})

_return_
: [float]({{ '/reference/value.html#float' | relative_url }}) or [integer]({{ '/reference/value.html#integer' | relative_url }})

Return the difference of current row and following row.
If _IGNORE NULLS_ keywords are specified, then rows that _expr_ values are nulls are skipped. 


### COUNT
{: #count}

```
COUNT([DISTINCT] expr) OVER ([partition_clause] [order by clause])
```

_expr_
: [value]({{ '/reference/value.html' | relative_url }})

_partition_clause_
: [Partition Clause](#syntax)

_order_by_clause_
: [Order By Clause]({{ '/reference/select-query.html#order_by_clause' | relative_url }})

_return_
: [integer]({{ '/reference/value.html#integer' | relative_url }})

Returns the number of non-null values of _expr_.

```
COUNT([DISTINCT] *) OVER ([partition_clause] [order by clause])
```

_partition_clause_
: [Partition Clause](#syntax)

_order_by_clause_
: [Order By Clause]({{ '/reference/select-query.html#order_by_clause' | relative_url }})

_return_
: [integer]({{ '/reference/value.html#integer' | relative_url }})

Return the number of all values including null values.


### MIN
{: #min}

```
MIN([DISTINCT] expr) OVER ([partition_clause] [order by clause])
```

_expr_
: [value]({{ '/reference/value.html' | relative_url }})

_partition_clause_
: [Partition Clause](#syntax)

_order_by_clause_
: [Order By Clause]({{ '/reference/select-query.html#order_by_clause' | relative_url }})

_return_
: [primitive value]({{ '/reference/value.html#primitive_types' | relative_url }})

Returns the minimum value of non-null values of _expr_.
If all values are null, return null.


### MAX
{: #max}

```
MAX([DISTINCT] expr) OVER ([partition_clause] [order by clause])
```

_expr_
: [value]({{ '/reference/value.html' | relative_url }})

_partition_clause_
: [Partition Clause](#syntax)

_order_by_clause_
: [Order By Clause]({{ '/reference/select-query.html#order_by_clause' | relative_url }})

_return_
: [primitive value]({{ '/reference/value.html#primitive_types' | relative_url }})

Returns the maximum value of non-null values of _expr_.
If all values are null, return null.


### SUM
{: #sum}

```
SUM([DISTINCT] expr) OVER ([partition_clause] [order by clause])
```

_expr_
: [value]({{ '/reference/value.html' | relative_url }})

_partition_clause_
: [Partition Clause](#syntax)

_order_by_clause_
: [Order By Clause]({{ '/reference/select-query.html#order_by_clause' | relative_url }})

_return_
: [float]({{ '/reference/value.html#float' | relative_url }}) or [integer]({{ '/reference/value.html#integer' | relative_url }})

Returns the sum of non-null values of _expr_.
If all values are null, return null.


### AVG
{: #avg}

```
AVG([DISTINCT] expr) OVER ([partition_clause] [order by clause])
```

_expr_
: [value]({{ '/reference/value.html' | relative_url }})

_partition_clause_
: [Partition Clause](#syntax)

_order_by_clause_
: [Order By Clause]({{ '/reference/select-query.html#order_by_clause' | relative_url }})

_return_
: [float]({{ '/reference/value.html#float' | relative_url }}) or [integer]({{ '/reference/value.html#integer' | relative_url }})

Returns the average of non-null values of _expr_.
If all values are null, return null.


### LISTAGG
{: #listagg}

```
LISTAGG([DISTINCT] expr [, separator]) OVER ([partition_clause] [order by clause])
```

_expr_
: [value]({{ '/reference/value.html' | relative_url }})

_separator_
: [string]({{ '/reference/value.html#string' | relative_url }})

_partition_clause_
: [Partition Clause](#syntax)

_order_by_clause_
: [Order By Clause]({{ '/reference/select-query.html#order_by_clause' | relative_url }})

Return the string result with the concatenated non-null values of _expr_.
If all values are null, return null.

Separator string _separator_ is placed between values. Empty string is the default.
