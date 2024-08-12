FROM golang:alpine

RUN mkdir /go/src/dory
WORKDIR /go/src/dory

RUN apk add --no-cache git build-base

COPY . .

RUN go mod download
RUN go build -o /go/bin/dory ./cmd

FROM alpine

RUN mkdir /app

COPY --from=0 --chmod=0755 /go/bin/dory /app/dory
COPY --chmod=0644 templates/* /app/templates/

WORKDIR /app
ENTRYPOINT ["./dory"]

EXPOSE 8000
