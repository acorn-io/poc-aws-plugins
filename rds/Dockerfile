FROM cgr.dev/chainguard/go as build
ARG MAIN
WORKDIR /src
COPY . .
RUN --mount=type=cache,target=/root/go/pkg \
    --mount=type=cache,target=/root/.cache/go-build \
    go build -o rds ./${MAIN}

FROM cgr.dev/chainguard/mariadb as user
WORKDIR /app
COPY ./scripts ./scripts
ENTRYPOINT ["/app/scripts/create_and_grant_users.sh"]

FROM cgr.dev/chainguard/wolfi-base
RUN apk add -U --no-cache nodejs bash busybox aws-cli jq curl zip && \
    apk del --no-cache wolfi-base apk-tools
RUN curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip" && \
    unzip awscliv2.zip && \
    ./aws/install
RUN npm install -g aws-cdk
WORKDIR /app
COPY ./cdk.json ./
COPY ./scripts ./scripts
COPY --from=build /src/rds .
CMD [ "/app/scripts/apply.sh" ]
