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
| [ROW_NUMBER](#row_number)     | Returns the sequential numbers |
| [RANK](#rank)                 | Returns the ranks |
| [DENSE_RANK](#dense_rank)     | Returns the ranks without any gaps in the ranking |
| [CUME_DIST](#cume_dist)       | Returns the cumulative distributions |
| [PERCENT_RANK](#percent_rank) | Returns the relative ranks |
| [NTILE](#ntile)               | Returns the number of the group |
| [FIRST_VALUE](#first_value)   | Returns the first value |
| [LAST_VALUE](#last_value)     | Returns the last value |
| [NTH_VALUE](#nth_value)       | Returns the n-th value |
| [LAG](#lag)                   | Returns the value in a previous row |
| [LEAD](#lead)                 | Returns the value in a following row |
| [COUNT](#count)               | Returns the number of values |
| [MIN](#min)                   | Returns the minimum value |
| [MAX](#max)                   | Returns the maximum value |
| [SUM](#sum)                   | Returns the sum of values |
| [AVG](#avg)                   | Returns the average of values |
| [MEDIAN](#median)             | Returns the median of values |
| [LISTAGG](#listagg)           | Returns the concatenated string of values |
| [JSON_AGG](#json_agg)         | Returns the string formatted in JSON array |

## Basic Syntax
{: #syntax}

```sql
analytic_function
  : function_name([args]) OVER ([partition_clause] [order_by_clause [windowing_clause]])

args
  : value [, value ...]

partition_clause
  : PARTITION BY value [, value ...]

windowing_clause
  : ROWS window_position
  | ROWS BETWEEN window_frame_low AND window_frame_high

window_position
  : {UNBOUNDED PRECEDING|offset PRECEDING|CURRENT ROW}

window_frame_low
  : {UNBOUNDED PRECEDING|offset PRECEDING|offset FOLLOWING|CURRENT_ROW}

window_frame_high
  : {UNBOUNDED FOLLOWING|offset PRECEDING|offset FOLLOWING|CURRENT_ROW}

```

_value_
: [value]({{ '/reference/value.html' | relative_url }})

_order_by_clause_
: [Order By Clause]({{ '/reference/select-query.html#order_by_clause' | relative_url }})

_offset_
: [integer]({{ '/reference/value.html#integer' | relative_url }})

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

Returns the sequential numbers of records in a group.


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

Returns the ranks of records in a group.


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

Returns the ranks of records without any gaps in the ranking in a group.


### CUME_DIST
{: #cume_dist}

```
CUME_DIST() OVER ([partition_clause] [order_by_clause])
```

_partition_clause_
: [Partition Clause](#syntax)

_order_by_clause_
: [Order By Clause]({{ '/reference/select-query.html#order_by_clause' | relative_url }})

_return_
: [float]({{ '/reference/value.html#float' | relative_url }})

Returns the cumulative distributions in a group.
The return value is greater than 0 and less than or equal to 1.


### PERCENT_RANK
{: #percent_rank}

```
PERCENT_RANK() OVER ([partition_clause] [order_by_clause])
```

_partition_clause_
: [Partition Clause](#syntax)

_order_by_clause_
: [Order By Clause]({{ '/reference/select-query.html#order_by_clause' | relative_url }})

_return_
: [float]({{ '/reference/value.html#float' | relative_url }})

Returns the relative ranks in a group.
The return value is greater than or equal to 0 and less than or equal to 1.


### NTILE
{: #ntile}

```
NTILE(number_of_groups) OVER ([partition_clause] [order_by_clause])
```

_number_of_groups_
: [integer]({{ '/reference/value.html#integer' | relative_url }})

_partition_clause_
: [Partition Clause](#syntax)

_order_by_clause_
: [Order By Clause]({{ '/reference/select-query.html#order_by_clause' | relative_url }})

_return_
: [integer]({{ '/reference/value.html#integer' | relative_url }})

The NTILE function splits the records into _number_of_groups_ groups, then returns the sequential numbers of the groups.


### FIRST_VALUE
{: #first_value}

```
FIRST_VALUE(expr) [IGNORE NULLS] OVER ([partition_clause] [order_by_clause [windowing_clause]])
```

_expr_
: [value]({{ '/reference/value.html' | relative_url }})

_partition_clause_
: [Partition Clause](#syntax)

_order_by_clause_
: [Order By Clause]({{ '/reference/select-query.html#order_by_clause' | relative_url }})

_return_
: [primitive type]({{ '/reference/value.html#primitive_types' | relative_url }})

Returns the first value in a group.
If _IGNORE NULLS_ keywords are specified, then returns the first value that is not a null.


### LAST_VALUE
{: #last_value}

```
LAST_VALUE(expr) [IGNORE NULLS] OVER ([partition_clause] [order_by_clause [windowing_clause]])
```

_expr_
: [value]({{ '/reference/value.html' | relative_url }})

_partition_clause_
: [Partition Clause](#syntax)

_order_by_clause_
: [Order By Clause]({{ '/reference/select-query.html#order_by_clause' | relative_url }})

_return_
: [primitive type]({{ '/reference/value.html#primitive_types' | relative_url }})

Returns the last value in a group.
If _IGNORE NULLS_ keywords are specified, then returns the last value that is not a null.


### NTH_VALUE
{: #nth_value}

```
NTH_VALUE(expr, n) [IGNORE NULLS] OVER ([partition_clause] [order_by_clause [windowing_clause]])
```

_expr_
: [value]({{ '/reference/value.html' | relative_url }})

_n_
: [integer]({{ '/reference/value.html#integer' | relative_url }})

_partition_clause_
: [Partition Clause](#syntax)

_order_by_clause_
: [Order By Clause]({{ '/reference/select-query.html#order_by_clause' | relative_url }})

_return_
: [primitive type]({{ '/reference/value.html#primitive_types' | relative_url }})

Returns the _n_-th value in a group.
If _IGNORE NULLS_ keywords are specified, then returns the _n_-th value excluding null values.


### LAG
{: #lag}

```
LAG(expr [, offset [, default]]) [IGNORE NULLS] OVER ([partition_clause] [order by clause])
```

_expr_
: [value]({{ '/reference/value.html' | relative_url }})

_offset_
: [integer]({{ '/reference/value.html#integer' | relative_url }})
  
  The number of rows from current row. The default is 1.

_default_
: [value]({{ '/reference/value.html' | relative_url }})

  The value to set when the offset row does not exist.
  The default is NULL.

_partition_clause_
: [Partition Clause](#syntax)

_order_by_clause_
: [Order By Clause]({{ '/reference/select-query.html#order_by_clause' | relative_url }})

_return_
: [float]({{ '/reference/value.html#float' | relative_url }}) or [integer]({{ '/reference/value.html#integer' | relative_url }})

Returns the value in a previous row.
If _IGNORE NULLS_ keywords are specified, then rows that _expr_ values are nulls are skipped. 


### LEAD
{: #lead}

```
LEAD(expr [, offset [, default]]) [IGNORE NULLS] OVER ([partition_clause] [order by clause])
```

_expr_
: [value]({{ '/reference/value.html' | relative_url }})

_offset_
: [integer]({{ '/reference/value.html#integer' | relative_url }})
  
  The number of rows from current row. The default is 1.

_default_
: [value]({{ '/reference/value.html' | relative_url }})

  The value to set when the offset row does not exist.
  The default is NULL.

_partition_clause_
: [Partition Clause](#syntax)

_order_by_clause_
: [Order By Clause]({{ '/reference/select-query.html#order_by_clause' | relative_url }})

_return_
: [float]({{ '/reference/value.html#float' | relative_url }}) or [integer]({{ '/reference/value.html#integer' | relative_url }})

Returns the value in a following row.
If _IGNORE NULLS_ keywords are specified, then rows that _expr_ values are nulls are skipped. 


### COUNT
{: #count}

```
COUNT([DISTINCT] expr) OVER ([partition_clause] [order_by_clause [windowing_clause]])
```

_expr_
: [value]({{ '/reference/value.html' | relative_url }})

_partition_clause_
: [Partition Clause](#syntax)

_return_
: [integer]({{ '/reference/value.html#integer' | relative_url }})

Returns the number of non-null values of _expr_.

```
COUNT([DISTINCT] *) OVER ([partition_clause])
```

_partition_clause_
: [Partition Clause](#syntax)

_return_
: [integer]({{ '/reference/value.html#integer' | relative_url }})

Returns the number of all values including null values.


### MIN
{: #min}

```
MIN(expr) OVER ([partition_clause] [order_by_clause [windowing_clause]])
```

_expr_
: [value]({{ '/reference/value.html' | relative_url }})

_partition_clause_
: [Partition Clause](#syntax)

_return_
: [primitive type]({{ '/reference/value.html#primitive_types' | relative_url }})

Returns the minimum value of non-null values of _expr_.
If all values are null, then returns a null.


### MAX
{: #max}

```
MAX(expr) OVER ([partition_clause] [order_by_clause [windowing_clause]])
```

_expr_
: [value]({{ '/reference/value.html' | relative_url }})

_partition_clause_
: [Partition Clause](#syntax)

_return_
: [primitive type]({{ '/reference/value.html#primitive_types' | relative_url }})

Returns the maximum value of non-null values of _expr_.
If all values are null, then returns a null.


### SUM
{: #sum}

```
SUM([DISTINCT] expr) OVER ([partition_clause] [order_by_clause [windowing_clause]])
```

_expr_
: [value]({{ '/reference/value.html' | relative_url }})

_partition_clause_
: [Partition Clause](#syntax)

_return_
: [float]({{ '/reference/value.html#float' | relative_url }}) or [integer]({{ '/reference/value.html#integer' | relative_url }})

Returns the sum of float values of _expr_.
If all values are null, then returns a null.


### AVG
{: #avg}

```
AVG([DISTINCT] expr) OVER ([partition_clause] [order_by_clause [windowing_clause]])
```

_expr_
: [value]({{ '/reference/value.html' | relative_url }})

_partition_clause_
: [Partition Clause](#syntax)

_return_
: [float]({{ '/reference/value.html#float' | relative_url }}) or [integer]({{ '/reference/value.html#integer' | relative_url }})

Returns the average of float values of _expr_.
If all values are null, then returns a null.


### MEDIAN
{: #median}

```
MEDIAN([DISTINCT] expr) OVER ([partition_clause] [order_by_clause [windowing_clause]])
```

_expr_
: [value]({{ '/reference/value.html' | relative_url }})

_partition_clause_
: [Partition Clause](#syntax)

_return_
: [float]({{ '/reference/value.html#float' | relative_url }}) or [integer]({{ '/reference/value.html#integer' | relative_url }})

Returns the median of float or datetime values of _expr_.
If all values are null, then returns a null.

Even if _expr_ values are datetime values, the _MEDIAN_ function returns a float or integer value.
The return value can be converted to a datetime value by using the [DATETIME function]({{ '/reference/cast-functions.html#datetime' | relative_url }}).


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

Returns the string result with the concatenated non-null values of _expr_.
If all values are null, then returns a null.

Separator string _separator_ is placed between values. Empty string is the default.



### JSON_AGG
{: #json_agg}

```
JSON_AGG([DISTINCT] expr) OVER ([partition_clause] [order by clause])
```

_expr_
: [value]({{ '/reference/value.html' | relative_url }})

_partition_clause_
: [Partition Clause](#syntax)

_order_by_clause_
: [Order By Clause]({{ '/reference/select-query.html#order_by_clause' | relative_url }})

Returns the string formatted in JSON array.
