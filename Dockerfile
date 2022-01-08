FROM golang:1.17 AS builder
WORKDIR /go/src 
ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn
ENV CGO_ENABLED=0  
ADD . .
RUN go mod tidy && go build -tags netgo -o /go/bin/replacer main.go

FROM alpine:3.10
COPY --from=builder /go/bin/lank8swebhook /
ENTRYPOINT ["/replacer"]