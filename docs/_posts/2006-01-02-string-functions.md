---
layout: default
title: String Functions - Reference Manual - csvq
category: reference
---

# String Functions

| name | description |
| :- | :- |
| [TRIM](#trim) | Returns the string with all leading and trailing characters removed |
| [LTRIM](#ltrim) | Returns the string with all leading characters removed |
| [RTRIM](#rtrim) | Returns the string with all trailing characters removed |
| [UPPER](#upper) | Returns the string with all characters mapped to their upper case |
| [LOWER](#lower) | Returns the string with all characters mapped to their lower case |
| [BASE64_ENCODE](#base64_encode) | Return the base64 encoding of string |
| [BASE64_DECODE](#base64_decode) | Return the string represented by the base64 string |
| [HEX_ENCODE](#hex_encode) | Return the hexadecimal encoding of string |
| [HEX_DECODE](#hex_decode) | Return the string represented by the hexadecimal string |
| [LEN](#len) | Return the character length of the string |
| [BYTE_LEN](#byte_len) | Return the byte length in utf-8 encoding of the string |
| [LPAD](#lpad) | Return the string left-side padded |
| [RPAD](#rpad) | Return the string right-side padded |
| [SUBSTR](#substr) | Return a substring of the string |
| [REPLACE](#replace) | Return the string with substrings replaced another strings |

## Definitions

### TRIM
{: #trim}

```
TRIM(str)
```

_str_
: [string]({{ '/reference/value.html#string' | relative_url }})

_return_
: [string]({{ '/reference/value.html#string' | relative_url }})

Returns the string _str_ with all leading and trailing white space removed

```
TRIM(str, charset)
```

_str_
: [string]({{ '/reference/value.html#string' | relative_url }})

_charset_
: [string]({{ '/reference/value.html#string' | relative_url }})

_return_
: [string]({{ '/reference/value.html#string' | relative_url }})

Returns the string _str_ with all leading and trailing character contained in _charset_ removed

### LTRIM
{: #ltrim}

```
LTRIM(str)
```

_str_
: [string]({{ '/reference/value.html#string' | relative_url }})

_return_
: [string]({{ '/reference/value.html#string' | relative_url }})

Returns the string _str_ with all leading white space removed

```
LTRIM(str, charset)
```

_str_
: [string]({{ '/reference/value.html#string' | relative_url }})

_charset_
: [string]({{ '/reference/value.html#string' | relative_url }})

_return_
: [string]({{ '/reference/value.html#string' | relative_url }})

Returns the string _str_ with all leading character contained in _charset_ removed


### RTRIM
{: #rtrim}

{: #trim}

```
RTRIM(str)
```

_str_
: [string]({{ '/reference/value.html#string' | relative_url }})

_return_
: [string]({{ '/reference/value.html#string' | relative_url }})

Returns the string _str_ with all trailing white space removed

```
RTRIM(str, charset)
```

_str_
: [string]({{ '/reference/value.html#string' | relative_url }})

_charset_
: [string]({{ '/reference/value.html#string' | relative_url }})

_return_
: [string]({{ '/reference/value.html#string' | relative_url }})

Returns the string _str_ with all trailing character contained in _charset_ removed

### UPPER
{: #upper}

```
UPPER(str)
```

_str_
: [string]({{ '/reference/value.html#string' | relative_url }})

_return_
: [string]({{ '/reference/value.html#string' | relative_url }})

Returns the string _str_ with all characters mapped to their upper case. 

### LOWER
{: #lower}

```
LOWER(str)
```

_str_
: [string]({{ '/reference/value.html#string' | relative_url }})

_return_
: [string]({{ '/reference/value.html#string' | relative_url }})

Returns the string _str_ with all characters mapped to their lower case. 

### BASE64_ENCODE
{: #base64_encode}

```
BASE64_ENCODE(str)
```

_str_
: [string]({{ '/reference/value.html#string' | relative_url }})

_return_
: [string]({{ '/reference/value.html#string' | relative_url }})

Return the base64 encoding of string _str_.

### BASE64_DECODE
{: #base64_decode}

```
BASE64_DECODE(str)
```

_str_
: [string]({{ '/reference/value.html#string' | relative_url }})

_return_
: [string]({{ '/reference/value.html#string' | relative_url }})

Return the string represented by the base64 string _str_.

### HEX_ENCODE
{: #hex_encode}

```
HEX_ENCODE(str)
```

_str_
: [string]({{ '/reference/value.html#string' | relative_url }})

_return_
: [string]({{ '/reference/value.html#string' | relative_url }})

Return the hexadecimal encoding of string _str_.

### HEX_DECODE
{: #hex_decode}

```
HEX_DECODE(str)
```

_str_
: [string]({{ '/reference/value.html#string' | relative_url }})

_return_
: [string]({{ '/reference/value.html#string' | relative_url }})

Return the string represented by the hexadecimal string _str_.

### LEN
{: #len}

```
LEN(str)
```

_str_
: [string]({{ '/reference/value.html#string' | relative_url }})

_return_
: [integer]({{ '/reference/value.html#integer' | relative_url }})

Return the character length of the string _str_.

### BYTE_LEN
{: #byte_len}

```
BYTE_LEN(str)
```

_str_
: [string]({{ '/reference/value.html#string' | relative_url }})

_return_
: [integer]({{ '/reference/value.html#integer' | relative_url }})

Return the byte length in utf-8 encoding of the string _str_.

### LPAD
{: #lpad}

```
LPAD(str, len, padstr)
```

_str_
: [string]({{ '/reference/value.html#string' | relative_url }})

_len_
: [integer]({{ '/reference/value.html#integer' | relative_url }})

_padstr_
: [string]({{ '/reference/value.html#string' | relative_url }})

_return_
: [string]({{ '/reference/value.html#string' | relative_url }})

Return the string _str_ padded with leading _padstr_ to a length specified by _len_.

### RPAD
{: #rpad}

```
RPAD(str, len, padstr)
```

_str_
: [string]({{ '/reference/value.html#string' | relative_url }})

_len_
: [integer]({{ '/reference/value.html#integer' | relative_url }})

_padstr_
: [string]({{ '/reference/value.html#string' | relative_url }})

_return_
: [string]({{ '/reference/value.html#string' | relative_url }})

Return the string _str_ padded with trailing _padstr_ to a length specified by _len_.

### SUBSTR
{: #substr}

```
SUBSTR(str, position)
```

_str_
: [string]({{ '/reference/value.html#string' | relative_url }})

_position_
: [integer]({{ '/reference/value.html#integer' | relative_url }})

_return_
: [string]({{ '/reference/value.html#string' | relative_url }})

Return a substring of the string _str_ from at _position_ to the end.
If _position_ is negative, starting position is _position_ from the end of the str.

```
SUBSTR(str, position, len)
```

_str_
: [string]({{ '/reference/value.html#string' | relative_url }})

_position_
: [integer]({{ '/reference/value.html#integer' | relative_url }})

_len_
: [integer]({{ '/reference/value.html#integer' | relative_url }})

_return_
: [string]({{ '/reference/value.html#string' | relative_url }})

Return a _len_ characters substring of the string _str_ from at _position_.
if _len_ is less than the length from _position_ to the end, return a substring from _position_ to the end. 

### REPLACE
{: #replace}

```
REPLACE(str, old, new)
```

_str_
: [string]({{ '/reference/value.html#string' | relative_url }})

_old_
: [string]({{ '/reference/value.html#string' | relative_url }})

_new_
: [string]({{ '/reference/value.html#string' | relative_url }})

_return_
: [string]({{ '/reference/value.html#string' | relative_url }})

Return the string _str_ with all occurrences of the string _old_ replaced by the string _new_.
