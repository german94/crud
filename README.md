# crud

This is a simple CRUD REST API that has two entities: persons and pets. The following operations are allowed:

- Create a person (a POST to /persons)
- List all the persons (a GET to /persons)
- Find a person (a GET to /persons/:id)
- Create a pet (a POST to /pets)
- Find a pet (a GET to /pets/:id)
- Delete a pet (a DELETE to /pets/:id)
- Adopt a pet (a PUT to /adopt)


## Technical details

- The stack is composed of: sqlite, gorm and gin-gonic
- Only the controller layer is being tested. Only integration tests are provided, in a sense that we're working directly on the database and the requests. As this is an short example we don't have much "service" logic. In a real world case we need to have unit tests as well and mock storage.
- Only pets and adoption related requests are being tested. This is because I don't think it's worth for the exercise to also test persons. Again, in a real world case we need to test everything no matter how similar are the domains.

## Compiling and testing

In order to compile we can just type `make` and `make run`. Then the server will be accepting requests on port 8080. To run the tests we just need to type `make test`.

## Examples

Create person body:

`{
    "name": "german",
    "last_name": "pinzon",
    "age": 27
}`

Create pet body:

`{
    "kind": "Dog",
    "specie": "Labrador",
    "name": "Felipe",
    "age": 5
}`

Adopt pet body:

`{
    "owner_id": 1,
    "pet_id": 1
}`

## Deployment

This code was deployed to heroku as well and the URL of the app is: https://crud-challenge-rollee.herokuapp.com/persons Note that this is using the heroku free tier, so the app sleeps if it is not used for some time and it can be kind of slow to start up again. Once it's started it should work normally.

The deployment was added after the 2 days deadline because the Heroku build was failing, as it was using a Go version lower than 1.13, and [Gorm needs a greater version than that](https://stackoverflow.com/questions/64773490/error-reflect-valueofval-iszero-undefined-type-reflect-value-has-no-field-or) because of the `IsZero` `reflect.Value` method.