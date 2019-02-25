---
layout: default
title: External Command - Reference Manual - csvq
category: reference
---

# External Command

You can run an external command by placing a Dollar Sign(U+0024 `$`) at the beggining of the line.
The result is written to the standard output.

[CALL Function]({{ '/reference/system-functions.html#call' | relative_url }}) also runs an external command, but the result is returned as a string.

```
$ command [args ...]
```

Arguments are separated with spaces. 
Strings including spaces are treated as a single string by enclosing in Apostrophes(U+0027 `'`) or Quotation Marks(U+0022 `"`).
Embedded expressions including spaces are also treated as chunks. 
  
In the arguments, following expressions can be embedded.

* [Variable]({{ '/reference/variable.html' | relative_url }})
* [Environment Variable]({{ '/reference/environment-variable.html' | relative_url }})
* [Runtime Information]({{ '/reference/runtime-information.html' | relative_url }})
* Embedded Expression written in the following format

## Embedded Expression
{: #embedded-expression}

```
${value_expression}
```

Expressions can be embedded in arguments by enclosing in a Dollar Sign(U+0024 `$`) and Curly Brackets(U+007B, U+007D `{}`).
The result of evaluation of the expression must be a single [value]({{ '/reference/value.html' | relative_url }}).


### Examples
```bash
csvq > $echo 'abc'
abc
csvq > VAR @ARG := 'argstr'
csvq > $echo @ARG
argstr
csvq > $echo @%HOME
/home/mithrandie
csvq > $echo ${@%HOME || '/docs'}
/home/mithrandie/docs
csvq > $sh -c 'echo ${@%HOME || "/docs"} | wc'
      1       1      22
csvq >
```
