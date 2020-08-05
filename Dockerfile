ARG BASE_ARCH
ARG BINARY_ARCH

FROM ${BASE_ARCH}/debian:9.12-slim

EXPOSE 9436

COPY scripts/start.sh /app/

COPY dist/mikrotik-exporter_linux_${BINARY_ARCH} /app/mikrotik-exporter

RUN chmod 755 /app/*

ENTRYPOINT ["/app/start.sh"]
