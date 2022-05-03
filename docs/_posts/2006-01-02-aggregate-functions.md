---
layout: default
title: Aggregate Functions - Reference Manual - csvq
category: reference
---

# Aggregate Functions

Aggregate functions calculate groupd records retrieved by a select query.
If records are not grouped, all records are dealt with as one group.

If distinct option is specified, aggregate functions calculate only unique values.

Aggregate Functions can be used only in [Select Clause]({{ '/reference/select-query.html#select_clause' | relative_url }}), [Having Clause]({{ '/reference/select-query.html#having_clause' | relative_url }}) and [Order By Clause]({{ '/reference/select-query.html#order_by_clause' | relative_url }})


| name | description |
| :- | :- |
| [COUNT](#count)       | Return the number of values |
| [MIN](#min)           | Return the minimum value |
| [MAX](#max)           | Return the maximum value |
| [SUM](#sum)           | Return the sum of values |
| [AVG](#avg)           | Return the average of values |
| [STDEV](#stdev)       | Return the sample standard deviation of values |
| [STDEVP](#stdevp)     | Return the population standard deviation of values |
| [VAR](#var)           | Return the sample variance of values |
| [VARP](#varp)         | Return the population variance of values |
| [MEDIAN](#median)     | Return the median of values |
| [LISTAGG](#listagg)   | Return the concatenated string of values |
| [JSON_AGG](#json_agg) | Return the string formatted in JSON array |

## Definitions

### COUNT
{: #count}

```
COUNT([DISTINCT] expr)
```

_expr_
: [value]({{ '/reference/value.html' | relative_url }})

_return_
: [integer]({{ '/reference/value.html#integer' | relative_url }})

Returns the number of non-null values of _expr_.

```
COUNT([DISTINCT] *)
```

_return_
: [integer]({{ '/reference/value.html#integer' | relative_url }})

Returns the number of all values including null values.

### MIN
{: #min}

```
MIN(expr)
```

_expr_
: [value]({{ '/reference/value.html' | relative_url }})

_return_
: [primitive type]({{ '/reference/value.html#primitive_types' | relative_url }})

Returns the minimum value of non-null values of _expr_.
If all values are null, then returns null.

### MAX
{: #max}

```
MAX(expr)
```

_expr_
: [value]({{ '/reference/value.html' | relative_url }})

_return_
: [primitive type]({{ '/reference/value.html#primitive_types' | relative_url }})

Returns the maximum value of non-null values of _expr_.
If all values are null, then return a null.

### SUM
{: #sum}

```
SUM([DISTINCT] expr)
```

_expr_
: [value]({{ '/reference/value.html' | relative_url }})

_return_
: [float]({{ '/reference/value.html#float' | relative_url }})

Returns the sum of float values of _expr_.
If all values are null, then returns a null.

### AVG
{: #avg}

```
AVG([DISTINCT] expr)
```

_expr_
: [value]({{ '/reference/value.html' | relative_url }})

_return_
: [float]({{ '/reference/value.html#float' | relative_url }})

Returns the average of float values of _expr_.
If all values are null, then returns a null.

### STDEV
{: #stdev}

```
STDEV([DISTINCT] expr)
```

_expr_
: [value]({{ '/reference/value.html' | relative_url }})

_return_
: [float]({{ '/reference/value.html#float' | relative_url }})

Returns the sample standard deviation of float values of _expr_.
If all values are null, then returns a null.

### STDEVP
{: #stdevp}

```
STDEVP([DISTINCT] expr)
```

_expr_
: [value]({{ '/reference/value.html' | relative_url }})

_return_
: [float]({{ '/reference/value.html#float' | relative_url }})

Returns the population standard deviation of float values of _expr_.
If all values are null, then returns a null.

### VAR
{: #var}

```
VAR([DISTINCT] expr)
```

_expr_
: [value]({{ '/reference/value.html' | relative_url }})

_return_
: [float]({{ '/reference/value.html#float' | relative_url }})

Returns the sample variance of float values of _expr_.
If all values are null, then returns a null.


### VARP
{: #varp}

```
VARP([DISTINCT] expr)
```

_expr_
: [value]({{ '/reference/value.html' | relative_url }})

_return_
: [float]({{ '/reference/value.html#float' | relative_url }})

Returns the population variance of float values of _expr_.
If all values are null, then returns a null.


### MEDIAN
{: #median}

```
MEDIAN([DISTINCT] expr)
```

_expr_
: [value]({{ '/reference/value.html' | relative_url }})

_return_
: [float]({{ '/reference/value.html#float' | relative_url }})

Returns the median of float or datetime values of _expr_.
If all values are null, then returns a null.

Even if _expr_ represents datetime values, this function returns a float or integer value.
The return value can be converted to a datetime value by using the [DATETIME function]({{ '/reference/cast-functions.html#datetime' | relative_url }}).

### LISTAGG
{: #listagg}

```
LISTAGG([DISTINCT] expr [, separator]) [WITHIN GROUP (order_by_clause)]
```

_expr_
: [value]({{ '/reference/value.html' | relative_url }})

_separator_
: [string]({{ '/reference/value.html#string' | relative_url }})

_order_by_clause_
: [Order By Clause]({{ '/reference/select-query.html#order_by_clause' | relative_url }})

_return_
: [string]({{ '/reference/value.html#string' | relative_url }})

Returns the string result with the concatenated non-null values of _expr_.
If all values are null, then returns a null.

_separator_ is placed between values. Empty string is the default.
By using _order_by_clause_, you can sort values.

### JSON_AGG
{: #json_agg}

```
JSON_AGG([DISTINCT] expr) [WITHIN GROUP (order_by_clause)]
```

_expr_
: [value]({{ '/reference/value.html' | relative_url }})

_order_by_clause_
: [Order By Clause]({{ '/reference/select-query.html#order_by_clause' | relative_url }})

_return_
: [string]({{ '/reference/value.html#string' | relative_url }})

Returns the string formatted in JSON array of _expr_.
