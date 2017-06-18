---
layout: default
title: Alter Table Query - Reference Manual - csvq
category: reference
---

# Alter Table Query

## Add Columns

```
ALTER TABLE table_name ADD column_name [DEFAULT value] [FIRST|LAST|AFTER column_name|BEFORE column_name]
ALTER TABLE table_name ADD (column_name [DEFAULT value], [column_name [DEFAULT value], ...]) [FIRST|LAST|AFTER column_name|BEFORE column_name]
```

## Drop Columns

```
ALTER TABLE table_name DROP column_name
ALTER TABLE table_name DROP (column_name, [column_name, ...])
```

## Rename Column

```
ALTER TABLE table_name RENAME old_column_name TO new_column_name
```