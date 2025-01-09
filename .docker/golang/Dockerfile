FROM golang:1.23-bookworm

RUN apt-get update

RUN mkdir -p /opt/bin
RUN curl -sSfL https://raw.githubusercontent.com/air-verse/air/master/install.sh | sh -s -- -b /opt/bin
RUN chmod ugo+rx /opt/bin/air
RUN echo 'export PATH=$PATH:/opt/bin' >> /root/.bashrc

WORKDIR /var/www/html

ENTRYPOINT ["/opt/bin/air"]
