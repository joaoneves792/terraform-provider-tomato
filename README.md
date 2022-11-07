# Terraform Provider Tomato

This repo contains a provider for manipulating Tomato router settings through terraform.

This was tested only on FreshTomato instalations on an ASUS RT-N12 and an ASUS RT-AC68U.

To use this change `OS_ARCH=linux_amd64` to your architecture and run `make` (alternatively run `make release` to build for all architectures)

In your terraform code include the provider as:

```
terraform {
  required_providers {
    tomato = {
      version = ">=0.0.3"
      source  = "warpenguin.dev/tomato/tomato"
    }
  }
}

provider "tomato" {
  url = "https://10.0.0.1"
  #Provide username and password with TOMATO_USERNAME and TOMATO_PASSWORD env variables, or alternatively set them here
  username = "xxxxx" 
  password = "xxxxx"
}

```

This provider exposes the following datasources and resources (see examples in ./examples/main.tf):


---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "tomato Provider"
subcategory: ""
description: |-
  
---

# tomato Provider





<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `password` (String, Sensitive)
- `url` (String, Sensitive)
- `username` (String)
---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "tomato_nvram Data Source - terraform-provider-tomato"
subcategory: ""
description: |-
  
---

# tomato_nvram (Data Source)





<!-- schema generated by tfplugindocs -->
## Schema

### Read-Only

- `id` (String) The ID of this resource.
- `nvram` (Map of String)


---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "tomato_static_ip Resource - terraform-provider-tomato"
subcategory: ""
description: |-
  
---

# tomato_static_ip (Resource)





<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `ip` (String)
- `mac` (String)

### Optional

- `bind` (Boolean)
- `hostname` (String)
- `mac2` (String)

### Read-Only

- `id` (String) The ID of this resource.


---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "tomato_dns_entry Resource - terraform-provider-tomato"
subcategory: ""
description: |-
  
---

# tomato_dns_entry (Resource)





<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String)
- `record` (String)

### Read-Only

- `id` (String) The ID of this resource.


---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "tomato_generic Resource - terraform-provider-tomato"
subcategory: ""
description: |-
  
---

# tomato_generic (Resource)





<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `key` (String)
- `services` (List of String)
- `value` (String)

### Read-Only

- `id` (String) The ID of this resource.


