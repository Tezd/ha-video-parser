FROM alpine:3.17.3

COPY docker/features /features

ENV PATH="$PATH:/usr/local/go/bin"

RUN adduser -D -g parser parser

RUN  /features/go/install.sh \
     && /features/ffmpeg/install.sh \
     && /features/imagemagick/install.sh
