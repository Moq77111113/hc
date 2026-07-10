# Build a static hc and ship it alone on scratch, so any image can grab it with:
#   COPY --from=ghcr.io/moq77111113/hc /hc /hc
FROM golang:1.26 AS build
WORKDIR /src
COPY go.mod ./
COPY *.go ./
COPY internal/ ./internal/
ARG TARGETOS TARGETARCH
ARG HC_TAGS=""
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build -tags "${HC_TAGS}" -trimpath -ldflags="-s -w" -o /hc .

FROM scratch
COPY --from=build /hc /hc
ENTRYPOINT ["/hc"]
