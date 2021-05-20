# todo-app
Production ready service with Clean Architecture approach

## Questions
- Now Create operation returns Item's ID if created. Should it return the whole Item?
What about Update operation?
- Database layer should check if item exists before updating.
Which layer should check if id is not empty string? (transport?)