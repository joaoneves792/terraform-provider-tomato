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

#output "dnsmasq_config" {
#  value = data.tomato_nvram.nvram.nvram.dnsmasq_custom
#}

output "rebind" {
  value = data.tomato_nvram.nvram.nvram.dns_norebind
}

resource "tomato_dns_entry" "batatas" {
  name = "batatas"
  record = "127.0.0.2"
}
resource "tomato_dns_entry" "cebola" {
  name = "cebola"
  record = "8.8.6.4"
}

