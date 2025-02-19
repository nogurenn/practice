# We'll use alpine as the base image for build, because
# the alpine runtime image cannot detect the binary as an
# executable. This is because the binary built on default
# golang images is linked against the debian glibc, which
# is incompatible with alpine's muslc.
# 
# A golang binary from a debian image (dynamically linked)
# won't readily run on an alpine image, but a binary built
# on an alpine image will run on a debian image. Setting
# CGO_ENABLED=0 always felt like a hack to me, so finally
# I found a better solution.

FROM golang:1.24-alpine3.21 AS build

ARG TARGET_APP="producer"

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/reference/dockerfile/#copy
COPY cmd/ ./cmd/
COPY internal/ ./internal/

RUN GOOS=linux go build -o /${TARGET_APP} ./cmd/${TARGET_APP}

FROM alpine:3.21 AS run

ARG TARGET_APP="producer"

# Copy the application executable from the build image
COPY --from=build /${TARGET_APP} /${TARGET_APP}

EXPOSE 3000

# Run
ENV CMD_RUN="./${TARGET_APP}"
CMD $CMD_RUN
