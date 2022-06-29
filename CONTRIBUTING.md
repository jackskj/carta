# Contributing

Thank you for contributing to carta.

Please follow the code of conduct when interacting with the carta project.

## Pull Request Process

The recommended PR process is as follows.

1. Please first discuss the change you wish to make via issue.
2. Upgrade with `go get -u github.com/jackskj/carta`
3. cd `$GOPATH/src/github.com/jackskj/carta`
4. Reflect your changes in example files as well as the README. 
5. Make sure that the build succedes with `go build` or `make build`
   and all tests pass with `make test`
6. Make sure to follow [Effective Go](https://golang.org/doc/effective_go.html)
   as well as [Go Code Review Comments](https://golang.org/wiki/CodeReviewComments)


## Testing
Go tests can be run via make by executing `make test`.

Prior to running unit tests, ensure postgres and mysql database containers are running by executing `make testdbs`.

Before running `make test` for the first time, you will need to create the database structure by running `make initdb`.

## Code of Conduct

Please do not make aggressive, harassing or condescending comments based on someone's
age, body size, disability, ethnicity, gender identity and expression, level of experience,
nationality, personal appearance, race, religion, or sexual identity and
orientation.

Please do not troll, or make insulting/derogatory comments  pr personal/political attacks
