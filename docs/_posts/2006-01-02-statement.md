---
layout: default
title: Statements - Reference Manual - csvq
category: reference
---

# Statements

* [Basics](#basics)
* [Parsing](#parsing)
* [Comments](#comments)
* [Reserved Words](#reserved_words)

## Basics
{: #basics}

You can pass a query or statements(it's also called procedure) as a csvq command argument or a source file.
Statements have to be encoded in UTF-8.

A statements is terminated with a semicolon. 
Stetaments are processed sequentially for each statement.
In statements, character cases of keywords are ignored.

If you want to execute a single query, you can omit the terminal semicolon.  

```bash
# Execute a single query
$ csvq "SELECT id, name FROM user"

# Execute statements
$ csvq "VAR @id := 0; SELECT @id := @id + 1 AS id, name FROM user;"

# Load statements from file
$ cat statements.sql
VAR @id := 0;
SELECT @id := @id + 1 AS id,
       name
  FROM user;

$ csvq -s statements.sql
```
## Parsing
{: #parsing}

You can use following types in statements.

Identifier
: A identifier is a word starting with any unicode letter or a low line(U+005F '\_') and followed by a character string that contains any unicode letters, any digits or low lines(U+005F '\_').
  You cannot use [reserved words](#reserved_words) as a identifier.

  Notwithstanding above naming restriction, you can use most character strings as a identifier by enclosing in back quotes.
  Back quotes are escaped by back slashes.
  
  Identifiers represent tables, columns, functions or cursors.
  These character cases are insensitive except file paths, and whether file paths are case insensitive or not depends on your file system.
  
String
: A string is a character string enclosed in single quotes or double quotes.
  In a string, single quotes or double quotes are escaped by back slashes.

Integer
: An integer is a word that contains only \[0-9\].

Float
: A float is a word that contains only \[0-9\] with a decimal point.

Ternary
: A ternary is represented by any one keyword of TRUE, FALSE or UNKNOWN.

Datetime
: A datetime is a string formatted as datetime.

  Strings of the form passed by the ["datetime-format" option]({{ '/reference/command.html#global_options' | relative_url }}) or the following forms can be converted to datetime values.
  
  | DateFormat | Example |
  | :- | :- |
  | YYYY-MM-DD | 2012-03-15 |
  | YYYY/MM/DD | 2012/03/15 |
  | YYYY-M-D   | 2012-3-15 |
  | YYYY/M/D   | 2012/3/15 |

  | DatetimeFormat | Example |
  | :- | :- |
  | DateFormat | 2012-03-15 |
  | DateFormat hh:mm:ss(.NanoSecods) | 2012-03-15 12:03:01<br />2012-03-15 12:03:01.123456789 |
  | DateFormat hh:mm:ss(.NanoSecods) Zhh:mm | 2012-03-15 12:03:01 -07:00 |
  | DateFormat hh:mm:ss(.NanoSecods) Zhhmm | 2012-03-15 12:03:01 -0700 |
  | DateFormat hh:mm:ss(.NanoSecods) TZ | 2012-03-15 12:03:01 PST |
  | RFC3339 | 2012-03-15T12:03:01-07:00 |
  | RFC3339 with Nano Seconds | 2012-03-15T12:03:01.123456789-07:00 |
  | RFC822 | 03 Mar 12 12:03 PST |
  | RFC822 with Numeric Zone | 03 Mar 12 12:03 -0700 |
  
  > Timezone abbreviations such as "PST" may not work properly depending on your environment, 
  > so you should use timezone offset such as "-07:00" as possible.

Null
: A null is represented by a keyword NULL.

Variable
: A variable is a word starting with "@" and followed by a character string that contains any unicode letters, any digits or low lines(U+005F '\_').

Flag
: A flag is a word starting with "@@" and followed by a character string that contains any unicode letters, any digits or low lines(U+005F '\_').

```sql
abcde                 -- identifier
識別子                 -- identifier
`ab+c\`de`            -- identifier
'abcd\'e'             -- string
"abcd\"e"             -- string
123                   -- integer
123.456               -- float
true                  -- ternary
'2012-03-15 12:03:01' -- datetime
null                  -- null
@var                  -- variable
@@flag                -- flag
```

## Comments
{: #comments}

Line Comment
: A single line comment starts with a string "--" and ends with a line-break character. 

Block Comment
: A block comment starts with a string "/\*" and ends with a string "\*/".


```sql
/*
 * Multi Line Comment
 */
VAR @id /* In Line Comment */ := 0;

-- Line Comment
SELECT @id := @id + 1 AS id, -- Line Comment
       name
  FROM user;
```

## Reserved Words
{: #reserved_words}

ABSOLUTE ADD AFTER AGGREGATE ALTER ALL AND ANY AS ASC
BEFORE BEGIN BETWEEN BREAK BY
CASE COMMIT CREATE CLOSE CONTINUE CROSS CURSOR
DECLARE DEFAULT DELETE DESC DISPOSE DISTINCT DO DROP DUAL
ELSE ELSEIF END EXCEPT EXISTS EXIT
FETCH FIRST FOR FROM FULL FUNCTION
GROUP
HAVING
IF IGNORE IN INNER INSERT INTERSECT INTO IS
JOIN
LAST LEFT LIKE LIMIT
NATURAL NEXT NOT NULL
OFFSET ON OPEN OR ORDER OUTER OVER
PARTITION PERCENT PRINT PRINTF PRIOR
RANGE RECURSIVE RELATIVE RENAME RETURN RIGHT ROLLBACK
SELECT SET SEPARATOR SOURCE STDIN
TABLE THEN TO TRIGGER
UNION UPDATE USING
VALUES VAR
WHEN WHERE WHILE WITH
