provider "oryhydra" {
  url = "http://localhost:4445"
}

resource "oryhydra_oauth2_client" "example_1" {
  client_name = "Example 1"
}

resource "oryhydra_oauth2_client" "example_2" {
  client_id     = "example_2"
  client_name   = "Example 2"
  client_secret = "super-secret!"
}

resource "oryhydra_oauth2_client" "example_3" {
  client_name = "Example 3"
  client_metadata = {
    "Foo" = "Bar"
  }

  scopes         = ["offline", "openid"]
  grant_types    = ["refresh_token", "authorization_code"]
  response_types = ["code"]

  audience      = ["http://localhost:8080"]
  redirect_uris = ["http://localhost:8080/redirect.html"]

  subject_type               = "public"
  token_endpoint_auth_method = "none"
}
