# ORY Hydra Provider

Terraform provider for [ORY Hydra](https://github.com/ory/hydra).

The provider is used to interact with and manage resources like OAuth2 Clients supported by Hydra.

## Example Usage

```hcl
# Configure the ORY Hydra Provider
provider "oryhydra" {
  url = "http://localhost:4445"
}
```

## Argument Reference

In addition to [generic `provider` arguments](https://www.terraform.io/docs/configuration/providers.html) (e.g. `alias` and `version`),
the following arguments are supported in the ORY Hydra `provider` block:

* `url` - (Optional) URL for Hydra [administrative API](https://www.ory.sh/hydra/docs/reference/api/#administrative-endpoints).
It must be provided, but it can also be sourced from the `ORY_HYDRA_URL` environment variable.

### Authentication

If the Hydra administrative API is protected with the OAuth2.0 "client credentials" token flow,
the following arguments can be set to obtain a bearer token beforehand.

* `oauth2_token_url` - (Optional) Token URL to use for OAuth2.0 flow. Can also be sourced from the `ORY_HYDRA_OAUTH2_TOKEN_URL` environment variable.
* `oauth2_client_id` - (Optional) Client ID used for OAuth2.0 flow. Can also be sourced from the `ORY_HYDRA_OAUTH2_CLIENT_ID` environment variable.
* `oauth2_client_secret` - (Optional) Client secret used for OAuth2.0 flow. Can also be sourced from the `ORY_HYDRA_OAUTH2_CLIENT_SECRET` environment variable.
