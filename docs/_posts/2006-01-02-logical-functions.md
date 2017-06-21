---
layout: default
title: Logical Functions - Reference Manual - csvq
category: reference
---

# Logical Functions

| name | description |
| :- | :- |
| [COALESCE](#coalesce) | Return the first non-null value in arguments |
| [IF](#if) | Return a value by condition |
| [IFNULL](#ifnull) | Return a value whether passed value is null |
| [NULLIF](#nullif) | Return null wheter passed values are equal |
| [CASE](#case) |  |

## Definitions

### COALESCE
{: #coalesce}

```
COALESCE(value, ...) return value
```

Return the first non-null value in arguments. If there are no non-null value, return null.

### IF
{: #if}

```
IF(condition, value1, value2) return value
```

If condition is TRUE, returns value1. Otherwise returns value2.

### IFNULL
{: #ifnull}

```
IFNULL(value1, value2) return value
```

If value1 is null, return value2. Otherwise return value1.

### NULLIF
{: #nullif}

```
NULLIF(value1, value2) return value
```

If value1 is equal to value2, return null. Otherwise return value1.

### CASE
{: #case}

```
CASE value WHEN value THEN value [WHEN value THE value ...] [ELSE value] END
```

```
CASE WHEN condition THEN value [WHEN condition THE value ...] [ELSE value] END
```