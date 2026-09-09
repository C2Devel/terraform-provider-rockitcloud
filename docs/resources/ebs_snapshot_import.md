---
subcategory: "EBS (EC2)"
layout: "aws"
page_title: "aws_ebs_snapshot_import"
description: |-
  Imports a disk image from object storage as a snapshot.
---

[default-tags]: https://registry.terraform.io/providers/hashicorp/aws/latest/docs#default_tags-configuration-block
[timeouts]: https://developer.hashicorp.com/terraform/plugin/framework/resources/timeouts

# Resource: aws_ebs_snapshot_import

Imports a disk image from object storage as a snapshot.

## Example usage

```terraform
resource "aws_ebs_snapshot_import" "example" {
  disk_container {
    format = "VHD"
    user_bucket {
      s3_bucket = "disk-images"
      s3_key    = "source.vhd"
    }
  }

  tags = {
    Name = "HelloWorld"
  }
}
```

## Argument reference

The following arguments are required:

* `disk_container` - (Required, Forces new resource, [Block](#disk_container)) Information about the disk container.

The following arguments are optional:

* `description` - (Optional, Editable, String) The description string for the snapshot import task.
* `tags` - (Optional, Editable, Map of strings) Key-value pairs to assign to the snapshot. If the [`default_tags` configuration block][default-tags] is used within a provider configuration, the tags with matching keys will overwrite those defined at the provider level.

### disk_container

The following arguments are required:

* `format` - (Required, Forces new resource, String) The format of the disk image being imported.
    * _Valid values:_  `QCOW2`, `RAW`, `VHD`, `VMDK`
* `user_bucket` - (Required, Forces new resource, [Block](#user_bucket)) The bucket for the disk image.

The following arguments are optional:

* `description` - (Optional, Editable, String) The description of the disk image being imported.

#### user_bucket

The following arguments are required:

* `s3_bucket` - (Required, Forces new resource, String) The name of the bucket where the disk image is located.
* `s3_key` - (Required, Forces new resource, String) The name of the disk image file.

## Attribute reference

### Supported attributes

In addition to all arguments above, the following attributes are exported:

* `arn` - (String) The Amazon Resource Name (ARN) of the snapshot.
* `id` - (String) The ID of the snapshot.
    * _Example:_ `snap-12345678`
* `owner_alias` - (String) The alias of the snapshot owner.
* `owner_id` - (String) The ID of the project that owns the resource.
* `volume_size` - (Integer) The size of the volume in GiB.
* `tags_all` - (Map of strings) Key-value pairs assigned to the resource, including any tags inherited from the [`default_tags` configuration block][default-tags] if used within a provider configuration.

### Unsupported attributes

~> **Note** These attributes may be present in the `terraform.tfstate` file, but they have preset values and cannot be specified in configuration files.

The following attributes are not currently supported:

`client_data`, `data_encryption_key_id`, `disk_container.url`, `encrypted`, `kms_key_id`, `outpost_arn`, `permanent_restore`, `role_name`, `storage_tier`, `temporary_restore_days`.

## Timeouts

The `timeouts` block allows you to specify [timeouts] for certain actions:

- `create` - (Default `60 minutes`) Used for importing the snapshot.
- `delete` - (Default `10 minutes`) Used for deleting the snapshot.

## Import

To import a snapshot, use the [aws_ebs_snapshot](ebs_snapshot.md#import) resource.
