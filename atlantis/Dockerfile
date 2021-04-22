FROM runatlantis/atlantis
ENV DEFAULT_TERRASCAN_VERSION=1.5.0
ENV PLANFILE tfplan
ADD setup.sh terrascan.sh launch-atlantis.sh entrypoint.sh /usr/local/bin/
RUN touch ${PLANFILE} && mkdir -p /etc/atlantis/ && /
    chmod +x /usr/local/bin/setup.sh /usr/local/bin/terrascan.sh /usr/local/bin/launch-atlantis.sh /usr/local/bin/entrypoint.sh && /
    /usr/local/bin/setup.sh
ADD terrascan-workflow.yaml /etc/atlantis/workflow.yaml
ENTRYPOINT ["/bin/bash", "entrypoint.sh"]
CMD ["server"]
