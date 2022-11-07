package tomato

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"sync"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var staticIpLock = &sync.Mutex{}

func resourceStaticIp() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceStaticIpCreate,
		ReadContext:   resourceStaticIpRead,
		UpdateContext: resourceStaticIpUpdate,
		DeleteContext: resourceStaticIpDelete,
		Schema: map[string]*schema.Schema{
			"mac": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"mac2": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"hostname": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"ip": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"bind": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceStaticIpCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	staticIpLock.Lock()
	defer staticIpLock.Unlock()

	n, err := c.getNVRAM()
	if err != nil {
		return diag.FromErr(err)
	}

	dhcp_static, err := url.QueryUnescape(n["dhcpd_static"])
	if err != nil {
		return diag.FromErr(err)
	}

	ip := d.Get("ip").(string)
	mac := d.Get("mac").(string)
	mac2 := d.Get("mac2").(string)
	hostname := d.Get("hostname").(string)
	bind := d.Get("bind").(bool)

	d.SetId(mac)

	nr := 0
	if bind {
		nr = 1
	}

	entry := ""
	if mac2 != "" {
		entry = fmt.Sprintf("%s,%s<%s<%s<%d>", mac, mac2, ip, hostname, nr)
	} else {
		entry = fmt.Sprintf("%s<%s<%s<%d>", mac, ip, hostname, nr)
	}

	tflog.Debug(ctx, "Entry dhcpconfig:\n"+entry)

	dhcpconfig := url.QueryEscape(fmt.Sprintf("%s%s", dhcp_static, entry))

	tflog.Debug(ctx, "Apply dhcpconfig:\n"+dhcpconfig)

	b, err := c.applyChange("dhcpd-restart%2Carpbind-restart%2Ccstats-restart%2Cdnsmasq-restart", "dhcpd_static="+dhcpconfig)

	tflog.Debug(ctx, b)

	if err != nil {
		return diag.FromErr(err)
	}

	resourceStaticIpRead(ctx, d, m)

	return diags
}

func resourceStaticIpRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	mac := d.Id()

	n, err := c.getNVRAM()
	if err != nil {
		return diag.FromErr(err)
	}

	dhcp_static, err := url.QueryUnescape(n["dhcpd_static"])
	if err != nil {
		return diag.FromErr(err)
	}

	_, emac, mac2, ip, hostname, options := staticIpFindEntry(mac, dhcp_static)
	if err := d.Set("mac", emac); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("mac2", mac2); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("ip", ip); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("hostname", hostname); err != nil {
		return diag.FromErr(err)
	}
	if options == "0" {
		if err := d.Set("bind", false); err != nil {
			return diag.FromErr(err)
		}
	} else if options == "1" {
		if err := d.Set("bind", true); err != nil {
			return diag.FromErr(err)
		}
	}

	if len(emac) == 0 {
		d.SetId("")
	} else {
		d.SetId(emac)
	}
	return diags
}

func staticIpFindEntry(smac, dhcp_static string) (string, string, string, string, string, string) {
	var re = regexp.MustCompile(`((?:[[:alnum:]]{2}:){5}[[:alnum:]]{2})(,((?:[[:alnum:]]{2}:){5}[[:alnum:]]{2}))?<([0-9]+\.[0-9]+\.[0-9]+\.[0-9]+)<([^<]+)<([0-2])>`)
	const (
		mac      = 1
		mac2     = 3
		ip       = 4
		hostname = 5
		bind     = 6
	)

	matches := re.FindAllStringSubmatch(dhcp_static, -1)
	for i := range matches {
		if matches[i][mac] == smac {
			return matches[i][0], matches[i][mac], matches[i][mac2], matches[i][ip], matches[i][hostname], matches[i][bind]
		}
	}
	return "", "", "", "", "", ""
}

func resourceStaticIpUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	c := m.(*Client)

	staticIpLock.Lock()
	defer staticIpLock.Unlock()

	ip := d.Get("ip").(string)
	mac := d.Get("mac").(string)
	mac2 := d.Get("mac2").(string)
	hostname := d.Get("hostname").(string)
	bind := d.Get("bind").(bool)

	n, err := c.getNVRAM()
	if err != nil {
		return diag.FromErr(err)
	}

	dhcpd_static, err := url.QueryUnescape(n["dhcpd_static"])
	if err != nil {
		return diag.FromErr(err)
	}

	eentry, _, _, _, _, _ := staticIpFindEntry(d.Id(), dhcpd_static)

	if eentry == "" {
		return diag.FromErr(errors.New("ID Not Found"))
	}

	tflog.Debug(ctx, "Existing entry:\n"+eentry)

	nr := 0
	if bind {
		nr = 1
	}

	entry := ""
	if mac2 != "" {
		entry = fmt.Sprintf("%s,%s<%s<%s<%d>", mac, mac2, ip, hostname, nr)
	} else {
		entry = fmt.Sprintf("%s<%s<%s<%d>", mac, ip, hostname, nr)
	}

	dhcpd_static = strings.Replace(dhcpd_static, eentry, entry, -1)

	dhcpconfig := url.QueryEscape(dhcpd_static)

	tflog.Debug(ctx, "Apply dnsconfig:\n"+dhcpd_static)

	b, err := c.applyChange("dhcpd-restart%2Carpbind-restart%2Ccstats-restart%2Cdnsmasq-restart", "dhcpd_static="+dhcpconfig)

	tflog.Debug(ctx, b)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(mac)
	return resourceStaticIpRead(ctx, d, m)
}

func resourceStaticIpDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	c := m.(*Client)

	staticIpLock.Lock()
	defer staticIpLock.Unlock()

	n, err := c.getNVRAM()
	if err != nil {
		return diag.FromErr(err)
	}

	dhcpd_static, err := url.QueryUnescape(n["dhcpd_static"])
	if err != nil {
		return diag.FromErr(err)
	}

	entry, _, _, _, _, _ := staticIpFindEntry(d.Id(), dhcpd_static)

	if len(entry) == 0 {
		return diags
	}

	dhcpd_static = strings.Replace(dhcpd_static, entry, "", -1)

	dhcpconfig := url.QueryEscape(dhcpd_static)

	tflog.Debug(ctx, "Apply dhcpconfig:\n"+dhcpconfig)

	b, err := c.applyChange("dhcpd-restart%2Carpbind-restart%2Ccstats-restart%2Cdnsmasq-restart", "dhcpd_static="+dhcpconfig)

	tflog.Debug(ctx, b)

	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}
