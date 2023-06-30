<p align="center">
  <img alt="App" src="https://github.com/fiufit/app/assets/86434696/2dc48884-cd7c-4aca-ad99-e9adf2f4410d" height="200" />
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
