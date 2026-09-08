---
subcategory: "EC2 (Elastic Compute Cloud)"
layout: "aws"
page_title: "aws_instance"
description: |-
  Provides information about an instance.
---

[base64decode-function]: https://developer.hashicorp.com/terraform/language/functions/base64decode
[describe-instances]: https://docs.k2.cloud/en/api/ec2/actions/instances/DescribeInstances.html

# Data Source: aws_instance

Provides information about an instance.

## Example usage

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

## Argument reference

* `filter` - (Optional, [Block](#filter)) One or more name/value pairs to use as filters.
    * _Valid values:_ See supported names and values in the [EC2 API documentation][describe-instances]
* `get_user_data` - (Optional, Boolean) Indicates whether to retrieve Base64 encoded user data contents into the `user_data_base64` attribute.
  A SHA-1 hash of the user data contents will always be present in the `user_data` attribute.
    * _Default value:_ `false`
* `instance_id` - (Optional, String) Specify the exact instance ID with which to populate the data source.
* `instance_tags` - (Optional, Map of strings) Key-value pairs. Must exactly match pairs on the desired instance.

~> **Note** At least one of the arguments `filter`, `instance_tags`, or `instance_id` must be specified.

~> **Note** If anything other than a single match is returned by the search, Terraform will fail.
Ensure that your search is specific enough to return a single instance ID only.

### filter

* `name` - (Required, String) The name of the filter.
    * _Constraints:_ Filter names are case-sensitive
* `values` - (Required, List of strings) One or more filter values.
    * _Constraints:_ Filter values are case-sensitive

## Attribute reference

### Supported attributes

In addition to all arguments above, the following attributes are exported:

* `affinity` - (String) The affinity setting for an instance on a dedicated host.
* `ami` - (String) The ID of the image used to launch the instance.
* `arn` - (String) The Amazon Resource Name (ARN) of the instance.
* `associate_public_ip_address` - (Boolean) Indicates whether the instance has an associated a public IP address.
* `availability_zone` - (String) The availability zone of the instance.
* `disable_api_termination` - (Boolean) Indicates whether API termination is disabled for the instance.
* `ebs_block_device` - ([Block](#ebs_block_device)) The EBS block device mappings of the instance.
* `host_id` - (String) The ID of the dedicated host where the instance will run.
* `id` - (String) The ID of the instance.
* `instance_state` - (String) The state of the instance.
    * _Valid values:_ `pending`, `running`, `shutting-down`, `stopped`, `stopping`, `terminated`
* `instance_type` - (String) The type of the instance.
* `key_name` - (String) The key name of the instance.
* `monitoring` - (Boolean) Indicates whether detailed monitoring is enabled for the instance.
* `network_interface_id` - (String) The ID of the network interface that was created with the instance.
* `placement_group` - (String) The placement group of the instance.
* `private_dns` - (String) The private DNS name assigned to the instance.
* `private_ip` - (String) The private IP address associated with the instance.
* `public_dns` - (String) The public DNS name assigned to the instance.
* `public_ip` - (String) The public IP address associated with the instance, if applicable.

    ~> **Note** If you are using an [`aws_eip`](../resources/eip.md) with your instance, you should refer to the EIP's address directly and not use `public_ip`, as this field will change after the EIP is attached.

* `root_block_device` - ([Block](#root_block_device)) The root block device mappings of the instance.
* `secondary_private_ips` - (Set of strings) The secondary private IPv4 addresses associated with the instance's primary network interface in a VPC.
* `security_groups` - (Set of strings) Security groups associated with the instance.
* `source_dest_check` - (Boolean) Indicates whether the network interface performs source/destination checking.
* `subnet_id` - (String) The ID of the subnet.
* `tags` - (Map of strings) Key-value pairs assigned to the instance.
* `tenancy` - (String) The placement type.
* `user_data` - (String) SHA-1 hash of user data supplied to the instance.
* `user_data_base64` - (String) Base64 encoded contents of user data supplied to the instance.
  Valid UTF-8 contents can be decoded with the [`base64decode` function][base64decode-function].
    * _Constraints:_ This attribute is only exported if `get_user_data` is true
* `vpc_security_group_ids` - (Set of strings) Security groups associated with the instance in a non-default VPC.

#### ebs_block_device

The `ebs_block_device` block has the following structure:

* `delete_on_termination` - (Boolean) Indicates whether the EBS volume will be deleted on instance termination.
* `device_name` - (String) The physical name of the device.
* `iops` - (Integer) The number of I/O operations per second for the volume.
    * _Constraints:_ `0` if the EBS volume is not a provisioned IOPS image, otherwise the supported IOPS count.
* `snapshot_id` - (String) The ID of the snapshot.
* `tags` - (Map of strings) Key-value pairs assigned to the volume.
* `volume_id` - (String) The ID of the EBS volume.
* `volume_size` - (Integer) The size of the volume in GiB.
* `volume_type` - (String) The volume type.

#### root_block_device

The `root_block_device` block has the following structure:

* `delete_on_termination` - (Boolean) Indicates whether the root block device will be deleted on instance termination.
* `device_name` - (String) The physical name of the device.
* `iops` - (Integer) The number of I/O operations per second for the volume.
    * _Constraints:_ `0` if the volume is not a provisioned IOPS image, otherwise the supported IOPS count.
* `tags` - (Map of strings) Key-value pairs assigned to the volume.
* `volume_id` - (String) The ID of the root block device volume.
* `volume_size` - (Integer) The size of the volume in GiB.
* `volume_type` - (String) The type of the volume.

### Unsupported attributes

~> **Note** These attributes may be present in the `terraform.tfstate` file, but they have preset values and cannot be specified in configuration files.

The following attributes are not currently supported:

`credit_specification`, `ebs_block_device.encrypted`, `ebs_block_device.kms_key_id`, `ebs_block_device.throughput`, `ebs_optimized`, `enclave_options`, `ephemeral_block_device`, `get_password_data`, `iam_instance_profile`, `ipv6_addresses`, `maintenance_options`, `metadata_options`, `outpost_arn`, `password_data`, `placement_partition_number`, `root_block_device.encrypted`, `root_block_device.kms_key_id`, `root_block_device.throughput`.
