---
layout: default
title: String Operators - Reference Manual - csvq
category: reference
---

# String Operators

| operator | description |
| :- | :- |
| \|\| | Concatnation |

## Syntax

```sql
string operator string
```

_string_
: [value]({{ '/reference/value.html' | relative_url }})

A string operator concatnate string values, and return a string value.
If each of operands is not a string value, the value is converted to a string value.

If either of operands is null or conversion to string failed, return null.
