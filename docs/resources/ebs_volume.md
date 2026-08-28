---
subcategory: "EBS (EC2)"
layout: "aws"
page_title: "aws_ebs_volume"
description: |-
  Manages a single volume.
---

[default-tags]: https://www.terraform.io/docs/providers/aws/index.html#default_tags-configuration-block
[timeouts]: https://developer.hashicorp.com/terraform/plugin/framework/resources/timeouts

# Resource: aws_ebs_volume

Manages a single volume.

## Example usage

```terraform
resource "aws_ebs_volume" "example" {
  availability_zone = "ru-msk-vol52"
  size              = 40

  tags = {
    Name = "HelloWorld"
  }
}
```

## Argument reference

The following arguments are required:

* `availability_zone` - (Required, Forces new resource, String) The availability zone to create the volume in.

The following arguments are optional:

* `iops` - (Optional, Editable, Integer) The amount of IOPS to provision for the volume.
    * _Constrains:_ Can be applied to volumes only with `type` set to `io2`
* `size` - (Optional, Editable, Integer) The size of the volume in GiB.
    * _Constraints:_ This argument is required if the `snapshot_id` argument is not specified
* `snapshot_id` (Optional, Forces new resource, String) The ID of the snapshot.
    * _Constraints:_ This argument is required if the `size` argument is not specified
* `tags` - (Optional, Editable, Map of tags) Key-value pairs to assign to the volume. If the [`default_tags` configuration block][default-tags] is used within a provider configuration, the tags with matching keys will overwrite those defined at the provider level.
* `type` - (Optional, Editable, String) The type of the volume.

## Attribute reference

### Supported attributes

In addition to all arguments above, the following attributes are exported:

* `arn` - (String) The Amazon Resource Name (ARN) of the volume.
* `encrypted` - (Boolean) Indicates whether the EBS volume is encrypted.
* `id` - (String) The ID of the volume.
    * _Example:_ `vol-12345678`
* `tags_all` - (Map of strings) Key-value pairs assigned to the volume, including any tags inherited from the [`default_tags` configuration block][default-tags] if used within a provider configuration.
* `throughput` - (Integer) The throughput that the volume supports in MiB/s.

### Unsupported attributes

~> **Note** These attributes may be present in the `terraform.tfstate` file, but they have preset values and cannot be specified in configuration files.

The following attributes are not currently supported:

`kms_key_id`, `multi_attach_enabled`, `outpost_arn`.

## Timeouts

The `timeouts` block allows you to specify [timeouts] for certain actions:

- `create` - (Default `5 minutes`) Used for creating the volume. This includes the time required for the volume to become available.
- `update` - (Default `5 minutes`) Used for updating the volume.
- `delete` - (Default `5 minutes`) Used for destroying the volume.

## Import

The volume can be imported using `id`, for example:

```
$ terraform import aws_ebs_volume.id vol-12345678
```
