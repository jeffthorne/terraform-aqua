package aqua

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	aquaSDK "github.com/jeffthorne/aqua-go/aqua"
	"log"
	"strings"
)



func resourceRegistry() *schema.Resource{
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Friendly name for resource in Aqua Enterprise",
			},
			"description": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Registry description",
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Registry Type: AWS, ",
			},
			"url": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Registry Type: AWS, ",
			},
			"username": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "registry authentication",
			},
			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "registry authentication",
			},
			"prefixes": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "registry authentication",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"auto_pull": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "registry authentication",
			},
			"auto_pull_max": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "registry authentication",
			},
			"auto_pull_time": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:  	 "",
				Description: "registry authentication",
			},
		},
		CreateContext: resourceCreateRegistry,
		UpdateContext: resourceUpdateRegistry,
		ReadContext: resourceGetRegistry,
		DeleteContext: resourceDeleteRegistry,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceCreateRegistry(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics
	aquaClient := m.(*aquaSDK.Aqua)
	name := d.Get("name").(string)
	description:= d.Get("description").(string)
	registryType := d.Get("type").(string)
	url := d.Get("url").(string)
	prefixes:= d.Get("prefixes").([]interface{})
	username := d.Get("username").(string)
	password := d.Get("password").(string)
	autoPull := d.Get("auto_pull").(bool)
	autoPullMax := 0
	autoPullTime := ""
	log.Printf("[DEBUG] In CREATE Registry Name: %s", name)

	if autoPull{
		autoPullMax = d.Get("auto_pull_max").(int)
		autoPullTime = d.Get("auto_pull_time ").(string);
	}


	err := aquaClient.CreateRegistry(name, description, registryType,url,username, password, convertIterface(prefixes), autoPull, int64(autoPullMax), autoPullTime)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(name)
	return diags
}

func resourceDeleteRegistry(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	aquaClient := m.(*aquaSDK.Aqua)
	itemId := d.Id()
	err := aquaClient.DeleteRegistry(itemId)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return diags
}


func resourceUpdateRegistry(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	aquaClient := m.(*aquaSDK.Aqua)

	id := d.Id()
	name := d.Get("name").(string)
	log.Printf("[DEBUG] Name: %s, ID: %s", name, id)

	if name != id {
		err := aquaClient.DeleteRegistry(id)
		log.Printf("[DEBUG] In UPDATE deleteing current registry %s: ", id)

		if err != nil {
			log.Println(err)
		}else{
			d.SetId(name)
		}

	}



	registryType := d.Get("type").(string)
	url := d.Get("url").(string)
	prefixes:= d.Get("prefixes").([]interface{})
	username := d.Get("username").(string)
	password := d.Get("password").(string)
	autoPull := d.Get("auto_pull").(bool)
	autoPullMax := d.Get("auto_pull_max").(int)
	autoPullTime := d.Get("auto_pull_time").(string)
	description := d.Get("description").(string)
	log.Printf("[DEBUG] In UPDATE re-creating current registry*************************************************************** NAme %s:",name)

	if name != id {
		aquaClient.CreateRegistry(name, description, registryType,url,username, password, convertIterface(prefixes), autoPull, int64(autoPullMax), autoPullTime)

	}else {
		aquaClient.UpdateRegistry(name, description, registryType, url, username, password, convertIterface(prefixes), autoPull, int64(autoPullMax), autoPullTime)
	}

	return diags
}

func resourceGetRegistry(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	aquaClient := m.(*aquaSDK.Aqua)
	itemId := d.Id()
	registry, err := aquaClient.GetRegistry(itemId)

	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			d.SetId("")
		} else {
			return diag.FromErr(err)//fmt.Errorf("error finding Registry with name %s", itemId)
		}
	}

	d.SetId(registry.Name)
	d.Set("name", registry.Name)
	d.Set("type", registry.Type)
	d.Set("url", registry.URL)
	d.Set("prefixes", registry.Prefixes)
	d.Set("username", registry.Username)
	d.Set("password", registry.Password)
	d.Set("auto_pull", registry.AutoPull)
	d.Set("auto_pull_max", registry.AutoPullMax)
	d.Set("auto_pull_time", registry.AutoPullTime)
	d.Set("description", registry.Description)
	return diags
}

func convertIterface(prefixes []interface{}) []string {
	s := make([]string, len(prefixes))

	for i, v := range prefixes {
		s[i] = fmt.Sprint(v)
	}

	return s
}
