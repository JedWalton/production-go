############################
# STEP 1 build executable binary
############################
FROM golang@sha256:dd8888bb7f1b0b05e1e625aa29483f50f38a9b64073a4db00b04076cec52b71c as builder
# Install git + SSL ca certificates.
RUN apk update && apk add --no-cache git ca-certificates tzdata && update-ca-certificates
# Create appuser
ENV USER=appuser
ENV UID=10001
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    "${USER}"
WORKDIR /go/src/mypackage/myapp/
COPY . .
RUN go mod download
RUN go mod verify
RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /go/bin/production-go

############################
# STEP 2 build a small image
############################
FROM alpine
# Import from builder.
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
COPY --from=builder /go/bin/production-go /go/bin/production-go

# Add entrypoint script
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

# Use an unprivileged user
USER appuser:appuser

# Set entrypoint to the script
ENTRYPOINT ["/entrypoint.sh"]
