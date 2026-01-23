---
subcategory: "VPC (Virtual Private Cloud)"
layout: "aws"
page_title: "aws_security_groups"
description: |-
  Provides information about a set of security groups.
---

[describe-security-groups]: https://docs.k2.cloud/en/api/ec2/actions/security_groups/DescribeSecurityGroups.html

# Data Source: aws_security_groups

Provides information about IDs and VPC membership of security groups.

## Example Usage

```terraform
data "aws_security_groups" "test" {
  tags = {
    Application = "k8s"
    Environment = "dev"
  }
}
```

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

## Argument Reference

* `tags` - (Optional) Map of tags, each pair of which must exactly match for desired security groups.
* `filter` - (Optional) One or more name/value pairs to use as filters.
    * _Valid values:_ See supported names and values in [EC2 API documentation][describe-security-groups]

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `arns` - The Amazon Resource Names (ARNs) of the matched security groups.
* `id` - The region.
* `ids` - IDs of the matched security groups.
* `vpc_ids` - The VPC IDs of the matched security groups. The data source's tag or filter *will span VPCs* unless the `vpc-id` filter is also used.
