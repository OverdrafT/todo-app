# todo-app
Work in progress

Production ready service with Clean Architecture approach

## Improvements
- Add metadata like version, build, commit

## Questions
- Now Create operation returns Item's ID if created. Should it return the whole Item?
What about Update operation?
- Database layer should check if item exists before updating.
Which layer should check if id is not empty string? (transport?)