---
layout: default
title: Numeric Functions - Reference Manual - csvq
category: reference
---

# Numeric Functions

| name | description |
| :- | :- |
| [CEIL](#ceil) | Round a number up |
| [FLOOR](#floor) | Round a number down |
| [ROUND](#round) | Round a number |
| [ABS](#abs) | Return the absolute value of a number |
| [ACOS](#acos) | Return the arc cosine of a number |
| [ASIN](#asin) | Return the arc sine of a number |
| [ATAN](#atan) | Return the arc tangent of a number |
| [ATAN2](#atan2) | Return the arc tangent of two numbers |
| [COS](#cos) | Return the cosine of a number |
| [SIN](#sin) | Return the sine of a number |
| [TAN](#tan) | Return the tangent of a number |
| [EXP](#exp) | Return the value of base _e_ raised to the power of a number |
| [EXP2](#exp2) | Return the value of base _2_ raised to the power of a number |
| [EXPM1](#expm1) | Return the value of base _e_ rised to the power of a number minus 1 |
| [LOG](#log) | Return the natural logarithm of a number |
| [LOG10](#log10) | Return the decimal logarithm of a number |
| [LOG2](#log2) | Return the binary logarithm of a number |
| [LOG1P](#log1p) | Return the natural logarithm of 1 plus a number |
| [SQRT](#sqrt) | Return the square root of a number |
| [POW](#pow) | Returns the value of a number raised to the power of another number |
| [BIN_TO_DEC](#bin_to_dec) | Convert a string representing a binary number to an integer |
| [OCT_TO_DEC](#oct_to_dec) | Convert a string representing a octal number to an integer |
| [HEX_TO_DEC](#hex_to_dec) | Convert a string representing a hexadecimal number to an integer |
| [ENOTATION_TO_DEC](#enotation_to_dec) | Convert a string representing a number with exponential notation to an integer or a float |
| [BIN](#bin) | Convert an integer to a string representing the bynary number |
| [OCT](#oct) | Convert an integer to a string representing the octal number |
| [HEX](#hex) | Convert an integer to a string representing the hexadecimal number |
| [ENOTATION](#enotation) | Convert a float to a string representing the number with exponential notation |
| [NUMBER_FORMAT](#number_format) | Convert a number to a string representing the number with separators |
| [RAND](#rand) | Return a pseudo-random number |

> _e_ is the base of natural logarithms

## Definitions

### CEIL
{: #ceil}

```
CEIL(number)
```

_number_
: [float]({{ '/reference/value.html#float' | relative_url }}) or [integer]({{ '/reference/value.html#integer' | relative_url }})

_return_
: [float]({{ '/reference/value.html#float' | relative_url }})

Rounds _number_ up to an integer value.

```
CEIL(number, place)
```

_number_
: [float]({{ '/reference/value.html#float' | relative_url }}) or [integer]({{ '/reference/value.html#integer' | relative_url }})

_place_
: [integer]({{ '/reference/value.html#integer' | relative_url }})

_return_
: [float]({{ '/reference/value.html#float' | relative_url }})

Rounds _number_ up to _place_ decimal place.
If _place_ is a negative number, _place_ represents the place in the integer part. 

### FLOOR
{: #floor}

```
FLOOR(number)
```

_number_
: [float]({{ '/reference/value.html#float' | relative_url }}) or [integer]({{ '/reference/value.html#integer' | relative_url }})

_return_
: [float]({{ '/reference/value.html#float' | relative_url }})

Rounds _number_ down to an integer value.

```
FLOOR(number, place)
```

_number_
: [float]({{ '/reference/value.html#float' | relative_url }}) or [integer]({{ '/reference/value.html#integer' | relative_url }})

_place_
: [integer]({{ '/reference/value.html#integer' | relative_url }})

_return_
: [float]({{ '/reference/value.html#float' | relative_url }})

Rounds _number_ down to _place_ decimal place.
If _place_ is a negative number, _place_ represents the place in the integer part. 

### ROUND
{: #round}

```
ROUND(number)
```

_number_
: [float]({{ '/reference/value.html#float' | relative_url }}) or [integer]({{ '/reference/value.html#integer' | relative_url }})

_return_
: [float]({{ '/reference/value.html#float' | relative_url }})

Rounds _number_ to an integer value.

```
ROUND(number, place)
```

_number_
: [float]({{ '/reference/value.html#float' | relative_url }}) or [integer]({{ '/reference/value.html#integer' | relative_url }})

_place_
: [integer]({{ '/reference/value.html#integer' | relative_url }})

_return_
: [float]({{ '/reference/value.html#float' | relative_url }})

Rounds _number_ to _place_ decimal place.
If _place_ is a negative number, _place_ represents the place in the integer part. 

### ABS
{: #abs}

```
ABS(number)
```

_number_
: [float]({{ '/reference/value.html#float' | relative_url }}) or [integer]({{ '/reference/value.html#integer' | relative_url }})

_return_
: [float]({{ '/reference/value.html#float' | relative_url }})

Returns the absolute value of _number_

### ACOS
{: #acos}

```
ACOS(number)
```

_number_
: [float]({{ '/reference/value.html#float' | relative_url }}) or [integer]({{ '/reference/value.html#integer' | relative_url }})

_return_
: [float]({{ '/reference/value.html#float' | relative_url }})

Returns the arc cosine of _number_.

### ASIN
{: #asin}

```
ASIN(number)
```

_number_
: [float]({{ '/reference/value.html#float' | relative_url }}) or [integer]({{ '/reference/value.html#integer' | relative_url }})

_return_
: [float]({{ '/reference/value.html#float' | relative_url }})

Returns the arc sine of _number_.

### ATAN
{: #atan}

```
ATAN(number)
```

_number_
: [float]({{ '/reference/value.html#float' | relative_url }}) or [integer]({{ '/reference/value.html#integer' | relative_url }})

_return_
: [float]({{ '/reference/value.html#float' | relative_url }})

Returns the arc tangent of _number_.

### ATAN2
{: #atan2}

```
ATAN2(number2, number1)
```

_number2_
: [float]({{ '/reference/value.html#float' | relative_url }}) or [integer]({{ '/reference/value.html#integer' | relative_url }})

_number1_
: [float]({{ '/reference/value.html#float' | relative_url }}) or [integer]({{ '/reference/value.html#integer' | relative_url }})

_return_
: [float]({{ '/reference/value.html#float' | relative_url }})

Returns the arc tangent of _number2_ / _number1_, using the signs of the two to determine the quadrant of the returns value.

### COS
{: #cos}

```
COS(number)
```

_number_
: [float]({{ '/reference/value.html#float' | relative_url }}) or [integer]({{ '/reference/value.html#integer' | relative_url }})

_return_
: [float]({{ '/reference/value.html#float' | relative_url }})

Returns the cosine of _number_.

### SIN
{: #sin}

```
SIN(number)
```

_number_
: [float]({{ '/reference/value.html#float' | relative_url }}) or [integer]({{ '/reference/value.html#integer' | relative_url }})

_return_
: [float]({{ '/reference/value.html#float' | relative_url }})

Returns the sine of _number_.

### TAN
{: #tan}

```
TAN(number)
```

_number_
: [float]({{ '/reference/value.html#float' | relative_url }}) or [integer]({{ '/reference/value.html#integer' | relative_url }})

_return_
: [float]({{ '/reference/value.html#float' | relative_url }})

Returns the tangent of _number_.

### EXP
{: #exp}

```
EXP(number)
```

_number_
: [float]({{ '/reference/value.html#float' | relative_url }}) or [integer]({{ '/reference/value.html#integer' | relative_url }})

_return_
: [float]({{ '/reference/value.html#float' | relative_url }})

Returns the value of base _e_ raised to the power of _number_.

### EXP2
{: #exp2}

```
EXP2(number)
```

_number_
: [float]({{ '/reference/value.html#float' | relative_url }}) or [integer]({{ '/reference/value.html#integer' | relative_url }})

_return_
: [float]({{ '/reference/value.html#float' | relative_url }})

Returns the value of base _2_ raised to the power of _number_.

### EXPM1
{: #expm1}

```
EXPM1(number)
```

_number_
: [float]({{ '/reference/value.html#float' | relative_url }}) or [integer]({{ '/reference/value.html#integer' | relative_url }})

_return_
: [float]({{ '/reference/value.html#float' | relative_url }})

Returns the value of base _e_ rised to the power of _number_ minus 1.

### LOG
{: #log}

```
LOG(number)
```

_number_
: [float]({{ '/reference/value.html#float' | relative_url }}) or [integer]({{ '/reference/value.html#integer' | relative_url }})

_return_
: [float]({{ '/reference/value.html#float' | relative_url }})

Returns the natural logarithm of _number_.

### LOG10
{: #log10}

```
LOG10(number)
```

_number_
: [float]({{ '/reference/value.html#float' | relative_url }}) or [integer]({{ '/reference/value.html#integer' | relative_url }})

_return_
: [float]({{ '/reference/value.html#float' | relative_url }})

Returns the decimal logarithm of _number_.

### LOG2
{: #log2}

```
LOG2(number)
```

_number_
: [float]({{ '/reference/value.html#float' | relative_url }}) or [integer]({{ '/reference/value.html#integer' | relative_url }})

_return_
: [float]({{ '/reference/value.html#float' | relative_url }})

Returns the binary logarithm of _number_.

### LOG1P
{: #log1p}

```
LOG1P(number)
```

_number_
: [float]({{ '/reference/value.html#float' | relative_url }}) or [integer]({{ '/reference/value.html#integer' | relative_url }})

_return_
: [float]({{ '/reference/value.html#float' | relative_url }})

Returns the natural logarithm of 1 plus _number_.

### SQRT
{: #sqrt}

```
SQRT(number)
```

_number_
: [float]({{ '/reference/value.html#float' | relative_url }}) or [integer]({{ '/reference/value.html#integer' | relative_url }})

_return_
: [float]({{ '/reference/value.html#float' | relative_url }})

Returns the square root of _number_.
 
### POW
{: #pow}

```
POW(base, exponent)
```

_base_
: [float]({{ '/reference/value.html#float' | relative_url }}) or [integer]({{ '/reference/value.html#integer' | relative_url }})

_exponent_
: [float]({{ '/reference/value.html#float' | relative_url }}) or [integer]({{ '/reference/value.html#integer' | relative_url }})

_return_
: [float]({{ '/reference/value.html#float' | relative_url }})

Returns the value of _base_ raised to the power of _exponent_.

### BIN_TO_DEC
{: #bin_to_dec}

```
BIN_TO_DEC(bin)
```

_bin_
: [string]({{ '/reference/value.html#string' | relative_url }})

_return_
: [integer]({{ '/reference/value.html#integer' | relative_url }})

Converts _bin_ representing a binary number to an integer.

### OCT_TO_DEC
{: #oct_to_dec}

```
OCT_TO_DEC(oct)
```

_oct_
: [string]({{ '/reference/value.html#string' | relative_url }})

_return_
: [integer]({{ '/reference/value.html#integer' | relative_url }})

Converts _hex_ representing a octal number to an integer.

### HEX_TO_DEC
{: #hex_to_dec}

```
HEX_TO_DEC(hex)
```

_hex_
: [string]({{ '/reference/value.html#string' | relative_url }})

_return_
: [integer]({{ '/reference/value.html#integer' | relative_url }})

Converts _hex_ representing a hexadecimal number to an integer.

### ENOTATION_TO_DEC
{: #enotation_to_dec}

```
ENOTATION_TO_DEC(enotation)
```

_enotation_
: [string]({{ '/reference/value.html#string' | relative_url }})

_return_
: [float]({{ '/reference/value.html#float' | relative_url }})

Converts _enotation_ representing a number with exponential notation to a float.

### BIN
{: #bin}

```
BIN(integer)
```

_integer_
: [integer]({{ '/reference/value.html#integer' | relative_url }})

_return_
: [string]({{ '/reference/value.html#string' | relative_url }})

Converts _integer_ to a string representing the binary number.

### OCT
{: #oct}

```
OCT(integer)
```

_integer_
: [integer]({{ '/reference/value.html#integer' | relative_url }})

_return_
: [string]({{ '/reference/value.html#string' | relative_url }})

Converts _integer_ to a string representing the octal number.

### HEX
{: #hex}

```
HEX(integer)
```

_integer_
: [integer]({{ '/reference/value.html#integer' | relative_url }})

_return_
: [string]({{ '/reference/value.html#string' | relative_url }})

Converts _integer_ to a string representing the hexadecimal number.

### ENOTATION
{: #enotation}

```
ENOTATION(float)
```

_float_
: [float]({{ '/reference/value.html#float' | relative_url }})

_return_
: [string]({{ '/reference/value.html#string' | relative_url }})

Converts _float_ to a string representing the number with exponential notation.


### NUMBER_FORMAT
{: #number_format}

```
NUMBER_FORMAT(number [, precision, decimalPoint, thousandsSeparator, decimalSeparator])
```

_number_
: [float]({{ '/reference/value.html#float' | relative_url }}) or [integer]({{ '/reference/value.html#integer' | relative_url }})

_precision_
: [integer]({{ '/reference/value.html#integer' | relative_url }})

  The default is -1.
  -1 is the special precision to determine the number of digits automatically.

_decimalPoint_
: [string]({{ '/reference/value.html#string' | relative_url }})

  The default is ".".

_thousandsSeparator_
: [string]({{ '/reference/value.html#string' | relative_url }})

  The default is ",".

_decimalSeparator_
: [string]({{ '/reference/value.html#string' | relative_url }})

  The default is empty string.

_return_
: [string]({{ '/reference/value.html#string' | relative_url }})

Converts _number_ to a string representing the number with separators.


### RAND
{: #rand}

```
RAND()
```

_return_
: [float]({{ '/reference/value.html#float' | relative_url }})

Returns a random float greater than or equal to 0.0 and less than 1.0. 

```
RAND(min, max)
```

_min_
: [integer]({{ '/reference/value.html#integer' | relative_url }})

_max_
: [integer]({{ '/reference/value.html#integer' | relative_url }})

_return_
: [integer]({{ '/reference/value.html#integer' | relative_url }})

Returns a random integer between _min_ and _max_.