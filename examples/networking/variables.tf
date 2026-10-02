variable "vpc_cidr_blocks" {
  description = "The CIDR block of the VPC in every region."
  type        = map(string)
  default = {
    ru-msk = "10.0.0.0/16"
    ru-spb = "10.1.0.0/16"
  }
}
