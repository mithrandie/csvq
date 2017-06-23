---
layout: default
title: Logic Operators - Reference Manual - csvq
category: reference
---

# Logic Operators

| name | description |
| :- | :- |
| [AND](#and) | Logical AND |
| [OR](#or)   | Logical OR |
| [NOT](#not) | Logical NOT |

A logic operator returns a ternary value.

## AND
{: #and}

```sql
condition AND condition
```

_condition_
: [value]({{ '/reference/value.html' | relative_url }})

If either of conditions is FALSE, then return FALSE.
If both of conditions are not FALSE, and either of conditions is UNKNOWN, then return UNKNOWN.
Otherwise return TRUE.

## OR
{: #or}

```sql
condition OR condition
```

_condition_
: [value]({{ '/reference/value.html' | relative_url }})

If either of conditions is TRUE, then return TRUE.
If both of conditions are not TRUE, and either of conditions is UNKNOWN, then return UNKNOWN.
Otherwise return FALSE.

## NOT
{: #not}

```sql
NOT condition
```

_condition_
: [value]({{ '/reference/value.html' | relative_url }})

If the condition is TRUE, return FALSE.
If the condition is FALSE, return TRUE.
IF the condition is UNKNOWN, return UNKNOWN.