terraform {
  required_version = ">= 1.11.0, < 2.0.0"
  required_providers {
    aws = {
      source  = "hc-registry.website.k2.cloud/c2devel/rockitcloud"
      version = ">= 25.5.5"
    }
    local = {
      source  = "hc-registry.website.k2.cloud/hashicorp/local"
      version = "~> 2.5"
    }
    http = {
      source  = "hc-registry.website.k2.cloud/hashicorp/http"
      version = "~> 3.4"
    }
  }
}

data "local_file" "credentials" {
  filename = "./c2rc.sh"
}

locals {
  loc_creds = [
    for i in split("\n", data.local_file.credentials.content) :
    split("=", replace(replace(i, "export ", ""), "\"", ""))
    if can(regex("export", i))
  ]
  dict_creds = { for i in local.loc_creds : i[0] => i[1] }
}

provider "aws" {
  access_key = "${local.dict_creds["C2_PROJECT"]}:${local.dict_creds["BASE_ACCESS_KEY"]}"
  secret_key = local.dict_creds["EC2_SECRET_KEY"]

  # For K2 Cloud, specify one of the supported regions.
  # For other cloud platforms, enter a non-empty string,
  # for example, "region-1", and API endpoints.
  region = var.region
}

# Not required: currently used in conjunction with using
# icanhazip.com to determine local workstation external IP
# to open EC2 Security Group access to the Kubernetes cluster.
# See workstation-external-ip.tf for additional information.
provider "http" {}
