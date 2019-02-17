---
layout: default
title: Change Log - csvq
---

# Change Log

## Release v1.8.3 (2019-02-17)

- Fix the following bugs.
  - RETURN statement does not return a value in IF and WHILE statements.
  - NOW Function returns different time from the specification in user-defined functions.

## Release v1.8.2 (2019-02-13)

- Fix the following bug.
  - Panic occurs when an empty environment variable is passed. (Github Pull Request #4)

## Release v1.8.1 (2019-01-05)

- Improve completer.
- Fix some bugs of completer.

## Release v1.8.0 (2018-12-31)

- Support LTSV Format.

## Release v1.7.3 (2018-12-26)

- Improve logics of parallel routine.

## Release v1.7.2 (2018-12-24)

- Implement Syntax subcommand.

## Release v1.7.1 (2018-12-15)

- Fix the following bugs.
  - TableObject does not accept identifier as an argument.

## Release v1.7.0 (2018-12-14)

- Enhance the interactive shell.
  - Completion (default: true)
  - Kill Whole Line (default: false)
  - Vi-mode (default: false)

## Release v1.6.7 (2018-12-01)

- Fork github.com/chzyer/readline and change dependency to github.com/mithrandie/readline-csvq to use the latest update that is not versioned. 

## Release v1.6.6 (2018-11-27)

- Fix fatal error of variable substitution in multithreading.

## Release v1.6.5 (2018-11-25)

- Implement Flag Related Statements.
  - ADD FLAG ELEMENT
  - REMOVE FLAG ELEMENT
- Fix a bug of datetime formats configuration.

## Release v1.6.4 (2018-11-24)

- Implement Identical Operator ("==").

## Release v1.6.3 (2018-11-24)

- Fix the following bug.
  - Color output is on by default. (Github Issue #3)

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