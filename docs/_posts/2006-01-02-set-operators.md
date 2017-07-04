---
layout: default
title: Set Operators - Reference Manual - csvq
category: reference
---

# Set Operators

| operator | description |
| :- | :- |
| [UNION](#union) | Return the union of result sets |
| [EXCEPT](#except)   | Return the relative complement of result sets  |
| [INTERSECT](#intersect) | Return the intersection of result sets |

A set operation combines result sets retrieved by select queries into a single result set.
If the ALL keyword is not specified, the result is distinguished.

## UNION
{: #union}

```sql
select_query UNION [ALL] select_query
```

_select_query_
: [select_set_entity]({{ '/reference/select-query.html' | relative_url }})

Return all records of both result sets.

## EXCEPT
{: #except}

```sql
select_query EXCEPT [ALL] select_query
```

_select_query_
: [select_set_entity]({{ '/reference/select-query.html' | relative_url }})

Return records of the result set of the left-hand side query that do not appear in the result set of the right-hand side query.

## INTERSECT
{: #intersect}

```sql
select_query INTERSECT [ALL] select_query
```

_select_query_
: [select_set_entity]({{ '/reference/select-query.html' | relative_url }})

Return only records that appear in both result sets.