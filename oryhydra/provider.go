package oryhydra

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	client "github.com/ory/hydra-client-go/v25"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"url": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("ORY_HYDRA_URL", nil),
			},
			"oauth2_token_url": {
				Type:         schema.TypeString,
				Optional:     true,
				RequiredWith: []string{"oauth2_client_id", "oauth2_client_secret"},
				DefaultFunc:  schema.EnvDefaultFunc("ORY_HYDRA_OAUTH2_TOKEN_URL", nil),
			},
			"oauth2_client_id": {
				Type:         schema.TypeString,
				Optional:     true,
				RequiredWith: []string{"oauth2_token_url", "oauth2_client_secret"},
				DefaultFunc:  schema.EnvDefaultFunc("ORY_HYDRA_OAUTH2_CLIENT_ID", nil),
			},
			"oauth2_client_secret": {
				Type:         schema.TypeString,
				Optional:     true,
				RequiredWith: []string{"oauth2_client_id", "oauth2_token_url"},
				DefaultFunc:  schema.EnvDefaultFunc("ORY_HYDRA_OAUTH2_CLIENT_SECRET", nil),
			},
			"header_authorization": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("HEADER_AUTHORIZATION", nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"oryhydra_oauth2_client": resourceOAuth2Client(),
		},
		ConfigureFunc: configure,
	}
}

func configure(data *schema.ResourceData) (interface{}, error) {
	adminURL := data.Get("url").(string)

	httpClient := cleanhttp.DefaultClient()

	if tokenURL, ok := data.GetOk("oauth2_token_url"); ok {
		config := clientcredentials.Config{
			TokenURL:     tokenURL.(string),
			ClientID:     data.Get("oauth2_client_id").(string),
			ClientSecret: data.Get("oauth2_client_secret").(string),
		}
		ctx := context.WithValue(context.TODO(), oauth2.HTTPClient, httpClient)
		httpClient = config.Client(ctx)
	} else if header, ok := data.GetOk("header_authorization"); ok {
		httpClient.Transport = &authHeaderTransport{
			origin: httpClient.Transport,
			header: header.(string),
		}
	}

	return newHydraClient(adminURL, httpClient)
}

// newHydraClient returns a new configured hydra OAuth2 API client.
func newHydraClient(hydraAdminURL string, httpClient *http.Client) (client.OAuth2API, error) {
	if hydraAdminURL == "" {
		return nil, fmt.Errorf("hydra admin URL must not be empty")
	}

	cfg := client.NewConfiguration()
	cfg.Servers = client.ServerConfigurations{
		{URL: hydraAdminURL},
	}
	cfg.HTTPClient = httpClient

	return client.NewAPIClient(cfg).OAuth2API, nil
}

type authHeaderTransport struct {
	origin http.RoundTripper
	header string
}

func (a *authHeaderTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", a.header)
	return a.origin.RoundTrip(req)
}
