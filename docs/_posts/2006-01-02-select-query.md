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
      [offset_clause]

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

_offset_clause_
: [Offset Clause](#offsetma_clause)

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
```

_value_
: [value]({{ '/reference/value.html' | relative_url }})

_alias_
: [identifier]({{ '/reference/statement.html#parsing' | relative_url }})

## From Clause
{: #from_clause}

```sql
FROM table [, table ...]
```

If multiple tables have been enumerated, tables are joined using cross join.

### table syntax

```sql
table
  : table_name
  | table_name alias 
  | table_name AS alias
  | subquery
  | subquery alias
  | subquery AS alias
  | join
```

_table_name_
: [identifier]({{ '/reference/statement.html#parsing' | relative_url }})
  or [special tables](#special_tables)
  
  A identifier represents a csv file path.
  You can use absolute path or relative path from the directory specified by the ["--repository" option]({{ '/reference/command.html#global_options' | relative_url }}).
  
  If a file name extension is ".csv" or ".tsv", you can omit it. 
  
  ```sql
  FROM `user.csv`          -- Relative path
  FROM `/path/to/user.csv` -- Absolute path
  FROM user                -- Relative path without file extension
  ```

_alias_
: [identifier]({{ '/reference/statement.html#parsing' | relative_url }})

  If _alias_ is not specified, _table_name_ stripped it's directory path and extension is used as alias.

  ```sql
  -- Following expressions are the same
  FROM `/path/to/user.csv`
  FROM `/path/to/user.csv` AS user
  ```

_subquery_
: A select query enclosed in parentheses.

#### Special Tables
{: #special_tables}

DUAL
: The dual table has one column and one record, and the only field is empty.
  This table is used to retrieve pseudo columns.

STDIN
: The stdin table loads data from the standard input as a csv data.

### join syntax

```sql
join
  : table CROSS JOIN table
  | table NATURAL [INNER] JOIN table
  | table [INNER] JOIN table [join_condition]
  | table NATURAL [LEFT|RIGHT|FULL] OUTER JOIN table
  | table NATURAL {LEFT|RIGHT|FULL} [OUTER] JOIN table
  | table [LEFT|RIGHT|FULL] OUTER JOIN table [join_condition]
  | table {LEFT|RIGHT|FULL} [OUTER] JOIN table [join_condition]
```

### join condition syntax

```sql
join_condition
  : ON condition
  | USING (column_name [, column_name, ...])
```

_condition_
: [value]({{ '/reference/value.html' | relative_url }})

_column_name_
: [identifier]({{ '/reference/statement.html#parsing' | relative_url }})

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
GROUP BY value [, value ...] 
```

_value_
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
  : field_name [order_direction] [null_position]
  
order_direction
  : {ASC|DESC}
  
null_position
  : NULLS {FIRST|LAST}
```

_field_name_
: [value]({{ '/reference/value.html' | relative_url }})
  
  If DISTINCT keyword is specified in the select clause, you can use only enumerated fields in the select clause as _field_name_.

_order_direction_
: _ASC_ sorts records in ascending order. _DESC_ sorts in descending order. _ASC_ is the default.

_null_position_
: _FIRST_ puts null values first. _LAST_ puts null values last. 
  If _order_direction_ is specified as _ASC_ then _FIRST_ is the default, otherwise _LAST_ is the default.


## Limit Clause
{: #limit_clause}

The Limit clause is used to specify the maximum number of records to return.

```sql
limit_clause
  : LIMIT number_of_records [WITH TIES]
  | LIMIT percent PERCENT [WITH TIES]
```

_number_of_records_
: [integer]({{ '/reference/value.html#integer' | relative_url }})

_percent_
: [float]({{ '/reference/value.html#integer' | relative_url }})

If _PERCENT_ keyword is specified, maximum number of records is _percent_ percent of the result set that includes the excluded records by _Offset Clause_. 

If _WITH TIES_ keywords are specified, all records that have the same sort keys specified by _Order By Clause_ as the last record of the limited records are included in the records to return.
If there is no _Order By Clause_ in the query, _WITH TIES_ keywords are ignored.

## Offset Clause
{: #offset_clause}

The Offset clause is used to exclude the first set of records.

```sql
OFFSET number
```

_number_
: [integer]({{ '/reference/value.html#integer' | relative_url }})
