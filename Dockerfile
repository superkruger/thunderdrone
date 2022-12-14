# backend build stage
FROM golang:alpine as backend-builder
RUN apk --no-cache add ca-certificates
ENV GO111MODULE=on
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . /app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build cmd/thunderdrone/thunderdrone.go

# final stage
FROM alpine
COPY --from=backend-builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=backend-builder ./app/wait-for-it.sh /app/
COPY --from=backend-builder /app/thunderdrone /app/
COPY --from=backend-builder /app/database/migrations /app/database/migrations
RUN apk add --no-cache bash busybox-extras

ENV GIN_MODE=release
WORKDIR /app

ENTRYPOINT ["./wait-for-it.sh", "thunderdrone-db:5432", "--", "./thunderdrone"]