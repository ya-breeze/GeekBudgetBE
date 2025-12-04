FROM golang:alpine AS builder
RUN apk add --no-cache gcc musl-dev sqlite-dev
WORKDIR /app
COPY backend/go.mod backend/go.sum ./
RUN go mod download
COPY backend/ ./
ENV CGO_ENABLED=1
RUN cd cmd && go build -o geekbudget .

FROM alpine:latest
RUN apk add --no-cache sqlite-libs ca-certificates tzdata
# Set the timezone (optional)
ENV TZ=Europe/Prague
WORKDIR /root/
COPY --from=builder /app/cmd/geekbudget .
COPY --from=builder /app/webapp ./webapp
EXPOSE 8080
CMD ["./geekbudget", "server"]
