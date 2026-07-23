---
subcategory: "VPC (Virtual Private Cloud)"
layout: "aws"
page_title: "aws_network_acls"
description: |-
  Provides a list of network ACL IDs for a VPC.
---

[describe-network-acls]: https://docs.k2.cloud/en/api/ec2/actions/network_acls/DescribeNetworkAcls.html

# Data Source: aws_network_acls

Provides a list of network ACL IDs for a VPC.

## Example usage

### Basic example

The following example shows all network ACL IDs in a VPC.

```terraform
variable vpc_id {}

data "aws_network_acls" "example" {
  vpc_id = var.vpc_id
}

output "example" {
  value = data.aws_network_acls.example.ids
}
```

### Specific examples

The following example retrieves a list of all network ACLs associated with a VPC with a custom `Tier` tag set to `Private`.

```terraform
variable vpc_id {}

data "aws_network_acls" "example" {
  vpc_id = var.vpc_id

  tags = {
    Tier = "Private"
  }
}
```

The following example retrieves the ID of a network ACL which is associated with a specific subnet in a VPC.

```terraform
variable vpc_id {}
variable subnet_id {}

data "aws_network_acls" "example" {
  vpc_id = var.vpc_id

  filter {
    name   = "association.subnet-id"
    values = [var.subnet_id]
  }
}
```

## Argument reference

The following arguments are supported:

* `filter` - (Optional, [Block](#filter)) One or more name/value pairs to use as filters.
    * _Valid values:_ See supported names and values in [EC2 API documentation][describe-network-acls]
* `tags` - (Optional, Map of strings) Key-value pairs.
    Must exactly match pairs on the required resources.
* `vpc_id` - (Optional, String) The ID of the VPC that you want to filter from.

### filter

* `name` - (Required, String) The name of the filter.
    * _Constraints:_ Filter names are case-sensitive
* `values` - (Required, List of strings) One or more filter values.
    * _Constraints:_ Filter values are case-sensitive

## Attribute reference

In addition to all arguments above, the following attributes are exported:

* `id` - (String) The region.
    * _Example:_ `ru-spb`
* `ids` - (List of strings) The list of all the network ACL IDs found.
