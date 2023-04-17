#!/usr/bin/env sh

apk add -Ut build-deps curl

apk add -Ut deps libc6-compat=1.2.3-r4
ln -s /lib/libc.so.6 /usr/lib/libresolv.so.2 #fix ldd missing dependency

curl -sSL https://go.dev/dl/go1.20.3.linux-amd64.tar.gz > /tmp/go.tar.gz
  rm -rf /usr/local/go && tar -C /usr/local -xzf /tmp/go.tar.gz

apk del build-deps
