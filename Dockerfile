FROM golang:1.20-buster AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN go build -o /api

## Deploy
FROM cgr.dev/chainguard/static

WORKDIR /
USER nonroot:nonroot

COPY --from=build /api /api

EXPOSE 4000
ENTRYPOINT ["/api"]
