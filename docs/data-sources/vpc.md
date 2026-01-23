---
subcategory: "VPC (Virtual Private Cloud)"
layout: "aws"
page_title: "aws_vpc"
description: |-
  Provides information about a VPC.
---

[describe-vpcs]: https://docs.k2.cloud/en/api/ec2/actions/vpcs/DescribeVpcs.html

# Data Source: aws_vpc

Provides information about a VPC.

This resource can be useful when a module accepts the ID of a VPC as
an input variable and needs to, for example, determine the CIDR block of that VPC.

## Example Usage

The following example shows how one might accept the ID of a VPC as a variable
and use this data source to obtain the data necessary to create a subnet
within it.

```terraform
variable "vpc_id" {}

data "aws_vpc" "selected" {
  id = var.vpc_id
}

resource "aws_subnet" "example" {
  vpc_id            = data.aws_vpc.selected.id
  availability_zone = "ru-msk-vol52"
  cidr_block        = cidrsubnet(data.aws_vpc.selected.cidr_block, 4, 1)
}
```

## Argument Reference

The arguments of this data source act as filters for querying the available
VPCs in the current region. The given filters must match exactly one
VPC whose data will be exported as attributes.

* `cidr_block` - (Optional) The CIDR block of the desired VPC.
* `dhcp_options_id` - (Optional) The ID of the DHCP options for the desired VPC.
* `filter` - (Optional) One or more name/value pairs to use as filters.
    * _Valid values:_ See supported names and values in [EC2 API documentation][describe-vpcs]
* `id` - (Optional) The ID of the specific VPC to retrieve.
* `state` - (Optional) The current state of the desired VPC.
  Can be either `pending` or `available`.
* `tags` - (Optional) Map of tags, each pair of which must exactly match
  a pair on the desired VPC.

## Attribute Reference

### Supported attributes

All argument attributes except `filter` blocks are also exported as
result attributes. This data source will complete the data by populating
any fields that are not included in the configuration with the data for
the selected VPC.

In addition to all arguments above, the following attributes are exported:

* `arn` - The Amazon Resource Name (ARN) of VPC.
* `enable_dns_support` - Whether the VPC has DNS support.
* `main_route_table_id` - ID of the main route table associated with this VPC.

`cidr_block_associations` is also exported with the following attributes:

* `association_id` - The association ID for the IPv4 CIDR block.
* `cidr_block` - The CIDR block for the association.
* `state` - The state of the association.

### Unsupported attributes

~> **Note** These attributes may be present in the `terraform.tfstate` file, but they have preset values and cannot be specified in configuration files.

The following attributes are not currently supported:

`enable_dns_hostnames`, `instance_tenancy`, `ipv6_association_id`, `ipv6_cidr_block`, `owner_id`.
