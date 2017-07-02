---
layout: default
title: Arithmetic Operators - Reference Manual - csvq
category: reference
---

# Arithmetic Operators

| operator | description |
| :- | :- |
| +  | Addition |
| \- | Subtraction |
| *  | Multiplication |
| /  | Division |
| %  | Modulo |

## Syntax

```sql
float operator float
```

_float_
: [value]({{ '/reference/value.html' | relative_url }})

An arithmetic operator calculate float values, and return a integer or float value.
If each of operands is not a float value, the value is converted to a float value.

If either of operands is null or conversions to float failed, return null.
