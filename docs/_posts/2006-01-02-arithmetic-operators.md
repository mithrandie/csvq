---
layout: default
title: Arithmetic Operators - Reference Manual - csvq
category: reference
---

# Arithmetic Operators

## Binary Operators
{: #binary}

| operator | description |
| :- | :- |
| +  | Addition |
| \- | Subtraction |
| *  | Multiplication |
| /  | Division |
| %  | Modulo |

### Syntax

```sql
float operator float
```

_float_
: [value]({{ '/reference/value.html' | relative_url }})

An binary arithmetic operator calculate float values, and return a integer or float value.
If each of operands is not a float value, the value is converted to a float value.

If either of operands is null or conversions to float failed, return null.

## Unary Operators
{: #unary}

| operator | description |
| :- | :- |
| +  | Plus |
| \- | Minus |

### Syntax

```sql
operator float
```

_float_
: [value]({{ '/reference/value.html' | relative_url }})

