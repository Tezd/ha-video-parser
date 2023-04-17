FROM alpine:3.17.3

COPY docker/features /features

ENV PATH="$PATH:/usr/local/go/bin"

RUN adduser -D -g parser parser

RUN  sh /features/go/install.sh
RUN  sh /features/ffmpeg/install.sh
