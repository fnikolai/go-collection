# syntax = docker/dockerfile:1-experimental

# In the first stage, create a dummy container with all the dependencies needed to compile timon
FROM  golang:1.15.2-alpine AS build

# Build parameters
ARG SrcDir="."
ARG TARGETOS
ARG TARGETARCH

# Execution parameters
ENV CGO_ENABLED=0

WORKDIR /src

# Feature1: separate download of our dependencies from our build
# Depending on the builder, git may need to be install
COPY ./go.* ./
RUN go mod download

COPY ${SrcDir} .

RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o /bin/exporter ./cmd

# Feature2: leverage Go's compiler cache using BuildKit's Dockerfile
#RUN --mount=type=cache,target=/root/.cache/go-build \
#    GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o /bin/exporter ./cmd


# In the second step, copy the compiled binary to another tiny container
FROM alpine AS bin
COPY --from=build /bin/exporter /

#EXPOSE 6666/tcp

ENTRYPOINT ["./exporter"]
