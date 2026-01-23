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

The following example retrieves a list of VPC IDs with a custom tag of `service` set to a value of "production".

## Example Usage

The following shows outputing all VPC IDs.

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

## Argument Reference

* `filter` - (Optional) One or more name/value pairs to use as filters.
    * _Valid values:_ See supported names and values in [EC2 API documentation][describe-vpcs]
* `tags` - (Optional) Map of tags, each pair of which must exactly match
  a pair on the desired VPCs.

## Attribute Reference

Provides a list of VPC IDs in a region.

* `id` - A region.
* `ids` - A list of all the VPC IDs found.
