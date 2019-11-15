---
layout: default
title: Cast Functions - Reference Manual - csvq
category: reference
---

# Cast Functions

| name | description |
| :- | :- |
| [STRING](#string) | Convert a value to a string |
| [INTEGER](#integer) | Convert a value to an integer |
| [FLOAT](#float) | Convert a value to a float |
| [DATETIME](#datetime) | Convert a value to a datetime |
| [BOOLEAN](#boolean) | Convert a value to a boolean |
| [TERNARY](#ternary) | Convert a value to a ternary |

## Definitions

### STRING
{: #string}

```
STRING(value)
```

_value_
: [value]({{ '/reference/value.html' | relative_url }})

_return_
: [string]({{ '/reference/value.html#string' | relative_url }})

Convert _value_ to a string.

| value type | description |
| :- | :- |
| Integer  | An integer value is converted to a string representing a decimal integer. |
| Float    | A float value is converted to a string representing a floating-point decimal. |
| Datetime | A datetime value is converted to a string formatted with RFC3339 with Nano Seconds. |
| Boolean  | A boolean value is converted to either 'true' or 'false'. |
| Ternary  | A ternaly value is converted to any one string of 'TRUE', 'FALSE' and 'UNKNOWN'. |
| Null     | A null value is kept as it is. |


### INTEGER
{: #integer}

```
INTEGER(value)
```

_value_
: [value]({{ '/reference/value.html' | relative_url }})

_return_
: [integer]({{ '/reference/value.html#integer' | relative_url }})

Convert _value_ to an integer.

| value type | description |
| :- | :- |
| String   | If a string is a representation of a decimal integer or its exponential notation, then it is converted to an integer. If a string is a representation of a floating-point decimal or its exponential notation, then it is converted and rounded to an integer. Otherwise it is converted to a null. |
| Float    | A float value is rounded to an integer. |
| Datetime | A datetime value is converted to an integer representing its unix time. |
| Boolean  | A boolean value is converted to a null. |
| Ternary  | A ternaly value is converted to a null. |
| Null     | A null value is kept as it is. |

### FLOAT
{: #float}

```
FLOAT(value)
```

_value_
: [value]({{ '/reference/value.html' | relative_url }})

_return_
: [float]({{ '/reference/value.html#float' | relative_url }})

Convert _value_ to a float.

| value type | description |
| :- | :- |
| String   | If a string is a representation of a floating-point decimal or its exponential notation, then it is converted to a float. Otherwise it is converted to a null. |
| Integer  | An integer value is converted to a float. |
| Datetime | A datetime value is converted to a float representing its unix time. |
| Boolean  | A boolean value is converted to a null. |
| Ternary  | A ternary value is converted to a null. |
| Null     | A null value is kept as it is. |

### DATETIME
{: #datetime}

```
DATETIME(value)
```

_value_
: [value]({{ '/reference/value.html' | relative_url }})

_return_
: [datetime]({{ '/reference/value.html#datetime' | relative_url }})

Convert _value_ to a datetime.

| value type | description |
| :- | :- |
| String   | If a string value is a representation of an integer or a float value, then it is converted to a datetime represented by the number as a unix time. If a string value is formatted as a datetime, then it is convered to a datetime. Otherwise it is converted to a null. |
| Integer  | An integer value is converted to a datetime represented by the integer value as a unix time. |
| Float    | A float value is converted to a datetime represented by the float value as a unix time. |
| Boolean  | A boolean value is converted to a null. |
| Ternary  | A ternaly value is converted to a null. |
| Null     | A null value is kept as it is. |

### BOOLEAN
{: #boolean}

```
BOOLEAN(value)
```

_value_
: [value]({{ '/reference/value.html' | relative_url }})

_boolean_
: [boolean]({{ '/reference/value.html#boolean' | relative_url }})

Convert _value_ to a boolean.

| value type | description |
| :- | :- |
| String   | If a string value is any of '1', 't', 'T', 'TRUE', 'true' and 'True', then it is converted to true. If a string value is any of '0', 'f', 'F', 'FALSE' and 'false', then it is converted to false. Otherwise it is converted to a null. |
| Integer  | If an integer value is 1, then it is converted to true. If an integer value is 0, then it is converted to false. Otherwise it is converted to a null. |
| Float    | If a float value is 1, then it is converted to true. If a float value is 0, then it is converted to false. Otherwise it is converted to a null. |
| Datetime | A datetime value is converted to a null. |
| Ternary  | If a ternary value is TRUE, then it is converted to true. If a ternary value is FALSE, then it is converted to false. Otherwise it is converted to a null. |
| Null     | A null value is kept as it is. |

### TERNARY
{: #ternary}

```
TERNARY(value)
```

_value_
: [value]({{ '/reference/value.html' | relative_url }})

_return_
: [ternary]({{ '/reference/value.html#ternary' | relative_url }})

Convert _value_ to a ternary.

| value type | description |
| :- | :- |
| String   | If a string value is any of '1', 't', 'T', 'TRUE', 'true' and 'True', then it is converted to TRUE. If a string value is any of '0', 'f', 'F', 'FALSE' and 'false', then it is converted to FALSE. Otherwise it is converted to UNKNOWN. |
| Integer  | If an integer value is 1, then it is converted to TRUE. If an integer value is 0, then it is converted to FALSE. Otherwise it is converted to UNKNOWN. |
| Float    | If a float value is 1, then it is converted to TRUE. If a float value is 0, then it is converted to FALSE. Otherwise it is converted to UNKNOWN. |
| Datetime | A datetime value is converted to UNKNOWN. |
| Boolean  | If a boolean value is true, then it is converted to TRUE. If a boolean value is false, then it is converted to FALSE. |
| Null     | A null value is converted to UNKNOWN. |
