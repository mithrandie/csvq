---
layout: default
title: Comparison Operators - Reference Manual - csvq
category: reference
---

# Comparison Operators

| \= |
| < |
| <\= |
| > |
| >\= |
| <> |
| !\= |

```
value IS [NOT] ternary
value IS [NOT] null
```

```
value [NOT] BETWEEN value AND value
```

```
value [NOT] IN (value, ...)
value [NOT] IN (select_query)
```

```
value [NOT] LIKE string
```

```
value comparison_operator ANY (select_query)
```

```
value comparison_operator ALL (select_query)
```

```
EXISTS (select_query)
```