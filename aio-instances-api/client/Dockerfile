FROM ubuntu:latest

ENV INSTANCE_ID=0

RUN \
  sed -i 's/# \(.*multiverse$\)/\1/g' /etc/apt/sources.list && \
  apt-get update && \
  apt-get -y upgrade && \
  apt-get install -y build-essential && \
  apt-get install -y software-properties-common && \
  rm -rf /var/lib/apt/lists/*

ADD client-local-v3 /bin/client-local-v3
RUN chmod +x /bin/client-local-v3

ENTRYPOINT ["client-local-v3"]