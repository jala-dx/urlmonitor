FROM ubuntu:latest
COPY ./urlmonitor /bin/urlmonitor
COPY ./config.json /tmp/config.json

RUN apt-get update && apt-get install -y curl
RUN apt-get install iputils-ping -y
CMD /bin/bash
ENTRYPOINT ["/bin/urlmonitor"]
