---
layout: default
title: Comparison Operators - Reference Manual - csvq
category: reference
---

# Comparison Operators

| name | description |
| :-: | :- |
| [Relational Operators](#relational_operators) | Compare values |
| [IS](#is)           | Compare a value with ternary value |
| [BETWEEN](#between) | Check if a value is with in a range of values |
| [IN](#in)           | Check if a value is within a set of values |
| [LIKE](#like)       | Check if a string matches a pattern |
| [ANY](#any)         | Check if any of values fulfill conditions |
| [ALL](#all)         | Check if all of values fulfill conditions |
| [EXISTS](#exists)   | Check if a subquery returns at least one row |

A comparison operator returns a ternary value.

## Relational Operators
{: #relational_operators}

| name | description |
| :-: | :- |
| \=      | Equal |
| <       | Less than |
| <\=     | Less than or equal |
| >       | Greater than |
| >\=     | Greater than or equal |
| <>, !\= | Not equal |

```sql
value operator value
```

_value_
: [value]({{ '/reference/value.html' | relative_url }})

At first, a relational operator attempt to convert both of operands to float values, and if convertions is successful then compare them.
If conversions failed, next a relational operater attempt to convert to datetime, and next to boolean, at last to string.

If either of operands is null or all conversion failed, return UNKNOWN.

## IS
{: #is}

```sql
value IS [NOT] ternary
value IS [NOT] NULL
```

_value_
: [value]({{ '/reference/value.html' | relative_url }})

_ternary_
: [ternary]({{ '/reference/value.html#ternary' | relative_url }})

Evaluate the ternary value of _value_ and check if the ternary value is equal to _ternary_.

## BETWEEN
{: #between}

```sql
value [NOT] BETWEEN low AND high
```

_value_
: [value]({{ '/reference/value.html' | relative_url }})

_low_
: [value]({{ '/reference/value.html' | relative_url }})

_high_
: [value]({{ '/reference/value.html' | relative_url }})

Check a _value_ is greater than or equal to _low_ and less than or equal to _high_.

The BETWEEN operation is equivalent to followings.
```sql
low <= value AND value <= high
NOT (low <= value AND value <= high)
```

## IN
{: #in}

```sql
value [NOT] IN (value [, value ...])
value [NOT] IN (select_query)
```

_value_
: [value]({{ '/reference/value.html' | relative_url }})

_select_query_
: [Select Query]({{ '/reference/select-query.html' | relative_url }})

Check if a _value_ is in within a set of _values_ or a result set of _select_query_.

The IN operation with subquery is equivalent to following ANY operation.
```sql
value = ANY (select_query)
```

## LIKE
{: #like}

```sql
string [NOT] LIKE pattern
```

_string_
: [string]({{ '/reference/value.html#string' | relative_url }})

_pattern_
: [string]({{ '/reference/value.html#string' | relative_url }})

Return TRUE if a _string_ matches a _pattern_, otherwise return FALSE.
If _string_ is null, return UNKNOWN. 

In pattern, following special characters are used.

%
: any number of characters.

_
: exactly one character

## ANY
{: #any}

```sql
value relational_operator ANY (select_query)
```

_value_
: [value]({{ '/reference/value.html' | relative_url }})

_relational_operator_
: [relational operator](#relational_operators)

_select_query_
: [Select Query]({{ '/reference/select-query.html' | relative_url }})

Compare a _value_ to each values retrieved by _select_query_.
If any of comparison results is TRUE, return TRUE.
If there is no TRUE result and there is at least one UNKNOWN result, return UNKNOWN.
Otherwise return FALSE.

If _select_query_ returns no values, return FALSE.

## ALL
{: #all}

```sql
value relational_operator ALL (select_query)
```

_value_
: [value]({{ '/reference/value.html' | relative_url }})

_relational_operator_
: [relational operator](#relational_operators)

_select_query_
: [Select Query]({{ '/reference/select-query.html' | relative_url }})

Compare a _value_ to every values retrieved by _select_query_.
If any of comparison results is FALSE, return FALSE.
If there is no FALSE result and there is at least one UNKNOWN result, return UNKNOWN.
Otherwise return TRUE.

If _select_query_ returns no values, return TRUE.

## Exists
{: #exists}

```sql
EXISTS (select_query)
```

_select_query_
: [Select Query]({{ '/reference/select-query.html' | relative_url }})

Return TRUE if a _select_query_ returns at least one row, otherwise return FALSE.
