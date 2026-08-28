---
subcategory: "EBS (EC2)"
layout: "aws"
page_title: "aws_ebs_snapshot"
description: |-
  Provides information about a snapshot.
---

[describe-snapshots]: https://docs.k2.cloud/en/api/ec2/actions/snapshots/DescribeSnapshots.html

# Data Source: aws_ebs_snapshot

Provides information about a snapshot.

## Example usage

```terraform
data "aws_ebs_snapshot" "ebs_snapshot" {
  most_recent = true
  owners      = ["self"]

  filter {
    name   = "volume-size"
    values = ["40"]
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
    * _Valid values:_ See supported names and values in [EC2 API documentation][describe-snapshots]
* `most_recent` - (Optional, Boolean) Indicates whether to use the most recent snapshot if more than one result is returned.
    * _Default value:_ `false`
* `owners` - (Optional, List of strings) The list of the snapshot owners.
    * _Valid values:_ `project@customer`, `self`
* `restorable_by_user_ids` - (Optional, List of strings) The list of the project IDs (`project@customer`), in which volumes can be created from a snapshot.
* `snapshot_ids` - (Optional, List of strings) The list of the snapshot IDs.

### filter

* `name` - (Required, String) The name of the filter.
    * _Constraints:_ Filter names are case-sensitive
* `values` - (Required, List of strings) One or more filter values.
    * _Constraints:_ Filter values are case-sensitive

## Attribute reference

### Supported attributes

In addition to all arguments above, the following attributes are exported:

* `arn` - (String) The Amazon Resource Name (ARN) of the snapshot.
* `description` - (String) The description for the snapshot.
* `encrypted` - (Boolean) Indicates whether the snapshot is encrypted.
* `id` - (String) The ID of the snapshot.
    * _Example:_ `snap-12345678`
* `owner_alias` - (String) The alias of the snapshot owner.
* `owner_id` - (String) The ID of the project that owns the snapshot.
* `snapshot_id` - (String) The ID of the snapshot.
    * _Example:_ `snap-12345678`
* `state` - (String) The state of the snapshot.
* `tags` - (Map of strings) Key-value pairs assigned to the snapshot.
* `volume_id` - (String) The ID of the volume.
    * _Example:_ `vol-12345678`
* `volume_size` - (Integer) The size of the volume in GiB.

### Unsupported attributes

~> **Note** These attributes may be present in the `terraform.tfstate` file, but they have preset values and cannot be specified in configuration files.

The following attributes are not currently supported:

`data_encryption_key_id`, `kms_key_id`, `outpost_arn`, `storage_tier`.
