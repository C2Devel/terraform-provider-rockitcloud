---
subcategory: "VPC (Virtual Private Cloud)"
layout: "aws"
page_title: "aws_route_tables"
description: |-
  Provides a list of route table IDs.
---

[describe-route-tables]: https://docs.k2.cloud/en/api/ec2/actions/routes/DescribeRouteTables.html

# Data Source: aws_route_tables

Provides a list of route table IDs matching the specified criteria.

## Example usage

### Basic example

```terraform
variable vpc_id {}

data "aws_route_tables" "rts" {
  vpc_id = var.vpc_id

  filter {
    name   = "tag:kubernetes.io/kops/role"
    values = ["private", "public"]
  }
}
```

## Argument reference

The following arguments are supported:

* `filter` - (Optional, [Block](#filter)) One or more name/value pairs to use as filters.
    * _Valid values:_ See supported names and values in [EC2 API documentation][describe-route-tables]
* `tags` - (Optional, Map of strings) Key-value pairs.
    Must exactly match pairs on the required resources.
* `vpc_id` - (Optional, String) The VPC ID for filtering the route tables related to this VPC.

### filter

* `name` - (Required, String) The name of the filter.
    * _Constraints:_ Filter names are case-sensitive
* `values` - (Required, List of strings) One or more filter values.
    * _Constraints:_ Filter values are case-sensitive

## Attribute reference

* `id` - (String) The region.
    * _Example:_ `ru-spb`
* `ids` - (List of strings) The list of all the route table IDs found.
