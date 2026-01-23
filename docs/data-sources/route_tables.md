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

## Example Usage

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

## Argument Reference

* `filter` - (Optional) One or more name/value pairs to use as filters.
    * _Valid values:_ See supported names and values in [EC2 API documentation][describe-route-tables]
* `vpc_id` - (Optional) The VPC ID that you want to filter from.
* `tags` - (Optional) Map of tags, each pair of which must exactly match
  a pair on the desired route tables.

## Attribute Reference

* `id` - The region.
* `ids` - List of all the route table IDs found.
