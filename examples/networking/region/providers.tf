terraform {
  required_version = ">= 1.11.0, < 2.0.0"
  required_providers {
    aws = {
      source  = "hc-registry.website.k2.cloud/c2devel/rockitcloud"
      version = ">= 25.5.5"
    }
  }
}

provider "aws" {
  access_key = var.access_key
  secret_key = var.secret_key

  # For K2 Cloud, specify one of the supported regions.
  # For other cloud platforms, enter a non-empty string,
  # for example, "region-1", and API endpoints.
  region = var.region
}
