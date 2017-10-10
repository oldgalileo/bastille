FROM ubuntu
RUN apt-get update
RUN apt-get install -y golang




ADD . /stuff

CMD ["go","run","/stuff/relay.go"]