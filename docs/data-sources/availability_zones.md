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

## Example usage

### By state

```terraform
data "aws_availability_zones" "available" {
  state = "available"
}
```


## Argument reference

The following arguments are supported:

* `all_availability_zones` - (Optional, Boolean) Indicates whether to include availability zones that are not currently available.
* `exclude_names` - (Optional, Set of strings) List of availability zone names to exclude from the results.
* `exclude_zone_ids` - (Optional, Set of strings) List of availability zone IDs to exclude from the results.
* `filter` - (Optional, [Block](#filter)) One or more name/value pairs to use as filters.
    * _Valid values:_ See supported names and values in [EC2 API documentation][describe-azs]
* `state` - (Optional, String) Filters the list of availability zones based on their
current state.
    * _Valid values:_ `available`, `impaired`, `information`, `unavailable`

### filter

* `name` - (Required, String) The name of the filter.
    * _Constraints:_ Filter names are case-sensitive
* `values` - (Required, List of strings) One or more filter values.
    * _Constraints:_ Filter values are case-sensitive

## Attribute reference

### Supported attributes

In addition to all arguments above, the following attributes are exported:

* `id` - (String) The region of the availability zones.
* `group_names` - (Set of strings) The list of group names of the availability zones.
* `names` - (List of strings) The list of availability zone names available to the account.
* `zone_ids` - (List of strings) The list of availability zone IDs available to the account.

### Unsupported attributes

~> **Note** These attributes may be present in the `terraform.tfstate` file, but they have preset values and cannot be specified in configuration files.

The following attributes are not currently supported:

_(None)_
