FROM golang:1.19 as builder
WORKDIR /go/src
RUN git clone https://github.com/watsonserve/scaner.git && export GOPROXY=goproxy.cn && cd /go/src/scaner && go install

FROM ubuntu
COPY --from=builder /go/bin/scaner /usr/bin/
EXPOSE 80
COPY ./scaner.conf /etc/scaner.conf

CMD  scaner --config=/etc/scaner.conf :80
