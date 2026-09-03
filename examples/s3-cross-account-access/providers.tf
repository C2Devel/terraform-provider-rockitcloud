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

# The example works with two accounts, so it needs two credentials files.
data "local_file" "first_credentials" {
  filename = "./c2rc-first.sh"
}

data "local_file" "second_credentials" {
  filename = "./c2rc-second.sh"
}

locals {
  first_creds = {
    for i in split("\n", data.local_file.first_credentials.content) :
    split("=", replace(replace(i, "export ", ""), "\"", ""))[0] =>
    split("=", replace(replace(i, "export ", ""), "\"", ""))[1]
    if can(regex("export", i))
  }

  second_creds = {
    for i in split("\n", data.local_file.second_credentials.content) :
    split("=", replace(replace(i, "export ", ""), "\"", ""))[0] =>
    split("=", replace(replace(i, "export ", ""), "\"", ""))[1]
    if can(regex("export", i))
  }
}

provider "aws" {
  alias = "first"

  access_key = "${local.first_creds["C2_PROJECT"]}:${local.first_creds["BASE_ACCESS_KEY"]}"
  secret_key = local.first_creds["EC2_SECRET_KEY"]

  # For K2 Cloud, specify one of the supported regions.
  # For other cloud platforms, enter a non-empty string,
  # for example, "region-1", and API endpoints.
  region = var.region
}

provider "aws" {
  alias = "second"

  access_key = "${local.second_creds["C2_PROJECT"]}:${local.second_creds["BASE_ACCESS_KEY"]}"
  secret_key = local.second_creds["EC2_SECRET_KEY"]

  # For K2 Cloud, specify one of the supported regions.
  # For other cloud platforms, enter a non-empty string,
  # for example, "region-1", and API endpoints.
  region = var.region
}
