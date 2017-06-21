---
layout: default
title: Operator Precedence - Reference Manual - csvq
category: reference
---

# Operator Precedence

The following table list operators from highest precedence to lowest.

| precedence | operators | Associativity |
| :- | :- | :- |
| 1 | [*]({{ '/reference/arithmetic-operators.html' | relative_url }})       | Left-to-right | 
|   | [/]({{ '/reference/arithmetic-operators.html' | relative_url }})       | Left-to-right | 
|   | [%]({{ '/reference/arithmetic-operators.html' | relative_url }})       | Left-to-right | 
| 2 | [+]({{ '/reference/arithmetic-operators.html' | relative_url }})       | Left-to-right | 
|   | [-]({{ '/reference/arithmetic-operators.html' | relative_url }})       | Left-to-right | 
| 3 | [\|\|]({{ '/reference/string-operators.html' | relative_url }})    | Left-to-right | 
| 4 | [\=]({{ '/reference/comparison-operators.html#relational_operators' | relative_url }})  | nonassoc | 
|   | [<]({{ '/reference/comparison-operators.html#relational_operators' | relative_url }})   | nonassoc | 
|   | [<\=]({{ '/reference/comparison-operators.html#relational_operators' | relative_url }}) | nonassoc | 
|   | [>]({{ '/reference/comparison-operators.html#relational_operators' | relative_url }})   | nonassoc | 
|   | [>\=]({{ '/reference/comparison-operators.html#relational_operators' | relative_url }}) | nonassoc | 
|   | [<>]({{ '/reference/comparison-operators.html#relational_operators' | relative_url }})  | nonassoc | 
|   | [!\=]({{ '/reference/comparison-operators.html#relational_operators' | relative_url }}) | nonassoc | 
|   | [IS]({{ '/reference/comparison-operators.html#is' | relative_url }})           | nonassoc | 
|   | [BETWEEN]({{ '/reference/comparison-operators.html#between' | relative_url }}) | nonassoc | 
|   | [IN]({{ '/reference/comparison-operators.html#in' | relative_url }})           | nonassoc | 
|   | [LIKE]({{ '/reference/comparison-operators.html#like' | relative_url }})       | nonassoc | 
| 5 | [NOT]({{ '/reference/logic-operators.html#not' | relative_url }})     | Right-to-left | 
| 6 | [AND]({{ '/reference/logic-operators.html#and' | relative_url }})     | Left-to-right | 
| 7 | [OR]({{ '/reference/logic-operators.html#or' | relative_url }})       | Left-to-right | 

