FROM golang:1.20-alpine as build

ADD . /src/
WORKDIR /src

RUN go build -o s3 .

FROM node:18-alpine3.17

RUN npm install -g aws-cdk
ADD . /app/
WORKDIR /app
COPY --from=build /src/s3 .

CMD [ "/app/cdk-synth.sh" ]