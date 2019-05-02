FROM gcr.io/gcp-runtimes/go1-builder:1.12 as builder
RUN apt-get install -y gcc
WORKDIR /go/src/app
COPY . .
RUN \
	rm -rf /root/go/src; \
	mkdir -p /root/go/src; \
	tar -xvf dep.tar -C /root/go/src; \
	rm dep.tar
RUN /usr/local/go/bin/go build -o /usr/local/bin/app .
CMD ["/usr/local/bin/app"]
