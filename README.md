<p align="center">
  <img alt="App" src="https://firebasestorage.googleapis.com/v0/b/fiufit.appspot.com/o/fiufit-logo.png?alt=media&token=39f3ae3f-34d1-4fb3-96ca-8707adf2bc37" height="200" />
</p>

# users
Microservice for managing fiufit's users, admins and profiles

[![codecov](https://codecov.io/github/fiufit/users/branch/main/graph/badge.svg?token=PDT69DRER8)](https://codecov.io/github/fiufit/users)
[![test](https://github.com/fiufit/users/actions/workflows/test.yml/badge.svg)](https://github.com/fiufit/users/actions/workflows/test.yml)
[![Lint Go Code](https://github.com/fiufit/users/actions/workflows/lint.yml/badge.svg)](https://github.com/fiufit/users/actions/workflows/lint.yml)
[![Fly Deploy](https://github.com/fiufit/users/actions/workflows/fly.yml/badge.svg)](https://github.com/fiufit/users/actions/workflows/fly.yml)

### Usage

#### With docker:
* Edit .example-env with your own secret credentials and rename it to .env
* `docker build -t fiufit-users .`
* `docker run -p PORT:PORT --env-file=.env fiufit-users`

#### Natively: 
* `go mod tidy`
* set your environvent variables to imitate the .env-example
* `go run main.go` or `go build` and run the executable


#### Running tests:
* `go test ./...`


### Links
* Fly.io deploy dashboard: https://fly.io/apps/fiufit-users
* Swagger docs: https://fiufit-users.fly.dev/v1/docs/index.html
* Coverage report: https://app.codecov.io/github/fiufit/users
