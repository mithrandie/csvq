# csvq

SQL-like query language for csv

[![Build Status](https://travis-ci.org/mithrandie/csvq.svg?branch=master)](https://travis-ci.org/mithrandie/csvq)
[![codecov](https://codecov.io/gh/mithrandie/csvq/branch/master/graph/badge.svg)](https://codecov.io/gh/mithrandie/csvq)

csvq is a command line tool to operate CSV files. 
You can read, update, delete CSV records with SQL-like query.

You can also execute multiple operations sequentially in managed transactions by passing a procedure or using the interactive shell.
In the multiple operations, you can use variables, cursors, temporary tables, and other features. 


## Features

* CSV File Operation
  * Select Query
  * Insert Query
  * Update Query
  * Replace Query
  * Delete Query
  * Create Table Query
  * Alter Table Query
* Cursor
* Temporary Table
* Transaction Management
* Support loading data from Standard Input
* Support following file formats
  * CSV
  * TSV
  * LTSV
  * Fixed-Length Format
  * JSON
  * JSON Lines
* Support following file encodings
  * UTF-8
  * UTF-16
  * Shift_JIS

  > JSON Format supports only UTF-8.

## Reference Manual

[Reference Manual - csvq](https://mithrandie.github.io/csvq/reference)

## Installation

### Install executable binary

1. Download an archive file from [release page](https://github.com/mithrandie/csvq/releases).
2. Extract the downloaded archive and add a binary file in it to your path.

### Build from source

#### Requirements

Go 1.16 or later (cf. [Getting Started - The Go Programming Language](https://golang.org/doc/install))

#### Build command

```$ env GO111MODULE=on go get github.com/mithrandie/csvq```

### Install using package manager

Installing using a package manager does not ensure that you always get the latest version, but it may make installation and updating easier.

#### Ubuntu

1. ```$ sudo add-apt-repository ppa:mithrandie/csvq```
2. ```$ sudo apt update```
3. ```$ sudo apt install csvq```

#### macOS (unofficial)

1. Install homebrew (cf. [The missing package manager for macOS (or Linux) — Homebrew](https://brew.sh))
2. ```$ brew install csvq```

## Usage

```shell
# Simple query
csvq 'select id, name from `user.csv`'
csvq 'select id, name from user'

# Specify data delimiter as tab character
csvq -d '\t' 'select count(*) from `user.csv`'

# Load no-header-csv
csvq --no-header 'select c1, c2 from user'

# Load from redirection or pipe
csvq 'select * from stdin' < user.csv
cat user.csv | csvq 'select *'

# Load from Fixed-Length Format
cat /var/log/syslog | csvq -n -i fixed -m '[15, 24, 124]' 'select *'

# Split lines with spaces automatically
ps | csvq -i fixed -m spaces 'select * from stdin'

# Output in JSON format
csvq -f json 'select integer(id) as id, name from user'

# Output to a file
csvq -o new_user.csv 'select id, name from user'

# Load statements from file
$ cat statements.sql
VAR @id := 0;
SELECT @id := @id + 1 AS id,
       name
  FROM user;

$ csvq -s statements.sql

# Execute statements in the interactive shell
$ csvq
csvq > UPDATE users SET name = 'Mildred' WHERE id = 2;
1 record updated on "/home/mithrandie/docs/csv/users.csv".
csvq > COMMIT;
Commit: file "/home/mithrandie/docs/csv/users.csv" is updated.
csvq > EXIT;

# Show help
csvq -h
```

More details >> [https://mithrandie.github.io/csvq](https://mithrandie.github.io/csvq)

## Execute csvq statements in Go

[csvq-driver](https://github.com/mithrandie/csvq-driver)

## Example of cooperation with other applications

- [csvq emacs extension](https://github.com/mithrandie/csvq-emacs-extension)
