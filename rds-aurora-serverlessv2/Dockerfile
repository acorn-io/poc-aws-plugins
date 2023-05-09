FROM golang:1.20-alpine as build

ADD . /src/
WORKDIR /src

RUN go build -o rds .

FROM node:18-alpine

RUN npm install -g aws-cdk
ADD . /app/
WORKDIR /app
COPY --from=build /src/rds .

CMD [ "/app/scripts/cdk-synth.sh" ]