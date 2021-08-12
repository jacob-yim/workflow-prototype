##
## Build
##

FROM golang:1.16-buster AS build

COPY . /usr/src/app
WORKDIR /usr/src/app
RUN go mod download

RUN go build -o /workflow-app

##
## Deploy
##

FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /workflow-app /workflow-app

USER nonroot:nonroot

ENTRYPOINT ["/workflow-app"]