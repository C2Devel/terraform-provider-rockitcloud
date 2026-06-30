---
subcategory: "EC2 (Elastic Compute Cloud)"
layout: "aws"
page_title: "aws_ami"
description: |-
  Manages an Amazon Machine Image (AMI).
---

[default-tags]: https://registry.terraform.io/providers/hashicorp/aws/latest/docs#default_tags-configuration-block
[images]: https://docs.k2.cloud/en/services/storage/images.html
[timeouts]: https://developer.hashicorp.com/terraform/plugin/framework/resources/timeouts

# Resource: aws_ami

Manages an Amazon Machine Image (AMI).
If you just want to share an existing image with another project,
it's better to use [`aws_ami_launch_permission`](ami_launch_permission.md) instead.

For more information about images, see [user documentation][images].

## Example usage

```terraform
# Creates an image that will start a machine whose root device is backed by
# an EBS volume populated from a snapshot. It is assumed that such a snapshot
# already exists with the ID "snap-12345678".
resource "aws_ami" "example" {
  name                = "tf-ami"
  virtualization_type = "hvm"
  root_device_name    = "disk1"

  ebs_block_device {
    device_name = "disk1"
    snapshot_id = "snap-12345678"
    volume_size = 8
  }
}
```

## Argument reference

The following arguments are required:

* `name` - (Required, Forces new resource, String) A unique name for the image.

The following arguments are optional:

* `architecture` - (Optional, Forces new resource, String) The machine architecture for created instances.
    * _Default value:_ `x86_64`
* `description` - (Optional, Editable, String) A longer, human-readable description for the image.
* `ebs_block_device` - (Optional, Editable, [Block](#ebs_block_device)) A list of EBS block devices that should be attached to created instances.
* `ephemeral_block_device` - (Optional, Forces new resource, [Block](#ephemeral_block_device)) A list of ephemeral block devices that should be attached to created instances.
* `root_device_name` - (Optional, Forces new resource, String) The name of the root device.
    * _Valid values:_ `cdrom<N>`, `disk<N>`, `floppy<N>`, `menu`, where `<N>` is a disk number
* `sriov_net_support` - (Optional, Editable, String) The value of the SriovNetSupport parameter.
    * _Default value:_ `simple`
* `tags` - (Optional, Editable, Map of strings) Key-value pairs to assign to the image. If the [`default_tags` configuration block][default-tags] is used within a provider configuration, the tags with matching keys will overwrite those defined at the provider level.
* `virtualization_type` - (Optional, Forces new resource, String) A keyword to choose what virtualization mode created instances will use.
    * _Valid values:_ `hvm`, `hvm-legacy`
    * _Default value:_ `hvm`

### ebs_block_device

The `ebs_block_device` block has the following structure:

* `device_name` - (Required, Forces new resource, String) The device name of one or more block device mapping entries.
    * _Valid values:_ `cdrom<N>`, `disk<N>`, `floppy<N>`, where `<N>` is a disk number
* `delete_on_termination` - (Optional, Forces new resource, Boolean) Indicates whether the EBS volumes will be deleted once the instance for which they were created is terminated.
    * _Default value:_ `true`
* `iops` - (Optional, Forces new resource, Integer) The number of I/O operations per second the created volumes will support.
    * _Constraints:_ Required if `volume_type` is `io2`
* `snapshot_id` - (Optional, Forces new resource, String) The ID of an EBS snapshot that will be used to initialize the created EBS volumes.
    * _Constraints:_ If set, the `volume_size` attribute must be at least as large as the referenced snapshot
* `volume_size` - (Optional, Forces new resource, Integer) The size of created volumes in GiB.
    * _Constraints:_ Required unless `snapshot_id` is set. If `snapshot_id` is set and `volume_size` is omitted then the volume will have the same size as the selected snapshot
* `volume_type` - (Optional, Forces new resource, String) The type of EBS volume to create.
    * _Default value:_ `st2`

### ephemeral_block_device

The `ephemeral_block_device` block has the following structure:

* `device_name` - (Required, Forces new resource, String) The device name of one or more block device mapping entries.
    * _Valid values:_ `cdrom<N>`, `floppy<N>`, where `<N>` is a disk number
* `virtual_name` - (Required, Forces new resource, String) A name for the ephemeral device.
    * _Constraints:_ Must match the device name

## Attribute reference

### Supported attributes

In addition to all arguments above, the following attributes are exported:

* `arn` - (String) The Amazon Resource Name (ARN) of the image.
* `id` - (String) The ID of the created image.
* `image_owner_alias` - (String) The alias of the image owner.
* `image_type` - (String) The type of the image.
* `owner_id` - (String) The ID of the image owner.
* `platform` - (String) The platform of the image.
* `public` - (Boolean) Indicates whether the image has public launch permissions.
* `root_snapshot_id` - (String) The ID of the snapshot for the root volume (for EBS-backed images).
* `tags_all` - (Map of strings) Key-value pairs assigned to the image, including any tags inherited from the [`default_tags` configuration block][default-tags] if used within a provider configuration.

### Unsupported attributes

~> **Note** These attributes may be present in the `terraform.tfstate` file, but they have preset values and cannot be specified in configuration files.

The following attributes are not currently supported:

`boot_mode`, `deprecation_time`, `ebs_block_device.encrypted`, `ebs_block_device.kms_key_id`, `ebs_block_device.outpost_arn`, `ebs_block_device.throughput`, `ena_support`, `hypervisor`, `image_location`, `kernel_id`, `platform_details`, `ramdisk_id`, `usage_operation`.

## Timeouts

The `timeouts` block allows you to specify [timeouts] for certain actions:

* `create` - (Default `40 minutes`) Used when creating the image.
* `update` - (Default `40 minutes`) Used when updating the image.
* `delete` - (Default `90 minutes`) Used when deregistering the image.

## Import

`aws_ami` can be imported using the ID of the image, for example:

```
$ terraform import aws_ami.example cmi-12345678
```
