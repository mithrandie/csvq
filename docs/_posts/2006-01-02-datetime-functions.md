---
layout: default
title: Datetime Functions - Reference Manual - csvq
category: reference
---

# Datetime Functions

| name | description |
| :- | :- |
| [NOW](#now) | Return current date and time |
| [DATETIME_FORMAT](#datetime_format) | Format a datetime |
| [YEAR](#year) | Return year of a datetime |
| [MONTH](#month) | Return month of a datetime |
| [DAY](#day) | Return day of a datetime |
| [HOUR](#hour) | Return hour of a datetime |
| [MINUTE](#minute) | Return minute of a datetime |
| [SECOND](#second) | Return second of a datetime |
| [MILLISECOND](#millisecond) | Return millisecond of a datetime |
| [MICROSECOND](#microsecond) | Return microsecond of a datetime |
| [NANOSECOND](#nanosecond) | Return nanosecond of a datetime |
| [WEEKDAY](#weekday) | Return weekday number of a datetime |
| [UNIX_TIME](#unix_time) | Return Unix time of a datetime |
| [UNIX_NANO_TIME](#unix_nano_time) | Return Unix nano time of a datetime |
| [DAY_OF_YEAR](#day_of_year) | Return day of year of a datetime |
| [WEEK_OF_YEAR](#week_of_year) | Return week number of year of a datetime |
| [ADD_YEAR](#add_year) | Add years to a datetime |
| [ADD_MONTH](#add_month) | Add monthes to a datetime |
| [ADD_DAY](#add_day) | Add days to a datetime |
| [ADD_HOUR](#add_hour) | Add hours to a datetime |
| [ADD_MINUTE](#add_minute) | Add minutes to a datetime |
| [ADD_SECOND](#add_second) | Add seconds to a datetime |
| [ADD_MILLI](#add_milli) | Add milliseconds to a datetime |
| [ADD_MICRO](#add_micro) | Add microseconds to a datetime |
| [ADD_NANO](#add_nano) | Add nanoseconds to a datetime |
| [DATE_DIFF](#date_diff) | Return the difference of days between two datetime values |
| [TIME_DIFF](#time_diff) | Return the difference of time between two datetime values |

## Definitions

### NOW
{: #now}

```
NOW()
```

_return_
: [datetime]({{ '/reference/value.html#datetime' | relative_url }})

Return current date and time.

### DATETIME_FORMAT
{: #datetime_format}

```
DATETIME_FORMAT(datetime, format)
```

_datetime_
: [datetime]({{ '/reference/value.html#datetime' | relative_url }})

_format_
: [string]({{ '/reference/value.html#string' | relative_url }})

_return_
: [string]({{ '/reference/value.html#string' | relative_url }})

Format the _datetime_ according to the string _format_. 

#### Format Placeholders

| placeholder | description |
| :- | :- |
| %a | Abbreviation of week name (Sun, Mon, ...) |
| %b | Abbreviation of month name (Jan, Feb, ...) |
| %c | Month number (0 - 12) |
| %d | Day of month in two digits (01 - 31) |
| %E | Day of month padding with a underscore (_1 - 31) |
| %e | Day of month (1 - 31) |
| %F | Microseconds that drops trailing zeros (empty - .999999) |
| %f | Microseconds (.000000 - .999999) |
| %H | Hour in 24-hour (00 - 23) |
| %h | Hour in two digits 12-hour (01 - 12) |
| %i | Minute in two digits (00 - 59) |
| %l | Hour in 12-hour (1 - 12) |
| %M | Month name (January, February, ...) |
| %m | Month number with two digits (01 - 12) |
| %N | Nanoseconds that drops trailing zeros (empty - .999999999) |
| %n | Nanoseconds (.000000000 - .999999999) |
| %p | Period in a day (AM or PM) |
| %r | Time with a period (%H:%i:%s %p) |
| %s | Second in two digits (00 - 59) |
| %T | Time (%H:%i:%s) |
| %W | Week name (Sunday, Monday, ...) |
| %Y | Year in four digits |
| %y | Year in two digits |
| %Z | Time zone in time difference |
| %z | Abbreviation of Time zone name |

> You can also use [the Time Layout of the Go Lang](https://golang.org/pkg/time/#Time.Format) as a format.

### YEAR
{: #year}

```
YEAR(datetime)
```

_datetime_
: [datetime]({{ '/reference/value.html#datetime' | relative_url }})

_return_
: [integer]({{ '/reference/value.html#integer' | relative_url }})

Return year of the _datetime_ as integer.

### MONTH
{: #month}

```
MONTH(datetime)
```

_datetime_
: [datetime]({{ '/reference/value.html#datetime' | relative_url }})

_return_
: [integer]({{ '/reference/value.html#integer' | relative_url }})

Return month number of the _datetime_ as integer.

### DAY
{: #day}

```
DAY(datetime)
```

_datetime_
: [datetime]({{ '/reference/value.html#datetime' | relative_url }})

_return_
: [integer]({{ '/reference/value.html#integer' | relative_url }})

Return day of month of the _datetime_ as integer.

### HOUR
{: #hour}

```
HOUR(datetime)
```

_datetime_
: [datetime]({{ '/reference/value.html#datetime' | relative_url }})

_return_
: [integer]({{ '/reference/value.html#integer' | relative_url }})

Return hour of the _datetime_ as integer.

### MINUTE
{: #minute}

```
MINUTE(datetime)
```

_datetime_
: [datetime]({{ '/reference/value.html#datetime' | relative_url }})

_return_
: [integer]({{ '/reference/value.html#integer' | relative_url }})

Return minute of the _datetime_ as integer.

### SECOND
{: #second}

```
SECOND(datetime)
```

_datetime_
: [datetime]({{ '/reference/value.html#datetime' | relative_url }})

_return_
: [integer]({{ '/reference/value.html#integer' | relative_url }})

Return seconds of the _datetime_ as integer.

### MILLISECOND
{: #millisecond}

```
MILLISECOND(datetime)
```

_datetime_
: [datetime]({{ '/reference/value.html#datetime' | relative_url }})

_return_
: [integer]({{ '/reference/value.html#integer' | relative_url }})

Return millisecond of the _datetime_ as integer.

### MICROSECOND
{: #microsecond}

```
MICROSECOND(datetime)
```

_datetime_
: [datetime]({{ '/reference/value.html#datetime' | relative_url }})

_return_
: [integer]({{ '/reference/value.html#integer' | relative_url }})

Return microsecond of the _datetime_ as integer.

### NANOSECOND
{: #nanosecond}

```
NANOSECOND(datetime)
```

_datetime_
: [datetime]({{ '/reference/value.html#datetime' | relative_url }})

_return_
: [integer]({{ '/reference/value.html#integer' | relative_url }})

Return nanosecond of the _datetime_ as integer.

### WEEKDAY
{: #weekday}

```
WEEKDAY(datetime)
```

_datetime_
: [datetime]({{ '/reference/value.html#datetime' | relative_url }})

_return_
: [integer]({{ '/reference/value.html#integer' | relative_url }})

Return weekday number of the _datetime_ as integer.

#### Weekday number

| weekday | number |
| :- | :- |
| Sunday | 0 |
| Monday | 1 |
| Tuesday | 2 |
| Wednesday | 3 |
| Thursday | 4 |
| Friday | 5 |
| Saturday | 6 |

### UNIX_TIME
{: #unix_time}

```
UNIX_TIME(datetime)
```

_datetime_
: [datetime]({{ '/reference/value.html#datetime' | relative_url }})

_return_
: [integer]({{ '/reference/value.html#integer' | relative_url }})

Return the number of seconds elapsed since January 1, 1970 UTC of the _datetime_ as integer.

### UNIX_NANO_TIME
{: #unix_nano_time}

```
UNIX_NANO_TIME(datetime)
```

_datetime_
: [datetime]({{ '/reference/value.html#datetime' | relative_url }})

_return_
: [integer]({{ '/reference/value.html#integer' | relative_url }})

Return the number of nanoseconds elapsed since January 1, 1970 UTC of the _datetime_ as integer.

### DAY_OF_YEAR
{: #day_of_year}

```
DAY_OF_YEAR(datetime)
```

_datetime_
: [datetime]({{ '/reference/value.html#datetime' | relative_url }})

_return_
: [integer]({{ '/reference/value.html#integer' | relative_url }})

Return day of year of the _datetime_ as integer.

### WEEK_OF_YEAR
{: #week_of_year}

```
WEEK_OF_YEAR(datetime)
```

_datetime_
: [datetime]({{ '/reference/value.html#datetime' | relative_url }})

_return_
: [integer]({{ '/reference/value.html#integer' | relative_url }})

Return week number of year of the _datetime_ as integer.
The week number is in the range from 1 to 53.
Jan 01 to Jan 03 of year might return week 52 or 53 of the last year, and Dec 29 to Dec 31 might return week 1 of next year.

### ADD_YEAR
{: #add_year}

```
ADD_YEAR(datetime, duration)
```

_datetime_
: [datetime]({{ '/reference/value.html#datetime' | relative_url }})

_duration_
: [integer]({{ '/reference/value.html#integer' | relative_url }})

_return_
: [datetime]({{ '/reference/value.html#datetime' | relative_url }})

Add _duration_ years to the _datetime_.

### ADD_MONTH
{: #add_month}

```
ADD_MONTH(datetime, duration)
```

_datetime_
: [datetime]({{ '/reference/value.html#datetime' | relative_url }})

_duration_
: [integer]({{ '/reference/value.html#integer' | relative_url }})

_return_
: [datetime]({{ '/reference/value.html#datetime' | relative_url }})

Add _duration_ monthes to the _datetime_.

### ADD_DAY
{: #add_day}

```
ADD_DAY(datetime, duration)
```

_datetime_
: [datetime]({{ '/reference/value.html#datetime' | relative_url }})

_duration_
: [integer]({{ '/reference/value.html#integer' | relative_url }})

_return_
: [datetime]({{ '/reference/value.html#datetime' | relative_url }})

Add _duration_ days to the _datetime_.

### ADD_HOUR
{: #add_hour}

```
ADD_HOUR(datetime, duration)
```

_datetime_
: [datetime]({{ '/reference/value.html#datetime' | relative_url }})

_duration_
: [integer]({{ '/reference/value.html#integer' | relative_url }})

_return_
: [datetime]({{ '/reference/value.html#datetime' | relative_url }})

Add _duration_ hours to the _datetime_.

### ADD_MINUTE
{: #add_minute}

```
ADD_MINUTE(datetime, duration)
```

_datetime_
: [datetime]({{ '/reference/value.html#datetime' | relative_url }})

_duration_
: [integer]({{ '/reference/value.html#integer' | relative_url }})

_return_
: [datetime]({{ '/reference/value.html#datetime' | relative_url }})

Add _duration_ minutes to the _datetime_.

### ADD_SECOND
{: #add_second}

```
ADD_SECOND(datetime, duration)
```

_datetime_
: [datetime]({{ '/reference/value.html#datetime' | relative_url }})

_duration_
: [integer]({{ '/reference/value.html#integer' | relative_url }})

_return_
: [datetime]({{ '/reference/value.html#datetime' | relative_url }})

Add _duration_ seconds to the _datetime_.

### ADD_MILLI
{: #add_milli}

```
ADD_MILLI(datetime, duration)
```

_datetime_
: [datetime]({{ '/reference/value.html#datetime' | relative_url }})

_duration_
: [integer]({{ '/reference/value.html#integer' | relative_url }})

_return_
: [datetime]({{ '/reference/value.html#datetime' | relative_url }})

Add _duration_ milliseconds to the _datetime_.

### ADD_MICRO
{: #add_micro}

```
ADD_MICRO(datetime, duration)
```

_datetime_
: [datetime]({{ '/reference/value.html#datetime' | relative_url }})

_duration_
: [integer]({{ '/reference/value.html#integer' | relative_url }})

_return_
: [datetime]({{ '/reference/value.html#datetime' | relative_url }})

Add _duration_ microseconds to the _datetime_.

### ADD_NANO
{: #add_nano}

```
ADD_NANO(datetime, duration)
```

_datetime_
: [datetime]({{ '/reference/value.html#datetime' | relative_url }})

_duration_
: [integer]({{ '/reference/value.html#integer' | relative_url }})

_return_
: [datetime]({{ '/reference/value.html#datetime' | relative_url }})

Add _duration_ nanoseconds to the _datetime_.

### DATE_DIFF
{: #date_diff}

```
DATE_DIFF(datetime1, datetime2)
```

_datetime1_
: [datetime]({{ '/reference/value.html#datetime' | relative_url }})

_duration2_
: [datetime]({{ '/reference/value.html#datetime' | relative_url }})

_return_
: [integer]({{ '/reference/value.html#integer' | relative_url }})

Return the difference of days between two _datetime_ values.
The time parts are ignored in the calculation. 

### TIME_DIFF
{: #time_diff}

```
TIME_DIFF(datetime1, datetime2)
```

_datetime1_
: [datetime]({{ '/reference/value.html#datetime' | relative_url }})

_duration2_
: [datetime]({{ '/reference/value.html#datetime' | relative_url }})

_return_
: [float]({{ '/reference/value.html#float' | relative_url }})

Return the difference of time between two _datetime_ values.
In the return value, the integer part representing seconds and the fractional part representing nanoseconds.