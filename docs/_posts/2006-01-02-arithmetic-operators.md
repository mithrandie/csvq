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
value operator value
```

_value_
: [value]({{ '/reference/value.html' | relative_url }})

A binary arithmetic operator calculate integer or float values, and return the result.

If either of operands is null or the conversions to integer or float failed, return null.

## Unary Operators
{: #unary}

| operator | description |
| :- | :- |
| +  | Plus |
| \- | Minus |

### Syntax

```sql
operator value
```

_value_
: [value]({{ '/reference/value.html' | relative_url }})

