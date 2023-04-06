FROM node:19-alpine3.16

RUN npm install -g aws-cdk
ADD . /app/

WORKDIR /app

CMD [ "/app/cdk-synth.sh" ]