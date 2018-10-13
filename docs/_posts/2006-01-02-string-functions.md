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
| [BASE64_ENCODE](#base64_encode) | Returns the base64 encoding of string |
| [BASE64_DECODE](#base64_decode) | Returns the string represented by the base64 encoding |
| [HEX_ENCODE](#hex_encode) | Returns the hexadecimal encoding of string |
| [HEX_DECODE](#hex_decode) | Returns the string represented by the hexadecimal encoding |
| [LEN](#len) | Returns the character length of the string |
| [BYTE_LEN](#byte_len) | Returns the byte length in utf-8 encoding of the string |
| [LPAD](#lpad) | Returns the string left-side padded |
| [RPAD](#rpad) | Returns the string right-side padded |
| [SUBSTR](#substr) | Returns a substring of the string |
| [INSTR](#instr) | Returns the index of the first occurrence of the substring |
| [LIST_ELEM](#list_elem) | Returns the element of the list |
| [REPLACE](#replace) | Returns the string with substrings replaced another strings |
| [FORMAT](#format) | Returns the formatted string |
| [JSON_VALUE](#json_value) | Returns a value from json |
| [JSON_OBJECT](#json_object) | Returns a string formatted in json object |

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

Returns the string _str_ with all leading and trailing white space removed.

```
TRIM(str, charset)
```

_str_
: [string]({{ '/reference/value.html#string' | relative_url }})

_charset_
: [string]({{ '/reference/value.html#string' | relative_url }})

_return_
: [string]({{ '/reference/value.html#string' | relative_url }})

Returns the string _str_ with all leading and trailing character contained in _charset_ removed.

### LTRIM
{: #ltrim}

```
LTRIM(str)
```

_str_
: [string]({{ '/reference/value.html#string' | relative_url }})

_return_
: [string]({{ '/reference/value.html#string' | relative_url }})

Returns the string _str_ with all leading white space removed.

```
LTRIM(str, charset)
```

_str_
: [string]({{ '/reference/value.html#string' | relative_url }})

_charset_
: [string]({{ '/reference/value.html#string' | relative_url }})

_return_
: [string]({{ '/reference/value.html#string' | relative_url }})

Returns the string _str_ with all leading character contained in _charset_ removed.


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

Returns the string _str_ with all trailing white space removed.

```
RTRIM(str, charset)
```

_str_
: [string]({{ '/reference/value.html#string' | relative_url }})

_charset_
: [string]({{ '/reference/value.html#string' | relative_url }})

_return_
: [string]({{ '/reference/value.html#string' | relative_url }})

Returns the string _str_ with all trailing character contained in _charset_ removed.

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

### INSTR
{: #instr}

```
INSTR(str, substr)
```

_str_
: [string]({{ '/reference/value.html#string' | relative_url }})

_substr_
: [string]({{ '/reference/value.html#string' | relative_url }})

_return_
: [string]({{ '/reference/value.html#string' | relative_url }})

Return the index of the first occurrence of _substr_ in _str_,
or null if _substr_ is not present in _str_, returns null.

### LIST_ELEM
{: #list_elem}

```
LIST_ELEM(str, sep, index)
```

_str_
: [string]({{ '/reference/value.html#string' | relative_url }})

_sep_
: [string]({{ '/reference/value.html#string' | relative_url }})

_index_
: [integer]({{ '/reference/value.html#integer' | relative_url }})

_return_
: [string]({{ '/reference/value.html#string' | relative_url }})

Return the string at _index_ in the list generated by splitting with _sep_ from _str_.

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

### FORMAT
{: #format}

```
FORMAT(format [, replace ... ])
```

_format_
: [string]({{ '/reference/value.html#string' | relative_url }})

_replace_
: [value]({{ '/reference/value.html' | relative_url }})

Return the formatted string.

#### Format Placeholder

```
%[flags][width][.precision]specifier
```

flags
: | flag | description |
  | :- | :- |
  | + | print a plus sign for numeric values |
  | '&nbsp;' (U+0020 Space) | print a space instead of a plus sign |
  | - | pad on the right |
  | 0 | pad with zeros |

width
: [integer]({{ '/reference/value.html#integer' | relative_url }})

precision
: [integer]({{ '/reference/value.html#integer' | relative_url }})

  Number of digits after the decimal point.

specifier
: | specifier | description |
  | :- | :- |
  | b | base 2 integer |
  | o | base 8 integer |
  | d | base 10 integer |
  | x | base 16 integer with lower cases |
  | X | base 16 integer with upper cases |
  | e | exponential notation with lower cases |
  | E | exponential notation with upper cases |
  | f | floating point decimal number |
  | s | string representing the value |
  | q | string representing the value with quotes |
  | % | '%' |

### JSON_VALUE
{: #json_value}

```
JSON_VALUE(query, json)
```

_query_
: [string]({{ '/reference/value.html#string' | relative_url }})

  [JSON Query]({{ '/reference/json.html#query' |relative_url }}) to uniquely specify data.

_json_
: [string]({{ '/reference/value.html#string' | relative_url }})

_return_
: [value]({{ '/reference/value.html' | relative_url }})

Return a value from json.

A JSON values are converted to following types.

| JSON value | csvq value |
|:-|:-|
| object | string |
| array  | string |
| number | integer or float |
| string | string |
| true   | boolean |
| false  | boolean |
| null   | null |


### JSON_OBJECT
{: #json_object}

```
JSON_OBJECT([field [, field ...]])

field
  : value
  : value AS alias
```

_value_
: [value]({{ '/reference/value.html' | relative_url }})

_alias_
: [identifier]({{ '/reference/statement.html#parsing' | relative_url }})

_return_
: [string]({{ '/reference/value.html#string' | relative_url }})


Returns a string formatted in json object.

If no arguments are passed, then the object include all fields in the view.
