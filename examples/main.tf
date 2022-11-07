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

#output "dnsrebind" {
#  value = data.tomato_nvram.nvram.nvram.dns_norebind
#}

#resource "tomato_generic" "norebind" {
#  key = "dns_norebind"
#  value = "0"
#  services = ["dnsmasq-restart"]
#}


#output "static" {
#  value = data.tomato_nvram.nvram.nvram.dhcpd_static
#}

#resource "tomato_dns_entry" "potato" {
#  name = "potato.com"
#  record = "127.0.0.1"
#}



#resource "tomato_static_ip" "desktop"{
#  ip = "10.6.4.6"
#  hostname = "Desktop"
#  mac = "56:56:56:56:56:56"
#  mac2 = "6C:60:6D:67:68:65"
#  bind = false
#}
