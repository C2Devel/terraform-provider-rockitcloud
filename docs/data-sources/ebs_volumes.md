---
subcategory: "EBS (EC2)"
layout: "aws"
page_title: "aws_ebs_volumes"
description: |-
  Provides a list of EBS volume IDs.
---

[describe-volumes]: https://docs.k2.cloud/en/api/ec2/actions/volumes/DescribeVolumes.html

# Data Source: aws_ebs_volumes

Provides a list of EBS volume IDs matching the specified criteria.
This data source can be used to get a list of volume IDs with (for example) matching tags.

## Example Usage

The following demonstrates obtaining a map of availability zone to EBS volume ID for volumes with a given tag value.

```terraform
data "aws_ebs_volumes" "example" {
  tags = {
    Name = "Example"
  }
}

data "aws_ebs_volume" "example" {
  for_each = toset(data.aws_ebs_volumes.example.ids)
  filter {
    name   = "volume-id"
    values = [each.value]
  }
}

output "availability_zone_to_volume_id" {
  value = { for s in data.aws_ebs_volume.example : s.id => s.availability_zone }
}
```

### Filter example

If matching against the `size` filter, use:

```terraform
data "aws_ebs_volumes" "ten_or_twenty_gb_volumes" {
  filter {
    name   = "size"
    values = ["10", "20"]
  }
}
```

## Argument Reference

In addition to all arguments above, the following attributes are exported:

* `filter` - (Optional) One or more name/value pairs to use as filters.
    * _Valid values:_ See supported names and values in [EC2 API documentation][describe-volumes]
* `tags` - (Optional) Map of tags, each pair of which must exactly match
  a pair on the desired volumes.

## Attribute Reference

* `id` - The region.
* `ids` - A set of all the EBS volume IDs found.
