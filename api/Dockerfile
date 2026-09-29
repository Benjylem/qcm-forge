FROM golang:1.26.8-alpine3.24 AS build
WORKDIR /src
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /api ./cmd/api

FROM alpine:3.24.2
COPY --from=build /api /api
USER nobody
CMD ["/api"]
