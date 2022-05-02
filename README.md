# Go-do
A todo cli application based on the [todo.txt](https://github.com/todotxt/todo.txt) format.

## Commands

| Command              | Args             | Flags                                                                                   | Description                                             |
|----------------------|------------------|-----------------------------------------------------------------------------------------|---------------------------------------------------------|
| create, c            | TODO             | -                                                                                       | Create and add a new todo based on the todo.txt format. |
| show                 | -                | --tag, -t tag type \| --value, -v value <br /> --complete, -c <br /> --incomplete, -inc | Show all todos.                                         |
| delete, d            | Todo description | -                                                                                       | Delete a todo by providing part of its description.     |
| delete-by-select, ds | -                | -                                                                                       | Lists all todos. The selected todo is deleted.          |
| help, h              | -                | -                                                                                       | Show the list of available commands.                    |

## Trello integration (in progress)
Generate API key and API token here: [Trello API](https://developer.atlassian.com/cloud/trello/guides/rest-api/api-introduction/).

Rename the file `.env.example` to `.env` and copy generated key and token values to that file.