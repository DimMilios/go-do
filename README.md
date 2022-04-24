# Go-do
A todo cli application based on the [todo.txt](https://github.com/todotxt/todo.txt) format.

## Features
### General
- [x] Display all todos command (`show`)
    - [] Pass a filter (e.g. show todos of context A or project B etc.) 
- [x] Create todo command (`create` or `c`)
- [] Export todos to `.txt` file specifying a file name (`export` or `e` followed by `--name` or `-n` for the file name)
- [] Help command (lists available commands + a description for the todo.txt format)

### Todos
- [x] Add a todo
- [] Update a todo
- [x] Delete a todo
- [x] List all todos
- [] Sort todos
  - [] By priority
  - [] By completion status
  - [] By tags
- [] Filter todos
  - [x] Complete or incomplete
  - [] By completion/creation date
  - [x] By tags
    - [x] By project tag
    - [x] By context tag
    - [x] By key value tag

## Project todo
- [] Completion date can't come before creation dates
- [] Marking a todo as done automatically should add a completion date if one wasn 't provided
- [] Providing a completion date should also mark the todo as done

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