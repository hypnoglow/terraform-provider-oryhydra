package oryhydra

import (
	"fmt"
	"net/url"

	httptransport "github.com/go-openapi/runtime/client"
	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	hydra "github.com/ory/hydra-client-go/client"
	"github.com/ory/hydra-client-go/client/admin"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"url": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("ORY_HYDRA_URL", nil),
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
	client, err := newHydraClient(adminURL)
	return client, err
}

// newHydraClient returns a new configured hydra client.
func newHydraClient(hydraAdminURL string) (admin.ClientService, error) {
	u, err := url.Parse(hydraAdminURL)
	if err != nil {
		return nil, fmt.Errorf("parse hydra url: %v", err)
	}

	config := hydra.DefaultTransportConfig()
	config.Schemes = []string{u.Scheme}
	config.Host = u.Host
	if u.Path != "" {
		config.BasePath = u.Path
	}

	transport := httptransport.NewWithClient(
		config.Host,
		config.BasePath,
		config.Schemes,
		cleanhttp.DefaultClient(),
	)

	client := hydra.New(transport, nil)
	return client.Admin, nil
}
