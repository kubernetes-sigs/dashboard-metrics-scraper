# Accept the Go version for the image to be set as a build argument.
ARG GO_VERSION=1.18

# First stage: build the executable.
FROM --platform=$BUILDPLATFORM golang:${GO_VERSION}-stretch AS builder

# Expose global args
ARG TARGETARCH
ARG TARGETOS

# What arch is it?
ARG GOARCH=$TARGETARCH
ARG GOOS=$TARGETOS

# Install the Certificate-Authority certificates for the app to be able to make
# calls to HTTPS endpoints.
RUN apt-get update && \
    apt-get install -y ca-certificates git gcc libc-dev libncurses5-dev sqlite3

# Set the working directory outside $GOPATH to enable the support for modules.
WORKDIR /src

# Import the code from the context.
COPY ./ ./

# Build the executable to `/app`. Mark the build as statically linked.
RUN echo "Building for $GOARCH" \
    && mkdir -p ${GOPATH}/src/github.com/kubernetes-sigs \
    && ln -sf `pwd` ${GOPATH}/src/github.com/kubernetes-sigs/dashboard-metrics-scraper \
    && GOARCH=${GOARCH} hack/build.sh 

# Create a nonroot user for final image
RUN useradd -u 10001 nonroot 

# Final stage: the running container.
FROM scratch AS final

# Import the compiled executable from the second stage.
COPY --from=builder /metrics-sidecar /metrics-sidecar

# Copy nonroot user
COPY --from=builder /etc/passwd /etc/passwd

# Declare the port on which the webserver will be exposed.
EXPOSE 8080
USER nonroot

# Run the compiled binary.
ENTRYPOINT ["/metrics-sidecar"]
