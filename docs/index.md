---
layout: default
title: csvq - SQL-like query language for csv
---

## Overview

csvq is a command line tool to operate CSV files. 
You can read, update, delete CSV records with SQL-like query.
 
## Features

* CSV File Operation
  * [Select Query]({{ '/reference/select-query.html' | relative_url }})
  * [Insert Query]({{ '/reference/insert-query.html' | relative_url }})
  * [Update Query]({{ '/reference/update-query.html' | relative_url }})
  * [Delete Query]({{ '/reference/delete-query.html' | relative_url }})
  * [Create Table Query]({{ '/reference/create-table-query.html' | relative_url }})
  * [Alter Table Query]({{ '/reference/alter-table-query.html' | relative_url }})
* [Cursor]({{ '/reference/cursor.html' | relative_url }})
* [Temporary Table]({{ '/reference/temporary-table.html' | relative_url }})
* [Transaction Management]({{ '/reference/transaction.html' | relative_url }})
* Support loading from standard input as a CSV
* Support output a result of select query in JSON format 
* Support following file encodings
  * UTF-8
  * Shift-JIS
  * (Now supported encodings are very few. If you want to use other encodings, please submit an issue on [github](https://github.com/mithrandie/csvq).)

## Install

[Install - Reference Manual - csvq]({{ '/reference/install.html' | relative_url }})

## Command Usage

[Command Usage - Reference Manual - csvq]({{ '/reference/command.html' | relative_url }})

## Reference Manual

[Reference Manual - csvq]({{ '/reference.html' | relative_url }})

## License

csvq is released under [the MIT License]({{ '/license.html' | relative_url }})