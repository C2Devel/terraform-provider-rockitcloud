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

## Example Usage

```terraform
data "aws_availability_zone" "example" {
  name = "ru-msk-vol52"
}

output "availability_zone_to_region" {
  value = data.aws_availability_zone.example.id
}
```

## Argument Reference

The arguments of this data source act as filters for querying the available
availability zones. The given filters must match exactly one availability
zone whose data will be exported as attributes.

* `filter` - (Optional) One or more name/value pairs to use as filters.
    * _Valid values:_ See supported names and values in [EC2 API documentation][describe-azs]
* `name` - (Optional) The full name of the availability zone to select.
* `state` - (Optional) A specific availability zone state to require.
    * _Valid values:_ `available`, `information`, `impaired`

## Attribute Reference

### Supported attributes

In addition to all arguments above, the following attributes are exported:

* `region` - The region where the selected availability zone resides.

### Unsupported attributes

~> **Note** These attributes may be present in the `terraform.tfstate` file, but they have preset values and cannot be specified in configuration files.

The following attributes are not currently supported:

`all_availability_zones`, `group_name`, `network_border_group`, `opt_in_status`, `parent_zone_id`, `parent_zone_name`, `zone_id`, `zone_type`.
