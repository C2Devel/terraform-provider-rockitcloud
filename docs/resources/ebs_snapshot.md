---
subcategory: "EBS (EC2)"
layout: "aws"
page_title: "aws_ebs_snapshot"
description: |-
  Creates a snapshot of a volume.
---

[default-tags]: https://registry.terraform.io/providers/hashicorp/aws/latest/docs#default_tags-configuration-block
[timeouts]: https://developer.hashicorp.com/terraform/plugin/framework/resources/timeouts

# Resource: aws_ebs_snapshot

Creates a snapshot of a volume.

## Example usage

```terraform
resource "aws_ebs_volume" "example" {
  availability_zone = "ru-msk-vol52"
  size              = 40

  tags = {
    Name = "HelloWorld"
  }
}

resource "aws_ebs_snapshot" "example_snapshot" {
  volume_id = aws_ebs_volume.example.id

  tags = {
    Name = "HelloWorld_snap"
  }
}
```

## Argument reference

The following arguments are required:

* `volume_id` - (Required, Forces new resource, String) The ID of the volume.

The following arguments are optional:

* `description` - (Optional, Editable, String) The description of the snapshot.
* `tags` - (Optional, Editable, Map of strings) Key-value pairs to assign to the snapshot. If the [`default_tags` configuration block][default-tags] is used within a provider configuration, the tags with matching keys will overwrite those defined at the provider level.

## Attribute reference

### Supported attributes

In addition to all arguments above, the following attributes are exported:

* `arn` - (String) The Amazon Resource Name (ARN) of the snapshot.
* `id` - (String) The ID of the snapshot.
    * _Example:_ `snap-12345678`
* `owner_alias` - (String) The alias of the snapshot owner.
* `owner_id` - (String) The ID of the project that owns the snapshot.
* `tags_all` - (Map of strings) Key-value pairs assigned to the snapshot, including any tags inherited from the [`default_tags` configuration block][default-tags] if used within a provider configuration.
* `volume_size` - (Integer) The size of the volume in GiB.

### Unsupported attributes

~> **Note** These attributes may be present in the `terraform.tfstate` file, but they have preset values and cannot be specified in configuration files.

The following attributes are not currently supported:

`data_encryption_key_id`, `encrypted`, `kms_key_id`, `outpost_arn`, `permanent_restore`, `storage_tier`, `temporary_restore_days`.

## Timeouts

The `timeouts` block allows you to specify [timeouts] for certain actions:

- `create` - (Default `10 minutes`) Used for creating the snapshot.
- `delete` - (Default `10 minutes`) Used for deleting the snapshot.

## Import

The snapshot can be imported using `id`, for example:

```
$ terraform import aws_ebs_snapshot.id snap-12345678
```
