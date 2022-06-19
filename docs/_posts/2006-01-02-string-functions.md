---
layout: default
title: String Functions - Reference Manual - csvq
category: reference
---

# String Functions

| name | description |
| :- | :- |
| [TRIM](#trim) | Return a string with all the specified leading and trailing characters removed |
| [LTRIM](#ltrim) | Return a string with all the specified leading characters removed |
| [RTRIM](#rtrim) | Return a string with all the specified trailing characters removed |
| [UPPER](#upper) | Return a string with all characters mapped to their upper case |
| [LOWER](#lower) | Return a string with all characters mapped to their lower case |
| [BASE64_ENCODE](#base64_encode) | Return a base64 encoding of a string |
| [BASE64_DECODE](#base64_decode) | Return a string represented by a base64 encoding |
| [HEX_ENCODE](#hex_encode) | Return a hexadecimal encoding of a string |
| [HEX_DECODE](#hex_decode) | Return a string represented by a hexadecimal encoding |
| [LEN](#len) | Return the number of characters of a string |
| [BYTE_LEN](#byte_len) | Return the byte length of a string |
| [WIDTH](#width) | Return the string width of a string |
| [LPAD](#lpad) | Return a string left-side padded |
| [RPAD](#rpad) | Return a string right-side padded |
| [SUBSTRING](#substring) | Return the substring of a string |
| [SUBSTR](#substr) | Return the substring of a string using zero-based indexing |
| [INSTR](#instr) | Return the index of the first occurrence of a substring |
| [LIST_ELEM](#list_elem) | Return a element of a list |
| [REPLACE](#replace) | Return a string replaced the substrings with another string |
| [REGEXP_MATCH](#regexp_match) | Verify a string matches with a regular expression |
| [REGEXP_FIND](#regexp_find) | Return a string that matches a regular expression |
| [REGEXP_FIND_SUBMATCHES](#regexp_find_submatches) | Return a string representing an array that matches a regular expression |
| [REGEXP_FIND_ALL](#regexp_all) | Return a string representing a nested array that matches a regular expression |
| [REGEXP_REPLACE](#regexp_replace) | Return a string replaced substrings that match a regular expression with another strings |
| [TITLE_CASE](#title_case) | Returns a string converted to Title Case |
| [FORMAT](#format) | Return a formatted string |
| [JSON_VALUE](#json_value) | Return a value from json |
| [JSON_OBJECT](#json_object) | Return a string formatted in json object |

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

Returns the string value that is removed all leading and trailing white spaces from _str_.

```
TRIM(str, charset)
```

_str_
: [string]({{ '/reference/value.html#string' | relative_url }})

_charset_
: [string]({{ '/reference/value.html#string' | relative_url }})

_return_
: [string]({{ '/reference/value.html#string' | relative_url }})

Returns the string value that is removed all leading and trailing characters contained in _charset_ from _str_.

### LTRIM
{: #ltrim}

```
LTRIM(str)
```

_str_
: [string]({{ '/reference/value.html#string' | relative_url }})

_return_
: [string]({{ '/reference/value.html#string' | relative_url }})

Returns the string value that is removed all leading white spaces from _str_.

```
LTRIM(str, charset)
```

_str_
: [string]({{ '/reference/value.html#string' | relative_url }})

_charset_
: [string]({{ '/reference/value.html#string' | relative_url }})

_return_
: [string]({{ '/reference/value.html#string' | relative_url }})

Returns the string value that is removed all leading characters contained in _charset_ from _str_.


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

Returns the string value that is removed all trailing white spaces from _str_.

```
RTRIM(str, charset)
```

_str_
: [string]({{ '/reference/value.html#string' | relative_url }})

_charset_
: [string]({{ '/reference/value.html#string' | relative_url }})

_return_
: [string]({{ '/reference/value.html#string' | relative_url }})

Returns the string value that is removed all trailing characters contained in _charset_ from _str_.

### UPPER
{: #upper}

```
UPPER(str)
```

_str_
: [string]({{ '/reference/value.html#string' | relative_url }})

_return_
: [string]({{ '/reference/value.html#string' | relative_url }})

Returns the string value replaced _str_ with characters mapped to their upper case.

### LOWER
{: #lower}

```
LOWER(str)
```

_str_
: [string]({{ '/reference/value.html#string' | relative_url }})

_return_
: [string]({{ '/reference/value.html#string' | relative_url }})

Returns the string value replaced _str_ with characters mapped to their upper case.

### BASE64_ENCODE
{: #base64_encode}

```
BASE64_ENCODE(str)
```

_str_
: [string]({{ '/reference/value.html#string' | relative_url }})

_return_
: [string]({{ '/reference/value.html#string' | relative_url }})

Returns the Base64 encoding of _str_.

### BASE64_DECODE
{: #base64_decode}

```
BASE64_DECODE(str)
```

_str_
: [string]({{ '/reference/value.html#string' | relative_url }})

_return_
: [string]({{ '/reference/value.html#string' | relative_url }})

Returns the string value represented by _str_ that is encoded with Base64.

### HEX_ENCODE
{: #hex_encode}

```
HEX_ENCODE(str)
```

_str_
: [string]({{ '/reference/value.html#string' | relative_url }})

_return_
: [string]({{ '/reference/value.html#string' | relative_url }})

Returns the hexadecimal encoding of _str_.

### HEX_DECODE
{: #hex_decode}

```
HEX_DECODE(str)
```

_str_
: [string]({{ '/reference/value.html#string' | relative_url }})

_return_
: [string]({{ '/reference/value.html#string' | relative_url }})

Returns the string value represented by _str_ that is encoded with hexadecimal.

### LEN
{: #len}

```
LEN(str)
```

_str_
: [string]({{ '/reference/value.html#string' | relative_url }})

_return_
: [integer]({{ '/reference/value.html#integer' | relative_url }})

Returns the number of characters of _str_.

### BYTE_LEN
{: #byte_len}

```
BYTE_LEN(str [, encoding])
```

_str_
: [string]({{ '/reference/value.html#string' | relative_url }})

_encoding_
: [string]({{ '/reference/value.html#string' | relative_url }})

  "UTF8", "UTF16" or "SJIS". The default is "UTF8".

_return_
: [integer]({{ '/reference/value.html#integer' | relative_url }})

Returns the byte length of _str_.

### WIDTH
{: #width}

```
WIDTH(str)
```

_str_
: [string]({{ '/reference/value.html#string' | relative_url }})

_return_
: [integer]({{ '/reference/value.html#integer' | relative_url }})

Returns the string width of _str_. Half-width characters are counted as 1, and full-width characters are counted as 2.

### LPAD
{: #lpad}

```
LPAD(str, len, padstr [, pad_type, encoding])
```

_str_
: [string]({{ '/reference/value.html#string' | relative_url }})

_len_
: [integer]({{ '/reference/value.html#integer' | relative_url }})

_padstr_
: [string]({{ '/reference/value.html#string' | relative_url }})

_pad_type_
: [string]({{ '/reference/value.html#string' | relative_url }})

  "LEN", "BYTE" or "WIDTH". The default is "LEN".

_encoding_
: [string]({{ '/reference/value.html#string' | relative_url }})

  "UTF8", "UTF16" or "SJIS". The default is "UTF8".

_return_
: [string]({{ '/reference/value.html#string' | relative_url }})

Returns the string value of _str_ padded with leading _padstr_ to a length specified by _len_.

### RPAD
{: #rpad}

```
RPAD(str, len, padstr [, pad_type, encoding])
```

_str_
: [string]({{ '/reference/value.html#string' | relative_url }})

_len_
: [integer]({{ '/reference/value.html#integer' | relative_url }})

_padstr_
: [string]({{ '/reference/value.html#string' | relative_url }})

_pad_type_
: [string]({{ '/reference/value.html#string' | relative_url }})

  "LEN", "BYTE" or "WIDTH". The default is "LEN".

_encoding_
: [string]({{ '/reference/value.html#string' | relative_url }})

  "UTF8", "UTF16" or "SJIS". The default is "UTF8".

_return_
: [string]({{ '/reference/value.html#string' | relative_url }})

Returns the string value of _str_ padded with trailing _padstr_ to a length specified by _len_.


### SUBSTRING
{: #substring}

```
SUBSTRING(str FROM position [FOR len])
SUBSTRING(str, position [, len])
```

_str_
: [string]({{ '/reference/value.html#string' | relative_url }})

_position_
: [integer]({{ '/reference/value.html#integer' | relative_url }})

_len_
: [integer]({{ '/reference/value.html#integer' | relative_url }})

_return_
: [string]({{ '/reference/value.html#string' | relative_url }})

Returns the _len_ characters in _str_ starting from the _position_-th character using one-based positional indexing.

If _position_ is 0, then it is treated as 1.<br />
if _len_ is not specified or _len_ is longer than the length from _position_ to the end, then returns the substring from _position_ to the end.<br /> 
If _position_ is negative, then starting position is _position_ from the end of the _str_.


### SUBSTR
{: #substr}

```
SUBSTR(str, position [, len])
```

_str_
: [string]({{ '/reference/value.html#string' | relative_url }})

_position_
: [integer]({{ '/reference/value.html#integer' | relative_url }})

_len_
: [integer]({{ '/reference/value.html#integer' | relative_url }})

_return_
: [string]({{ '/reference/value.html#string' | relative_url }})

This function behaves the same as [SUBSTRING](#substring), but uses zero-based positional indexing.


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
: [integer]({{ '/reference/value.html#integer' | relative_url }})

Returns the index of the first occurrence of _substr_ in _str_, 
or null if _substr_ is not present in _str_.

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

Returns the string at _index_ in the list generated by splitting with _sep_ from _str_.

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

Returns the string that is replaced all occurrences of _old_ with _new_ in _str_.

### REGEXP_MATCH
{: #regexp_match}

```
REGEXP_MATCH(str, regexp [, flags])
```

_str_
: [string]({{ '/reference/value.html#string' | relative_url }})

_regexp_
: [string]({{ '/reference/value.html#string' | relative_url }})

_flags_
: [string]({{ '/reference/value.html#string' | relative_url }})

  A string including the [flags of regular expressions](#flags-of-regular-expressions)

_return_
: [ternary]({{ '/reference/value.html#ternary' | relative_url }})

Verifies the string _str_ matches with the regular expression _regexp_.

### REGEXP_FIND
{: #regexp_find}

```
REGEXP_FIND(str, regexp [, flags])
```

_str_
: [string]({{ '/reference/value.html#string' | relative_url }})

_regexp_
: [string]({{ '/reference/value.html#string' | relative_url }})

_flags_
: [string]({{ '/reference/value.html#string' | relative_url }})

  A string including the [flags of regular expressions](#flags-of-regular-expressions)

_return_
: [string]({{ '/reference/value.html#string' | relative_url }})

Returns the string that matches the regular expression _regexp_ in _str_.

#### Examples

```shell
# Return the matched string.
$ csvq "SELECT REGEXP_FIND('ABCDEFG abcdefg', 'cdef')"
+----------------------------------------+
| REGEXP_FIND('ABCDEFG abcdefg', 'cdef') |
+----------------------------------------+
| cdef                                   |
+----------------------------------------+

# Return the submatched string if there is a submatch expression.
$ csvq "SELECT REGEXP_FIND('ABCDEFG abcdefg', 'c(de)f')"
+------------------------------------------+
| REGEXP_FIND('ABCDEFG abcdefg', 'c(de)f') |
+------------------------------------------+
| de                                       |
+------------------------------------------+

# Return the first matched string if there are multiple matched strings.
$ csvq "SELECT REGEXP_FIND('ABCDEFG abcdefg', 'cdef', 'i')"
+---------------------------------------------+
| REGEXP_FIND('ABCDEFG abcdefg', 'cdef', 'i') |
+---------------------------------------------+
| CDEF                                        |
+---------------------------------------------+
```

### REGEXP_FIND_SUBMATCHES
{: #regexp_find_submatches}

```
REGEXP_FIND_SUBMATCHES(str, regexp [, flags])
```

_str_
: [string]({{ '/reference/value.html#string' | relative_url }})

_regexp_
: [string]({{ '/reference/value.html#string' | relative_url }})

_flags_
: [string]({{ '/reference/value.html#string' | relative_url }})

  A string including the [flags of regular expressions](#flags-of-regular-expressions)

_return_
: [string]({{ '/reference/value.html#string' | relative_url }})

  A string representing a JSON array.

Returns the string representing an array that matches the regular expression _regexp_ in _str_.

#### Examples

```shell
# Return all the first matched strings including submatches.
$ csvq "SELECT REGEXP_FIND_SUBMATCHES('ABCDEFG abcdefg', 'c(de)f', 'i')"
+----------------------------------------------------------+
| REGEXP_FIND_SUBMATCHES('ABCDEFG abcdefg', 'c(de)f', 'i') |
+----------------------------------------------------------+
| ["CDEF","DE"]                                            |
+----------------------------------------------------------+

# Return only the submatched string.
$ csvq "SELECT JSON_VALUE('[1]', REGEXP_FIND_SUBMATCHES('ABCDEFG abcdefg', 'c(de)f', 'i'))"
+-----------------------------------------------------------------------------+
| JSON_VALUE('[1]', REGEXP_FIND_SUBMATCHES('ABCDEFG abcdefg', 'c(de)f', 'i')) |
+-----------------------------------------------------------------------------+
| DE                                                                          |
+-----------------------------------------------------------------------------+
```

### REGEXP_FIND_ALL
{: #regexp_all}

```
REGEXP_FIND_ALL(str, regexp [, flags])
```

_str_
: [string]({{ '/reference/value.html#string' | relative_url }})

_regexp_
: [string]({{ '/reference/value.html#string' | relative_url }})

_flags_
: [string]({{ '/reference/value.html#string' | relative_url }})

  A string including the [flags of regular expressions](#flags-of-regular-expressions)

_return_
: [string]({{ '/reference/value.html#string' | relative_url }})

  A string representing a nested JSON array.

Returns the string representing a nested array that matches the regular expression _regexp_ in _str_.

#### Examples

```shell
# Return all the matched strings.
$ csvq "SELECT REGEXP_FIND_ALL('ABCDEFG abcdefg', 'c(de)f', 'i')"
+---------------------------------------------------+
| REGEXP_FIND_ALL('ABCDEFG abcdefg', 'c(de)f', 'i') |
+---------------------------------------------------+
| [["CDEF","DE"],["cdef","de"]]                     |
+---------------------------------------------------+

# Return only submatched strings as an array.
$ csvq "SELECT JSON_VALUE('[][1]', REGEXP_FIND_ALL('ABCDEFG abcdefg', 'c(de)f', 'i'))"
+------------------------------------------------------------------------+
| JSON_VALUE('[][1]', REGEXP_FIND_ALL('ABCDEFG abcdefg', 'c(de)f', 'i')) |
+------------------------------------------------------------------------+
| ["DE","de"]                                                            |
+------------------------------------------------------------------------+
```

### REGEXP_REPLACE
{: #regexp_replace}

```
REGEXP_REPLACE(str, regexp, replacement_value [, flags])
```

_str_
: [string]({{ '/reference/value.html#string' | relative_url }})

_regexp_
: [string]({{ '/reference/value.html#string' | relative_url }})

_replacement_value_
: [string]({{ '/reference/value.html#string' | relative_url }})

_flags_
: [string]({{ '/reference/value.html#string' | relative_url }})

  A string including the [flags of regular expressions](#flags-of-regular-expressions)

_return_
: [string]({{ '/reference/value.html#string' | relative_url }})


Returns the string replaced substrings that match the regular expression _regexp_ with _replacement_value_ in _str_.

### TITLE_CASE
{: #title_case}

```
TITLE_CASE(str)
```

_str_
: [string]({{ '/reference/value.html#string' | relative_url }})

_return_
: [string]({{ '/reference/value.html#string' | relative_url }})

Returns a string with the first letter of each word in _str_ capitalized.

### FORMAT
{: #format}

```
FORMAT(format [, replace_value ... ])
```

_format_
: [string]({{ '/reference/value.html#string' | relative_url }})

_replace_value_
: [value]({{ '/reference/value.html' | relative_url }})

_return_
: [string]({{ '/reference/value.html#string' | relative_url }})

Returns a formatted string replaced placeholders with _replace_ in _format_.

#### Format Placeholder

```
%[flag][width][.precision]specifier
```

flag
: | flag | description |
  | :- | :- |
  | + | Print a plus sign for numeric values |
  | '&nbsp;' (U+0020 Space) | Print a space instead of a plus sign |
  | - | Pad on the right |
  | 0 | Pad with zeros |

width
: [integer]({{ '/reference/value.html#integer' | relative_url }})

  Width of the replaced string.

precision
: [integer]({{ '/reference/value.html#integer' | relative_url }})

  Number of digits after the decimal point for a float value, 
  or max length for a string value.
  
specifier
: | specifier | description |
  | :- | :- |
  | b | Base 2 integer |
  | o | Base 8 integer |
  | d | Base 10 integer |
  | x | Base 16 integer with lower cases |
  | X | Base 16 integer with upper cases |
  | e | Exponential notation with lower cases |
  | E | Exponential notation with upper cases |
  | f | Floating point decimal number |
  | s | String representation of the value |
  | q | Quoted string representation of the value |
  | i | Quoted identifier representation of the value |
  | T | Type of the value |
  | % | '%' |

  > Quoted string and identifier representations are escaped for [special characters]({{ '/reference/command.html#special_characters' | relative_url }}).

### JSON_VALUE
{: #json_value}

```
JSON_VALUE(json_query, json_data)
```

_json_query_
: [string]({{ '/reference/value.html#string' | relative_url }})

  [JSON Query]({{ '/reference/json.html#query' |relative_url }}) to uniquely specify a value.

_json_data_
: [string]({{ '/reference/value.html#string' | relative_url }})

_return_
: [value]({{ '/reference/value.html' | relative_url }})

Returns a value in _json_data.

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


Returns a string formatted in JSON.

If no arguments are passed, then the object include all fields in the view.

## Flags of regular expressions
{: #flags-of-regular-expressions}

| flag | description |
|:-|:-|
| i | case-insensitive |
| m  | multi-line mode |
| s | let . match \n |
| U | swap meaning of x* and x*?, x+ and x+?, etc. |
