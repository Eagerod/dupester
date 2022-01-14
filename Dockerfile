FROM golang:1.16 AS build

WORKDIR /usr/src/app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN make


FROM debian:10 AS app

COPY --from=build /usr/src/app/build/* /usr/bin/dupester

ENTRYPOINT ["/usr/bin/dupester"]
