---
subcategory: "EC2 (Elastic Compute Cloud)"
layout: "aws"
page_title: "aws_instance"
description: |-
  Provides an EC2 instance resource. This allows instances to be created, updated, and deleted. Instances also support provisioning.
---

[default-tags]: https://www.terraform.io/docs/providers/aws/index.html#default_tags-configuration-block
[provisioning]: https://www.terraform.io/docs/provisioners/index.html
[timeouts]: https://www.terraform.io/docs/configuration/blocks/resources/syntax.html#operation-timeouts

# Resource: aws_instance

Provides an EC2 instance resource. This allows instances to be created, updated, and deleted.
Instances also support [provisioning].

## Example Usage

### Basic Example Using Image Lookup

```terraform
data "aws_ami" "selected" {
  most_recent = true
  owners      = ["self"]

  filter {
    name   = "name"
    values = ["ubuntu/images/hvm-ssd/ubuntu-focal-20.04-amd64-server"]
  }

  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }
}

resource "aws_instance" "example" {
  ami           = data.aws_ami.selected.id
  instance_type = "m1.micro"

  tags = {
    Name = "tf-instance"
  }
}
```

### Network Example

```terraform
resource "aws_vpc" "example" {
  cidr_block = "172.16.0.0/16"

  tags = {
    Name = "tf-vpc"
  }
}

resource "aws_subnet" "example" {
  vpc_id            = aws_vpc.example.id
  cidr_block        = "172.16.10.0/24"
  availability_zone = "ru-msk-vol52"

  tags = {
    Name = "tf-subnet"
  }
}

resource "aws_network_interface" "example" {
  subnet_id   = aws_subnet.example.id
  private_ips = ["172.16.10.100"]

  tags = {
    Name = "tf-primary-network-interface"
  }
}

resource "aws_instance" "example" {
  ami           = "cmi-12345678" # add image id, change instance type if needed
  instance_type = "m1.micro"

  network_interface {
    network_interface_id = aws_network_interface.example.id
    device_index         = 0
  }
}
```

## Argument Reference

The following arguments are supported:

* `affinity` - (Optional) The affinity setting for an instance on a dedicated host. Valid values are `default`, `host`. The parameter could be set to `host` only if `tenancy` is `host`.
* `ami` - (Optional) An image to use for the instance. Required unless `launch_template` is specified.
  If an image is specified in the launch template, setting `ami` will override it.
* `associate_public_ip_address` - (Optional, Conflicts with `network_interface`) Whether to associate a public IP address with an instance in a VPC.
  The address will be assigned to the `eth0` interface if there are free allocated external addresses.
  This operation is available only for instances running in the VPC and for new network interfaces.
* `availability_zone` - (Optional) AZ to start the instance in.
* `disable_api_termination` - (Optional) If `true`, disables the possibility to terminate an instance via API.
* `ebs_block_device` - (Optional) One or more configuration blocks with additional EBS block devices to attach to the instance. Block device configurations only apply on resource creation. See [Block Devices](#ebs-ephemeral-and-root-block-devices) below for details on attributes and drift detection. When accessing this as an attribute reference, it is a set of objects.
* `ephemeral_block_device` - (Optional) One or more configuration blocks to customize ephemeral volumes on the instance. See [Block Devices](#ebs-ephemeral-and-root-block-devices) below for details. When accessing this as an attribute reference, it is a set of objects.
* `host_id` - (Optional) The ID of the dedicated host that the instance will be assigned to.
* `instance_initiated_shutdown_behavior` - (Optional) Shutdown behavior for the instance. Valid values are `stop`, `terminate`.
* `instance_type` - (Optional) The instance type to use for the instance. Updates to this field will trigger a stop/start of the EC2 instance.
* `key_name` - (Optional) Key name of the key pair to use for the instance; which can be managed using [the `aws_key_pair` resource](key_pair.html.markdown).
* `launch_template` - (Optional) Specifies a launch template to configure the instance. Parameters configured on this resource will override the corresponding parameters in the Launch Template.
  See [Launch Template Specification](#launch-template-specification) below for more details.
* `monitoring` - (Optional) If `true`, the launched EC2 instance will have detailed monitoring enabled.
* `network_interface` - (Optional, Conflicts with `associate_public_ip_address`, `private_ip`, `secondary_private_ips`, `subnet_id`, `vpc_security_group_ids`) Customize network interfaces to be attached at instance boot time. See [Network Interfaces](#network-interfaces) below for more details.
* `placement_group` - (Optional) Placement group to start the instance in.
* `private_ip` - (Optional, Conflicts with `network_interface`) Private IP address to associate with the instance in a VPC.
* `root_block_device` - (Optional) Configuration block to customize details about the root block device of the instance. See [Block Devices](#ebs-ephemeral-and-root-block-devices) below for details. When accessing this as an attribute reference, it is a list containing one object.
* `secondary_private_ips` - (Optional, Conflicts with `network_interface`) A list of secondary private IPv4 addresses to assign to the instance's primary network interface in a VPC. Currently, only specifying the primary private IP address is supported.
* `source_dest_check` - (Optional) Controls if traffic is routed to the instance when the destination address does not match the instance. Defaults to `true`.
* `subnet_id` - (Optional) VPC subnet ID to launch in.
* `tags` - (Optional) A map of tags to assign to the resource. Note that these tags apply to the instance and not block storage devices. If configured with a provider [`default_tags` configuration block][default-tags] present, tags with matching keys will overwrite those defined at the provider-level.
* `tenancy` - (Optional) The placement type. Valid values are `default`, `host`.

~> **NOTE:** If you use the `host` value, you may encounter the `NotEnoughResourcesForInstanceType` error when running an instance. To avoid this, it is recommended to provide either the `subnet_id` argument or the `availability_zone` argument.

* `user_data` - (Optional, Conflicts with `user_data_base64`) User data to provide when launching the instance. Do not pass gzip-compressed data via this argument; see `user_data_base64` instead. Updates to this field will trigger a stop/start of the EC2 instance by default. If the `user_data_replace_on_change` is set then updates to this field will trigger a destroy and recreate.
* `user_data_base64` - (Optional, Conflicts with `user_data`) Can be used instead of `user_data` to pass base64-encoded binary data directly. Use this instead of `user_data` whenever the value is not a valid UTF-8 string. For example, gzip-encoded user data must be base64-encoded and passed via this argument to avoid corruption. Updates to this field will trigger a stop/start of the EC2 instance by default. If the `user_data_replace_on_change` is set then updates to this field will trigger a destroy and recreate.
* `user_data_replace_on_change` - (Optional) When used in combination with `user_data` or `user_data_base64` will trigger a destroy and recreate when set to `true`. Defaults to `false`.
* `volume_tags` - (Optional) A map of tags to assign, at instance-creation time, to root and EBS volumes.

~> **NOTE:** Do not use `volume_tags` if you plan to manage block device tags outside the `aws_instance` configuration, such as using `tags` in an [`aws_ebs_volume`](ebs_volume.html.markdown) resource attached via [`aws_volume_attachment`](volume_attachment.html.markdown). Doing so will result in resource cycling and inconsistent behavior.

* `vpc_security_group_ids` - (Optional, Conflicts with `network_interface`) A list of security group IDs to associate with.

### EBS, Ephemeral, and Root Block Devices

Each of the `*_block_device` attributes control a portion of the EC2 instance's "Block Device Mapping".

The `root_block_device` block supports the following:

* `delete_on_termination` - (Optional) Whether the volume should be destroyed on instance termination. Defaults to `true`.
* `iops` - (Optional) Amount of provisioned IOPS. Only valid for volume_type of `io2`.
* `tags` - (Optional) A map of tags to assign to the device.
* `volume_size` - (Optional) Size of the volume in GiB.
* `volume_type` - (Optional) Type of volume. Valid values are `st2`, `gp2`, `io2`.

Each `ebs_block_device` block supports the following:

* `delete_on_termination` - (Optional) Whether the volume should be destroyed on instance termination. Defaults to `true`.
* `device_name` - (Required) Name of the device to mount.
* `iops` - (Optional) Amount of provisioned IOPS. Only valid for volume_type of `io2`.
* `snapshot_id` - (Optional) Snapshot ID to mount.
* `tags` - (Optional) A map of tags to assign to the device.
* `volume_size` - (Optional) Size of the volume in gibibytes (GiB).
* `volume_type` - (Optional) Type of volume. Valid values are `st2`, `gp2`, `io2`.

~> **NOTE:** Currently, changes to the `ebs_block_device` configuration of _existing_ resources cannot be automatically detected by Terraform.
To manage changes and attachments of an EBS block to an instance, use the [`aws_ebs_volume`](ebs_volume.html.markdown) and [`aws_volume_attachment`](volume_attachment.html.markdown) resources instead.
If you use `ebs_block_device` on an `aws_instance`, Terraform will assume management over the full set of non-root EBS block devices for the instance, treating additional block devices as drift.
For this reason, `ebs_block_device` cannot be mixed with external `aws_ebs_volume` and `aws_volume_attachment` resources for a given instance.

Each `ephemeral_block_device` block supports the following:

* `device_name` - (Required) The name of the block device to mount on the instance.
* `no_device` - (Optional) Suppresses the specified device included in the block device mapping.
* `virtual_name` - (Optional) A name for the ephemeral device. Must match with the device name.

### Network Interfaces

Each of the `network_interface` blocks attach a network interface to an EC2 instance during boot time.
However, because the network interface is attached at boot-time, replacing/modifying the network interface **WILL** trigger a recreation of the EC2 Instance.
If you should need at any point to detach/modify/re-attach a network interface to the instance, use the [`aws_network_interface`](network_interface.html.markdown) or [`aws_network_interface_attachment`](network_interface_attachment.html.markdown) resources instead.

The `network_interface` configuration block _does_, however, allow users to supply their own network interface to be used as the default network interface on an EC2 instance, attached at `eth0`.

Each `network_interface` block supports the following:

* `delete_on_termination` - (Optional) Whether to delete the network interface on instance termination. Defaults to `false`.
  Currently, the only valid value is `false`, as this option is only supported when creating new network interfaces during instance launching.
* `device_index` - (Required) Integer index of the network interface attachment.
* `network_interface_id` - (Required) ID of the network interface to attach.

### Launch Template Specification

-> **Note:** Launch template parameters will be used only once during instance creation. If you want to update existing instance you need to change parameters
directly. Updating Launch Template Specification will force a new instance.

Any other instance parameters that you specify will override the same parameters in the launch template.

The `launch_template` block supports the following:

* `id` - The ID of the launch template.
* `name` - The name of the launch template.
* `version` - Template version. Valid values are a specific version number, `$Latest`, `$Default`. Defaults to `$Default`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `arn` - The ARN of the instance.
* `instance_state` - The state of the instance. One of: `pending`, `running`, `shutting-down`, `terminated`, `stopping`, `stopped`.
* `primary_network_interface_id` - The ID of the instance's primary network interface.
* `private_dns` - The private DNS name assigned to the instance. For EC2-VPC, this is only available if you've enabled DNS hostnames for your VPC.
* `public_dns` - The public DNS name assigned to the instance. For EC2-VPC, this is only available if you've enabled DNS hostnames for your VPC.
* `public_ip` - The public IP address assigned to the instance, if applicable. **NOTE**: If you are using an [`aws_eip`](eip.html.markdown) with your instance, you should refer to the EIP's address directly and not use `public_ip` as this field will change after the EIP is attached.
* `security_groups` - The list of security group names associated with the instance.
* `tags_all` - A map of tags assigned to the resource, including those inherited from the provider [`default_tags` configuration block][default-tags].

For `ebs_block_device`, in addition to the arguments above, the following attribute is exported:

* `volume_id` - ID of the volume. For example, the ID can be accessed like this, `aws_instance.web.ebs_block_device.2.volume_id`.

For `root_block_device`, in addition to the arguments above, the following attributes are exported:

* `volume_id` - ID of the volume. For example, the ID can be accessed like this, `aws_instance.web.root_block_device.0.volume_id`.
* `device_name` - Device name, e.g., `disk1`.

->  **Unsupported attributes**
These exported attributes are currently unsupported:

* `capacity_reservation_specification` - Describes an instance's Capacity Reservation targeting option. Always empty.
    * `capacity_reservation_preference` - Indicates the instance's Capacity Reservation preferences.
    * `capacity_reservation_target` - Information about the target Capacity Reservation.
        * `capacity_reservation_id` - The ID of the Capacity Reservation in which to run the instance.
        * `capacity_reservation_resource_group_arn` - The ARN of the Capacity Reservation resource group in which to run the instance.
* `cpu_core_count` - Sets the number of CPU cores for an instance. Always empty.
* `cpu_threads_per_core` - If set to `1`, hyperthreading is disabled on the launched instance. Always empty.
* `credit_specification` - Configuration block for customizing the credit specification of the instance. Always empty.
    * `cpu_credits` - The credit option for CPU usage.
* `ebs_block_device`:
    * `encrypted` - Whether to enable volume encryption. Always `false`.
    * `kms_key_id` - The ARN of the KMS Key to use when encrypting the volume. Always `""`.
    * `throughput` - Throughput to provision for a volume in mebibytes per second (MiB/s). Always `0`.
* `ebs_optimized` - If true, the launched EC2 instance will be EBS-optimized. Always `false`.
* `enclave_options` - Enable Nitro Enclaves on launched instances. Always empty.
    * `enabled` - Whether Nitro Enclaves will be enabled on the instance.
* `get_password_data` - If true, wait for password data to become available and retrieve it. Always `false`.
* `hibernation` - If true, the launched EC2 instance will support hibernation. Always empty.
* `iam_instance_profile` - IAM Instance Profile to launch the instance with. Always `""`.
* `ipv6_address_count`- A number of IPv6 addresses to associate with the primary network interface. Always `0`.
* `ipv6_addresses` - Specify one or more IPv6 addresses from the range of the subnet to associate with the primary network interface. Always empty.
* `maintenance_options` - The maintenance and recovery options for the instance. Always empty.
    * `auto_recovery` - The automatic recovery behavior of the Instance.
* `metadata_options` - Customize the metadata options of the instance. Always empty.
    * `http_endpoint` - Whether the metadata service is available.
    * `http_tokens` - Whether the metadata service requires session tokens, also referred to as _Instance Metadata Service Version 2 (IMDSv2)_.
    * `http_put_response_hop_limit` - Desired HTTP PUT response hop limit for instance metadata requests.
    * `instance_metadata_tags` - Enables or disables access to instance tags from the instance metadata service.
* `network_interface`:
    * `network_card_index` - Integer index of the network card. Limited by instance type. Always `0`.
* `outpost_arn` - The ARN of the Outpost the instance is assigned to. Always `""`.
* `password_data` - Base-64 encoded encrypted password data for the instance. Always `""`.
* `placement_partition_number` - The number of the partition the instance is in. Always empty.
* `root_block_device`:
    * `encrypted` - Whether to enable volume encryption. Always `false`.
    * `kms_key_id` - The ARN of the KMS Key to use when encrypting the volume. Always `""`.
    * `throughput` - Throughput to provision for a volume in mebibytes per second (MiB/s). Always `0`.

### Timeouts

The `timeouts` block allows you to specify [timeouts] for certain actions:

* `create` - (Default `10 minutes`) Used when launching the instance (until it reaches the initial `running` state).
* `update` - (Default `10 minutes`) Used when stopping and starting the instance when necessary during update - e.g., when changing instance type.
* `delete` - (Default `20 minutes`) Used when terminating the instance.

## Import

Instances can be imported using the `id`, e.g.,

```
$ terraform import aws_instance.web i-12345678
```
