package tomato

import (
	"context"
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"sync"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var lock = &sync.Mutex{}

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
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceDNSEntryCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	lock.Lock()
	defer lock.Unlock()

	n, err := c.getNVRAM()
	if err != nil {
		return diag.FromErr(err)
	}

	dnsmasq_custom := n["dnsmasq_custom"]

	name := d.Get("name").(string)
	record := d.Get("record").(string)

	d.SetId(name)

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

	name := recordID

	n, err := c.getNVRAM()
	if err != nil {
		return diag.FromErr(err)
	}

	dnsmasq_custom := n["dnsmasq_custom"]

	entry, name, record := findEntry(name, dnsmasq_custom)
	if err := d.Set("name", name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("record", record); err != nil {
		return diag.FromErr(err)
	}

	if len(entry) == 0 {
		d.SetId("")
	}
	return diags
}

func findEntry(name, dnsmasq_config string) (string, string, string) {
	re := regexp.MustCompile(`(?m)address=/([a-zA-Z-\.]+)/([0-9]+\.[0-9]+\.[0-9]+\.[0-9]+)$`)

	matches := re.FindAllStringSubmatch(dnsmasq_config, -1)
	for i := range matches {
		if matches[i][1] == name {
			return matches[i][0], matches[i][1], matches[i][2]
		}
	}
	return "", "", ""
}

func resourceDNSEntryUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	c := m.(*Client)

	lock.Lock()
	defer lock.Unlock()

	oname := d.Get("name").(string)
	orecord := d.Get("record").(string)

	n, err := c.getNVRAM()
	if err != nil {
		return diag.FromErr(err)
	}

	dnsmasq_custom := n["dnsmasq_custom"]

	eentry, ename, erecord := findEntry(oname, dnsmasq_custom)

	//Nothing has changed
	if oname == ename && orecord == erecord {
		return resourceDNSEntryRead(ctx, d, m)
	}

	entry := fmt.Sprintf("address=/%s/%s", oname, orecord)

	dnsmasq_custom = strings.Replace(dnsmasq_custom, eentry, entry, -1)

	dnsconfig := url.QueryEscape(dnsmasq_custom)

	tflog.Debug(ctx, "Apply dnsconfig:\n"+dnsconfig)

	b, err := c.applyChange("dnsmasq-restart", "dnsmasq_custom="+dnsconfig)

	tflog.Debug(ctx, b)

	if err != nil {
		return diag.FromErr(err)
	}

	return resourceDNSEntryRead(ctx, d, m)
}

func resourceDNSEntryDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	id := d.Id()
	c := m.(*Client)

	lock.Lock()
	defer lock.Unlock()

	n, err := c.getNVRAM()
	if err != nil {
		return diag.FromErr(err)
	}

	dnsmasq_custom := n["dnsmasq_custom"]

	entry, _, _ := findEntry(id, dnsmasq_custom)

	if len(entry) == 0 {
		return diags
	}

	dnsmasq_custom = strings.Replace(dnsmasq_custom, entry, "", -1)

	dnsconfig := url.QueryEscape(dnsmasq_custom)

	tflog.Debug(ctx, "Apply dnsconfig:\n"+dnsconfig)

	b, err := c.applyChange("dnsmasq-restart", "dnsmasq_custom="+dnsconfig)

	tflog.Debug(ctx, b)

	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}
