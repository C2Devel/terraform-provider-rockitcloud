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

This resource can be used when a module accepts a subnet ID as an input variable and needs to, for example, determine the ID of the VPC that the subnet belongs to.

## Example Usage

The following example shows how one might accept a subnet ID as a variable and use this data source to obtain the data necessary to create a security group that allows connections from hosts in that subnet.

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

### Filter Example

If you want to match against tag `Name`, use:

```terraform
data "aws_subnet" "selected" {
  filter {
    name   = "tag:Name"
    values = ["yakdriver"]
  }
}
```

## Argument Reference

The arguments of this data source act as filters for querying the available subnets in the current region. The given filters must match exactly one subnet whose data will be exported as attributes.

The following arguments are optional:

* `availability_zone` - (Optional) Availability zone where the subnet must reside.
* `default_for_az` - (Optional) Whether the desired subnet must be the default subnet for its associated availability zone.
* `filter` - (Optional) One or more name/value pairs to use as filters.
    * _Valid values:_ See supported names and values in [EC2 API documentation][describe-subnets]
* `id` - (Optional) ID of the specific subnet to retrieve.
* `state` - (Optional) State that the desired subnet must have.
* `tags` - (Optional) Map of tags, each pair of which must exactly match a pair on the desired subnet.
* `vpc_id` - (Optional) ID of the VPC that the desired subnet belongs to.

## Attribute Reference

### Supported attributes

In addition to all arguments above, the following attributes are exported:

* `map_public_ip_on_launch` - Indicates whether public IP addresses will be associated with instances created in this subnet. Addresses are associated only if there are available allocated Elastic IP addresses.

### Unsupported attributes

~> **Note** These attributes may be present in the `terraform.tfstate` file, but they have preset values and cannot be specified in configuration files.

The following attributes are not currently supported:

`arn`, `assign_ipv6_address_on_creation`, `availability_zone_id`, `customer_owned_ipv4_pool`, `enable_dns64`, `enable_resource_name_dns_aaaa_record_on_launch`, `enable_resource_name_dns_a_record_on_launch`, `ipv6_cidr_block_association_id`, `ipv6_native`, `map_customer_owned_ip_on_launch`, `outpost_arn`, `owner_id`, `private_dns_hostname_type_on_launch`.
