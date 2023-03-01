# syntax=docker/dockerfile:1

FROM golang:1.20-alpine AS development

VOLUME [ "/project" ]

WORKDIR /project

RUN go install github.com/cosmtrek/air@latest

ENTRYPOINT [ "air" ]

FROM golang:1.20-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

# Copy project files
# NOTE: For more efficient builds, copy only
# required files here.
# For example, for a Go project, you may only
# need go.mod, go.sum and the actual *.go files
# that make up your project
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o go-scaffold

FROM scratch AS production

COPY --from=builder /app/go-scaffold .

ENTRYPOINT [ "/go-scaffold" ]
