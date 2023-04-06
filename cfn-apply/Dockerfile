FROM amazon/aws-cli

ADD . /app/
WORKDIR /app

ENTRYPOINT ["/app/cfn-apply.sh"]
