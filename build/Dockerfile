# -------- builder stage -------- #
FROM golang:alpine AS builder

ARG CGO_ENABLED_VAL=0

WORKDIR $GOPATH/src/terrascan

# download go dependencies
COPY go.mod go.sum ./
RUN go mod download
RUN apk update && apk add --no-cache --update build-base git

# copy terrascan source
COPY . .

# build binary
RUN CGO_ENABLED=${CGO_ENABLED_VAL} go build -v -ldflags "-w -s -X google.golang.org/protobuf/reflect/protoregistry.conflictPolicy=ignore" -o /go/bin/terrascan ./cmd/terrascan


# -------- prod stage -------- #
FROM alpine:3

# create non root user
RUN addgroup --gid 101 terrascan && \
    adduser -S --uid 101 --ingroup terrascan terrascan && \
    apk add --no-cache git openssh openssl>=1.1.1u-r2

# create ~/.ssh & ~/bin folder and change owner to terrascan
RUN mkdir -p /home/terrascan/.ssh /home/terrascan/bin /data && \
    chown -R terrascan:terrascan /home/terrascan /data

# run as non root user
USER 101

ENV PATH /go/bin:$PATH

# copy terrascan binary from build
COPY --from=builder /go/bin/terrascan /go/bin/terrascan

EXPOSE 9010

ENTRYPOINT ["/go/bin/terrascan"]
CMD ["server", "--log-type", "json"]
