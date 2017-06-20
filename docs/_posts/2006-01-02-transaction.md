---
layout: default
title: Transaction Management - Reference Manual - csvq
category: reference
---

# Transaction Management

A transaction is a single logical unit of work.
A transaction is a atomic unit, so you get the result of either all of the changes are performed or none of them is performed.

## Usage Flow

### Start Transaction

A transaction is started automatically when the statements execution is started, or after a commit or rollback statement is executed.


### Terminate Transaction

A transaction is terminated when a commit or rollback statement is executed.

If the statements execution is normally terminated, commit all of the changes automatically.

If some errors occurred in the statements execution, roll all of the changes back automatically.

### Commit Statement

A commit statement writes all of the changes to files.

```sql
COMMIT;
```

### Rollback Statement

A rollback statement discards all of the changes.

```sql
ROLLBACK;
```

