FROM amazon/aws-cli

RUN yum install -y jq

ADD . /app/
WORKDIR /app

ENTRYPOINT ["/app/cfn-apply.sh"]
