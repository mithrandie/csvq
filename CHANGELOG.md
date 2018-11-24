## Release v1.6.2 (2018-11-24)

- Implement run-external-command statement.
- Add value expressions.
  - Runtime Information.
- Add built-in commands.
  - ECHO
  - CHDIR
  - PWD
  - RELOAD CONFIG
- Add configuration items to csvq_env.json.
  - interactive_shell.prompt
  - interactive_shell.continuous_prompt
- Bug Fixes

## Release v1.6.1 (2018-11-18)

- Fix a bug of colorize in JSON pretty print.
- Add --json-escape option.

## Release v1.6.0 (2018-11-17)

- Support Environment Variables.
- Support Configuration Files and Pre-Load Statements.
- Add command options.
  - enclose-all
  - east-asian-encoding
  - count-diacritical-sign
  - count-format-code

## Release v1.5.4 (2018-11-08)

- Fix a bug of string format.

## Release v1.5.3 (2018-11-08)

- Implement EXECUTE statement.
- Implement NUMBER_FORMAT function.
- Make FORMAT function to determine the number of digits automatically when precision is not specified.

## Release v1.5.2 (2018-11-05)

- Fix a bug in calculation of Shift-JIS byte length.

## Release v1.5.1 (2018-11-04)

- Fix a bug of interactive shell that hide query results when the --out option is specified.

## Release v1.5.0 (2018-11-04)

- Support Fixed-Length Format.
- Implement WIDTH function.
- Support operate with byte length in Shift-JIS encoding in the following string functions.
  - BYTE_LEN
  - LPAN
  - RPAD
- Implement ALTER TABLE SET ATTRIBUTE statement.


## Release v1.4.3 (2018-10-19)

- Fix return code when on-usage-error occurred.
- Add flags for write out.

## Release v1.4.2 (2018-10-17)

- Fix output format problems on the specifications.
  - Conversion to GigHub Flavored Markdown Format
    - Ternary -> bool or empty string
    - Null -> empty string
  - Conversion to Org-mode Table Format
    - Ternary -> bool or empty string
    - Null -> empty string

## Release v1.4.1 (2018-10-17)

- Fix a bug of datetime conversion.

## Release v1.4.0 (2018-10-16)

- Add output formats.
  - Text Table for GitHub Flavored Markdown
  - Text Table for Emacs Org-mode

## Release v1.3.1 (2018-10-13)

- Fix a bug of output ANSI escape sequence.

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