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

This data source can be used when a module accepts the ID of a VPC as an input variable and needs to, for example, determine the CIDR block of that VPC.

## Example usage

### Specific example

The following example shows how to accept the ID of a VPC as a variable and use this data source to obtain the data necessary to create a subnet within it.

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

## Argument reference

The arguments of this data source act as filters for querying available VPCs in the current region.

~> **Note** The given filters must exactly match the resource whose data will be exported as attributes.

The following arguments are supported:

* `cidr_block` - (Optional, String) The CIDR block of the desired VPC.
* `dhcp_options_id` - (Optional, String) The ID of the DHCP options set for the desired VPC.
* `filter` - (Optional, [Block](#filter)) One or more name/value pairs to use as filters.
    * _Valid values:_ See supported names and values in [EC2 API documentation][describe-vpcs]
* `id` - (Optional, String) The ID of the specific VPC to retrieve.
* `state` - (Optional, String) The current state of the desired VPC.
    * _Valid values:_ `available`, `pending`
* `tags` - (Optional, Map of strings) Key-value pairs.
    Must exactly match pairs on the required resource.

### filter

* `name` - (Required, String) The name of the filter.
    * _Constraints:_ Filter names are case-sensitive
* `values` - (Required, List of strings) One or more filter values.
    * _Constraints:_ Filter values are case-sensitive

## Attribute reference

### Supported attributes

This data source will complete the data by populating any fields that are not included in the configuration with the data for the selected VPC.

In addition to all arguments above, the following attributes are exported:

* `arn` - (String) The Amazon Resource Name (ARN) of VPC.
* `cidr_block_associations` - ([Block](#cidr_block_associations)) The IPv4 CIDR blocks associated with the VPC.
* `enable_dns_support` - (Boolean) Indicates whether the VPC has DNS support.
* `main_route_table_id` - (String) The ID of the main route table associated with this VPC.

#### cidr_block_associations

* `association_id` - (String) The association ID for the IPv4 CIDR block.
* `cidr_block` - (String) The CIDR block for the association.
* `state` - (String) The state of the association.

### Unsupported attributes

~> **Note** These attributes may be present in the `terraform.tfstate` file, but they have preset values and cannot be specified in configuration files.

The following attributes are not currently supported:

`enable_dns_hostnames`, `instance_tenancy`, `ipv6_association_id`, `ipv6_cidr_block`, `owner_id`.
