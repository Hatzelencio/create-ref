FROM golang:1.14-alpine
LABEL maintainer='Hatzel Renteria'

WORKDIR /action

CMD ["go", "run", "/action/main.go"]