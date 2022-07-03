---
layout: default
title: System Defined Constant - Reference Manual - csvq
category: reference
---

# System Defined Constant

A system-defined constant is specified by two words: the category and the constant name.
The two words are separated by "::".

e.g. `MATH::PI`

| Category | Name             | Type    | Value                                                             |
|:---------|:-----------------|:--------|:------------------------------------------------------------------|
| MATH     | E                | float   | 2.71828182845904523536028747135266249775724709369995957496696763  |
|          | PI               | float   | 3.14159265358979323846264338327950288419716939937510582097494459  |
|          | PHI              | float   | 1.61803398874989484820458683436563811772030917980576286213544862  |
|          | SQRT2            | float   | 1.41421356237309504880168872420969807856967187537694807317667974  |
|          | SQRTE            | float   | 1.64872127070012814684865078781416357165377610071014801157507931  |
|          | SQRTPI           | float   | 1.77245385090551602729816748334114518279754945612238712821380779  |
|          | SQRTPHI          | float   | 1.27201964951406896425242246173749149171560804184009624861664038  |
|          | LN2              | float   | 0.693147180559945309417232121458176568075500134360255254120680009 |
|          | LOG2E            | float   | 1 / Ln2                                                           |
|          | LN10             | float   | 2.30258509299404568401799145468436420760110148862877297603332790  |
|          | LOG10E           | float   | 1 / Ln10                                                          |
| FLOAT    | MAX              | float   | 0x1p1023 * (1 + (1 - 0x1p-52))                                    |
|          | SMALLEST_NONZERO | float   | 0x1p-1022 * 0x1p-52                                               |
| INTEGER  | MAX              | integer | 1<<63 - 1                                                         |
|          | MIN              | integer | -1 << 63                                                          |

