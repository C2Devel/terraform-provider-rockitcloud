---
subcategory: "VPC (Virtual Private Cloud)"
layout: "aws"
page_title: "aws_security_group"
description: |-
  Provides information about a security group.
---

[describe-security-groups]: https://docs.k2.cloud/en/api/ec2/actions/security_groups/DescribeSecurityGroups.html

# Data Source: aws_security_group

Provides information about a security group.

This data source can be used when a module accepts the ID of a security group as an input variable and needs to, for example, determine the ID of the VPC that owns the security group.

## Example usage

### Specific example

The following example shows how to accept the ID of a security group as a variable and use this data source to obtain the data necessary to create a subnet.

```terraform
variable "security_group_id" {}

data "aws_security_group" "selected" {
  id = var.security_group_id
}

resource "aws_subnet" "subnet" {
  vpc_id     = data.aws_security_group.selected.vpc_id
  cidr_block = "10.0.1.0/24"
}
```

## Argument reference

The arguments of this data source act as filters for querying the available security group in the current region.

~> **Note** The given filters must exactly match the resource whose data will be exported as attributes.

* `filter` - (Optional, [Block](#filter)) One or more name/value pairs to use as filters.
    * _Valid values:_ See supported names and values in [EC2 API documentation][describe-security-groups]
* `id` - (Optional, String) The ID of the specific security group to retrieve.
* `name` - (Optional, String) The name that the desired security group must have.

    -> **Info** The default security group for a VPC has the name `default`.

* `tags` - (Optional, Map of strings) Key-value pairs.
    Must exactly match pairs on the required resource.
* `vpc_id` - (Optional, String) The ID of the VPC that owns the desired security group.

### filter

* `name` - (Required, String) The name of the filter.
    * _Constraints:_ Filter names are case-sensitive
* `values` - (Required, List of strings) One or more filter values.
    * _Constraints:_ Filter values are case-sensitive

## Attribute reference

This data source will complete the data by populating any fields that are not included in the configuration with the data for the selected security group.

In addition to all arguments above, the following attributes are exported:

* `arn` - (String) The Amazon Resource Name (ARN) of the security group.
* `description` - (String) The description of the security group.

