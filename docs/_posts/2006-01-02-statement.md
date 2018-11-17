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

A statement is terminated with a semicolon. 
In statements, character cases of keywords are ignored.

If you want to execute a single query, you can omit the terminal semicolon.  

### Interactive Shell

When the csvq command is called with no argument and no "--source" (or "-s") option, the interactive shell is launched.
You can use the interactive shell in order to sequencial input and execution.
In the interactive shell, statements are executed when a line ends with a semicolon.

If you want to continue to input a statement on the next line even though the end of the line is a semicolon, you can use a backslash at the end of the line to continue.

#### Command options in the interactive shell

--out
: Ignored 

--stats
: Show only Query Execution Time

#### Line editor in the interactive shell

On the some systems, the interactive shell provides a more powerful line editor by using the package [https://github.com/chzyer/readline](https://github.com/chzyer/readline).
Not all features of the readline package are available, but its short cut keys, command history and features like that will help you.

The readline package is used on the following systems.
- darwin dragonfly freebsd linux netbsd openbsd solaris windows


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

# Execute in the interactive shell
$ csvq
csvq > SELECT id, name FROM users;
+----+-------+
| id | name  |
+----+-------+
| 1  | Louis |
| 2  | Sean  |
+----+-------+
csvq > UPDATE users SET name = 'Mildred' WHERE id = 2;
1 record updated on "/home/mithrandie/docs/csv/users.csv".
csvq > SELECT id, name FROM users;
+----+----------+
| id | name     |
+----+----------+
| 1  | Louis    |
| 2  | Mildred  |
+----+----------+
csvq > COMMIT;
Commit: file "/home/mithrandie/docs/csv/users.csv" is updated.
csvq > IF (SELECT name FROM users WHERE id = 2) = 'Mildred' THEN
     >   PRINT TRUE; \
     > ELSE
     >   PRINT FALSE; \
     > END IF;
TRUE
csvq > EXIT;
```

## Parsing
{: #parsing}

You can use following types in statements.

Identifier
: A identifier is a word starting with any unicode letter or a low line(U+005F '\_') and followed by a character string that contains any unicode letters, any digits or low lines(U+005F '\_').
  You cannot use [reserved words](#reserved_words) as a identifier.

  Notwithstanding above naming restriction, you can use most character strings as a identifier by enclosing in back quotes(U+0060 '`').
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

  &nbsp;

  | DatetimeFormat | Example |
  | :- | :- |
  | DateFormat | 2012-03-15 |
  | DateFormat hh:mm:ss(.NanoSecods) | 2012-03-15 12:03:01<br />2012-03-15 12:03:01.123456789 |
  | DateFormat hh:mm:ss(.NanoSecods) ±hh:mm | 2012-03-15 12:03:01 -07:00 |
  | DateFormat hh:mm:ss(.NanoSecods) ±hhmm | 2012-03-15 12:03:01 -0700 |
  | DateFormat hh:mm:ss(.NanoSecods) TimeZone | 2012-03-15 12:03:01 PST |
  | YYYY-MM-DDThh:mm:ss(.NanoSeconds) | 2012-03-15T12:03:01 |
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

ABSOLUTE ADD AFTER AGGREGATE ALTER ALL AND ANY AS ASC AVG
BEFORE BEGIN BETWEEN BREAK BY
CASE CLOSE COMMIT CONTINUE COUNT CREATE CROSS CUME_DIST CURRENT CURSOR
DECLARE DEFAULT DELETE DENSE_RANK DESC DISPOSE DISTINCT DO DROP DUAL
ELSE ELSEIF END EXCEPT EXECUTE EXISTS EXIT
FALSE FETCH FIRST FIRST_VALUE FOLLOWING FOR FROM FULL FUNCTION
GROUP
HAVING
IF IGNORE IN INNER INSERT INTERSECT INTO IS
JOIN JSON_AGG JSON_OBJECT JSON_ROW JSON_TABLE
LAG LAST LAST_VALUE LEAD LEFT LIKE LIMIT LISTAGG
MAX MEDIAN MIN
NATURAL NEXT NOT NTH_VALUE NTILE NULL
OFFSET ON OPEN OR ORDER OUTER OVER
PARTITION PERCENT PERCENT_RANK PRECEDING PRINT PRINTF PRIOR
RANGE RANK RECURSIVE RELATIVE RENAME RETURN RIGHT ROLLBACK ROW ROW_NUMBER
SELECT SEPARATOR SET SHOW SOURCE STDIN SUM
TABLE THEN TO TRIGGER TRUE
UNBOUNDED UNION UNKNOWN UNSET UPDATE USING
VALUES VAR VIEW
WHEN WHERE WHILE WITH WITHIN

