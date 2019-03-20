# Accept the Go version for the image to be set as a build argument.
# Default to Go 1.11
ARG GO_VERSION=1.11
ARG GOARCH=amd64

# First stage: build the executable.
FROM golang:${GO_VERSION}-alpine AS builder

# Install the Certificate-Authority certificates for the app to be able to make
# calls to HTTPS endpoints.
RUN apk add --no-cache ca-certificates git gcc libc-dev

# Set the environment variables for the go command:
# * GOFLAGS=-mod=vendor to force `go build` to look into the `/vendor` folder.
ENV GOFLAGS=-mod=vendor

# Set the working directory outside $GOPATH to enable the support for modules.
WORKDIR /src

# Import the code from the context.
COPY ./ ./

# Build the executable to `/app`. Mark the build as statically linked.
RUN mkdir -p ${GOPATH}/src/github.com/kubernetes-sigs \
    && ln -sf `pwd` ${GOPATH}/src/github.com/kubernetes-sigs/dashboard-metrics-scraper \
    && go build \
    -installsuffix 'static' \
    -ldflags '-extldflags "-static"' \
    -o /metrics-sidecar github.com/kubernetes-sigs/dashboard-metrics-scraper

# Final stage: the running container.
FROM scratch AS final

# Import the Certificate-Authority certificates for enabling HTTPS.
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Import the compiled executable from the second stage.
COPY --from=builder /metrics-sidecar /metrics-sidecar

# We need a tmp folder too
COPY --from=builder /tmp /tmp

# Declare the port on which the webserver will be exposed.
EXPOSE 8080

# Run the compiled binary.
ENTRYPOINT ["/metrics-sidecar"]