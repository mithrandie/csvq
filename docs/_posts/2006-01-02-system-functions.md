---
layout: default
title: System Functions - Reference Manual - csvq
category: reference
---

# System Functions

| name | description |
| :- | :- |
| [AUTO_INCREMENT](#auto_increment) | Return a sequence number in a single query |
| [CALL](#call) | Execute a external command |

## Definitions

### AUTO_INCREMENT
{: #auto_increment}

```
AUTO_INCREMENT()
```

Return a sequence number that starts with 1 in a single query.

```
AUTO_INCREMENT(initial_value)
```

_initial_value_
: [integer]({{ '/reference/value.html#integer' | relative_url }})

_return_
: [integer]({{ '/reference/value.html#integer' | relative_url }})

Return a sequence number that starts with _initial_value_ in a single query.

### CALL
{: #call}

```
CALL(command [, argument ...])
```

_command_
: [string]({{ '/reference/value.html#string' | relative_url }})

_argument_
: [string]({{ '/reference/value.html#string' | relative_url }})

_return_
: [string]({{ '/reference/value.html#string' | relative_url }})

Execute a external command then return the standard output as a string.
If the external command failed, the executing procedure is terminated with an error.