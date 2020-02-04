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
| [LIKE](#like)       | Check if a string matches a pattern |
| [IN](#in)           | Check if a value is within a set of values |
| [ANY](#any)         | Check if any of values fulfill conditions |
| [ALL](#all)         | Check if all of values fulfill conditions |
| [EXISTS](#exists)   | Check if a subquery returns at least one row |

A comparison operator returns a ternary value.

## Relational Operators
{: #relational_operators}

| operator | description |
| :------- | :---------- |
| \=       | LHS is equal to RHS |
| \=\=     | LHS and RHS are of the same type, and LHS is equal to RHS |
| <        | LHS is less than RHS |
| <\=      | LHS is less than or equal to RHS |
| >        | LHS is greater than RHS |
| >\=      | LHS is greater than or equal to RHS |
| <>, !\=  | LHS is not equal to RHS |

```sql
relational_operation
  : value operator value
  | row_value operator row_value
```

_value_
: [value]({{ '/reference/value.html' | relative_url }})

_row_value_
: [Row Value]({{ '/reference/row-value.html' | relative_url }})

Except for identical operator("=="), at first, the relational operator attempts to convert both of operands to integer values, and if both conversions are successful then compares them.
If conversions failed, next the relational operater attempts to convert the values to float, and next to datetime, boolean, at last to string.

If either of operands is null or all conversions failed, then the comparison returns UNKNOWN.

Identical operator does not perform automatic type conversion.
The result will be true only when both operands are of the same type.

In case of _row_values_ comparison, both of _row_values_ must be tha same lengths.
Values at the same indices are compared in order from left to right.


## IS
{: #is}

```sql
value IS [NOT] NULL
```

_value_
: [value]({{ '/reference/value.html' | relative_url }})

Check if a _value_ is a null value.

```sql
value IS [NOT] ternary
```

_value_
: [value]({{ '/reference/value.html' | relative_url }})

_ternary_
: [ternary]({{ '/reference/value.html#ternary' | relative_url }})

Evaluate the ternary value of a _value_ and check if the ternary value is equal to _ternary_.

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

## LIKE
{: #like}

```sql
string [NOT] LIKE pattern
```

_string_
: [string]({{ '/reference/value.html#string' | relative_url }})

_pattern_
: [string]({{ '/reference/value.html#string' | relative_url }})

Returns TRUE if _string_ matches _pattern_, otherwise returns FALSE.
If _string_ is a null, return UNKNOWN. 

In _pattern_, following special characters can be used.

%
: any number of characters

_ (U+005F Low Line)
: exactly one character

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

Check if _value_ or _row_value_ is in within the set of _values_ or the result set of _select_query_.

_IN_ is equivalent to [= ANY](#any).

_NOT IN_ is equivalent to [<> ALL](#all).

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

Compare _value_ or _row_value_ to each listed _values_ or each records retrieved by _select_query_.
If any of comparison results is TRUE, returns TRUE.
If there is no TRUE result and there is at least one UNKNOWN result, returns UNKNOWN.
Otherwise returns FALSE.

If _select_query_ returns no record, returns FALSE.

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

Compare _value_ or _row_value_ to every listed _values_ or each records retrieved by _select_query_.
If any of comparison results is FALSE, returns FALSE.
If there is no FALSE result and there is at least one UNKNOWN result, returns UNKNOWN.
Otherwise returns TRUE.

If _select_query_ returns no record, returns TRUE.

## Exists
{: #exists}

```sql
EXISTS (select_query)
```

_select_query_
: [Select Query]({{ '/reference/select-query.html' | relative_url }})

Returns TRUE if a _select_query_ returns at least one record, otherwise returns FALSE.
