---
subcategory: "EBS (EC2)"
layout: "aws"
page_title: "aws_ebs_snapshot"
description: |-
  Provides information about an EBS snapshot.
---

[describe-snapshots]: https://docs.k2.cloud/en/api/ec2/actions/snapshots/DescribeSnapshots.html

# Data Source: aws_ebs_snapshot

Provides information about an EBS snapshot.

## Example Usage

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

## Argument Reference

The following arguments are supported:

* `filter` - (Optional) One or more name/value pairs to use as filters.
    * _Valid values:_ See supported names and values in [EC2 API documentation][describe-snapshots]
* `most_recent` - (Optional) If more than one result is returned, use the most recent snapshot.
* `owners` - (Optional) List of the snapshot owners.
    * _Valid values:_ Project ID (`project@customer`) or `self`
* `restorable_by_user_ids` - (Optional) List of the project IDs (`project@customer`), in which volumes can be created from a snapshot.
* `snapshot_ids` - (Optional) Returns information on a snapshot ID.

## Attribute Reference

### Supported attributes

In addition to all arguments above, the following attributes are exported:

* `arn` - The Amazon Resource Name (ARN) of the EBS snapshot.
* `id` - The snapshot ID (e.g., snap-12345678).
* `snapshot_id` - The snapshot ID (e.g., snap-12345678).
* `description` - A description for the snapshot
* `owner_id` - The project ID.
* `owner_alias` - The alias of the EBS snapshot owner.
* `volume_id` - The volume ID (e.g., vol-12345678).
* `volume_size` - The size of the drive in GiB.
* `state` - The snapshot state.
* `tags` - Map of tags assigned to the snapshot.

### Unsupported attributes

~> **Note** These attributes may be present in the `terraform.tfstate` file, but they have preset values and cannot be specified in configuration files.

The following attributes are not currently supported:

`data_encryption_key_id`, `encrypted`, `kms_key_id`, `outpost_arn`, `storage_tier`.
