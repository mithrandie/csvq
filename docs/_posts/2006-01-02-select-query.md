---
layout: default
title: Select Query - Reference Manual - csvq
category: reference
---

# Select Query

Select query is used to retrieve data from csv files.

```
select_query
  : [with_clause]
      select_entity
      [order_by_clause]
      [limit_clause]
      [FOR UPDATE]

select_entity
  : select_clause
      [from_clause]
      [where_clause]
      [group_by_clause]
      [having_clause]
  | select_set_entity set_operator [ALL] select_set_entity 

select_set_entity
  : select_entity
  | (select_query)
```

_with_clause_
: [With Clause](#with_clause)

_select_clause_
: [Select Clause](#select_clause)

_from_clause_
: [From Clause](#from_clause)

_where_clause_
: [Where Clause](#where_clause)

_group_by_clause_
: [Group By Clause](#group_by_clause)

_having_clause_
: [Having Clause](#having_clause)

_order_by_clause_
: [Order By Clause](#order_by_clause)

_limit_clause_
: [Limit Clause](#limit_clause)

_set_operator_
: [Set Operators]({{ '/reference/set-operators.html' | relative_url }})

## With Clause
{: #with_clause}

```sql
WITH common_table_expression [, common_table_expression ...]
```
_common_table_expression_
: [Common Table Expression]({{ '/reference/common-table-expression.html' | relative_url }})

## Select Clause
{: #select_clause}

```sql
SELECT [DISTINCT] field [, field ...]
```

### Distinct

You can use DISTINCT keyword to retrieve only unique records.

### field syntax

```sql
field
  : value
  | value AS alias
  | *
  | table_name.*
```

_value_
: [value]({{ '/reference/value.html' | relative_url }})

_alias_
: [identifier]({{ '/reference/statement.html#parsing' | relative_url }})

_table_name_
: [identifier]({{ '/reference/statement.html#parsing' | relative_url }})

_*_
: Asterisk(U+002A `*`) denotes all columns. 

  When used alone, the asterisk selects all columns in all tables; when used with a table name, it selects all columns in that table.

## From Clause
{: #from_clause}

```sql
FROM table [, {table|LATERAL laterable_table} ...]
```

If multiple tables have been enumerated, tables are joined using cross join.

### table syntax

```sql
table
  : table_entity
  | table_entity alias 
  | table_entity AS alias
  | join
  | DUAL
  | laterable_table
  | (table)

table_entity
  : table_identifier
  | table_object
  | inline_table_object

table_identifier
  : table_name
  | url
  | table_identification_function
  | STDIN
  
laterable_table
  : subquery
  | subquery alias
  | subquery AS alias

subquery
  : (select_query)

join
  : table CROSS JOIN table
  | table [INNER] JOIN table join_condition
  | table {LEFT|RIGHT|FULL} [OUTER] JOIN table join_condition
  | table NATURAL [INNER] JOIN table
  | table NATURAL {LEFT|RIGHT|FULL} [OUTER] JOIN table
  | table CROSS JOIN LATERAL laterable_table
  | table [INNER] JOIN LATERAL laterable_table join_condition
  | table LEFT [OUTER] JOIN LATERAL laterable_table join_condition
  | table NATURAL [INNER] JOIN LATERAL laterable_table
  | table NATURAL LEFT [OUTER] JOIN LATERAL laterable_table

join_condition
  : ON condition
  | USING (column_name [, column_name, ...])

table_identification_function
  : FILE::(file_path)
  : INLINE::(file_path)
  : URL::(url_string)
  : DATA::(data_string)

table_object
  : CSV(delimiter, table_identifier [, encoding [, no_header [, without_null]]])
  | FIXED(delimiter_positions, table_identifier [, encoding [, no_header [, without_null]]])
  | JSON(json_query, table_identifier)
  | JSONL(json_query, table_identifier)
  | LTSV(table_identifier [, encoding [, without_null]])

inline_table_object  -- Deprecated. Table identification functions can be used instead.
  : CSV_INLINE(delimiter, inline_table_identifier [, encoding [, no_header [, without_null]]])
  | CSV_INLINE(delimiter, csv_data)
  | JSON_INLINE(json_query, inline_table_identifier [, encoding [, no_header [, without_null]]])
  | JSON_INLINE(json_query, json_data)

inline_table_identifier
  : table_name
  | url_identifier

```

_table_name_
: [identifier]({{ '/reference/statement.html#parsing' | relative_url }})
  
  A _table_name_ represents a file path, a [temporary table]({{ '/reference/temporary-table.html' | relative_url }}), or a [inline table]({{ '/reference/common-table-expression.html' | relative_url }}).
  You can use absolute path or relative path from the directory specified by the ["--repository" option]({{ '/reference/command.html#options' | relative_url }}) as a file path.
  
  When the file name extension is ".csv", ".tsv", ".json", ".jsonl" or ".txt", the format to be loaded is automatically determined by the file extension, and you can omit it. 
  
  ```sql
  FROM `user.csv`          -- Relative path
  FROM `/path/to/user.csv` -- Absolute path
  FROM user                -- Relative path without file extension
  ```
  
  The specifications of the command options are used as file attributes such as encoding to be loaded. 
  If you want to specify the different attributes for each file, you can use _table_object_ expressions for each file to load.

  Once a file is loaded, then the data is cached, and it can be loaded with only file name after that within the transaction.

_url_
: A string of characters representing URL starting with a schema name and a colon.

  "http", "https" and "file" schemes are available.

  ```sql
  https://example.com/files/data.csv       -- Remote resource downloaded using HTTP GET method
  file:///C:/Users/yourname/files/data.csv -- Local file specified by absolute path
  file:./data.csv                          -- Local file specified by relative path
  ```

  An inline table is created from remote resources.
  The downloaded data is cached until the transaction ends.

  The file format is automatically determined when the http response specifies the following content types.

| MIME type        | Format |
|:-----------------|:-------|
| text/csv         | CSV    |
| application/json | JSON   |

_table_identification_function_
: Function notation with a name followed by two colons.

  - FILE::(file_path)

    file_path: [string]({{ '/reference/value.html#string' | relative_url }})

    This is the same as specifying a file using _table_name_.
  
  - INLINE::(file_path)

    file_path: [string]({{ '/reference/value.html#string' | relative_url }})

    Files read by this function are not cached and cannot be updated.

  - URL::(url_string)

    url_string: [string]({{ '/reference/value.html#string' | relative_url }})

    When specifying a resource using _url_, the path must be encoded, but this function does not require encoding.

  - DATA::(data_string)

    file_path: [string]({{ '/reference/value.html#string' | relative_url }})

    This function creates an inline table from a string.

  Example of use in a query:
  
  ```sql
  SELECT id,
         tag_name,
         (SELECT COUNT(*) FROM JSON('', DATA::(assets))) AS number_of_assets,
         published_at
    FROM https://api.github.com/repos/mithrandie/csvq/releases
   WHERE prerelease = false
   ORDER BY published_at DESC
   LIMIT 10
  ```

_alias_
: [identifier]({{ '/reference/statement.html#parsing' | relative_url }})

  If _alias_ is not specified, _table_name_ stripped its directory path and extension is used as alias.

  ```sql
  -- Following expressions are equivalent
  FROM `/path/to/user.csv`
  FROM `/path/to/user.csv` AS user
  ```

_select_query_
: [Select Query]({{ '/reference/select-query.html' | relative_url }})

_condition_
: [value]({{ '/reference/value.html' | relative_url }})

_column_name_
: [identifier]({{ '/reference/statement.html#parsing' | relative_url }})

_delimiter_  
: [string]({{ '/reference/value.html#string' | relative_url }})

_json_query_
: [JSON Query]({{ '/reference/json.html#query' | relative_url }})

  Empty string is equivalent to "{}".

_delimiter_positions_  
: [string]({{ '/reference/value.html#string' | relative_url }})

  "SPACES" or JSON Array of integers

_encoding_
: [string]({{ '/reference/value.html#string' | relative_url }}) or [identifier]({{ '/reference/statement.html#parsing' | relative_url }})
  
  "AUTO", "UTF8", "UTF8M", "UTF16", "UTF16BE", "UTF16LE", "UTF16BEM", "UTF16LEM" or "SJIS".

_no_header_
: [boolean]({{ '/reference/value.html#boolean' | relative_url }})

_without_null_
: [boolean]({{ '/reference/value.html#boolean' | relative_url }})

_url_identifier_
: [identifier]({{ '/reference/statement.html#parsing' | relative_url }})

  A URL of the http or https scheme to refer to a resource.

_csv_data_
: [string]({{ '/reference/value.html#string' | relative_url }})

_json_data_
: [string]({{ '/reference/value.html#string' | relative_url }})

#### Special Tables
{: #special_tables}

DUAL
: The dual table has one column and one record, and the only field is empty.
  This table is used to retrieve pseudo columns.

STDIN
: The stdin table loads data from pipe or redirection as a csv data.
  The stdin table is one of [temporary tables]({{ '/reference/temporary-table.html' | relative_url }}) that is declared automatically.
  This table cannot to be used in the interactive shell.


## Where Clause
{: #where_clause}

The Where clause is used to filter records.

```sql
WHERE condition
```

_condition_
: [value]({{ '/reference/value.html' | relative_url }})

## Group By Clause
{: #group_by_clause}

The Group By clause is used to group records.

```sql
GROUP BY field [, field ...] 
```

_field_
: [value]({{ '/reference/value.html' | relative_url }})

## Having Clause
{: #having_clause}

The Having clause is used to filter grouped records.

```sql
HAVING condition
```

_condition_
: [value]({{ '/reference/value.html' | relative_url }})

## Order By Clause
{: #order_by_clause}

The Order By clause is used to sort records.

```sql
ORDER BY order_item [, order_item ...]
```

### order item

```sql
order_item
  : field [order_direction] [null_position]
  
order_direction
  : {ASC|DESC}
  
null_position
  : NULLS {FIRST|LAST}
```

_field_
: [value]({{ '/reference/value.html' | relative_url }})
  
  If DISTINCT keyword is specified in the select clause, you can use only enumerated fields in the select clause as _field_.

_order_direction_
: _ASC_ sorts records in ascending order. _DESC_ sorts in descending order. _ASC_ is the default.

_null_position_
: _FIRST_ puts null values first. _LAST_ puts null values last. 
  If _order_direction_ is specified as _ASC_ then _FIRST_ is the default, otherwise _LAST_ is the default.


## Limit Clause
{: #limit_clause}

The Limit clause is used to specify the maximum number of records to return and exclude the first set of records.

```sql
limit_clause
  : LIMIT number_of_records [{ROW|ROWS}] [{ONLY|WITH TIES}] [offset_clause]
  | LIMIT percentage PERCENT [{ONLY|WITH TIES}] [offset_clause]
  | [offset_clause] FETCH {FIRST|NEXT} number_of_records {ROW|ROWS} [{ONLY|WITH TIES}]
  | [offset_clause] FETCH {FIRST|NEXT} percentage PERCENT [{ONLY|WITH TIES}]
  | offset_clause

offset_clause
  : OFFSET number_of_records [{ROW|ROWS}]
```

_number_of_records_
: [integer]({{ '/reference/value.html#integer' | relative_url }})

_percent_
: [float]({{ '/reference/value.html#integer' | relative_url }})

_ROW_ and _ROWS_ after _number_of_records_, _FIRST_ and _NEXT_ after _FETCH_, and _ONLY_ keyword does not affect the result.

If _PERCENT_ keyword is specified, maximum number of records is _percentage_ percent of the result set that includes the excluded records by _offset_clause_. 

If _WITH TIES_ keywords are specified, all records that have the same sort keys specified by _Order By Clause_ as the last record of the limited records are included in the records to return.
If there is no _Order By Clause_ in the query, _WITH TIES_ keywords are ignored.

