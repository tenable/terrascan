# Dockerfile with a script to use terrascan's validating webhook
# configured in the kubernetes cluster, to scan a repo for violations
FROM alpine:3.19.4

#curl to send request to terrascan validating webhook
RUN apk add --no-cache curl

WORKDIR /home/terrascan

RUN mkdir bin

COPY scripts/argocd-terrascan-remote-scan.sh  bin/terrascan-remote-scan.sh

# create non root user
RUN addgroup --gid 101 terrascan && \
    adduser -S --uid 101 --ingroup terrascan terrascan && \
    chown -R terrascan:terrascan bin && \
    chmod u+x bin/terrascan-remote-scan.sh

USER 101

CMD ["sh"]
