---
layout: default
title: Runtime Information - Reference Manual - csvq
category: reference
---

# Runtime Information

Following values can be referred as runtime information.

| name | type | description |
| :- | :- | :- |
| @#UNCOMMITTED        | boolean | Whether there are tables or views that have not been committed. |
| @#CREATED            | integer | Number of uncommitted tables after creation |
| @#UPDATED            | integer | Number of uncommitted tables after update |
| @#UPDATED_VIEWS      | integer | Number of uncommitted views after update |
| @#LOADED_TABLES      | integer | Number of loaded tables |
| @#WORKING_DIRECTORY  | string  | Current working directory |
| @#VERSION            | string  | Version of csvq |

