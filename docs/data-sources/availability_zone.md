---
subcategory: "EC2 (Elastic Compute Cloud)"
layout: "aws"
page_title: "aws_availability_zone"
description: |-
  Provides information about an availability zone.
---

[describe-azs]: https://docs.k2.cloud/en/api/ec2/actions/placements/DescribeAvailabilityZones.html

# Data Source: aws_availability_zone

Provides information about an availability zone.
To get a list of the available zones, use the [`aws_availability_zones`](availability_zones.md) (plural) data source.

## Example usage

### Basic example

```terraform
data "aws_availability_zone" "example" {
  name = "ru-msk-vol52"
}

output "availability_zone_to_region" {
  value = data.aws_availability_zone.example.id
}
```

## Argument reference

The arguments of this data source act as filters for querying the available
availability zones. The given filters must match exactly one availability
zone whose data will be exported as attributes.

* `all_availability_zones` - (Optional, Boolean) Indicates whether to include availability zones that are not currently available.
* `filter` - (Optional, [Block](#filter)) One or more name/value pairs to use as filters.
    * _Valid values:_ See supported names and values in [EC2 API documentation][describe-azs]
* `name` - (Optional, String) The full name of the availability zone to select.
* `state` - (Optional, String) A specific availability zone state to require.
    * _Valid values:_ `available`, `impaired`, `information`
* `zone_id` - (Optional, String) The ID of the availability zone to select.

### filter

* `name` - (Required, String) The name of the filter.
    * _Constraints:_ Filter names are case-sensitive
* `values` - (Required, List of strings) One or more filter values.
    * _Constraints:_ Filter values are case-sensitive

## Attribute reference

### Supported attributes

In addition to all arguments above, the following attributes are exported:

* `id` - (String) The ID of the availability zone.
* `region` - (String) The region where the selected availability zone resides.
* `zone_id` - (String) The ID of the availability zone.

### Unsupported attributes

~> **Note** These attributes may be present in the `terraform.tfstate` file, but they have preset values and cannot be specified in configuration files.

The following attributes are not currently supported:

`group_name`, `network_border_group`, `opt_in_status`, `parent_zone_id`, `parent_zone_name`, `zone_type`.
