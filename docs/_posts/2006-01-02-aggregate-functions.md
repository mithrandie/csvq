---
layout: default
title: Aggregate Functions - Reference Manual - csvq
category: reference
---

# Aggregate Functions

Aggregate functions calculate groupd records retrieved by a select query.
If records are not grouped, all records are dealt with as one group.

If distinct option is specified, aggregate functions calculate only unique values.

| name | description |
| :- | :- |
| [COUNT](#count) | Return the number of rows |
| [MAX](#max) | |
| [MIN](#min) | |
| [SUM](#sum) | |
| [AVG](#avg) | |
| [GROUP_CONCAT](#group_concat) | |

## Definitions

### COUNT
{: #count}

```
COUNT([DISTINCT] expr) return integer
```

Returns a number of non-null values of expr.

```
COUNT([DISTINCT] *) return integer
```

Return the number of values of expr including null values.

### MAX
{: #max}

```
MAX([DISTINCT] expr) return value
```

Returns the maximum value of non-null values of expr.
If all values are null, return null.

### MIN
{: #min}

```
MIN([DISTINCT] expr) return value
```

Returns the minimum value of non-null values of expr.
If all values are null, return null.

### SUM
{: #sum}

```
SUM([DISTINCT] expr) return decimal
```

Returns the sum of non-null values of expr.
If all values are null, return null.

### AVG
{: #avg}

```
AVG([DISTINCT] expr) return decimal
```

Returns the average of non-null values of expr.
If all values are null, return null.

### GROUP_CONCAT
{: #group_concat}

```
GROUP_CONCAT([DISTINCT] expr [order_by_clause] [SEPARATOR string]) return string
```

Return the string result with the concatenated non-null values of expr.
If all values are null, return null.

Separator string is placed between values. Default separator string is empty string.

By using order by clause, you can sort values.