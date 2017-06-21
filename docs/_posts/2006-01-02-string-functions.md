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
TRIM(str string) return string
```

Returns the string str with all leading and trailing white space removed

```
TRIM(str string, charset string) return string
```

Returns the string str with all leading and trailing character contained in charset removed

### LTRIM
{: #ltrim}

```
LTRIM(str string) return string
```

Returns the string str with all leading white space removed

```
LTRIM(str string, charset string) return string
```

Returns the string str with all leading character contained in charset removed


### RTRIM
{: #rtrim}

{: #trim}

```
RTRIM(str string) return string
```

Returns the string str with all trailing white space removed

```
RTRIM(str string, charset string) return string
```

Returns the string str with all trailing character contained in charset removed

### UPPER
{: #upper}

```
UPPER(str string) return string
```

Returns the string str with all characters mapped to their upper case. 

### LOWER
{: #lower}

```
LOWER(str string) return string
```

Returns the string str with all characters mapped to their lower case. 

### BASE64_ENCODE
{: #base64_encode}

```
BASE64_ENCODE(str string) return string
```

Return the base64 encoding of string str.

### BASE64_DECODE
{: #base64_decode}

```
BASE64_DECODE(str string) return string
```

Return the string represented by the base64 string str.

### HEX_ENCODE
{: #hex_encode}

```
HEX_ENCODE(str string) return string
```

Return the hexadecimal encoding of string str.

### HEX_DECODE
{: #hex_decode}

```
HEX_DECODE(str string) return string
```

Return the string represented by the hexadecimal string str.

### LEN
{: #len}

```
LEN(str string) return integer
```

Return the character length of the string str.

### BYTE_LEN
{: #byte_len}

```
BYTE_LEN(str string) return integer
```

Return the byte length in utf-8 encoding of the string str.

### LPAD
{: #lpad}

```
LPAD(str string, len integer, padstr string) return string
```

Return the string str padded with leading string padstr to a length specified by len.

### RPAD
{: #rpad}

```
RPAD(str string, len integer, padstr string) return string
```

Return the string str padded with trailing string padstr to a length specified by len.

### SUBSTR
{: #substr}

```
SUBSTR(str string, pos integer) return string
```

Return a substring of the string str from at position pos to the end.
If pos is negative, starting position is pos from the end of the str.

```
SUBSTR(str string, pos integer, len integer) return string
```

Return a len characters substring of the string str from at position pos.
if len is less than the length from pos to the end, return a substring from pos to the end. 

### REPLACE
{: #replace}

```
REPLACE(str string, old string, new string) return string
```
Return the string str with all occurrences of the string old replaced by the string new.
