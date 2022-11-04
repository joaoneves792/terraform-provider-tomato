package tomato

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/hex"
	"strconv"
	"time"

  "github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

)

func dataSourceNVRAM() *schema.Resource {
	return &schema.Resource{
		ReadContext: nvramRead,
		Schema: map[string]*schema.Schema{
			"nvram": &schema.Schema{
				Type:     schema.TypeMap,
				Computed: true,
			},
		},
	}
}

func nvramRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	b, err := c.getNVRAM()
	if err != nil {
		return diag.FromErr(err)
	}

	pad := make([]byte, 1)
	pad[0] = 0x00
  
  tflog.Info(ctx,"NVRAM 8 byte Header: " + hex.EncodeToString(b[0:8]))

	length := int(binary.BigEndian.Uint32(append(pad, b[4:7]...)))
	magic := b[7]
  tflog.Info(ctx, "NVRAM dump length: " + strconv.Itoa(length))

	for i := 8; i < length; i++ {
		if b[i] > (0xfd - 0x1) {
			b[i] = 0x0
		} else {
			b[i] = 0xff + magic - b[i]
		}
	}

	n := make(map[string]string)

	se := 8
	for se < length+8 {
		nb := bytes.IndexByte(b[se:length+8], 0x00)
		cfg := string(b[se : se+nb+1])
		se = se + nb + 1
		eq := bytes.Index([]byte(cfg), []byte("="))
		if eq == -1 {
			break
		}
		n[string(cfg[0:eq])] = string(cfg[eq+1:])
	}

	//	nvram := make(map[string]interface{}, 0)
	//	nvram["nvram"] = r
	if err := d.Set("nvram", n); err != nil {
		return diag.FromErr(err)
	}

	// always run
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}
