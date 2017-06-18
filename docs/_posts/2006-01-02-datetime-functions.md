---
layout: default
title: Datetime Functions - Reference Manual - csvq
category: reference
---

# Datetime Functions

| name | description |
| :- | :- |
| [NOW](#now) | |
| [DATETIME_FORMAT](#datetime_format) | |
| [YEAR](#year) | |
| [MONTH](#month) | |
| [DAY](#day) | |
| [HOUR](#hour) | |
| [MINUTE](#minute) | |
| [SECOND](#second) | |
| [MILLISECOND](#millisecond) | |
| [MICROSECOND](#microsecond) | |
| [NANOSECOND](#nanosecond) | |
| [WEEKDAY](#weekday) | |
| [UNIX_TIME](#unix_time) | |
| [UNIX_NANO_TIME](#unix_nano_time) | |
| [DAY_OF_YEAR](#day_of_year) | |
| [WEEK_OF_YEAR](#week_of_year) | |
| [ADD_YEAR](#add_year) | |
| [ADD_MONTH](#add_month) | |
| [ADD_DAY](#add_day) | |
| [ADD_HOUR](#add_hour) | |
| [ADD_MINUTE](#add_minute) | |
| [ADD_SECOND](#add_second) | |
| [ADD_MILLI](#add_milli) | |
| [ADD_MICRO](#add_micro) | |
| [ADD_NANO](#add_nano) | |
| [DATE_DIFF](#date_diff) | |
| [TIME_DIFF](#time_diff) | |

## Definitions

### NOW
{: #now}

```
NOW() return datetime
```

### DATETIME_FORMAT
{: #datetime_format}

```
DATETIME_FORMAT(dt datetime, format string) return string
```

### YEAR
{: #year}

```
YEAR(dt datetime) return integer
```

### MONTH
{: #month}

```
MONTH(dt datetime) return integer
```

### DAY
{: #day}

```
DAY(dt datetime) return integer
```

### HOUR
{: #hour}

```
HOUR(dt datetime) return integer
```

### MINUTE
{: #minute}

```
MINUTE(dt datetime) return integer
```

### SECOND
{: #second}

```
SECOND(dt datetime) return integer
```

### MILLISECOND
{: #millisecond}

```
MILLISECOND(dt datetime) return integer
```

### MICROSECOND
{: #microsecond}

```
MICROSECOND(dt datetime) return integer
```

### NANOSECOND
{: #nanosecond}

```
NANOSECOND(dt datetime) return integer
```

### WEEKDAY
{: #weekday}

```
WEEKDAY(dt datetime) return integer
```

### UNIX_TIME
{: #unix_time}

```
UNIX_TIME(dt datetime) return integer
```

### UNIX_NANO_TIME
{: #unix_nano_time}

```
UNIX_NANO_TIME(dt datetime) return integer
```

### DAY_OF_YEAR
{: #day_of_year}

```
DAY_OF_YEAR(dt datetime) return integer
```

### WEEK_OF_YEAR
{: #week_of_year}

```
WEEK_OF_YEAR(dt datetime) return integer
```

### ADD_YEAR
{: #add_year}

```
ADD_YEAR(dt datetime, duration int) return datetime
```

### ADD_MONTH
{: #add_month}

```
ADD_MONTH(dt datetime, duration int) return datetime
```

### ADD_DAY
{: #add_day}

```
ADD_DAY(dt datetime, duration int) return datetime
```

### ADD_HOUR
{: #add_hour}

```
ADD_HOUR(dt datetime, duration int) return datetime
```

### ADD_MINUTE
{: #add_minute}

```
ADD_MINUTE(dt datetime, duration int) return datetime
```

### ADD_SECOND
{: #add_second}

```
ADD_SECOND(dt datetime, duration int) return datetime
```

### ADD_MILLI
{: #add_milli}

```
ADD_MILLI(dt datetime, duration int) return datetime
```

### ADD_MICRO
{: #add_micro}

```
ADD_MICRO(dt datetime, duration int) return datetime
```

### ADD_NANO
{: #add_nano}

```
ADD_NANO(dt datetime, duration int) return datetime
```

### DATE_DIFF
{: #date_diff}

```
DATE_DIFF(dt1 datetime, dt2 datetime) return integer
```

### TIME_DIFF
{: #time_diff}

```
TIME_DIFF(dt1 datetime, dt2 datetime) return float
```
