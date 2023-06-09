# Devfile Registry Index Server

<div id="header">

[![Apache2.0 License](https://img.shields.io/badge/license-Apache2.0-brightgreen.svg)](LICENSE)
</div>

Provides REST API support for devfile registries and serves [devfile registry viewer](https://github.com/devfile/devfile-web) client.

For more information on REST API docs: [registry-REST-API.adoc](registry-REST-API.adoc)

## Build

If you want to run the build scripts with Podman, set the environment variable
`export USE_PODMAN=true`

You can build a local binary for the index server using:
```sh
CGO_ENABLED=0 go build -mod=vendor -o index-server main.go
```

or

You can build the index server, `devfile-index-base`, container image by using:
```sh
bash ./build.sh
```

After building the container image you can push it up to an image registry running the following:
```sh
bash push.sh <registry>/<user>/devfile-index-base:<tag>
```

This also applies a new tag you can give the image rather than using the default `latest`.

For building the entire devfile registry mechanism refer back to the [base readme](../../README.md#build).

## Testing

Endpoint unit testing is defined under `pkg/server/endpoint_test.go` and can be performed by running the following:

```sh
go test pkg/server/endpoint_test.go
```

or by running all tests:

```sh
go test ./...
```

**Environment Variables**

- `DEVFILE_REGISTRY`: Optional environment variable for specifying testing registry path
    - default: `../../tests/registry`
