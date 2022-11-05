package tomato

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceDNSEntry() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDNSEntryCreate,
		ReadContext:   resourceDNSEntryRead,
		UpdateContext: resourceDNSEntryUpdate,
		DeleteContext: resourceDNSEntryDelete,
		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"record": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceDNSEntryCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	n, err := c.getNVRAM()
	if err != nil {
		return diag.FromErr(err)
	}

	dnsmasq_custom := n["dnsmasq_custom"]

	name := d.Get("name").(string)
	record := d.Get("record").(string)

	sEnc := base64.StdEncoding.EncodeToString([]byte(name + "/" + record))

	d.SetId(sEnc)

	entry := fmt.Sprintf("address=/%s/%s", name, record)

	dnsconfig := url.QueryEscape(fmt.Sprintf("%s\n%s", dnsmasq_custom, entry))

	tflog.Debug(ctx, "Apply dnsconfig:\n"+dnsconfig)

	b, err := c.applyChange("dnsmasq-restart", "dnsmasq_custom="+dnsconfig)

	tflog.Debug(ctx, b)

	if err != nil {
		return diag.FromErr(err)
	}

	resourceDNSEntryRead(ctx, d, m)

	return diags
}

func resourceDNSEntryRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	recordID := d.Id()

	tflog.Info(ctx, "Read ID: "+recordID)

	sDec, err := base64.StdEncoding.DecodeString(recordID)
	if err != nil {
		return diag.FromErr(err)
	}

	sp := strings.Split(string(sDec), "/")
	name := sp[0]

	n, err := c.getNVRAM()
	if err != nil {
		return diag.FromErr(err)
	}

	dnsmasq_custom := n["dnsmasq_custom"]

	re := regexp.MustCompile(`(?m)address=/([a-zA-Z-\.]+)/([0-9]+\.[0-9]+\.[0-9]+\.[0-9]+)$`)

	matches := re.FindAllStringSubmatch(dnsmasq_custom, -1)
	for i := range matches {
		tflog.Info(ctx, "Found: "+matches[i][0])
		if matches[i][1] == name {
			if err := d.Set("name", name); err != nil {
				return diag.FromErr(err)
			}
			if err := d.Set("record", matches[i][2]); err != nil {
				return diag.FromErr(err)
			}
			return diags
		}
	}

	if err := d.Set("name", ""); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("record", ""); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}

func resourceDNSEntryUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceDNSEntryRead(ctx, d, m)
}

func resourceDNSEntryDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	return diags
}
