package oryhydra

import (
	"context"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	client "github.com/ory/hydra-client-go/v25"
)

func resourceOAuth2Client() *schema.Resource {
	return &schema.Resource{
		Create: resourceOAuth2ClientCreate,
		Read:   resourceOAuth2ClientRead,
		Update: resourceOAuth2ClientUpdate,
		Delete: resourceOAuth2ClientDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"client_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
				// client_id must not be empty string. It must be either unspecified (to make hydra generate the id)
				// or specified as id string.
				ValidateFunc: validation.NoZeroValues,
			},
			"client_secret": {
				Type:      schema.TypeString,
				Optional:  true,
				Computed:  true,
				Sensitive: true,
			},
			"client_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"client_metadata": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"scopes": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"grant_types": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						"refresh_token", "authorization_code", "client_credentials", "implicit",
					}, false),
				},
			},
			"response_types": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						"token", "code", "id_token",
					}, false),
				},
			},
			"audience": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"post_logout_redirect_uris": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"redirect_uris": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"owner": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"policy_uri": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"allowed_cors_origins": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"tos_uri": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"client_uri": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"logo_uri": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"contacts": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"subject_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"public", "pairwise",
				}, false),
			},
			"token_endpoint_auth_method": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"none", "client_secret_basic", "client_secret_post", "private_key_jwt",
				}, false),
			},
			"backchannel_logout_session_required": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"backchannel_logout_uri": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"frontchannel_logout_session_required": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"frontchannel_logout_uri": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"skip_consent": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"skip_logout_consent": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func resourceOAuth2ClientCreate(d *schema.ResourceData, m interface{}) error {
	lg.Print("resourceOAuth2ClientCreate")

	cli := m.(client.OAuth2API)

	resp, _, err := cli.CreateOAuth2Client(context.Background()).
		OAuth2Client(*expandClient(d)).
		Execute()
	if err != nil {
		return err
	}

	d.SetId(resp.GetClientId())

	// NOTE: client secret is only returned on create/update, not read.
	d.Set("client_secret", resp.GetClientSecret())

	return resourceOAuth2ClientRead(d, m)
}

func resourceOAuth2ClientRead(d *schema.ResourceData, m interface{}) error {
	lg.Print("resourceOAuth2ClientRead")

	cli := m.(client.OAuth2API)

	resp, httpResp, err := cli.GetOAuth2Client(context.Background(), d.Id()).Execute()
	if err != nil {
		// NOTE: when client does not exist, Hydra returns 401 even if no authentication is required.
		if httpResp != nil && (httpResp.StatusCode == http.StatusNotFound || httpResp.StatusCode == http.StatusUnauthorized) {
			d.SetId("")
			return nil
		}
		return err
	}

	flattenClient(d, resp)

	return nil
}

func resourceOAuth2ClientUpdate(d *schema.ResourceData, m interface{}) error {
	lg.Print("resourceOAuth2ClientUpdate")

	cli := m.(client.OAuth2API)
	body := expandClient(d)

	_, _, err := cli.SetOAuth2Client(context.Background(), d.Id()).
		OAuth2Client(*body).
		Execute()
	if err != nil {
		return err
	}

	// NOTE: client secret is only returned on create/update, not read.
	d.Set("client_secret", body.GetClientSecret())

	return resourceOAuth2ClientRead(d, m)
}

func resourceOAuth2ClientDelete(d *schema.ResourceData, m interface{}) error {
	lg.Print("resourceOAuth2ClientDelete")

	cli := m.(client.OAuth2API)

	_, err := cli.DeleteOAuth2Client(context.Background(), d.Id()).Execute()
	return err
}

func expandClient(d *schema.ResourceData) *client.OAuth2Client {
	lg.Print("expandClient")

	c := client.NewOAuth2Client()

	if v := d.Get("client_id").(string); v != "" {
		c.SetClientId(v)
	}
	if v := d.Get("client_secret").(string); v != "" {
		c.SetClientSecret(v)
	}
	if v := d.Get("client_name").(string); v != "" {
		c.SetClientName(v)
	}
	if v := d.Get("client_metadata"); v != nil {
		c.Metadata = v
	}

	var scopeArray []string
	for _, sc := range d.Get("scopes").([]interface{}) {
		scopeArray = append(scopeArray, sc.(string))
	}
	c.SetScope(strings.Join(scopeArray, " "))

	var grantTypes []string
	for _, gt := range d.Get("grant_types").([]interface{}) {
		grantTypes = append(grantTypes, gt.(string))
	}
	c.GrantTypes = grantTypes

	var responseTypes []string
	for _, rt := range d.Get("response_types").([]interface{}) {
		responseTypes = append(responseTypes, rt.(string))
	}
	c.ResponseTypes = responseTypes

	var audience []string
	for _, aud := range d.Get("audience").([]interface{}) {
		audience = append(audience, aud.(string))
	}
	c.Audience = audience

	var postLogoutRedirectUris []string
	for _, ru := range d.Get("post_logout_redirect_uris").([]interface{}) {
		postLogoutRedirectUris = append(postLogoutRedirectUris, ru.(string))
	}
	c.PostLogoutRedirectUris = postLogoutRedirectUris

	var redirectUris []string
	for _, ru := range d.Get("redirect_uris").([]interface{}) {
		redirectUris = append(redirectUris, ru.(string))
	}
	c.RedirectUris = redirectUris

	if v := d.Get("owner").(string); v != "" {
		c.SetOwner(v)
	}
	if v := d.Get("policy_uri").(string); v != "" {
		c.SetPolicyUri(v)
	}

	var allowedCorsOrigins []string
	for _, aco := range d.Get("allowed_cors_origins").([]interface{}) {
		allowedCorsOrigins = append(allowedCorsOrigins, aco.(string))
	}
	c.AllowedCorsOrigins = allowedCorsOrigins

	if v := d.Get("tos_uri").(string); v != "" {
		c.SetTosUri(v)
	}
	if v := d.Get("client_uri").(string); v != "" {
		c.SetClientUri(v)
	}
	if v := d.Get("logo_uri").(string); v != "" {
		c.SetLogoUri(v)
	}

	var contacts []string
	for _, ct := range d.Get("contacts").([]interface{}) {
		contacts = append(contacts, ct.(string))
	}
	c.Contacts = contacts

	if v := d.Get("subject_type").(string); v != "" {
		c.SetSubjectType(v)
	}
	if v := d.Get("token_endpoint_auth_method").(string); v != "" {
		c.SetTokenEndpointAuthMethod(v)
	}

	c.SetBackchannelLogoutSessionRequired(d.Get("backchannel_logout_session_required").(bool))
	if v := d.Get("backchannel_logout_uri").(string); v != "" {
		c.SetBackchannelLogoutUri(v)
	}
	c.SetFrontchannelLogoutSessionRequired(d.Get("frontchannel_logout_session_required").(bool))
	if v := d.Get("frontchannel_logout_uri").(string); v != "" {
		c.SetFrontchannelLogoutUri(v)
	}

	c.SetSkipConsent(d.Get("skip_consent").(bool))
	c.SetSkipLogoutConsent(d.Get("skip_logout_consent").(bool))

	return c
}

func flattenClient(d *schema.ResourceData, c *client.OAuth2Client) {
	lg.Print("flattenClient")

	d.Set("client_id", c.GetClientId())
	d.Set("client_name", c.GetClientName())
	d.Set("client_metadata", c.Metadata)

	// NOTE: client secret is never returned from read operations, so don't set this field.

	if scope := c.GetScope(); scope != "" {
		d.Set("scopes", strings.Split(scope, " "))
	}
	d.Set("grant_types", c.GrantTypes)
	d.Set("response_types", c.ResponseTypes)
	d.Set("audience", c.Audience)
	d.Set("post_logout_redirect_uris", c.PostLogoutRedirectUris)
	d.Set("redirect_uris", c.RedirectUris)
	d.Set("owner", c.GetOwner())
	d.Set("policy_uri", c.GetPolicyUri())
	d.Set("allowed_cors_origins", c.AllowedCorsOrigins)
	d.Set("tos_uri", c.GetTosUri())
	d.Set("client_uri", c.GetClientUri())
	d.Set("logo_uri", c.GetLogoUri())
	d.Set("contacts", c.Contacts)
	d.Set("subject_type", c.GetSubjectType())
	d.Set("token_endpoint_auth_method", c.GetTokenEndpointAuthMethod())
	d.Set("backchannel_logout_session_required", c.GetBackchannelLogoutSessionRequired())
	d.Set("backchannel_logout_uri", c.GetBackchannelLogoutUri())
	d.Set("frontchannel_logout_session_required", c.GetFrontchannelLogoutSessionRequired())
	d.Set("frontchannel_logout_uri", c.GetFrontchannelLogoutUri())
	d.Set("skip_consent", c.GetSkipConsent())
	d.Set("skip_logout_consent", c.GetSkipLogoutConsent())
}
