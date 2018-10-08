FROM centos:7

ARG GO_VERSION=1.11.1
ARG GO_OS=linux
ARG GO_ARCH=amd64

RUN GO_TGZ=go$GO_VERSION.$GO_OS-$GO_ARCH.tar.gz && \
    curl -LO https://dl.google.com/go/$GO_TGZ && \
    tar -C /usr/local -xzf $GO_TGZ && \
    rm $GO_TGZ

ENV GOPATH /workspace
ENV PATH "$PATH:/usr/local/go/bin:$GOPATH/bin"

RUN mkdir -p "$GOPATH/src" "$GOPATH/bin"

# glide installation script need it
RUN yum install which git -y
RUN curl https://glide.sh/get | sh

ENV SRC_DIR "$GOPATH/src/junyan.org/test"
WORKDIR $SRC_DIR

COPY glide.yaml $SRC_DIR/
COPY glide.lock $SRC_DIR/

RUN glide install

COPY *.go $SRC_DIR

RUN go build  -o /usr/local/bin/test

EXPOSE 8888
ENTRYPOINT [ "/usr/local/bin/test" ]
