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

#output "static" {
#  value = data.tomato_nvram.nvram.nvram.dhcpd_static
#}

/*resource "tomato_dns_entry" "batatas" {
  name = "batatas"
  record = "127.0.0.2"
}
resource "tomato_dns_entry" "cebola" {
  name = "cebola"
  record = "8.8.6.4"
}
*/
/*resource "tomato_static_ip" "desktop"{
  ip = "10.6.4.6"
  hostname = "Lumber"
  mac = "56:56:56:56:56:56"
  mac2 = "6C:60:6D:67:68:65"
  bind = false
}*/
