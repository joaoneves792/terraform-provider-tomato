#! /bin/bash
cat > README.md << "EOT"
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


EOT

for file in $(find ./docs/ -name '*.md'); do 
  cat $file | tail --lines +8 >> ./README.md
done
