---
layout: default
title: Operator Precedence - Reference Manual - csvq
category: reference
---

# Operator Precedence

The following table list operators from the highest precedence to lowest.

| precedence | operators | associativity |
| :- | :- | :- |
| 1  | [+ (unary plus)]({{ '/reference/arithmetic-operators.html#unary' | relative_url }})  | Right-to-left | 
|    | [- (unary minus)]({{ '/reference/arithmetic-operators.html#unary' | relative_url }}) | Right-to-left | 
|    | [!]({{ '/reference/logic-operators.html#not' | relative_url }})                      | Right-to-left | 
| 2  | [*]({{ '/reference/arithmetic-operators.html' | relative_url }})       | Left-to-right | 
|    | [/]({{ '/reference/arithmetic-operators.html' | relative_url }})       | Left-to-right | 
|    | [%]({{ '/reference/arithmetic-operators.html' | relative_url }})       | Left-to-right | 
| 3  | [+]({{ '/reference/arithmetic-operators.html' | relative_url }})       | Left-to-right | 
|    | [-]({{ '/reference/arithmetic-operators.html' | relative_url }})       | Left-to-right | 
| 4  | [\|\|]({{ '/reference/string-operators.html' | relative_url }})    | Left-to-right | 
| 5  | [\=]({{ '/reference/comparison-operators.html#relational_operators' | relative_url }})  | nonassoc | 
|    | [\=\=]({{ '/reference/comparison-operators.html#relational_operators' | relative_url }})   | nonassoc | 
|    | [<]({{ '/reference/comparison-operators.html#relational_operators' | relative_url }})   | nonassoc | 
|    | [<\=]({{ '/reference/comparison-operators.html#relational_operators' | relative_url }}) | nonassoc | 
|    | [>]({{ '/reference/comparison-operators.html#relational_operators' | relative_url }})   | nonassoc | 
|    | [>\=]({{ '/reference/comparison-operators.html#relational_operators' | relative_url }}) | nonassoc | 
|    | [<>]({{ '/reference/comparison-operators.html#relational_operators' | relative_url }})  | nonassoc | 
|    | [!\=]({{ '/reference/comparison-operators.html#relational_operators' | relative_url }}) | nonassoc | 
|    | [IS]({{ '/reference/comparison-operators.html#is' | relative_url }})           | nonassoc | 
|    | [BETWEEN]({{ '/reference/comparison-operators.html#between' | relative_url }}) | nonassoc | 
|    | [IN]({{ '/reference/comparison-operators.html#in' | relative_url }})           | nonassoc | 
|    | [LIKE]({{ '/reference/comparison-operators.html#like' | relative_url }})       | nonassoc | 
| 6  | [NOT]({{ '/reference/logic-operators.html#not' | relative_url }})     | Right-to-left | 
| 7  | [AND]({{ '/reference/logic-operators.html#and' | relative_url }})     | Left-to-right | 
| 8  | [OR]({{ '/reference/logic-operators.html#or' | relative_url }})       | Left-to-right | 
| 9  | [INTERSECT]({{ '/reference/set-operators.html#intersect' | relative_url }}) | Left-to-right | 
| 10 | [UNION]({{ '/reference/set-operators.html#union' | relative_url }})         | Left-to-right | 
|    | [EXCEPT]({{ '/reference/set-operators.html#except' | relative_url }})       | Left-to-right | 
| 11 | [:=]({{ '/reference/variable.html#substitution' | relative_url }})         | Right-to-left | 

