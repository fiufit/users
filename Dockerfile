FROM golang:1.20-alpine

##BUILD
WORKDIR /app
COPY . ./
RUN go mod download
RUN go build -o /fiufit-users

##SECRETS

##DEPLOY
EXPOSE 8888
CMD [ "/fiufit-users" ]

