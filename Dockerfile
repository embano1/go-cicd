FROM golang:alpine3.6 AS build
WORKDIR /go/src/go-cicd
RUN apk add --no-cache git
COPY cmd/main.go .
RUN go get . 
RUN go install

FROM scratch
LABEL MAINTAINER=embano1@live.com
LABEL VERSION="1.1"
COPY --from=build /go/bin/go-cicd /
ENTRYPOINT [ "/go-cicd" ]
CMD ["-h"]