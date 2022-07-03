---
layout: default
title: Cast Functions - Reference Manual - csvq
category: reference
---

# Cast Functions

| name                  | description                   |
|:----------------------|:------------------------------|
| [STRING](#string)     | Convert a value to a string   |
| [INTEGER](#integer)   | Convert a value to an integer |
| [FLOAT](#float)       | Convert a value to a float    |
| [DATETIME](#datetime) | Convert a value to a datetime |
| [BOOLEAN](#boolean)   | Convert a value to a boolean  |
| [TERNARY](#ternary)   | Convert a value to a ternary  |

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

| Type     | Value after conversion                          |
|:---------|:------------------------------------------------|
| Integer  | String representing a decimal integer           |
| Float    | String representing a floating-point decimal    |
| Datetime | String formatted with RFC3339 with Nano Seconds |
| Boolean  | 'true' or 'false'                               |
| Ternary  | 'TRUE', 'FALSE' and 'UNKNOWN'                   |
| Null     | Null                                            |


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

| Type     | Value                                                                  | Value after conversion                   |
|:---------|:-----------------------------------------------------------------------|:-----------------------------------------|
| String   | Representation of a decimal integer                                    | Integer represented by the string        |
|          | Representation of a floating-point decimal or its exponential notation | Integer with decimal places rounded down |
|          | Other values                                                           | Null                                     |
| Float    | +Inf, -Inf, NaN                                                        | Null                                     |
|          | Other values                                                           | Integer with decimal places rounded down |
| Datetime |                                                                        | Integer representing its unix time       |
| Boolean  |                                                                        | Null                                     |
| Ternary  |                                                                        | Null                                     |
| Null     |                                                                        | Null                                     |

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

| Type     | Value                                                                  | Value after conversion           |
|:---------|:-----------------------------------------------------------------------|:---------------------------------|
| String   | Representation of a floating-point decimal or its exponential notation | Float represented by the string  |
|          | 'Inf', '+Inf'                                                          | +Inf                             |
|          | '-Inf'                                                                 | -Inf                             |
|          | 'NaN'                                                                  | NaN                              |
|          | Other values                                                           | Null                             |
| Integer  |                                                                        | Float equivalent to the integer  |
| Datetime |                                                                        | Float representing its unix time |
| Boolean  |                                                                        | Null                             |
| Ternary  |                                                                        | Null                             |
| Null     |                                                                        | Null                             |

### DATETIME
{: #datetime}

```
DATETIME(value [, timezone])
```

_value_
: [value]({{ '/reference/value.html' | relative_url }})

_timezone_
: [string]({{ '/reference/value.html#string' | relative_url }})

  _Local_, _UTC_ or a timezone name in the IANA TimeZone database(in the form of _"Area/Location"_. e.g. _"America/Los_Angeles"_).
  The default is the timezone set to the flag [_@@TIMEZONE_]({{ '/reference/flag.html' | relative_url }}).

  See: [--timezone option]({{ '/reference/command.html#options' | relative_url }})

_return_
: [datetime]({{ '/reference/value.html#datetime' | relative_url }})

Convert _value_ to a datetime.

| Type    | Value                                                                  | Value after conversion                                   |
|:--------|:-----------------------------------------------------------------------|:---------------------------------------------------------|
| String  | Datetime Formats                                                       | Datetime represented by the string                       |
|         | Representation of a decimal integer                                    | Datetime represented by the integer value as a unix time |
|         | Representation of a floating-point decimal or its exponential notation | Datetime represented by the float value as a unix time   |
|         | Other values                                                           | Null                                                     |
| Integer |                                                                        | Datetime represented by the integer value as a unix time |
| Float   | +Inf, -Inf, NaN                                                        | Null                                                     |
|         | Other values                                                           | Datetime represented by the float value as a unix time   |
| Boolean |                                                                        | Null                                                     |
| Ternary |                                                                        | Null                                                     |
| Null    |                                                                        | Null                                                     |

#### Format of string to be interpreted as datetime
{: #format-of-string-as-datetime}

Strings of the form passed by the [--datetime-format option]({{ '/reference/command.html#options' | relative_url }}) and defined in the [configuration files]({{ '/reference/command.html#configurations' | relative_url }}),  or the following forms can be converted to datetime values.

| DateFormat | Example    |
|:-----------|:-----------|
| YYYY-MM-DD | 2012-03-15 |
| YYYY/MM/DD | 2012/03/15 |
| YYYY-M-D   | 2012-3-15  |
| YYYY/M/D   | 2012/3/15  |

&nbsp;

| DatetimeFormat                            | Example                                                |
|:------------------------------------------|:-------------------------------------------------------|
| DateFormat                                | 2012-03-15                                             |
| DateFormat hh:mm:ss(.NanoSecods)          | 2012-03-15 12:03:01<br />2012-03-15 12:03:01.123456789 |
| DateFormat hh:mm:ss(.NanoSecods) ±hh:mm   | 2012-03-15 12:03:01 -07:00                             |
| DateFormat hh:mm:ss(.NanoSecods) ±hhmm    | 2012-03-15 12:03:01 -0700                              |
| DateFormat hh:mm:ss(.NanoSecods) TimeZone | 2012-03-15 12:03:01 PST                                |
| YYYY-MM-DDThh:mm:ss(.NanoSeconds)         | 2012-03-15T12:03:01                                    |
| RFC3339                                   | 2012-03-15T12:03:01-07:00                              |
| RFC3339 with Nano Seconds                 | 2012-03-15T12:03:01.123456789-07:00                    |
| RFC822                                    | 03 Mar 12 12:03 PST                                    |
| RFC822 with Numeric Zone                  | 03 Mar 12 12:03 -0700                                  |

> Timezone abbreviations such as "PST" may not work properly depending on your environment,
> so you should use timezone offset such as "-07:00" as possible.

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

| Type     | Value             | Value after conversion |
|:---------|:------------------|:-----------------------|
| String   | '1', 't', 'true'  | true                   |
| String   | '0', 'f', 'false' | false                  |
| String   | Other values      | Null                   |
| Integer  | 1                 | true                   |
| Integer  | 0                 | false                  |
| Integer  | Other values      | Null                   |
| Float    | 1                 | true                   |
| Float    | 0                 | false                  |
| Float    | Other values      | Null                   |
| Datetime |                   | Null                   |
| Ternary  | TRUE              | true                   |
|          | FALSE             | false                  |
|          | UNKNOWN           | Null                   |
| Null     |                   | Null                   |

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

| Type     | Value             | Value after conversion |
|:---------|:------------------|:-----------------------|
| String   | '1', 't', 'true'  | TRUE                   |
| String   | '0', 'f', 'false' | FALSE                  |
| String   | Other values      | UNKNOWN                |
| Integer  | 1                 | TRUE                   |
| Integer  | 0                 | FALSE                  |
| Integer  | Other values      | UNKNOWN                |
| Float    | 1                 | TRUE                   |
| Float    | 0                 | FALSE                  |
| Float    | Other values      | UNKNOWN                |
| Datetime |                   | UNKNOWN                |
| Boolean  | true              | TRUE                   |
|          | false             | FALSE                  |
| Null     |                   | UNKNOWN                |
