# Base stage for certificates and user setup
FROM --platform=$BUILDPLATFORM alpine:latest AS certs

# Install ca-certificates
RUN apk update && apk add --no-cache ca-certificates tzdata && rm -rf /var/cache/apk/*
RUN update-ca-certificates
RUN adduser -D -g '' appuser

# Final stage - scratch image for minimal size
FROM scratch

# Import certificates and user from certs stage
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=certs /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=certs /etc/passwd /etc/passwd

# Copy the pre-built binary (GoReleaser will handle this)
COPY aws-sso-config /usr/local/bin/aws-sso-config

# Use an unprivileged user
USER appuser

# Set working directory to /workspace for mounting host files
WORKDIR /workspace

# Set entrypoint
ENTRYPOINT ["/usr/local/bin/aws-sso-config"]
