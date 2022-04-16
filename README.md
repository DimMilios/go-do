# Go-do
A todo cli application based on the [todo.txt](https://github.com/todotxt/todo.txt) format.

## Features
### General
- [] Display all todos command (`--all` or `-a`)
- [] Create todo command (`--create` or `-c`)
- [] Help command (lists available commands + a description for the todo.txt format)

### Todos
- [] Add a todo
- [] Update a todo
- [] Delete a todo
- [] List all todos
- [] Filter todos
  - [] Complete or incomplete
  - [] By completion/creation date
  - [] By tags
    - [] By project tag
    - [] By context tag
    - [] By key value tag

## Grammar translation
This is a representation of the `todo.txt` format in `EBNF`.
```
todo
  : 
  (completed SPACE)?
  (priority SPACE)?
  (dates SPACE)?
  description{1}
  ;
description
  : 
  STRING{1}
  project_tag*
  context_tag*
  key_value_tag*
  ;
completed: 'x';
priority: '(' [A-Z]{1} ')';
dates: DATE_FORMAT SPACE DATE_FORMAT; // completion_date creation_date
project_tag: SPACE '+' STRING;
context_tag: SPACE '@' STRING;
key_value_tag: SPACE \S+:\S+;
SPACE: ' ';
DATE_FORMAT: /\d{4}-\d{2}-\d{2}/;
```