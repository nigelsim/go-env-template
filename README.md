# go-env-template

A statically compilable templating application written in Golang. It uses [Mustache](https://github.com/alexkappa/mustache) for the templating, and gets all the template variables from the environment (e.g., Docker container).

## The Goal

The purpose of this application is to be able to template a static website for use in a Docker container. Specifically, the Docker container will run using [goStatic](https://github.com/PierreZ/goStatic), however, some of the paths, or features of the site will need to be configured based on how it is deployed (e.g., Dev vs Prod, context path, etc). This is where **go-env-template** comes in. It can do a one-shot parse of the static site, and render the final product.

The key is that the application can be statically compiled, resulting in a 3MB file that can be dropped into a `FROM scratch` Docker container.

### Usage

```
./go-env-template --help
Usage of ./go-env-template:
  -ext string
        The template extension to look for (default ".tmpl")
  -path string
        The path for the static files (default "/srv/http")
```

### Building

For my own notes:

```
docker run -it -w /src/go-env-template -v $PWD:/src/go-env-template golang:latest /bin/bash
go get github.com/alexkappa/mustache
CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build
```

### Example Usage

This is an example that builds a Docker container for a static website. It pulls in this tool, along with x for process hosting, and goStatic for web serving. It then processes the web content using the environment variables

Dockerfile:
```
# stage 0
FROM golang:latest as builder
RUN go get github.com/nigelsim/go-env-template
RUN go get github.com/pablo-ruth/go-init

FROM pierrezemb/gostatic

COPY --from=builder /go/bin/go-static-template /
COPY --from=builder /go/bin/go-init /

COPY dist/ /srv/http
COPY frontend-static/ /srv/http

ENTRYPOINT ["/go-init"]

CMD ["-pre", "/go-env-template /srv/http", "-main", "/goStatic --fallback index.html"]
```
