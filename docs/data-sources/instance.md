---
subcategory: "EC2 (Elastic Compute Cloud)"
layout: "aws"
page_title: "aws_instance"
description: |-
  Provides information about an instance.
---

[base64decode-function]: https://www.terraform.io/docs/configuration/functions/base64decode.html
[describe-instances]: https://docs.k2.cloud/en/api/ec2/actions/instances/DescribeInstances.html

# Data Source: aws_instance

Provides information about an instance.

## Example Usage

```terraform
data "aws_instance" "selected" {
  instance_id = "i-12345678"

  filter {
    name   = "image-id"
    values = ["cmi-12345678"]
  }

  instance_tags = {
    type = "test"
  }

  filter {
    name   = "tag:Name"
    values = ["example"]
  }
}
```

## Argument Reference

* `filter` - (Optional) One or more name/value pairs to use as filters.
    * _Valid values:_ See supported names and values in [EC2 API documentation][describe-instances]
* `get_user_data` - (Optional) Retrieve Base64 encoded user data contents into the `user_data_base64` attribute.
  A SHA-1 hash of the user data contents will always be present in the `user_data` attribute.
    * _Default value:_ `false`
* `instance_id` - (Optional) Specify the exact instance ID with which to populate the data source.
* `instance_tags` - (Optional) Map of tags, each pair of which must exactly match a pair on the desired instance.

~> **Note** At least one of the arguments `filter`, `instance_tags`, or `instance_id` must be specified.

~> **Note** If anything other than a single match is returned by the search,
Terraform will fail. Ensure that your search is specific enough to return
a single instance ID only.

## Attribute Reference

### Supported attributes

In addition to all arguments above, the following attributes are exported:

* `affinity` - The affinity setting for an instance on a dedicated host.
* `ami` - The ID of the image used to launch the instance.
* `arn` - The Amazon Resource Name (ARN) of the instance.
* `associate_public_ip_address` - Whether the instance is associated with a public IP address or not.
* `availability_zone` - The availability zone of the instance.
* `ebs_block_device` - The EBS block device mappings of the instance.
* `ephemeral_block_device` - The ephemeral block device mappings of the instance.
  The structure of this block is [described below](#ephemeral_block_device).
* `host_id` - The ID of the dedicated host that the instance will be assigned to.
* `id` - The ID of the instance.
* `instance_state` - The state of the instance.
    * _Valid values:_ `pending`, `running`, `shutting-down`, `terminated`, `stopping`, `stopped`
* `instance_type` - The type of the instance.
* `key_name` - The key name of the instance.
* `monitoring` - Whether detailed monitoring is enabled or disabled for the instance.
* `network_interface_id` - The ID of the network interface that was created with the instance.
* `placement_group` - The placement group of the instance.
* `private_dns` - The private DNS name assigned to the instance.
* `private_ip` - The private IP address assigned to the instance.
* `secondary_private_ips` - The secondary private IPv4 addresses assigned to the instance's primary network interface in a VPC.
* `public_dns` - The public DNS name assigned to the instance.
* `public_ip` - The public IP address assigned to the instance, if applicable.
    ~> **Note** If you are using an [`aws_eip`](../resources/eip.md) with your instance, you should refer to the EIP's address directly and not use `public_ip`, as this field will change after the EIP is attached.
* `root_block_device` - The root block device mappings of the instance.
  The structure of this block is [described below](#root_block_device).
* `security_groups` - The associated security groups.
* `source_dest_check` - Whether the network interface performs source/destination checking.
* `subnet_id` - The ID of the subnet.
* `user_data` - SHA-1 hash of user data supplied to the instance.
* `user_data_base64` - Base64 encoded contents of user data supplied to the instance. Valid UTF-8 contents can be decoded with the [`base64decode` function][base64decode-function].
    * _Constraints:_ This attribute is only exported if `get_user_data` is true
* `tags` - Map of tags assigned to the instance.
* `tenancy` - The placement type.
* `vpc_security_group_ids` - The associated security groups in a non-default VPC.

#### ebs_block_device

The `ebs_block_device` block has the following structure:

* `delete_on_termination` - If the EBS volume will be deleted on termination.
* `device_name` - The physical name of the device.
* `iops` - `0` if the EBS volume is not a provisioned IOPS image, otherwise the supported IOPS count.
* `snapshot_id` - The ID of the snapshot.
* `volume_size` - The size of the volume, in GiB.
* `volume_type` - The volume type.

#### ephemeral_block_device

The `ephemeral_block_device` block has the following structure:

* `device_name` - The physical name of the device.
* `no_device` - Whether the specified device included in the device mapping was suppressed or not.
* `virtual_name` - The virtual device name.

#### root_block_device

The `root_block_device` block has the following structure:

* `device_name` - The physical name of the device.
* `delete_on_termination` - Indicates whether the root block device will be deleted on termination.
* `iops` - `0` if the volume is not a provisioned IOPS image, otherwise the supported IOPS count.
* `volume_size` - The size of the volume, in GiB.
* `volume_type` - The type of the volume.

### Unsupported attributes

~> **Note** These attributes may be present in the `terraform.tfstate` file, but they have preset values and cannot be specified in configuration files.

The following attributes are not currently supported:

`credit_specification`, `ebs_block_device.encrypted`, `ebs_block_device.kms_key_id`, `ebs_block_device.throughput`, `ebs_optimized`, `enclave_options`, `get_password_data`, `iam_instance_profile`, `ipv6_addresses`, `maintenance_options`, `metadata_options`, `outpost_arn`, `password_data`, `placement_partition_number`, `root_block_device.encrypted`, `root_block_device.kms_key_id`, `root_block_device.throughput`.
