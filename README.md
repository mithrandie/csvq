# csvq

SQL-like query language for csv

[![Test](https://github.com/mithrandie/csvq/actions/workflows/test.yml/badge.svg)](https://github.com/mithrandie/csvq/actions/workflows/test.yml)
[![codecov](https://codecov.io/gh/mithrandie/csvq/branch/master/graph/badge.svg)](https://codecov.io/gh/mithrandie/csvq)
[![License: MIT](https://img.shields.io/badge/License-MIT-lightgrey.svg)](https://opensource.org/licenses/MIT)

Csvq is a command line tool to operate CSV files. 
You can read, update, delete CSV records with SQL-like query.

You can also execute multiple operations sequentially in managed transactions by passing a procedure or using the interactive shell.
In the multiple operations, you can use variables, cursors, temporary tables, and other features. 

## Latest Release
[![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/mithrandie/csvq?color=%2320b2aa&label=GitHub%20Release&sort=semver)](https://github.com/mithrandie/csvq/releases/latest)
[![GitHub tag (latest SemVer)](https://img.shields.io/github/v/tag/qittu/csvq-deb?color=%2320b2aa&label=Launchpad%20PPA)](https://launchpad.net/~mithrandie/+archive/ubuntu/csvq)

## Intended Use
Csvq is intended for one-time queries and routine processing described in source files on the amount of data that can be handled by spreadsheet applications.

It is not suitable for handling very large data since all data is kept on memory when queries are executed.
There is no indexing, calculation order optimization, etc., and the execution speed is not fast due to the inclusion of mechanisms for updating data and handling various other features.

However, it can be run with a single executable binary, and you don't have to worry about troublesome dependencies during installation.
You can not only write and run your own queries, but also share source files with co-workers on multiple platforms.

This tool may be useful for those who want to handle data easily and roughly, without having to think about troublesome matters.

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
  * [CSV](https://datatracker.ietf.org/doc/html/rfc4180)
  * TSV
  * [LTSV](http://ltsv.org)
  * Fixed-Length Format
  * [JSON](https://datatracker.ietf.org/doc/html/rfc8259)
  * [JSON Lines](https://jsonlines.org)
* Support following file encodings
  * UTF-8
  * UTF-16
  * Shift_JIS

  > JSON and JSON Lines formats support only UTF-8.

## Reference Manual

[Reference Manual - csvq](https://mithrandie.github.io/csvq/reference)

## Installation

### Install executable binary

1. Download an archive file from [release page](https://github.com/mithrandie/csvq/releases).
2. Extract the downloaded archive and add a binary file in it to your path.

### Build from source

#### Requirements

Go 1.18 or later (cf. [Getting Started - The Go Programming Language](https://golang.org/doc/install))

#### Build command

```$ go install github.com/mithrandie/csvq```

### Install using package manager

Installing using a package manager does not ensure that you always get the latest version, but it may make installation and updating easier.

#### Ubuntu

1. ```$ sudo add-apt-repository ppa:mithrandie/csvq```
2. ```$ sudo apt update```
3. ```$ sudo apt install csvq```

#### Arch Linux (unofficial)

Install the [csvq-git](https://aur.archlinux.org/packages/csvq-git) or [csvq-bin](https://aur.archlinux.org/packages/csvq-bin) from the [Arch User Repository](https://wiki.archlinux.org/title/Arch_User_Repository) (e.g. `yay -S csvq-git`)

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
