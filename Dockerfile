FROM --platform=$BUILDPLATFORM debian:9.12-slim

ARG TARGETPLATFORM
ARG BUILDPLATFORM

RUN if [ "$TARGETPLATFORM" == "linux/amd64" ]; \
    then export ARCH=amd64; \
    elif [ "$TARGETPLATFORM" == "linux/arm/v7" ]; \
    then export ARCH=arm; \
    lif [ "$TARGETPLATFORM" == "linux/arm64" ]; \
    then export ARCH=arm64; \
    lif [ "$TARGETPLATFORM" == "linux/386" ]; \
    then export ARCH=386; \
    fi 

ARG ARCH
ENV BINARY=$ARCH

RUN echo "I am running on $BUILDPLATFORM"
RUN echo "Building for $TARGETPLATFORM"
RUN echo "Binary is $BINARY"

EXPOSE 9436

COPY scripts/start.sh /app/

COPY dist/mikrotik-exporter_linux_${BINARY} /app/mikrotik-exporter

RUN chmod 755 /app/*

ENTRYPOINT ["/app/start.sh"]
