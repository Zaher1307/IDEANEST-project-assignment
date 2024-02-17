FROM golang:1.21 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /IDEANEST cmd/*.go

# FROM build-stage AS run-test-stage
# RUN go test -v ./...

FROM ubuntu:latest AS build-release-stage

WORKDIR /

COPY --from=build-stage /IDEANEST /IDEANEST

EXPOSE 8080

ENTRYPOINT ["/IDEANEST"]
