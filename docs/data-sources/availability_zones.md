---
subcategory: "EC2 (Elastic Compute Cloud)"
layout: "aws"
page_title: "aws_availability_zones"
description: |-
  Provides a list of availability zone names.
---

[describe-azs]: https://docs.k2.cloud/en/api/ec2/actions/placements/DescribeAvailabilityZones.html

# Data Source: aws_availability_zones

Provides a list of availability zone names matching the specified criteria.
To get information about a specific availability zone, use the [`aws_availability_zone`](availability_zone.md) (singular) data source.

## Example Usage

### By State

```terraform
data "aws_availability_zones" "available" {
  state = "available"
}
```


## Argument Reference

The following arguments are supported:

* `filter` - (Optional) One or more name/value pairs to use as filters.
    * _Valid values:_ See supported names and values in [EC2 API documentation][describe-azs]
* `state` - (Optional) Filters the list of availability zones based on their
current state.
    * _Valid values:_ `available`, `information`, `impaired`, `unavailable`

## Attribute Reference

### Supported attributes

In addition to all arguments above, the following attributes are exported:

* `id` - The region of the availability zones.
* `names` - List of the availability zone names available to the account.

### Unsupported attributes

~> **Note** These attributes may be present in the `terraform.tfstate` file, but they have preset values and cannot be specified in configuration files.

The following attributes are not currently supported:

`all_availability_zones`, `exclude_names`, `exclude_zone_ids`, `group_names`, `zone_ids`.
