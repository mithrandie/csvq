---
layout: default
title: csvq - SQL like query language for csv
---

## Overview

csvq is a command line tool to operate CSV files. 
You can read, update, delete CSV records with SQL like query.
 
## Features

* CSV File Operation
  * Select Query
  * Insert Query
  * Update Query
  * Delete Query
  * Create Table Query
  * Alter Table Query
* Control Flow
  * Variable
  * Cursor
  * If Statement
  * While Statement
* Transaction Management
* Support loading from standard input as a CSV
* Support output a result of select query in JSON format 
* Support following file encodings
  * UTF-8
  * Shift-JIS
  * (Now supported encodings are very few. If you want to use other encodings, please submit an issue on [github](https://github.com/mithrandie/csvq).)

## Install

[Install - Reference Manual - csvq]({{ '/reference/install.html' | relative_url }})

## Reference Manual

[Reference Manual - csvq]({{ '/reference.html' | relative_url }})

## License

csvq is released under [the MIT License]({{ '/license.html' | relative_url }})