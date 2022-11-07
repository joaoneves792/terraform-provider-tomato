package tomato

import (
	"context"
	"errors"
	"strings"
	"sync"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var genericLock = &sync.Mutex{}

func resourceGeneric() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGenericCreate,
		ReadContext:   resourceGenericRead,
		UpdateContext: resourceGenericUpdate,
		DeleteContext: resourceGenericDelete,
		Schema: map[string]*schema.Schema{
			"key": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"value": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"services": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceGenericCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	genericLock.Lock()
	defer genericLock.Unlock()

	key := d.Get("key").(string)
	d.SetId(key)

	value := d.Get("value").(string)
	services := d.Get("services").(([]interface{}))

	services_restart := ""
	for i := range services {
		services_restart = services_restart + "%2C" + services[i].(string)
	}
	services_restart = strings.Replace(services_restart, "%2C", "", 1)

	b, err := c.applyChange(services_restart, key+"="+value)

	tflog.Debug(ctx, b)

	if err != nil {
		return diag.FromErr(err)
	}

	resourceGenericRead(ctx, d, m)

	return diags
}

func resourceGenericRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	key := d.Id()

	n, err := c.getNVRAM()
	if err != nil {
		return diag.FromErr(err)
	}

	value, found := n[key]
	if !found {
		d.SetId("")
		return diag.FromErr(errors.New("key:" + key + " not found in NVRAM"))
	}

	if err := d.Set("key", key); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("value", value); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(key)
	return diags
}

func resourceGenericUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	resourceGenericCreate(ctx, d, m)
	return resourceGenericRead(ctx, d, m)
}

func resourceGenericDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	return diags
}
