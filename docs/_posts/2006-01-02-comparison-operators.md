---
layout: default
title: Comparison Operators - Reference Manual - csvq
category: reference
---

# Comparison Operators

| operator | description |
| :- | :- |
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

| operator | description |
| :- | :- |
| \=      | Equal to |
| <       | Less than |
| <\=     | Less than or equal to |
| >       | Greater than |
| >\=     | Greater than or equal to |
| <>, !\= | Not equal to |

```sql
relational_operation
  : value operator value
  | row_value operator row_value
```

_value_
: [value]({{ '/reference/value.html' | relative_url }})

_row_value_
: [Row Value]({{ '/reference/row-value.html' | relative_url }})

At first, a relational operator attempt to convert both of operands to float values, and if both convertions are successful then compare them.
If conversions failed, next a relational operater attempt to convert to datetime, and next to boolean, at last to string.

If either of operands is null or all conversion failed, return UNKNOWN.

In case of _row_values_ comparison, both of _row_values_ must be tha same lengths.
Values at the same indices are compared, and the result is decided by [AND]({{ '/reference/logic-operators.html#and' | relative_url }}) condition of each comparisons

## IS
{: #is}

```sql
is_operation
  : value IS [NOT] ternary
  | value IS [NOT] NULL
```

_value_
: [value]({{ '/reference/value.html' | relative_url }})

_ternary_
: [ternary]({{ '/reference/value.html#ternary' | relative_url }})

Evaluate the ternary value of _value_ and check if the ternary value is equal to _ternary_.

## BETWEEN
{: #between}

```sql
between_operation
  : value [NOT] BETWEEN low AND high
  | row_value [NOT] BETWEEN row_value_low AND row_value_high
```

_value_
: [value]({{ '/reference/value.html' | relative_url }})

_low_
: [value]({{ '/reference/value.html' | relative_url }})

_high_
: [value]({{ '/reference/value.html' | relative_url }})

_row_value_
: [Row Value]({{ '/reference/row-value.html' | relative_url }})

_row_value_low_
: [Row Value]({{ '/reference/row-value.html' | relative_url }})

_row_value_high_
: [Row Value]({{ '/reference/row-value.html' | relative_url }})

Check a _value_ is greater than or equal to _low_ and less than or equal to _high_.

The BETWEEN operation is equivalent to followings.
```sql
low <= value AND value <= high
NOT (low <= value AND value <= high)
```

## IN
{: #in}

```sql
in_operation
  : value [NOT] IN (value [, value ...])
  | value [NOT] IN single_field_subquery
  | row_value [NOT] IN (row_value [, row_value ...])
  | row_value [NOT] IN multiple_fields_subquery
```

_value_
: [value]({{ '/reference/value.html' | relative_url }})

_row_value_
: [Row Value]({{ '/reference/row-value.html' | relative_url }})

_single_field_subquery_
: [subquery]({{ '/reference/value.html#subquery' | relative_url }})

_multiple_fields_subquery_
: [subquery]({{ '/reference/value.html#subquery' | relative_url }})

Check if a _value_ or _row_value_ is in within a set of _values_ or a result set of _select_query_.

A IN operation is equivalent to a [ANY](#any) operation that _relational_operator_ is specified as "=".

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
any_operation
  : value relational_operator ANY (value [, value ...])
  | value relational_operator ANY single_field_subquery
  | row_value relational_operator ANY (row_value [, row_value ...])
  | row_value relational_operator ANY multiple_fields_subquery
```

_value_
: [value]({{ '/reference/value.html' | relative_url }})

_row_value_
: [Row Value]({{ '/reference/row-value.html' | relative_url }})

_relational_operator_
: [relational operator](#relational_operators)

_single_field_subquery_
: [subquery]({{ '/reference/value.html#subquery' | relative_url }})

_multiple_fields_subquery_
: [subquery]({{ '/reference/value.html#subquery' | relative_url }})

Compare a _value_ or _row_value_ to each listed _values_ or each records retrieved by _select_query_.
If any of comparison results is TRUE, return TRUE.
If there is no TRUE result and there is at least one UNKNOWN result, return UNKNOWN.
Otherwise return FALSE.

If _select_query_ returns no record, return FALSE.

## ALL
{: #all}

```sql
all_operation
  : value relational_operator ALL (value [, value ...])
  | value relational_operator ALL single_field_subquery
  | row_value relational_operator ALL (row_value [, row_value ...])
  | row_value relational_operator ALL multiple_fields_subquery
```

_value_
: [value]({{ '/reference/value.html' | relative_url }})

_row_value_
: [Row Value]({{ '/reference/row-value.html' | relative_url }})

_relational_operator_
: [relational operator](#relational_operators)

_single_field_subquery_
: [subquery]({{ '/reference/value.html#subquery' | relative_url }})

_multiple_fields_subquery_
: [subquery]({{ '/reference/value.html#subquery' | relative_url }})

Compare a _value_ or _row_value_ to every listed _values_ or each records retrieved by _select_query_.
If any of comparison results is FALSE, return FALSE.
If there is no FALSE result and there is at least one UNKNOWN result, return UNKNOWN.
Otherwise return TRUE.

If _select_query_ returns no record, return TRUE.

## Exists
{: #exists}

```sql
EXISTS (select_query)
```

_select_query_
: [Select Query]({{ '/reference/select-query.html' | relative_url }})

Return TRUE if a _select_query_ returns at least one record, otherwise return FALSE.
