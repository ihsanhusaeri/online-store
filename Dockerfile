FROM golang:alpine

RUN apk update && apk add git ca-certificates

WORKDIR /src

COPY go.mod  .
RUN go mod download -x

COPY . .

RUN CGO_ENABLED=0 go build -a -o /src/app


FROM alpine

RUN apk update && apk add --no-cache tzdata

WORKDIR /app

COPY --from=0 /src/app .

EXPOSE 8080

ENTRYPOINT ["./app"]
