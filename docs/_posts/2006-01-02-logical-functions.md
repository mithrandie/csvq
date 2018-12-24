---
layout: default
title: Logical Functions - Reference Manual - csvq
category: reference
---

# Logical Functions

| name | description |
| :- | :- |
| [COALESCE](#coalesce) | Return the first non-null value in arguments |
| [IF](#if) | Return a value by condition |
| [IFNULL](#ifnull) | Return a value whether passed value is null |
| [NULLIF](#nullif) | Return null whether passed values are equal |

## Definitions

### COALESCE
{: #coalesce}

```
COALESCE(value [, value ...])
```

_value_
: [value]({{ '/reference/value.html' | relative_url }})

_return_
: [primitive type]({{ '/reference/value.html#primitive_types' | relative_url }})

Returns the first non-null _value_ in arguments. If there is no non-null _value_, then returns null.

### IF
{: #if}

```
IF(condition, value1, value2)
```

_condition_
: [value]({{ '/reference/value.html' | relative_url }})

_value1_
: [value]({{ '/reference/value.html' | relative_url }})

_value2_
: [value]({{ '/reference/value.html' | relative_url }})

_return_
: [primitive type]({{ '/reference/value.html#primitive_types' | relative_url }})

If _condition_ is TRUE, then returns _value1_. Otherwise returns _value2_.

### IFNULL
{: #ifnull}

```
IFNULL(value1, value2)
```

_value1_
: [value]({{ '/reference/value.html' | relative_url }})

_value2_
: [value]({{ '/reference/value.html' | relative_url }})

_return_
: [primitive type]({{ '/reference/value.html#primitive_types' | relative_url }})

If _value1_ is null, then returns _value2_. Otherwise returns _value1_.

### NULLIF
{: #nullif}

```
NULLIF(value1, value2)
```

_value1_
: [value]({{ '/reference/value.html' | relative_url }})

_value2_
: [value]({{ '/reference/value.html' | relative_url }})

_return_
: [primitive type]({{ '/reference/value.html#primitive_types' | relative_url }})

If _value1_ is equal to _value2_, returns null. Otherwise returns _value1_.

