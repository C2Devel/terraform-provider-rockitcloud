---
subcategory: "VPC (Virtual Private Cloud)"
layout: "aws"
page_title: "aws_network_interfaces"
description: |-
  Provides a list of network interface IDs.
---

[describe-network-interfaces]: https://docs.k2.cloud/en/api/ec2/actions/network_interfaces/DescribeNetworkInterfaces.html

# Data Source: aws_network_interfaces

Provides a list of network interface IDs matching the specified criteria.

## Example Usage

The following shows all network interface IDs.

```terraform
data "aws_network_interfaces" "example" {}

output "example" {
  value = data.aws_network_interfaces.example.ids
}
```

The following example retrieves a list of all network interface IDs with a custom tag of `Name` set to a value of `test`.

```terraform
data "aws_network_interfaces" "example1" {
  tags = {
    Name = "test"
  }
}

output "example1" {
  value = data.aws_network_interfaces.example.ids
}
```

The following example retrieves a network interface IDs which associated with specific subnet.

```terraform
data "aws_network_interfaces" "example2" {
  filter {
    name   = "subnet-id"
    values = ["subnet-xxxxxxxx"]
  }
}

output "example2" {
  value = data.aws_network_interfaces.example.ids
}
```

## Argument Reference

* `filter` - (Optional) One or more name/value pairs to use as filters.
    * _Valid values:_ See supported names and values in [EC2 API documentation][describe-network-interfaces]
* `tags` - (Optional) Map of tags, each pair of which must exactly match
  a pair on the desired network interfaces.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The region.
* `ids` - A list of all the network interface IDs found.
