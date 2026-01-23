---
subcategory: "EC2 (Elastic Compute Cloud)"
layout: "aws"
page_title: "aws_eips"
description: |-
  Provides a list of Elastic IPs.
---

[describe-addresses]: https://docs.k2.cloud/en/api/ec2/actions/addresses/DescribeAddresses.html

# Data Source: aws_eips

Provides a list of Elastic IPs.

## Example Usage

The following shows all Elastic IPs with the specific tag value.

```terraform
data "aws_eips" "example" {
  tags = {
    Env = "dev"
  }
}

output "allocation_ids" {
  value = data.aws_eips.example.allocation_ids
}

output "public_ips" {
  value = data.aws_eips.example.public_ips
}
```

## Argument Reference

* `filter` - (Optional) One or more name/value pairs to use as filters.
    * _Valid values:_ See supported names and values in [EC2 API documentation][describe-addresses]
* `tags` - (Optional) Map of tags, each pair of which must exactly match a pair on the desired Elastic IPs.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `allocation_ids` - List of all allocation IDs.
* `id` - The region.
* `public_ips` - List of all Elastic IP addresses.
