# ORY Hydra Terraform Provider

Terraform provider for [ORY Hydra](https://github.com/ory/hydra).

## Install

As this provider is custom, terraform won't be able to download it automatically.

To install this provider in a [recommended way](https://www.terraform.io/docs/configuration/providers.html#third-party-plugins)
e.g. for Linux:

```shell script
VERSION="v0.1.0"

cd /tmp
wget https://github.com/hypnoglow/terraform-provider-oryhydra/releases/download/${VERSION}/terraform-provider-oryhydra_${VERSION}_linux_amd64.tar.gz
tar xzf terraform-provider-oryhydra_${VERSION}_linux_amd64.tar.gz
mv terraform-provider-oryhydra_${VERSION} ~/.terraform.d/plugins/linux_amd64/terraform-provider-oryhydra_${VERSION}
```

Then you can set up a provider in your terraform project:

```hcl-terraform
provider "oryhydra" {
  version = "0.1.1"
}
```

## Usage

For example usage see [examples](examples/README.md).

## Implemented resources

- OAuth2 Client - `oryhydra_oauth2_client`
