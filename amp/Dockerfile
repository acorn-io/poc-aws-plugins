FROM golang:1.20-alpine as build

ADD . /src/
WORKDIR /src

RUN go build -o amp .

FROM node:19-alpine3.16

RUN npm install -g aws-cdk 
ADD . /app/
WORKDIR /app
COPY --from=build /src/amp .

CMD [ "/app/cdk-synth.sh" ]