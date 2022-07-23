---
layout: default
title: Datetime Functions - Reference Manual - csvq
category: reference
---

# Datetime Functions

| name | description |
| :- | :- |
| [NOW](#now) | Return a datetime value of current date and time |
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
| [TRUNC_MONTH](#trunc_month)   | Truncate time information less than 1 year from a datetime |
| [TRUNC_DAY](#trunc_day)       | Truncate time information less than 1 month from a datetime |
| [TRUNC_TIME](#trunc_time)     | Truncate time information less than 1 day from a datetime |
| [TRUNC_HOUR](#trunc_time)     | Alias for TRUNC_TIME |
| [TRUNC_MINUTE](#trunc_minute) | Truncate time information less than 1 hour from a datetime |
| [TRUNC_SECOND](#trunc_second) | Truncate time information less than 1 minute from a datetime |
| [TRUNC_MILLI](#trunc_milli)   | Truncate time information less than 1 second from a datetime |
| [TRUNC_MICRO](#trunc_micro)   | Truncate time information less than 1 millisecond from a datetime |
| [TRUNC_NANO](#trunc_nano)     | Truncate time information less than 1 microsecond from a datetime |
| [DATE_DIFF](#date_diff) | Return the difference of days between two datetime values |
| [TIME_DIFF](#time_diff) | Return the difference of time between two datetime values as seconds |
| [TIME_NANO_DIFF](#time_nano_diff) | Return the difference of time between two datetime values as nanoseconds |
| [UTC](#utc) | Return a datetime in UTC |
| [MILLI_TO_DATETIME](#milli_to_datetime) | Convert an integer representing Unix milliseconds to a datetime |
| [NANO_TO_DATETIME](#nano_to_datetime) | Convert an integer representing Unix nano time to a datetime |

## Definitions

### NOW
{: #now}

```
NOW()
```

_return_
: [datetime]({{ '/reference/value.html#datetime' | relative_url }})

Returns a datetime value of current date and time.
In a single query, every this function returns the same value. 

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

Formats _datetime_ according to _format_. 

#### Format Placeholders

| placeholder | replacement value |
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
| %% | '%' |

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

Returns the year of _datetime_ as an integer.

### MONTH
{: #month}

```
MONTH(datetime)
```

_datetime_
: [datetime]({{ '/reference/value.html#datetime' | relative_url }})

_return_
: [integer]({{ '/reference/value.html#integer' | relative_url }})

Returns the month number of _datetime_ as an integer.

### DAY
{: #day}

```
DAY(datetime)
```

_datetime_
: [datetime]({{ '/reference/value.html#datetime' | relative_url }})

_return_
: [integer]({{ '/reference/value.html#integer' | relative_url }})

Returns the day of month of _datetime_ as an integer.

### HOUR
{: #hour}

```
HOUR(datetime)
```

_datetime_
: [datetime]({{ '/reference/value.html#datetime' | relative_url }})

_return_
: [integer]({{ '/reference/value.html#integer' | relative_url }})

Returns the hour of _datetime_ as an integer.

### MINUTE
{: #minute}

```
MINUTE(datetime)
```

_datetime_
: [datetime]({{ '/reference/value.html#datetime' | relative_url }})

_return_
: [integer]({{ '/reference/value.html#integer' | relative_url }})

Returns the minute of _datetime_ as an integer.

### SECOND
{: #second}

```
SECOND(datetime)
```

_datetime_
: [datetime]({{ '/reference/value.html#datetime' | relative_url }})

_return_
: [integer]({{ '/reference/value.html#integer' | relative_url }})

Returns the seconds of _datetime_ as an integer.

### MILLISECOND
{: #millisecond}

```
MILLISECOND(datetime)
```

_datetime_
: [datetime]({{ '/reference/value.html#datetime' | relative_url }})

_return_
: [integer]({{ '/reference/value.html#integer' | relative_url }})

Returns the millisecond of _datetime_ as an integer.

### MICROSECOND
{: #microsecond}

```
MICROSECOND(datetime)
```

_datetime_
: [datetime]({{ '/reference/value.html#datetime' | relative_url }})

_return_
: [integer]({{ '/reference/value.html#integer' | relative_url }})

Returns the microsecond of _datetime_ as an integer.

### NANOSECOND
{: #nanosecond}

```
NANOSECOND(datetime)
```

_datetime_
: [datetime]({{ '/reference/value.html#datetime' | relative_url }})

_return_
: [integer]({{ '/reference/value.html#integer' | relative_url }})

Returns the nanosecond of _datetime_ as an integer.

### WEEKDAY
{: #weekday}

```
WEEKDAY(datetime)
```

_datetime_
: [datetime]({{ '/reference/value.html#datetime' | relative_url }})

_return_
: [integer]({{ '/reference/value.html#integer' | relative_url }})

Returns the weekday number of _datetime_ as an integer.

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

Returns the number of seconds elapsed since January 1, 1970 UTC of _datetime_ as an integer.

### UNIX_NANO_TIME
{: #unix_nano_time}

```
UNIX_NANO_TIME(datetime)
```

_datetime_
: [datetime]({{ '/reference/value.html#datetime' | relative_url }})

_return_
: [integer]({{ '/reference/value.html#integer' | relative_url }})

Returns the number of nanoseconds elapsed since January 1, 1970 UTC of _datetime_ as an integer.

### DAY_OF_YEAR
{: #day_of_year}

```
DAY_OF_YEAR(datetime)
```

_datetime_
: [datetime]({{ '/reference/value.html#datetime' | relative_url }})

_return_
: [integer]({{ '/reference/value.html#integer' | relative_url }})

Returns the day of the year of _datetime_ as an integer.

### WEEK_OF_YEAR
{: #week_of_year}

```
WEEK_OF_YEAR(datetime)
```

_datetime_
: [datetime]({{ '/reference/value.html#datetime' | relative_url }})

_return_
: [integer]({{ '/reference/value.html#integer' | relative_url }})

Returns the week number of the year of _datetime_ as an integer.
The week number is in the range from 1 to 53.
Jan 01 to Jan 03 of a year might returns week 52 or 53 of the last year, and Dec 29 to Dec 31 might returns week 1 of the next year.

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

Adds _duration_ years to _datetime_.

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

Adds _duration_ monthes to _datetime_.

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

Adds _duration_ days to _datetime_.

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

Adds _duration_ hours to _datetime_.

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

Adds _duration_ minutes to _datetime_.

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

Adds _duration_ seconds to _datetime_.

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

Adds _duration_ milliseconds to _datetime_.

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

Adds _duration_ microseconds to _datetime_.

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

Adds _duration_ nanoseconds to _datetime_.


### TRUNC_MONTH
{: #trunc_month}

```
TRUNC_MONTH(datetime)
```

_datetime_
: [datetime]({{ '/reference/value.html#datetime' | relative_url }})

_return_
: [datetime]({{ '/reference/value.html#datetime' | relative_url }})

Truncates time information less than 1 year from _datetime_.


### TRUNC_DAY
{: #trunc_day}

```
TRUNC_DAY(datetime)
```

_datetime_
: [datetime]({{ '/reference/value.html#datetime' | relative_url }})

_return_
: [datetime]({{ '/reference/value.html#datetime' | relative_url }})

Truncates time information less than 1 month from _datetime_.


### TRUNC_TIME
{: #trunc_time}

```
TRUNC_TIME(datetime)
```

_datetime_
: [datetime]({{ '/reference/value.html#datetime' | relative_url }})

_return_
: [datetime]({{ '/reference/value.html#datetime' | relative_url }})

Truncates time information less than 1 day from _datetime_.


### TRUNC_MINUTE
{: #trunc_minute}

```
TRUNC_MINUTE(datetime)
```

_datetime_
: [datetime]({{ '/reference/value.html#datetime' | relative_url }})

_return_
: [datetime]({{ '/reference/value.html#datetime' | relative_url }})

Truncates time information less than 1 hour from _datetime_.


### TRUNC_SECOND
{: #trunc_second}

```
TRUNC_SECOND(datetime)
```

_datetime_
: [datetime]({{ '/reference/value.html#datetime' | relative_url }})

_return_
: [datetime]({{ '/reference/value.html#datetime' | relative_url }})

Truncates time information less than 1 minute from _datetime_.


### TRUNC_MILLI
{: #trunc_milli}

```
TRUNC_MILLI(datetime)
```

_datetime_
: [datetime]({{ '/reference/value.html#datetime' | relative_url }})

_return_
: [datetime]({{ '/reference/value.html#datetime' | relative_url }})

Truncates time information less than 1 second from _datetime_.


### TRUNC_MICRO
{: #trunc_micro}

```
TRUNC_MICRO(datetime)
```

_datetime_
: [datetime]({{ '/reference/value.html#datetime' | relative_url }})

_return_
: [datetime]({{ '/reference/value.html#datetime' | relative_url }})

Truncates time information less than 1 millisecond from _datetime_.


### TRUNC_NANO
{: #trunc_nano}

```
TRUNC_NANO(datetime)
```

_datetime_
: [datetime]({{ '/reference/value.html#datetime' | relative_url }})

_return_
: [datetime]({{ '/reference/value.html#datetime' | relative_url }})

Truncates time information less than 1 microsecond from _datetime_.



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

Returns the difference of days between two _datetime_ values.
The time information less than 1 day are ignored in the calculation. 

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
: [float]({{ '/reference/value.html#float' | relative_url }}) or [integer]({{ '/reference/value.html#integer' | relative_url }})

Returns the difference of time between two _datetime_ values as seconds.
In the return value, the integer part represents seconds and the fractional part represents nanoseconds.

### TIME_NANO_DIFF
{: #time_nano_diff}

```
TIME_NANO_DIFF(datetime1, datetime2)
```

_datetime1_
: [datetime]({{ '/reference/value.html#datetime' | relative_url }})

_duration2_
: [datetime]({{ '/reference/value.html#datetime' | relative_url }})

_return_
: [integer]({{ '/reference/value.html#integer' | relative_url }})

Returns the difference of time between two _datetime_ values as nanoseconds.

### UTC
{: #utc}

```
UTC(datetime)
```

_datetime_
: [datetime]({{ '/reference/value.html#datetime' | relative_url }})

_return_
: [datetime]({{ '/reference/value.html#datetime' | relative_url }})

Returns the datetime value of _datetime_ in UTC.

### MILLI_TO_DATETIME
{: #milli_to_datetime}

```
MILLI_TO_DATETIME(unix_milliseconds)
```

_unix_milliseconds_
: [integer]({{ '/reference/value.html#integer' | relative_url }})

_return_
: [datetime]({{ '/reference/value.html#datetime' | relative_url }})

Converts an integer representing Unix milliseconds to a datetime.

### NANO_TO_DATETIME
{: #nano_to_datetime}

```
NANO_TO_DATETIME(unix_nano_time)
```

_unix_nano_time_
: [integer]({{ '/reference/value.html#integer' | relative_url }})

_return_
: [datetime]({{ '/reference/value.html#datetime' | relative_url }})

Converts an integer representing Unix nano time to a datetime.
