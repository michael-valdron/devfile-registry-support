# Registry Support

<div id="header">

[![Apache2.0 License](https://img.shields.io/badge/license-Apache2.0-brightgreen.svg)](LICENSE)
</div>

Provide support for devfile registries

Issue tracking repo: https://github.com/devfile/api with label area/registry

## Components

This repository contains the list of devfile registry components below:

- [Devfile Index Server](./index/server/README.md)
- [Devfile Index Generator](./index/generator/README.md)
- [Devfile Registry Library](./registry-library/README.md)
- [Devfile Registry Helm Chart](./deploy/chart/devfile-registry/README.md)

## Build

If you want to run the build scripts with Podman, set the environment variable
`export USE_PODMAN=true`

To build all of the components together (recommended) for dev/test, run `bash ./build_registry.sh` to build a Devfile Registry index image that is populated with the mock devfile registry data under `tests/registry/`.

Once the container has been pushed, you can push it to a container registry of your choosing with the following commands:

```
docker tag devfile-index <registry>/<username>/devfile-index:latest
```

followed by

```
docker push <registry>/<username>/devfile-index:latest
```

See the [readme](./build-tools/README.md) on build tools to learn more about the devfile registry build automation.

## Deploy

### Via the Devfile Registry Operator

We recommend using the [Devfile Registry Operator](https://github.com/devfile/registry-operator) to install a Devfile Registry on your Kubernetes or OpenShift cluster. Consult [its Readme for more information](https://github.com/devfile/registry-operator#running-the-controller-in-a-cluster).

### Via the Devfile Registry Helm Chart

Alternatively, a [Helm chart](./deploy/chart/devfile-registry) is also provided if you do not wish to use an operator. To install (with Helm 3) run:

```bash
$ helm install devfile-registry ./deploy/chart/devfile-registry \ 
    --set global.ingress.domain=<ingress-domain> \
	--set devfileIndex.image=<index-image> \
	--set devfileIndex.tag=<index-image-tag>
```

Where `<ingress-domain>` is the ingress domain for your cluster, `<index-image>` is the devfile index image you want to deploy, and `<index-image-tag>` is the corresponding image tag for the devfile index image.

For example, if you're installing your own custom devfile registry image for dev/test purposes on Minikube, you might run:

```bash
$ helm install devfile-registry ./deploy/chart/devfile-registry \
    --set global.ingress.domain="$(minikube ip).nip.io" \
	--set devfileIndex.image=quay.io/someuser/devfile-index \
	--set devfileIndex.tag=latest
```

You can deploy a devfile registry with a custom registry viewer image (uses `quay.io/devfile/registry-viewer:next` by default) by running the following:

```bash
$ helm install devfile-registry ./deploy/chart/devfile-registry \
    --set global.ingress.domain="$(minikube ip).nip.io" \
	--set devfileIndex.image=quay.io/someuser/devfile-index \
	--set devfileIndex.tag=latest \
	--set registryViewer.image=quay.io/someuser/registry-viewer \
	--set registryViewer.tag=latest
```

You can deploy a *headless* devfile registry (i.e. without the registry viewer) by specifying `--set global.headless=true` which will look like:

```bash
$ helm install devfile-registry ./deploy/chart/devfile-registry \
    --set global.ingress.domain="$(minikube ip).nip.io" \
	--set global.headless=true \
	--set devfileIndex.image=quay.io/someuser/devfile-index \
	--set devfileIndex.tag=latest
```

For more information on the Helm chart, consult [its readme](deploy/chart/devfile-registry/README.md).

## Testing

Please see the integration tests [readme](./tests/integration/README.md) for more information about devfile registry testing.

## Troubleshooting

Please see [TROUBLESHOOTING.md](./TROUBLESHOOTING.md) for more information on how to troubleshoot issues within the devfile registry service.

## Contributing

Please see our [contributing.md](./CONTRIBUTING.md).
