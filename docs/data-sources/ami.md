---
subcategory: "EC2 (Elastic Compute Cloud)"
layout: "aws"
page_title: "aws_ami"
description: |-
  Provides information about an Amazon Machine Image (AMI).
---

[describe-images]: https://docs.k2.cloud/en/api/ec2/actions/images/DescribeImages.html

# Data Source: aws_ami

Provides information about an Amazon Machine Image (AMI).

## Example usage

### Basic example

```terraform
data "aws_ami" "example" {
  executable_users = ["self"]
  most_recent      = true
  name_regex       = "^example\\d{1}"
  owners           = ["self"]

  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }
}
```

## Argument reference

The following arguments are supported:

* `owners` - (Required, List of strings) List of image owners to limit search. At least one value must be specified.
    * _Valid values:_ `project@customer` or `self`
* `executable_users` - (Optional, List of strings) Limits search to project with the _explicit_ launch permission on the image.
    * _Valid values:_ `all`, `project@customer`, or `self`
* `filter` - (Optional, [Block](#filter)) One or more name/value pairs to use as filters.
    * _Valid values:_ See supported names and values in [EC2 API documentation][describe-images]
* `most_recent` - (Optional, Boolean) If more than one result is returned, use the most recent image.
    * _Default value:_ `false`
* `name_regex` - (Optional, String) A regex string to apply to the image list returned by the EC2 API.
  It is recommended to combine this with other options to narrow down the list that the EC2 API returns.

~> **Note** The search must return a single match, otherwise Terraform will fail.
Ensure that your search is specific enough to return the ID of
a single image only, or use `most_recent` to choose the most recent one. If
you want to match multiple images, use the [`aws_ami_ids`](ami_ids.md) data source instead.

### filter

* `name` - (Required, String) The name of the filter.
    * _Constraints:_ Filter names are case-sensitive
* `values` - (Required, List of strings) One or more filter values.
    * _Constraints:_ Filter values are case-sensitive

## Attribute reference

### Supported attributes

In addition to all arguments above, the following attributes are exported:

* `arn` - (String) The Amazon Resource Name (ARN) of the image.
* `architecture` - (String) The OS architecture of the image.
* `block_device_mappings` - ([Block](#block_device_mappings)) Set of objects with block device mappings of the image.
  The structure of this block is [described below](#block_device_mappings).
* `description` - (String) The description of the image that was provided during image
  creation.
* `id` - (String) The ID of the image.
* `image_id` - (String) The ID of the image. Should be the same as the resource `id`.
* `image_owner_alias` - (String) The alias of the image owner.
* `image_type` - (String) The type of the image.
* `name` - (String) The name of the image that was provided during image creation.
* `owner_id` - (String) The ID of the image owner.
* `platform` - (String) The platform of the image.
* `public` - (Boolean) Indicates whether the image has public launch permissions.
* `root_device_name` - (String) The device name of the root device.
* `root_device_type` - (String) The type of the root device.
* `root_snapshot_id` - (String) The ID of the snapshot associated with the root device, if any
  (only applies to `ebs` root devices).
* `state` - (String) The current state of the image. If the state is `available`, the image
  is successfully registered and can be used to launch an instance.
* `tags` - (Map of strings) Key-value pairs assigned to the resource.
* `virtualization_type` - (String) The type of virtualization of the image.

#### block_device_mappings

The `block_device_mappings` block has the following structure:

* `device_name` - (String) The physical name of the device.
* `ebs` - (Map of strings) Map containing EBS information, if the device is EBS based.
  Unlike most object attributes, these are accessed directly (e.g., `ebs.volume_size` or `ebs["volume_size"]`) rather than accessed through the first element of a list (e.g., `ebs[0].volume_size`).
  The structure of this block is [described below](#ebs).
* `no_device` - (String) The value of `no_device`, if the device is not mapped.
* `virtual_name` - (String) The virtual device name (for instance stores).

##### ebs

The `ebs` block is a part of the [`block_device_mappings`](#block_device_mappings) block. It has the following structure:

* `delete_on_termination` - (String) Indicates whether the EBS volume will be deleted on termination.
* `encrypted` - (String) Indicates whether the EBS volume is encrypted.
* `iops` - (String) `0` if the EBS volume is not a provisioned IOPS image, otherwise the supported IOPS count.
* `snapshot_id` - (String) The ID of the snapshot.
* `throughput` - (String) The throughput of the volume in MiB/s.
* `volume_size` - (String) The size of the volume in GiB.
* `volume_type` - (String) The volume type.

### Unsupported attributes

~> **Note** These attributes may be present in the `terraform.tfstate` file, but they have preset values and cannot be specified in configuration files.

The following attributes are not currently supported:

`boot_mode`, `creation_date`, `deprecation_time`, `ena_support`, `hypervisor`, `image_location`, `kernel_id`, `platform_details`, `product_codes`, `ramdisk_id`, `sriov_net_support`, `state_reason`, `usage_operation`.
