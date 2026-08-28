---
subcategory: "EBS (EC2)"
layout: "aws"
page_title: "aws_ebs_volume"
description: |-
  Provides information about a volume.
---

[describe-volumes]: https://docs.k2.cloud/en/api/ec2/actions/volumes/DescribeVolumes.html

# Data Source: aws_ebs_volume

Provides information about a volume.

## Example usage

```terraform
data "aws_ebs_volume" "ebs_volume" {
  most_recent = true

  filter {
    name   = "volume-type"
    values = ["st2"]
  }

  filter {
    name   = "tag:Name"
    values = ["Example"]
  }
}
```

## Argument reference

The following arguments are supported:

* `filter` - (Optional, [Block](#filter)) One or more name/value pairs to use as filters.
    * _Valid values:_ See supported names and values in [EC2 API documentation][describe-volumes]
* `most_recent` - (Optional, Boolean) If more than one result is returned, use the most recent volume.
    * _Default value:_ `false`

### filter

* `name` - (Required, String) The name of the filter.
    * _Constraints:_ Filter names are case-sensitive
* `values` - (Required, List of strings) One or more filter values.
    * _Constraints:_ Filter values are case-sensitive

## Attribute reference

### Supported attributes

In addition to all arguments above, the following attributes are exported:

* `arn` - (String) The Amazon Resource Name (ARN) of the volume.
* `availability_zone` - (String) The availability zone where the volume will be located.
* `encrypted` - (Boolean) Indicates whether the EBS volume is encrypted.
* `id` - (String) The ID of the volume.
    * _Example:_ `vol-12345678`
* `iops` - (Integer) The amount of IOPS for the volume.
* `size` - (Integer) The size of the volume in GiB.
* `snapshot_id` - (String) The ID of the snapshot the volume is based on.
* `tags` - (Map of strings) Key-value pairs assigned to the volume.
* `throughput` - (Integer) The throughput that the volume supports in MiB/s.
* `volume_id` - (String) The ID of the volume.
    * _Example:_ `vol-12345678`
* `volume_type` - (String) The type of the volume.

### Unsupported attributes

~> **Note** These attributes may be present in the `terraform.tfstate` file, but they have preset values and cannot be specified in configuration files.

The following attributes are not currently supported:

`kms_key_id`, `multi_attach_enabled`, `outpost_arn`.
