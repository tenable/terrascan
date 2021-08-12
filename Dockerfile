# -------- builder stage -------- #
FROM golang:alpine AS builder

ARG GOOS_VAL=linux
ARG GOARCH_VAL=amd64
ARG CGO_ENABLED_VAL=1

WORKDIR $GOPATH/src/terrascan

# download go dependencies
COPY go.mod go.sum ./
RUN go mod download
RUN apk add -U build-base
RUN apk add --no-cache curl 


# copy terrascan source
COPY . .


# build binary
RUN CGO_ENABLED=${CGO_ENABLED_VAL} GOOS=${GOOS_VAL} GOARCH=${GOARCH_VAL} go build -v -ldflags "-w -s" -o /go/bin/terrascan ./cmd/terrascan


# -------- prod stage -------- #
FROM alpine:3.12.0

# create non root user
RUN addgroup --gid 101 terrascan && \
    adduser -S --uid 101 --ingroup terrascan terrascan && \
    apk add --no-cache git openssh


# create ~/.ssh & ~/bin folder and change owner to terrascan
RUN mkdir -p /home/terrascan/.ssh /home/terrascan/bin && \
    chown -R terrascan:terrascan /home/terrascan

COPY /go/bin/terrascan/integrations/argocd/scripts/argocd-terrascan-remote-scan.sh  /home/terrascan/bin/terrascan-remote-scan.sh

RUN chmod u+x /home/terrascan/bin/terrascan-remote-scan.sh

# run as non root user
USER 101

ENV PATH /go/bin:$PATH

# copy terrascan binary from build
COPY --from=builder /go/bin/terrascan /go/bin/terrascan


# Copy webhooks UI templates & assets
COPY ./pkg/http-server/templates /go/terrascan
COPY ./pkg/http-server/assets /go/terrascan/assets

EXPOSE 9010

ENTRYPOINT ["/go/bin/terrascan"]
CMD ["server", "--log-type", "json", "sh"]