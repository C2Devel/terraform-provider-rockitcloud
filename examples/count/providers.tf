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
