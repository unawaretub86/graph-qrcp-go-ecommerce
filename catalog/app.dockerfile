FROM golang:1.23-alpine AS build
RUN apk --no-cache add gcc g++ make ca-certificates
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY catalog catalog
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/app ./catalog/cmd

FROM alpine:3.18
WORKDIR /app
COPY --from=build /app/app .
RUN mkdir -p configs
EXPOSE 8080
CMD ["./app"]
