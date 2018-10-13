## Release v1.3.0 (2018-10-13)

- Support ANSI escape sequence.
- Enhance support for JSON.
  - Load data from a JSON file with the JSON_TABLE expression in From Clause.
  - Load data from a JSON data from standard input with the –json-query option.
  - Export a result of a select query in JSON format with the –format {JSON | JSONH | JSONA} option.
  - Load a value from a JSON data using functions.
    - JSON_VALUE
    - JSON_OBJECT
    - JSON_AGG (Aggregate Function)
    - JSON_AGG (Analytic Function)
  - Load a row value from a JSON data using the JSON_ROW expression.

## Release v1.2.0 (2018-09-25)

- Support for Go 1.11 Modules.

## Release v1.1.1 (2018-04-05)

- Implement string functions.
  - INSTR
  - LIST_ELEM

## Release v1.1.0 (2018-02-28)

- Support for Go 1.10

## Release v1.0.2 (2017-12-08)

- Fix some bugs of operetor precedence.

## Release v1.0.1 (2017-09-26)

- Implement DISPOSE FUNCTION statement.
- Implement windowing clause in analytic function.

## Release v1.0.0 (2017-09-19)

The first general release. 