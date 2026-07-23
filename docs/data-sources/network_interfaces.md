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

## Example usage

### Basic examples

The following example retrieves the IDs of all network interfaces.

```terraform
data "aws_network_interfaces" "example" {}

output "example" {
  value = data.aws_network_interfaces.example.ids
}
```

The following example retrieves the IDs of network interfaces with a custom `Name` tag set to `test`.

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

The following example retrieves the IDs of network interfaces associated with the specific subnet.

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

## Argument reference

The following arguments are supported:

* `filter` - (Optional, [Block](#filter)) One or more name/value pairs to use as filters.
    * _Valid values:_ See supported names and values in [EC2 API documentation][describe-network-interfaces]
* `tags` - (Optional, Map of strings) Key-value pairs.
    Must exactly match pairs on the required resources.

### filter

* `name` - (Required, String) The name of the filter.
    * _Constraints:_ Filter names are case-sensitive
* `values` - (Required, List of strings) One or more filter values.
    * _Constraints:_ Filter values are case-sensitive

## Attribute reference

In addition to all arguments above, the following attributes are exported:

* `id` - (String) The region.
    * _Example:_ `ru-spb`
* `ids` - (List of strings) The list of all the network interface IDs found.
