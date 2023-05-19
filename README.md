# users
Microservice for managing fiufit's users, admins and profiles

### Usage

#### With docker:
* Edit .example-env with your own secret credentials and rename it to .env
* `docker build -t fiufit-users`
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