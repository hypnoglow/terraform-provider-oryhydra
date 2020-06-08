package oryhydra

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/ory/hydra-client-go/client/admin"
	"github.com/ory/hydra-client-go/models"
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
		},
	}
}

func resourceOAuth2ClientCreate(d *schema.ResourceData, m interface{}) error {
	lg.Print("resourceOAuth2ClientCreate")

	cli := m.(*admin.Client)

	resp, err := cli.CreateOAuth2Client(
		admin.NewCreateOAuth2ClientParams().
			WithBody(expandClient(d)),
	)
	if err != nil {
		return err
	}

	client := resp.Payload

	d.SetId(client.ClientID)

	// NOTE: client secret is only returned on create/update, not read.
	d.Set("client_secret", client.ClientSecret)

	return resourceOAuth2ClientRead(d, m)
}

func resourceOAuth2ClientRead(d *schema.ResourceData, m interface{}) error {
	lg.Print("resourceOAuth2ClientRead")

	cli := m.(*admin.Client)

	resp, err := cli.GetOAuth2Client(
		admin.NewGetOAuth2ClientParams().
			WithID(d.Id()),
	)
	if err != nil {
		return err
	}

	client := resp.Payload

	flattenClient(d, client)

	return nil
}

func resourceOAuth2ClientUpdate(d *schema.ResourceData, m interface{}) error {
	lg.Print("resourceOAuth2ClientUpdate")

	client := expandClient(d)

	cli := m.(*admin.Client)

	_, err := cli.UpdateOAuth2Client(
		admin.NewUpdateOAuth2ClientParams().
			WithID(d.Id()).
			WithBody(client),
	)
	if err != nil {
		return err
	}

	// NOTE: client secret is only returned on create/update, not read.
	d.Set("client_secret", client.ClientSecret)

	return resourceOAuth2ClientRead(d, m)
}

func resourceOAuth2ClientDelete(d *schema.ResourceData, m interface{}) error {
	lg.Print("resourceOAuth2ClientDelete")

	cli := m.(*admin.Client)

	_, err := cli.DeleteOAuth2Client(
		admin.NewDeleteOAuth2ClientParams().
			WithID(d.Id()),
	)
	if err != nil {
		return err
	}

	return nil
}

func expandClient(d *schema.ResourceData) *models.OAuth2Client {
	lg.Print("expandClient")

	clientID := d.Get("client_id").(string)
	clientSecret := d.Get("client_secret").(string)
	clientName := d.Get("client_name").(string)
	clientMetadata := d.Get("client_metadata")
	lg.Printf("metadata: %T %v", clientMetadata, clientMetadata)

	var scopeArray []string
	for _, sc := range d.Get("scopes").([]interface{}) {
		scopeArray = append(scopeArray, sc.(string))
	}
	scope := strings.Join(scopeArray, " ")

	var grantTypes []string
	for _, gt := range d.Get("grant_types").([]interface{}) {
		grantTypes = append(grantTypes, gt.(string))
	}

	var responseTypes []string
	for _, rt := range d.Get("response_types").([]interface{}) {
		responseTypes = append(responseTypes, rt.(string))
	}

	var audience []string
	for _, aud := range d.Get("audience").([]interface{}) {
		audience = append(audience, aud.(string))
	}

	var postLogoutRedirectUris []string
	for _, ru := range d.Get("post_logout_redirect_uris").([]interface{}) {
		postLogoutRedirectUris = append(postLogoutRedirectUris, ru.(string))
	}

	var redirectUris []string
	for _, ru := range d.Get("redirect_uris").([]interface{}) {
		redirectUris = append(redirectUris, ru.(string))
	}

	owner := d.Get("owner").(string)
	policyURI := d.Get("policy_uri").(string)

	var allowedCorsOrigins []string
	for _, aco := range d.Get("allowed_cors_origins").([]interface{}) {
		allowedCorsOrigins = append(allowedCorsOrigins, aco.(string))
	}

	tosURI := d.Get("tos_uri").(string)
	clientURI := d.Get("client_uri").(string)
	logoURI := d.Get("logo_uri").(string)

	var contacts []string
	for _, c := range d.Get("contacts").([]interface{}) {
		contacts = append(contacts, c.(string))
	}

	subjectType := d.Get("subject_type").(string)
	tokenEndpointAuthMethod := d.Get("token_endpoint_auth_method").(string)

	return &models.OAuth2Client{
		ClientID:                clientID,
		ClientName:              clientName,
		ClientSecret:            clientSecret,
		Metadata:                clientMetadata,
		Scope:                   scope,
		GrantTypes:              grantTypes,
		ResponseTypes:           responseTypes,
		Audience:                audience,
		PostLogoutRedirectUris:  postLogoutRedirectUris,
		RedirectUris:            redirectUris,
		Owner:                   owner,
		PolicyURI:               policyURI,
		AllowedCorsOrigins:      allowedCorsOrigins,
		TosURI:                  tosURI,
		ClientURI:               clientURI,
		LogoURI:                 logoURI,
		Contacts:                contacts,
		SubjectType:             subjectType,
		TokenEndpointAuthMethod: tokenEndpointAuthMethod,
	}
}

func flattenClient(d *schema.ResourceData, client *models.OAuth2Client) {
	lg.Print("flattenClient")

	d.Set("client_id", client.ClientID)
	d.Set("client_name", client.ClientName)
	d.Set("client_metadata", client.Metadata)

	lg.Printf("metadata: %T %v", client.Metadata, client.Metadata)

	// NOTE: client secret is never returned from read operations, so don't set this field.

	d.Set("scopes", strings.Split(client.Scope, " "))
	d.Set("grant_types", client.GrantTypes)
	d.Set("response_types", client.ResponseTypes)
	d.Set("audience", client.Audience)
	d.Set("post_logout_redirect_uris", client.PostLogoutRedirectUris)
	d.Set("redirect_uris", client.RedirectUris)
	d.Set("owner", client.Owner)
	d.Set("policy_uri", client.PolicyURI)
	d.Set("allowed_cors_origins", client.AllowedCorsOrigins)
	d.Set("tos_uri", client.TosURI)
	d.Set("client_uri", client.ClientURI)
	d.Set("logo_uri", client.LogoURI)
	d.Set("contacts", client.Contacts)
	d.Set("subject_type", client.SubjectType)
	d.Set("token_endpoint_auth_method", client.TokenEndpointAuthMethod)
}
