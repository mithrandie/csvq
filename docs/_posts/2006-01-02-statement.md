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
In statements, character case of keywords is ignored.

If you want to execute a single query, you can omit the terminal semicolon.  

### Interactive Shell

When the csvq command is called with no argument and no "--source" (or "-s") option, the interactive shell is launched.
You can use the interactive shell in order to sequencial input and execution.

If you want to continue to input the statement on the next line, you can use Backslash(U+005C `\`) at the end of the line to continue.

#### Command options in the interactive shell

--out
: Ignored 

--stats
: Show only Query Execution Time

```bash
# Execute a single query
$ csvq 'SELECT id, name FROM user'

# Execute statements
$ csvq 'VAR @id := 0; SELECT @id := @id + 1 AS id, name FROM user;'

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
: An identifier is a word starting with any unicode letter or a Low Line(U+005F `_`) and followed by a character string that contains any unicode letters, any digits or Low Lines(U+005F `_`).
  You cannot use [reserved words](#reserved_words) as an identifier.

  Notwithstanding above naming restriction, you can use most character strings as an identifier by enclosing in Grave Accents(U+0060 \` ) or Quotation Marks(U+0022 `"`) if [--ansi-quotes]({{ '/reference/command.html#options' | relative_url }}) is specified. 
  Enclosure characters are escaped by back slashes or double enclosures.
  
  Identifiers represent tables, columns, functions or cursors.
  Character case is insensitive except file paths, and whether file paths are case insensitive or not depends on your file system.
  
String
: A string is a character string enclosed in Apostrophes(U+0027 `'`) or Quotation Marks(U+0022 `"`) if [--ansi-quotes]({{ '/reference/command.html#options' | relative_url }}) is not specified.
  In a string, enclosure characters are escaped by back slashes or double enclosures.

Integer
: An integer is a word that contains only \[0-9\].

Float
: A float is a word that contains only \[0-9\] with a decimal point, or its exponential notation.

Ternary
: A ternary is represented by any one keyword of TRUE, UNKNOWN or FALSE.

Null
: A null is represented by a keyword NULL.

Variable
: A [variable]({{ '/reference/variable.html' | relative_url }}) is a word starting with "@" and followed by a character string that contains any unicode letters, any digits or Low Lines(U+005F `_`).

Flag
: A [flag]({{ '/reference/flag.html' | relative_url }}) is a word starting with "@@" and followed by a character string that contains any unicode letters, any digits or Low Lines(U+005F `_`). Character case is ignored.

Environment Variable
: A [environment variable]({{ '/reference/environment-variable.html' | relative_url }}) is a word starting with "@%" and followed by a character string that contains any unicode letters, any digits or Low Lines(U+005F `_`).
  If a environment variable includes other characters, you can use the variable by enclosing in Back Quotes(U+0060 ` ).

Runtime Information
: A [runtime information]({{ '/reference/runtime-information.html' | relative_url }}) is a word starting with "@#" and followed by a character string that contains any unicode letters, any digits or Low Lines(U+005F `_`). Character case is ignored.

```
abcde                 -- identifier
識別子                 -- identifier
`abc\`de`             -- identifier
`abc``de`             -- identifier
'abcd\'e'             -- string
'abcd''e'             -- string
123                   -- integer
123.456               -- float
true                  -- ternary
null                  -- null
@var                  -- variable
@@FLAG                -- flag
@%ENV_VAR             -- environment variable
@%`ENV_VAR`           -- environment variable
@#INFO                -- runtime information

/* if --ansi-quotes is specified */
"abcd\"e"             -- identifier
"abcd""e"             -- identifier

/* if --ansi-quotes is not specified */
"abcd\"e"             -- string
"abcd""e"             -- string
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
CASE CHDIR CLOSE COMMIT CONTINUE COUNT CREATE CROSS CSV_INLINE CUME_DIST CURRENT CURSOR
DECLARE DEFAULT DELETE DENSE_RANK DESC DISPOSE DISTINCT DO DROP DUAL
ECHO ELSE ELSEIF END EXCEPT EXECUTE EXISTS EXIT
FALSE FETCH FIRST FIRST_VALUE FOLLOWING FOR FROM FULL FUNCTION
GROUP
HAVING
IF IGNORE IN INNER INSERT INTERSECT INTO IS
JOIN JSONL JSON_AGG JSON_INLINE JSON_OBJECT JSON_ROW JSON_TABLE
LAG LAST LAST_VALUE LATERAL LEAD LEFT LIKE LIMIT LISTAGG
MAX MEDIAN MIN
NATURAL NEXT NOT NTH_VALUE NTILE NULL
OFFSET ON ONLY OPEN OR ORDER OUTER OVER
PARTITION PERCENT PERCENT_RANK PRECEDING PREPARE PRINT PRINTF PRIOR PWD
RANGE RANK RECURSIVE RELATIVE RELOAD REMOVE RENAME REPLACE RETURN RIGHT ROLLBACK ROW ROW_NUMBER
SELECT SEPARATOR SET SHOW SOURCE STDEV STDEVP STDIN SUBSTRING SUM SYNTAX
TABLE THEN TO TRIGGER TRUE
UNBOUNDED UNION UNKNOWN UNSET UPDATE USING
VALUES VAR VARP VIEW
WHEN WHERE WHILE WITH WITHIN

