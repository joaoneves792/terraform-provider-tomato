terraform {
  required_providers {
    tomato = {
      version = ">=0.0.1"
      source  = "warpenguin.dev/tomato/tomato"
    }
  }
}

provider "tomato" {
  #Provide username and password with TOMATO_USERNAME and TOMATO_PASSWORD env variables
  url = "https://10.0.0.1"
}

data "tomato_nvram" "nvram" {}

output "dnsmasq_config" {
  value = data.tomato_nvram.nvram.nvram.dnsmasq_custom
}

output "rebind" {
  value = data.tomato_nvram.nvram.nvram.dns_norebind
}

resource "tomato_nvram" "dnsmasq" {
  items {
    service = "dnsmasq-restart"
    nvram = {
      dns_norebind = 0
    }
  }
}

