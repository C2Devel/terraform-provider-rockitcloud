---
subcategory: "VPC (Virtual Private Cloud)"
layout: "aws"
page_title: "aws_vpcs"
description: |-
  Provides a list of VPC IDs.
---

[describe-vpcs]: https://docs.k2.cloud/en/api/ec2/actions/vpcs/DescribeVpcs.html

# Data Source: aws_vpcs

Provides a list of VPC IDs.

## Example usage

### Basic example

The following example retrieves a list of VPC IDs with a custom `service` tag set to `production`.

```terraform
data "aws_vpcs" "foo" {
  tags = {
    service = "production"
  }
}

output "foo" {
  value = data.aws_vpcs.foo.ids
}
```

## Argument reference

The following arguments are supported:

* `filter` - (Optional, [Block](#filter)) One or more name/value pairs to use as filters.
    * _Valid values:_ See supported names and values in [EC2 API documentation][describe-vpcs]
* `tags` - (Optional, Map of strings) Key-value pairs.
    Must exactly match pairs on the required resources.

### filter

* `name` - (Required, String) The name of the filter.
    * _Constraints:_ Filter names are case-sensitive
* `values` - (Required, List of strings) One or more filter values.
    * _Constraints:_ Filter values are case-sensitive

## Attribute reference

Provides a list of VPC IDs in a region.

* `id` - (String) The region.
    * _Example_: `ru-spb`
* `ids` - (List of strings) The list of VPC IDs found.
