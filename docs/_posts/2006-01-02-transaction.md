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

