# -------- builder stage -------- #
FROM golang:alpine AS builder

ARG GOOS_VAL=linux
ARG GOARCH_VAL=amd64
ARG CGO_ENABLED_VAL=0

WORKDIR $GOPATH/src/terrascan

# download go dependencies
COPY go.mod go.sum ./
RUN go mod download

# copy terrascan source
COPY . .

# build binary
RUN CGO_ENABLED=${CGO_ENABLED_VAL} GOOS=${GOOS_VAL} GOARCH=${GOARCH_VAL} go build -v -ldflags "-w -s" -o /go/bin/terrascan ./cmd/terrascan


# -------- prod stage -------- #
FROM alpine:3.12.0

# create non root user
RUN addgroup --gid 101 terrascan && \
    adduser -S --uid 101 --ingroup terrascan terrascan && \
    apk add --no-cache git

# run as non root user
USER terrascan

# copy terrascan binary from build
COPY --from=builder /go/bin/terrascan /go/bin/terrascan

EXPOSE 9010

ENTRYPOINT ["/go/bin/terrascan"]
CMD ["server", "--log-type", "json"]
