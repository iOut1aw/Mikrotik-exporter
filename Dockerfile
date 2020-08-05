ARG TARGETPLATFORM

# FROM ${TARGETPLATFORM}/debian:9.12-slim
FROM amd64/debian:9.12-slim

RUN echo "I am running on $BUILDPLATFORM, building for $TARGETPLATFORM" > /log

EXPOSE 9436

COPY scripts/start.sh /app/

COPY dist/mikrotik-exporter_linux_${TARGETPLATFORM} /app/mikrotik-exporter

RUN chmod 755 /app/*

ENTRYPOINT ["/app/start.sh"]
