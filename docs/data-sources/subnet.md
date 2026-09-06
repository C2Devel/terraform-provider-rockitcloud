---
subcategory: "VPC (Virtual Private Cloud)"
layout: "aws"
page_title: "aws_subnet"
description: |-
  Provides information about a subnet.
---

[describe-subnets]: https://docs.k2.cloud/en/api/ec2/actions/subnets/DescribeSubnets.html

# Data Source: aws_subnet

Provides information about a subnet.

This data source can be used when a module accepts a subnet ID as an input variable and needs to, for example, determine the ID of the VPC that owns the subnet.

## Example usage

### Specific examples

The following example shows how to accept a subnet ID as a variable and use this data source to obtain the data necessary to create a security group that allows connections from hosts in that subnet.

```terraform
variable "subnet_id" {}

data "aws_subnet" "selected" {
  id = var.subnet_id
}

resource "aws_security_group" "subnet" {
  vpc_id = data.aws_subnet.selected.vpc_id

  ingress {
    cidr_blocks = [data.aws_subnet.selected.cidr_block]
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
  }
}
```

If you want to match against tag `Name`, use:

```terraform
data "aws_subnet" "selected" {
  filter {
    name   = "tag:Name"
    values = ["yakdriver"]
  }
}
```

## Argument reference

The arguments of this data source act as filters for querying the available subnets in the current project within the current region.

~> **Note** The given filters must exactly match the resource whose data will be exported as attributes.

The following arguments are optional:

* `availability_zone` - (Optional, String) The availability zone where the subnet must reside.
* `default_for_az` - (Optional, Boolean) Indicates whether the desired subnet must be the default subnet for its associated availability zone.
    * _Default value_: `false`

    ~> **Note** Setting this argument to `false` does not actually filter by it *not* being the default, because Terraform cannot distinguish between `false` and "not set".

* `filter` - (Optional, [Block](#filter)) One or more name/value pairs to use as filters.
    * _Valid values:_ See supported names and values in [EC2 API documentation][describe-subnets]
* `tags` - (Optional, Map of strings) Key-value pairs. Must exactly match pairs on the required resource.
* `vpc_id` - (Optional, String) The ID of the VPC that owns the desired subnet.

### filter

* `name` - (Required, String) The name of the filter.
    * _Constraints:_ Filter names are case-sensitive
* `values` - (Required, List of strings) One or more filter values.
    * _Constraints:_ Filter values are case-sensitive

## Attribute reference

### Supported attributes

In addition to all arguments above, the following attribute is exported:

* `available_ip_address_count` - (Integer) The number of remaining available IPv4 addresses in the subnet.
* `cidr_block` - (String) The IPv4 CIDR block for the subnet.
* `map_public_ip_on_launch` - (Boolean, **Deprecated**) Indicates whether public IP addresses will be associated with instances created in this subnet. Addresses are associated only if there are available allocated Elastic IP addresses.
* `state` - (String) The state that the desired subnet must have.

~> **Note** The attribute `map_public_ip_on_launch` is deprecated and will be removed in a future release.

### Unsupported attributes

~> **Note** These attributes may be present in the `terraform.tfstate` file, but they have preset values and cannot be specified in configuration files.

The following attributes are not currently supported:

`arn`, `assign_ipv6_address_on_creation`, `availability_zone_id`, `customer_owned_ipv4_pool`, `enable_dns64`, `enable_resource_name_dns_aaaa_record_on_launch`, `enable_resource_name_dns_a_record_on_launch`, `ipv6_cidr_block`, `ipv6_cidr_block_association_id`, `ipv6_native`, `map_customer_owned_ip_on_launch`, `outpost_arn`, `owner_id`, `private_dns_hostname_type_on_launch`.
