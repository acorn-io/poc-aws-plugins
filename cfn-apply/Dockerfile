FROM amazon/aws-cli

RUN yum install -y jq

ADD . /app/
WORKDIR /app

CMD ["/app/cfn-apply.sh"]