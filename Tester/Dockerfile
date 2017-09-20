FROM docker:dind

RUN ["apk", "add", "-U", "bash", "git", "go", "e2fsprogs", "curl", "zip", "sed"]

RUN curl https://gitlab-ci-multi-runner-downloads.s3.amazonaws.com/latest/binaries/gitlab-ci-multi-runner-linux-amd64 > /usr/bin/gitlab-runner && \
	chmod +x /usr/bin/gitlab-runner

ENV GOPATH=/go \
	PATH=$GOPATH/bin:/usr/local/go/bin:$PATH

COPY /install /go/src/app/
RUN ["/go/src/app/install"]
VOLUME ["/var/lib/docker"]

WORKDIR /go/src/app/
COPY /*.go /go/src/app/
RUN ["go", "build", "-o", "tester"]

COPY /run /go/src/app/
ENTRYPOINT ["/go/src/app/run"]
