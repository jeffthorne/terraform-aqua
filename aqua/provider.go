package aqua

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	aquaSDK "github.com/jeffthorne/aqua-go/aqua"
)


type AquaOpts struct {
	User string
	Password string
	Host string
	Port int
	Secure bool
	Verify bool
}



// Provider returns a terraform aqua for Aqua Enterprise
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"user": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("AQUA_USER", nil),
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("AQUA_PASSWORD", nil),
			},
			"host": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("AQUA_HOST", nil),
			},
			"port": {
				Type:        schema.TypeInt,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("AQUA_PORT", nil),
			},
			"secure": {
				Type:        schema.TypeBool,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("AQUA_SECURE", nil),
			},
			"verify": {
				Type:        schema.TypeBool,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("AQUA_VERIFY", nil),
			},
		},
		ResourcesMap:         map[string]*schema.Resource{
			"aqua_create_registry": resourceRegistry(),
		},
		ConfigureContextFunc:  providerConfigure,
	}
}


func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {

	var aquaAPI *aquaSDK.Aqua
	var err error
	user := d.Get("user").(string)
	password := d.Get("password").(string)
	host := d.Get("host").(string)
	port := d.Get("port").(int)
	secure := d.Get("secure").(bool)
	verify := d.Get("verify").(bool)


	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	if (user != "") && (password != "") {
		aquaAPI, err = aquaSDK.NewCSP(host, port, user, password, secure, verify)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Unable to create Aqua client",
				Detail:   err.Error(),
			})

			return nil, diags
		}

		return aquaAPI, diags
	}

	return aquaAPI, diags
}
