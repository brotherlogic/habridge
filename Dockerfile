# syntax=docker/dockerfile:1

FROM golang:1.24.0 AS build

WORKDIR $GOPATH/src/github.com/brotherlogic/habridge

COPY go.mod ./
COPY go.sum ./

RUN mkdir proto
COPY proto/*.go ./proto/

RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 go build -o /habridge

##
## Deploy
##
FROM gcr.io/distroless/base-debian11

WORKDIR /

COPY --from=build /habridge /habridge

EXPOSE 8080
EXPOSE 8081

USER nonroot:nonroot

ENTRYPOINT ["/habridge"]