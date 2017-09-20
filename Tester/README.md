## HackWAW Judge :)

### How it works?

1. It uses [bats](https://github.com/sstephenson/bats),
2. In `tests/` we do have behavior tests,
3. We build app using the Dockerfile from external repository. The app needs to be stored in `/app`,
4. We start the app from `docker-compose.yml`, defining the persistent volume (`/storage`)
and proxy settings (we define `HTTP_PROXY` and `HTTPS_PROXY`),
5. All outgoing traffic from APP should go through proxy server,
6. We provide a set of tests `tests/*.bats` testing behavior of app in different scenarios.

### The flow

1. We build the tester app docker image from this [Dockerfile](./Dockerfile),
2. We start docker container, passing the docker.sock to the container,
3. We execute all bats tests from [tests/](./tests/),
4. (In the future) we will post test results outside.

### How to create an app

1. Create `Dockerfile` in your directory, with `CMD` defined:

    ```
    FROM golang:alpine
    
    ADD . /go/src/hackwaw/app/
    WORKDIR /go/src/hackwaw/app/
    CMD ["go", "run", "main.go"]
    EXPOSE 8080
    ````

2. Create `.gitlab-ci.yml` with this content:

    ```yaml
    image: hackwaw-disrupt/hackwaw-tester
    test:
      script: ""
    ```

3. Your app needs to start an HTTP server listening on `:8080`.

## How to setup GitLab Runner (in secure way)

1. Install a runner on machine using this instruction: https://gitlab.com/gitlab-org/gitlab-ci-multi-runner/blob/master/docs/install/linux-repository.md,
1. Register runner: `sudo gitlab-ci-multi-runner register`,
1. Use credentials from: https://gitlab.com/hackwaw-disrupt/hackwaw-golang-example (subject to change: we ideally should use shared runners),
1. Select `docker` executor,
1. Finish the registration,
1. Edit `/etc/gitlab-runner/config.toml` and add `allowed_images` and `allowed_services` (security):

    ```yaml
    [[runners]]
    executor = "docker"
    ...
    [runners.docker]
    allowed_images = ["hackwaw-tester"]
    allowed_services = [""]
    ```
    
1. Download sources from this repository and build the docker image on machine:

    ```
    docker build -t hackwaw-tester .
    ```

1. You should see your registered runner in your project: https://gitlab.com/hackwaw-disrupt/hackwaw-golang-example/runners.
1. You should be ready to test everything now.

## Test locally

1. Create VirtualBox docker-machine (for OSX, on Linux you can use your Docker Engine):
`docker-machine create -d virtualbox hackwaw-disrupt`,
2. Use `eval $(docker-machine env hackwaw-disrupt)`,
3. Execute `./run-locally <path-to-your-app>`
