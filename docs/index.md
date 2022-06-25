---
layout: default
title: csvq - SQL-like query language for csv
---

## Overview

csvq is a command line tool to operate CSV files. 
You can read, update, delete CSV records with SQL-like query.

You can also execute multiple operations sequentially in managed transactions by passing a procedure or using the interactive shell.
In the multiple operations, you can use variables, cursors, temporary tables, and other features. 

## Latest Release

Version 1.17.4
: Released on Jun 25, 2022

  <a class="waves-effect waves-light btn" href="https://github.com/mithrandie/csvq/releases/tag/v1.17.4">
    <i class="material-icons left">file_download</i>download
  </a>

## Features

* CSV File Operation
  * [Select Query]({{ '/reference/select-query.html' | relative_url }})
  * [Insert Query]({{ '/reference/insert-query.html' | relative_url }})
  * [Update Query]({{ '/reference/update-query.html' | relative_url }})
  * [Replace Query]({{ '/reference/replace-query.html' | relative_url }})
  * [Delete Query]({{ '/reference/delete-query.html' | relative_url }})
  * [Create Table Query]({{ '/reference/create-table-query.html' | relative_url }})
  * [Alter Table Query]({{ '/reference/alter-table-query.html' | relative_url }})
* [Cursor]({{ '/reference/cursor.html' | relative_url }})
* [Temporary Table]({{ '/reference/temporary-table.html' | relative_url }})
* [Transaction Management]({{ '/reference/transaction.html' | relative_url }})
* Support loading data from Standard Input
* Support following file formats
  * CSV
  * TSV
  * LTSV
  * Fixed-Length Format
  * [JSON]({{ '/reference/json.html' | relative_url }})
  * JSON Lines
* Support following file encodings
  * UTF-8
  * UTF-16
  * Shift_JIS

  > JSON and JSON Lines formats support only UTF-8.

## Installation

[Installation - Reference Manual - csvq]({{ '/reference/install.html' | relative_url }})

## Command Usage

[Command Usage - Reference Manual - csvq]({{ '/reference/command.html' | relative_url }})

## Reference Manual

[Reference Manual - csvq]({{ '/reference.html' | relative_url }})

## Execute csvq statements in Go

[csvq-driver](https://github.com/mithrandie/csvq-driver)

## Example of cooperation with other applications

- [csvq emacs extension](https://github.com/mithrandie/csvq-emacs-extension)

## License

csvq is released under [the MIT License]({{ '/license.html' | relative_url }})