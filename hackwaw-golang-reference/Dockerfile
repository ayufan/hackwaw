FROM golang:alpine
RUN ["apk", "add", "-U", "git"]

ADD . /go/src/hackwaw/app/
WORKDIR /go/src/hackwaw/app/
RUN ["go", "get", "./..."]
RUN ["go", "install"]
CMD ["app"]
EXPOSE 8080
