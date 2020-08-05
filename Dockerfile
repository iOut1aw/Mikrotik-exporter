# FROM ${TARGETPLATFORM}/debian:9.12-slim
FROM --platform=$BUILDPLATFORM debian:9.12-slim

ARG TARGETPLATFORM
ARG BUILDPLATFORM

RUN echo "I am running on $BUILDPLATFORM"
RUN echo "Building for $TARGETPLATFORM"

EXPOSE 9436

COPY scripts/start.sh /app/

COPY dist/mikrotik-exporter_linux_${TARGETPLATFORM} /app/mikrotik-exporter

RUN chmod 755 /app/*

ENTRYPOINT ["/app/start.sh"]
