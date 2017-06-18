---
layout: default
title: Select Query - Reference Manual - csvq
category: reference
---

# Select Query

```
select_clause [from_clause] [where_clause] [group_by_clause] [having_clause] [order_by_clause] [limit_clause]
```

## Select Clause

```
SELECT [DISTINCT] field, [field, ...]
```

### field

```
value
value AS alias
```

## From Clause

```
FROM table, [table, ...]
```

### table

```
table_name
table_name alias 
table_name AS alias
subquery
subquery alias
subquery AS alias
table CROSS JOIN table
table NATURAL [INNER] JOIN table
table [INNER] JOIN table join_condition
table NATURAL [LEFT|RIGHT|FULL] OUTER JOIN table
table NATURAL {LEFT|RIGHT|FULL} [OUTER] JOIN table
table [LEFT|RIGHT|FULL] OUTER JOIN table join_condition
table {LEFT|RIGHT|FULL} [OUTER] JOIN table join_condition
```

### join condition

```
ON condition
USING (column, [column, ...])
```

## Where Clause

```
WHERE condition
```

## Group By Clause

```
GROUP BY value, [value, ...] 
```

## Having Clause

```
HAVING condition
```

## Order By Clause

```
ORDER BY order_item, [order_item, ...]
```

### order item

```
field_name
field_name [ASC|DESC]
```

## Limit Clause

```
LIMIT integer
```