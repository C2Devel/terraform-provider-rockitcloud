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

## Example usage

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

## Argument reference

The arguments of this data source act as filters for querying the available Elastic IPs.

* `filter` - (Optional, [Block](#filter)) One or more name/value pairs to use as filters.
    * _Valid values:_ See supported names and values in the [EC2 API documentation][describe-addresses]
* `tags` - (Optional, Map of strings) Key-value pairs. Must exactly match pairs on the desired Elastic IPs.

### filter

* `name` - (Required, String) The name of the filter.
    * _Constraints:_ Filter names are case-sensitive
* `values` - (Required, List of strings) One or more filter values.
    * _Constraints:_ Filter values are case-sensitive

## Attribute reference

In addition to all arguments above, the following attributes are exported:

* `allocation_ids` - (List of strings) List of all allocation IDs.
* `id` - (String) The region.
* `public_ips` - (List of strings) List of all Elastic IP addresses.
