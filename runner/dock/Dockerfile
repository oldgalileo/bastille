FROM ubuntu

RUN apt-get update
ENV DEBIAN_FRONTEND noninteractive

ENV INITRD No
ENV GOVERSION 1.9.1
ENV GOROOT /opt/go
ENV GOPATH /root/.go

RUN apt-get install -y wget

RUN cd /opt && wget https://storage.googleapis.com/golang/go${GOVERSION}.linux-amd64.tar.gz && \
    tar zxf go${GOVERSION}.linux-amd64.tar.gz && rm go${GOVERSION}.linux-amd64.tar.gz && \
    ln -s /opt/go/bin/go /usr/bin/ && \
    mkdir $GOPATH

RUN apt-get install -y nano

RUN apt-get update \
  && apt-get install -y python3-pip python3-dev \
  && cd /usr/local/bin \
  && ln -s /usr/bin/python3 python \
  && pip3 install --upgrade pip







ADD . /stuff

CMD ["go","run","/stuff/relay.go"]