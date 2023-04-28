# syntax=docker/dockerfile:1
FROM golang:1.20-alpine as build

ADD . /src/
WORKDIR /src

RUN go build -o cdk-context-render

FROM alpine as app
RUN apk --no-cache add ca-certificates
COPY --from=build /src/cdk-context-render /

ENV OUTPUT_FILE "/run/secrets/output"
ENV AWS_CONFIG_PATH "/acorn/aws"

CMD ["/cdk-context-render"]
