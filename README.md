# ORY Hydra Terraform Provider

[![main](https://github.com/hypnoglow/terraform-provider-oryhydra/actions/workflows/main.yml/badge.svg)](https://github.com/hypnoglow/terraform-provider-oryhydra/actions/workflows/main.yml)
[![release](https://github.com/hypnoglow/terraform-provider-oryhydra/actions/workflows/release.yml/badge.svg)](https://github.com/hypnoglow/terraform-provider-oryhydra/actions/workflows/release.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/hypnoglow/terraform-provider-oryhydra)](https://goreportcard.com/report/github.com/hypnoglow/terraform-provider-oryhydra)

Terraform provider for [ORY Hydra](https://github.com/ory/hydra).

## Using the provider

The provider is published to the [Terraform Registry](https://registry.terraform.io/providers/hypnoglow/oryhydra/latest).
Terraform v0.13.0+ supports automatic provider installation from the Terraform Registry. If you run older versions of terraform, you
will have to install the provider manually.

To install this provider, add it into your Terraform configuration. Then, run `terraform init`.

```hcl
terraform {
  required_providers {
    oryhydra = {
      source = "hypnoglow/oryhydra"
      version = "0.4.0"
    }
  }
}

provider "oryhydra" {}
```

To configure the provider, define the `url` argument:

```hcl
provider "oryhydra" {
  url = "https://admin.hydra.example.tld"
}
```
 
Alternatively you can set `ORY_HYDRA_URL` environment variable.

Refer to [documentation](https://registry.terraform.io/providers/hypnoglow/oryhydra/latest/docs) for available resources reference.

Examples can be found in [examples](examples/) directory.

## Developing the Provider

This section explains how to build the provider and test it out on the example.

First of all you need to build the provider and link it to the examples:

```shell script
make build
make prepare-examples
```

To test examples, you need a running Hydra instance. For demonstrational purposes,
we can run Hydra locally as a Docker container:

```shell script
docker container run -i -t --rm --name hydra \
  -p 4444:4444 \
  -p 4445:4445 \
  -e LOG_LEVEL=debug \
  -e DSN=memory \
  oryd/hydra:v1.4.10 serve all --dangerous-force-http
```

Then you can simply run examples as usual terraform project. Enter a particular example directory (e.g. `cd examples/oryhydra_oauth2_client`)
and run:

```shell script
terraform init

terraform plan

terraform apply
```

### Test release

Releases managed by [GoReleaser](https://goreleaser.com/).

Run goreleaser locally:

```shell script
goreleaser release --config .github/goreleaser.yaml --rm-dist --skip-publish --skip-validate
```
