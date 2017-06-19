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

You can pass a query or Statements as a csvq command argument or source file.

Statements are terminated with semicolons. 
Stetaments are processed sequentially for each statement.
In statements, character case is ignored.

If you want to execute a sigle query, you can omit a terminal semicolon.  

```bash
# Execute a single query
$ csvq "SELECT id, name FROM user"

# Execute statements
$ csvq "var @id := 0; SELECT @id := @id + 1 AS id, name FROM user;"

# Load statements from file
$ cat statements.sql
var @id := 0;
SELECT @id := @id + 1 AS id,
       name
  FROM user;

$ csvq -s statements.sql
```
## Parsing
{: #parsing}

You can use following types in statements.

Identifier
: A identifier is a word starting with a character \[a-z\|A-Z\|\_\], contains any of characters \[a-z\|A-Z\|0-9\|\_\].
  You cannot use [reserved words](#reserved_words) as identifier.

  Notwithstanding above naming restriction, you can use most characters as identifier by enclusing in back quotes.
  
  Identifiers represent tables, columns or cursors. 
  
String
: A string is a character string enclosed in single quotes or double quotes.
  In a string, single quotes or double quotes are escaped by back slash.

Integer
: A integer is a word contains only \[0-9\].

Float
: A float is a word contains only \[0-9\] with a decimal point.

Ternary
: A ternary is represented by any one keyword of \[TRUE\|FALSE\|UNKNOWN\].

Datetime
: A datetime is a string formatted as datetime.

  **Datetime format:**
  
  | Format | Example |
  | :- | :- |
  | YYYY-MM-DD | 2012-03-15 |
  | YYYY-MM-DD HH:mi:ss | 2012-03-15 12:03:01 |
  | YYYY-MM-DD HH:mi:ss.Nano | 2012-03-15 12:03:01.123456789 |
  | YYYY-MM-DD HH:mi:ss TZ | 2012-03-15 12:03:01 PST |
  | RFC3339 | 2012-03-15T12:03:01-07:00 |
  | RFC3339 with Nano Seconds | 2012-03-15T12:03:01.123456789-07:00 |
  | RFC822 | 03 Mar 12 12:03 PST |
  | RFC822 with Numeric Zone | 03 Mar 12 12:03 -0700 |

Null
: A null is represented by a keyword "NULL"

Variable
: A variable is a word starting with "@", contains any of characters \[a-z\|A-Z\|0-9\|\_\]

Flag
: A flag is a word starting with "@@", contains any of characters \[a-z\|A-Z\|0-9\|\_\]

## Comments
{: #comments}

Line Comment
: A single line comment starts with a string "--", ends with a line-break character. 

Block Comment
: A block comment starts with a string "/\*", ends with a string "\*/"


```sql
/*
 * Multi Line Comment
 */
var @id /* In-line Comment */ := 0;

-- Line Comment
SELECT @id := @id + 1 AS id, -- Line Comment
       name
  FROM user;
```

## Reserved Words
{: #reserved_words}

IDENTIFIER STRING INTEGER FLOAT BOOLEAN TERNARY DATETIME VARIABLE FLAG
SELECT FROM UPDATE SET DELETE WHERE INSERT INTO VALUES AS DUAL STDIN
CREATE ADD DROP ALTER TABLE FIRST LAST AFTER BEFORE DEFAULT RENAME TO
ORDER GROUP HAVING BY ASC DESC LIMIT JOIN INNER OUTER LEFT RIGHT FULL
CROSS ON USING NATURAL UNION ALL ANY EXISTS IN AND OR NOT BETWEEN LIKE
IS NULL DISTINCT WITH CASE IF ELSEIF WHILE WHEN THEN ELSE DO END
DECLARE CURSOR FOR FETCH OPEN CLOSE DISPOSE GROUP_CONCAT SEPARATOR
COMMIT ROLLBACK CONTINUE BREAK EXIT PRINT VAR