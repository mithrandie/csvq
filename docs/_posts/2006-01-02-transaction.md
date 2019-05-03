---
layout: default
title: Transaction Management - Reference Manual - csvq
category: reference
---

# Transaction Management

A transaction is a single logical unit of work.
A transaction is a atomic unit, so you get the result of either all of the changes are performed or none of them is performed.

* [Usage Flow in a Procedure](#usage_flow_in_prodecure)
* [Usage Flow in the Interactive Shell](#usage_flow_in_shell)
* [File Locking](#file_locking)
* [Commit Statement](#commit)
* [Rollback Statement](#rollback)

## Usage Flow in a Procedure
{: #usage_flow_in_prodecure}

### Start Transaction

A transaction is started automatically when the procedure execution is started, or after a commit or rollback statement is executed.

### Terminate Transaction

A transaction is terminated when a commit or rollback statement is executed.

When the procedure is normally terminated, then commit all of the changes automatically.

When some errors occurred in the procedure, then roll all of the changes back automatically.

When the procedure is exited by [EXIT statement]({{ '/reference/control-flow.html#exit' | relative_url }}), then roll all of the changes back automatically.

## Usage Flow in the Interactive Shell
{: #usage_flow_in_shell}

### Start Transaction

A transaction is started automatically when the interactive shell is launched, or after a commit or rollback statement is executed.

### Terminate Transaction

A transaction is terminated when a commit or rollback statement is executed.

When the interactive shell is terminated, then roll all of the changes back automatically.


## File Locking
{: #file_locking}

In a transaction, created files and updated files are locked by using lock files, so these files are protected from other csvq processes.

This locking does not guarantee that these files are protected from other applications.
System-provided file locking to protect them from other applications are used only on the systems supported by the package [github.com/mithrandie/go-file](https://github.com/mithrandie/go-file).

SELECT queries use shared locks. INSERT, UPDATE, DELETE, CREATE and ALTER TABLE queries use exclusive locks to update files.
Shared locks are unlocked immediately after reading, and exclusive locks remain until the termination of the transaction.

Once you load files, that data is cached until the termination of the transaction, so in a transaction, that data is basically unaffected by the other transactions.
However, as an exception, when trying to update a file that has been loaded by a SELECT query, the file will be reloaded.
In that case, there is a probability the data is changed in tha same transaction.
You can use [FOR UPDATE]({{ '/reference/select-query.html' | relative_url }}) keywords in SELECT queries to use exclusive locks and prevent the probability. 

### Recover file locking

Program panics and unterminated transactions remain lock files.
In that case, you must manually remove following hidden files created by csvq.

- ._FILE_NAME_.[0-9a-zA-Z]{12}.rlock 
- ._FILE_NAME_.lock 
- ._FILE_NAME_.temp


## Commit Statement
{: #commit}

A commit statement writes all of the changes to files.

```sql
COMMIT;
```

## Rollback Statement
{: #rollback}

A rollback statement discards all of the changes.

```sql
ROLLBACK;
```

