FROM ubuntu

RUN apt-get update
ENV DEBIAN_FRONTEND noninteractive

ENV INITRD No
ENV LANG en_US.UTF-8
ENV GOVERSION 1.9.1
ENV GOROOT /opt/go
ENV GOPATH /root/.go

RUN apt-get install -y wget

RUN cd /opt && wget https://storage.googleapis.com/golang/go${GOVERSION}.linux-amd64.tar.gz && \
    tar zxf go${GOVERSION}.linux-amd64.tar.gz && rm go${GOVERSION}.linux-amd64.tar.gz && \
    ln -s /opt/go/bin/go /usr/bin/ && \
    mkdir $GOPATH

RUN apt-get install -y nano




ADD . /stuff

CMD ["go","run","/stuff/relay.go"]