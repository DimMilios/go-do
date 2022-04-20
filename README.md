# Go-do
A todo cli application based on the [todo.txt](https://github.com/todotxt/todo.txt) format.

## Features
### General
- [] Display all todos command (`--all` or `-a`)
- [] Create todo command (`--create` or `-c`)
- [] Export todos to `.txt` file specifying a file name (`--export` or `-e` followed by `--name` or `-n` for the file name)
- [] Help command (lists available commands + a description for the todo.txt format)

### Todos
- [] Add a todo
- [] Update a todo
- [] Delete a todo
- [] List all todos
- [] Sort todos
  - [] By priority
  - [] By completion status
  - [] By tags
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
  (complete SPACE)?
  (priority SPACE)?
  (completion_date SPACE)?
  description{1}
  ;
description
  : 
  STRING{1}
  project_tag*
  context_tag*
  key_value_tag*
  ;
complete: 'x';
priority: '(' [A-Z]{1} ')';
completion_date: DATE_FORMAT?;
project_tag: SPACE '+' STRING;
context_tag: SPACE '@' STRING;
key_value_tag: SPACE \S+:\S+;
SPACE: ' ';
DATE_FORMAT: /\d{4}-\d{2}-\d{2}/; // YYYY-MM-DD
```