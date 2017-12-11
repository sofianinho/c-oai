FROM golang:1.8

COPY . /go/src/github.com/sofianinho/vnf-api-golang

WORKDIR /go/src/github.com/sofianinho/vnf-api-golang

RUN apt-get update \
  && apt-get install -y ca-certificates curl \
  && cd /tmp && curl -L https://glide.sh/get -O -J && sh ./get \
  && rm /tmp/get

RUN glide install

RUN CGO_ENABLED=0 go build -a -installsuffix nocgo -o /go/bin/vnf-api .


FROM alpine

LABEL atom="vnf-api-golang"
LABEL REPO="https://gitlab.forge.orange-labs.fr/lucy/vnf-api-golang"

ENV PATH=/app:$PATH

WORKDIR /app


RUN apk --update add ca-certificates
RUN mkdir -p /app/

COPY --from=0 /go/bin/vnf-api /app/vnf-api
COPY ./templates /app/templates/
COPY ./config/config.json /app/config/


CMD ["/app/vnf-api"]

EXPOSE 1337


