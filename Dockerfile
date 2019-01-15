# Compiles the binary
FROM golang:1.11.4 as builder
WORKDIR /app
COPY . .
# We need to disable cgo in order for this binary to run in the other
# container:
# https://stackoverflow.com/questions/36279253/go-compiled-binary-wont-run-in-an-alpine-docker-container-on-ubuntu-host
RUN CGO_ENABLED=0 go build -o /main cmd/api/main.go

# Executes the binary
FROM scratch
EXPOSE 8080
COPY --from=builder /main /
ENTRYPOINT ["/main"]
