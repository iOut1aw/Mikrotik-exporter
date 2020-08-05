ARG BASE_ARCH

FROM ${BASE_ARCH}/debian:9.12-slim

ARG BINARY_ARCH

EXPOSE 9436

COPY scripts/start.sh /app/

COPY dist/mikrotik-exporter_linux_${BINARY_ARCH} /app/mikrotik-exporter

ENTRYPOINT ["/app/start.sh"]
