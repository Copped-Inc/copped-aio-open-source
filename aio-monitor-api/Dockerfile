FROM golang:1.19-alpine

RUN apk update && apk add git
RUN apk add gcc

WORKDIR /build

RUN --mount=type=secret,id=gh_access_token \
   export GH_ACCESS_TOKEN=$(cat /run/secrets/gh_access_token) && \
   git config --global url."https://$GH_ACCESS_TOKEN@github.com".insteadOf https://github.com

RUN mkdir /root/.ssh && echo "StrictHostKeyChecking no " > /root/.ssh/config


COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY . .

RUN chmod +x ./ai-model/ai-model-cli
RUN go build -o /main

CMD [ "/main" ]
