FROM node:18-alpine

RUN npm install -g aws-cdk
ADD . /app/
WORKDIR /app

CMD [ "/app/scripts/cdk-synth.sh" ]