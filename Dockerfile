ARG TARGETPLATFORM

RUN echo "I am running on $BUILDPLATFORM, building for $TARGETPLATFORM" > /log

FROM ${TARGETPLATFORM}/debian:9.12-slim

EXPOSE 9436

COPY scripts/start.sh /app/

COPY dist/mikrotik-exporter_linux_${TARGETPLATFORM} /app/mikrotik-exporter

RUN chmod 755 /app/*

ENTRYPOINT ["/app/start.sh"]
