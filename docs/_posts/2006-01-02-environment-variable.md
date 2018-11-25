---
layout: default
title: Environment Variable - Reference Manual - csvq
category: reference
---

# Environment Variable

* [SET ENVIRONMENT_VARIABLE](#set_env)
* [UNSET ENVIRONMENT_VARIABLE](#unset_env)


### Set Environment Variable
{: #set_env}

```sql
SET @%env_var_name TO value;
SET @%env_var_name = value;
```

_value_
: [string]({{ '/reference/value.html#string' | relative_url }}) or [identifier]({{ '/reference/statement.html#parsing' | relative_url }})

A set environment variable statement is used to set the value to the environment variable. 


### Unset Environment Variable
{: #unset_env}

```sql
UNSET @%env_var_name;
```

A unset environment variable statement is remove the environment variable. 


