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

# Returns all coffees
output "dnsmasq_config" {
  value = data.tomato_nvram.nvram.nvram.dnsmasq_custom
}


