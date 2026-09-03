variable "region" {
  description = "The region to set up a network within."
  type        = string
}

variable "vpc_cidr" {
  description = "The CIDR block of the VPC."
  type        = string
}

variable "access_key" {
  description = "The access key taken from c2rc.sh."
  type        = string
  sensitive   = true
}

variable "secret_key" {
  description = "The secret key taken from c2rc.sh."
  type        = string
  sensitive   = true
}
