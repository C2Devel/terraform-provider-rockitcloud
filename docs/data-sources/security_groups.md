---
subcategory: "VPC (Virtual Private Cloud)"
layout: "aws"
page_title: "aws_security_groups"
description: |-
  Provides information about a set of security groups.
---

[describe-security-groups]: https://docs.k2.cloud/en/api/ec2/actions/security_groups/DescribeSecurityGroups.html

# Data Source: aws_security_groups

Provides information about IDs of security groups and their association to any VPC.

## Example usage

### Basic example

```terraform
data "aws_security_groups" "test" {
  tags = {
    Application = "k8s"
    Environment = "dev"
  }
}
```

### Specific example

```terraform
variable vpc_id {}

data "aws_security_groups" "test" {
  filter {
    name   = "group-name"
    values = ["nodes"]
  }

  filter {
    name   = "vpc-id"
    values = [var.vpc_id]
  }
}
```

## Argument reference

The following arguments are supported:

* `filter` - (Optional, [Block](#filter)) One or more name/value pairs to use as filters.
    * _Valid values:_ See supported names and values in [EC2 API documentation][describe-security-groups]
* `tags` - (Optional, Map of strings) Key-value pairs.
    Must exactly match pairs on the required resources.

### filter

* `name` - (Required, String) The name of the filter.
    * _Constraints:_ Filter names are case-sensitive
* `values` - (Required, List of strings) One or more filter values.
    * _Constraints:_ Filter values are case-sensitive

## Attribute reference

In addition to all arguments above, the following attributes are exported:

* `arns` - (String) The Amazon Resource Names (ARNs) of the matched security groups.
* `id` - (String) The region.
    * _Example:_ `ru-spb`
* `ids` - (List of strings) The IDs of the matched security groups.
* `vpc_ids` - (String) The VPC IDs of the matched security groups.
    The data source's tag or filter *will span VPCs* unless the `vpc-id` filter is also used.
