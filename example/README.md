# Basic gqlgen + @validate Example

This example shows how to use the `@validate` directive and the validator-aware
model generator in a small gqlgen project. In addition to the required `rule`
argument, you can provide an optional `message` to customise the error returned
to the client—for instance the `age` field in this demo uses
`@validate(rule: "omitempty,gte=18", message: "Age must be 18+ or left blank")`
so the validation error mirrors product copy instead of the default formatter.

## Generate code

```bash
go run ./cmd/gqlgen
```

## Run the server

```bash
go run server.go
```

Once the server is running open http://localhost:8080 to access the GraphQL
playground. A sample mutation:

```graphql
mutation RegisterUser {
  registerUser(
    input: {
      email: "john.doe@example.com"
      password: "secret123"
      confirmPassword: "secret123"
      age: 18
      termsAndConditions: [1]
      questionnaireAnswers: [
        { questionId: 1, answerId: 1 }
        { questionId: 2, answerText: "Yes" }
      ]
    }
  ) {
    id
    email
    password
    age
    termsAndConditions
    questionnaireAnswers {
      questionId
      answerId
      answerText
    }
  }
}
```

If you submit an invalid payload (for example an empty email, a short password
or an underage `age` value) the playground displays the validation errors
produced by the runtime directive using the customised message text where
provided.
