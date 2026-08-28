---
subcategory: "EC2 (Elastic Compute Cloud)"
layout: "aws"
page_title: "aws_instance"
description: |-
  Manages an EC2 instance.
---

[default-tags]: https://www.terraform.io/docs/providers/aws/index.html#default_tags-configuration-block
[provisioning]: https://www.terraform.io/docs/provisioners/index.html
[timeouts]: https://developer.hashicorp.com/terraform/plugin/framework/resources/timeouts

# Resource: aws_instance

Manages an EC2 instance. This allows instances to be created, updated, and deleted.
Instances also support [provisioning].

## Example Usage

### Basic example: using image lookup

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

* `affinity` - (Optional) The affinity setting for an instance on a dedicated host.
    * _Valid values:_ `default`, `host`
    * _Constraints:_ The parameter could be set to `host` only if `tenancy` is `host`
* `ami` - (Optional) An image to use for the instance.
  If an image is specified in the launch template, the `ami` setting will override it.
    * _Constraints:_ Required unless `launch_template` is specified
* `associate_public_ip_address` - (Optional) Indicates whether to associate a public IP address with an instance in a VPC.
    * _Constraints:_ Conflicts with the `network_interface` argument

    ~> **Note** The address will be assigned to the `eth0` interface only if there are free allocated external addresses.
    This operation is available only for instances running in the VPC and for new network interfaces.

    ~> **Important** The address is picked arbitrarily among all free allocated addresses of the project, so the instance can take an address that is intended for another resource, for example, for an [`aws_eip`](eip.md) association or for an [`aws_nat_gateway`](nat_gateway.md) that allocates its addresses automatically.
    Terraform doesn't manage the creation order of such resources on its own, so it's recommended to set it explicitly via `depends_on`: the instance must be created after the resources that need their own external addresses, for example, `depends_on = [aws_nat_gateway.example]`.

* `availability_zone` - (Optional) An availability zone to start the instance in.
* `disable_api_termination` - (Optional) If `true`, disables the possibility to terminate an instance via API.
* `ebs_block_device` - (Optional) One or more configuration blocks with additional EBS block devices to attach to the instance. The structure of this block and details on drift detection are [described below](#ebs_block_device). When accessing this as an attribute reference, it is a set of objects.
    * _Constraints:_ Block device configurations are applied only when the resource is created
* `ephemeral_block_device` - (Optional) One or more configuration blocks to customize ephemeral volumes on the instance. The structure of this block is [described below](#ephemeral_block_device). When accessing this as an attribute reference, it is a set of objects.
* `host_id` - (Optional) The ID of the dedicated host that the instance will be assigned to.
* `instance_initiated_shutdown_behavior` - (Optional) Shutdown behavior for the instance.
    * _Valid values:_ `stop`, `terminate`
* `instance_type` - (Optional) The instance type to use for the instance. Updates to this field will trigger a stop/start of the EC2 instance.
* `key_name` - (Optional) Key name of the key pair to use for the instance; which can be managed using [the `aws_key_pair` resource](key_pair.md).
* `launch_template` - (Optional) Specifies a launch template to configure the instance. Parameters configured on this resource will override the corresponding parameters in the launch template.
  The structure of this block is [described below](#launch_template).
* `monitoring` - (Optional) If `true`, the launched EC2 instance will have detailed monitoring enabled.
* `network_interface` - (Optional) Customize network interfaces to be attached at instance boot time. The structure of this block is [described below](#network_interface).
    * _Constraints:_ Conflicts with `associate_public_ip_address`, `private_ip`, `secondary_private_ips`, `subnet_id`, `vpc_security_group_ids` arguments
* `placement_group` - (Optional) Placement group to start the instance in.
* `private_ip` - (Optional) Private IP address to associate with the instance in a VPC.
    * _Constraints:_ Conflicts with the `network_interface` argument
* `root_block_device` - (Optional) Root block device of the instance. The structure of this block is [described below](#root_block_device). When accessing this as an attribute reference, it is a list containing one object.
* `secondary_private_ips` - (Optional) List of secondary private IPv4 addresses to assign to the instance's primary network interface in a VPC.
    * _Constraints:_
        * Conflicts with the `network_interface` argument
        * Only the primary private IP address can be specified
* `source_dest_check` - (Optional) Controls if traffic is routed to the instance when the destination address does not match the instance.
    * _Default value:_ `true`
* `subnet_id` - (Optional) The ID of a subnet to launch in.
* `tags` - (Optional) Map of tags to assign to the instance. If a provider [`default_tags` configuration block][default-tags] is used, tags with matching keys will overwrite those defined at the provider level.
    * _Constraints:_ These tags apply to the instance and not block storage devices
* `tenancy` - (Optional) The placement type.
    * _Valid values:_ `default`, `host`

    ~> **Note** If you use the `host` value, you may encounter the `NotEnoughResourcesForInstanceType` error when running an instance. To avoid this, it is recommended to provide either the `subnet_id` argument or the `availability_zone` argument.

* `user_data` - (Optional) User data to provide when launching the instance. Do not pass gzip-compressed data via this argument; see `user_data_base64` instead. Updates to this field will trigger a stop/start of the EC2 instance by default. If the `user_data_replace_on_change` is set then updates to this field will trigger a destroy and recreate.
    * _Constraints:_ Conflicts with the `user_data_base64` argument
* `user_data_base64` - (Optional) Can be used instead of `user_data` to pass base64-encoded binary data directly. Use this instead of `user_data` whenever the value is not a valid UTF-8 string. For example, gzip-encoded user data must be base64-encoded and passed via this argument to avoid corruption. Updates to this field will trigger a stop/start of the EC2 instance by default. If the `user_data_replace_on_change` is set then updates to this field will trigger a destroy and recreate.
    * _Constraints:_ Conflicts with the `user_data` argument
* `user_data_replace_on_change` - (Optional) When used in combination with `user_data` or `user_data_base64` will trigger a destroy and recreate when set to `true`.
    * _Default value:_ `false`
* `volume_tags` - (Optional) Map of tags to assign to root and EBS volumes when the instance is created.

    ~>**Note** Do not use `volume_tags` if you plan to manage block device tags outside the `aws_instance` configuration, such as using `tags` in an [`aws_ebs_volume`](ebs_volume.md) resource attached via [`aws_volume_attachment`](volume_attachment.md). Doing so will result in resource cycling and inconsistent behavior.

* `vpc_security_group_ids` - (Optional) List of security group IDs to associate with.
    * _Constraints:_ Conflicts with the `network_interface` argument

### ebs_block_device

The following arguments are required:

* `device_name` - (Required) Name of the device to mount.

The following arguments are optional:

* `delete_on_termination` - (Optional) Whether the volume should be destroyed on instance termination.
    * _Default value:_ `true`
* `encrypted` - (Optional) Indicates whether the EBS volume is encrypted.
* `iops` - (Optional) Amount of provisioned IOPS.
    * _Constraints:_ Only valid for the volume type `io2`
* `snapshot_id` - (Optional) The ID of the snapshot to mount.
* `tags` - (Optional) Map of tags to assign to the device.
* `volume_size` - (Optional) Size of the volume in GiB.
* `volume_type` - (Optional) Type of volume.

~> **Note** Currently, changes to the `ebs_block_device` configuration of _existing_ resources cannot be automatically detected by Terraform.
To manage changes and attachments of an EBS block to an instance, use the [`aws_ebs_volume`](ebs_volume.md) and [`aws_volume_attachment`](volume_attachment.md) resources instead.
If you use `ebs_block_device` on an `aws_instance`, Terraform will assume management over the full set of non-root EBS block devices for the instance, treating additional block devices as drift.
For this reason, `ebs_block_device` cannot be mixed with external `aws_ebs_volume` and `aws_volume_attachment` resources for a given instance.

### ephemeral_block_device

The following arguments are required:

* `device_name` - (Required) The name of the block device to mount on the instance.

The following arguments are optional:

* `no_device` - (Optional) Suppresses the specified device included in the block device mapping.
* `virtual_name` - (Optional) A name for the ephemeral device. Must match with the device name.

### network_interface

Each of the `network_interface` blocks attaches a network interface to an EC2 instance during boot time.
However, because the network interface is attached at boot-time, replacing/modifying the network interface **WILL** trigger a recreation of the EC2 instance.
If you should need at any point to detach/modify/re-attach a network interface to the instance, use the [`aws_network_interface`](network_interface.md) or [`aws_network_interface_attachment`](network_interface_attachment.md) resources instead.

The `network_interface` configuration block _does_, however, allow users to supply their own network interface to be used as the default network interface on an EC2 instance, attached at `eth0`.

The following arguments are required:

* `device_index` - (Required) Integer index of the network interface attachment.
* `network_interface_id` - (Required) The ID of the network interface to attach.

The following arguments are optional:

* `delete_on_termination` - (Optional) Whether to delete the network interface on instance termination.
    * _Default value:_ `false`
    * _Constraints:_ Currently, the only valid value is `false`, as this option is only supported when creating new network interfaces during instance launching

### launch_template

~> **Note** Launch template parameters will be used only once when the instance is created.
If you want to update existing instance you need to change parameters directly.
Updating the launch template specification will force a new instance.

Any other instance parameters that you specify will override the same parameters in the launch template.

The `launch_template` block has the following structure:

* `id` - The ID of the launch template.
* `name` - The name of the launch template.
* `version` - Template version.
    * _Valid values:_ A version number, `$Latest`, `$Default`
    * _Default value:_ `$Default`

### root_block_device

The `root_block_device` block has the following structure:

* `delete_on_termination` - (Optional) Whether the volume should be destroyed on instance termination.
    * _Default value:_ `true`
* `encrypted` - (Optional) Indicates whether the EBS volume is encrypted.
* `iops` - (Optional) Amount of provisioned IOPS.
    * _Constraints:_ Only valid for volume_type of `io2`
* `tags` - (Optional) Map of tags to assign to the device.
* `volume_size` - (Optional) Size of the volume in GiB.
* `volume_type` - (Optional) Type of volume.

## Attribute Reference

### Supported attributes

In addition to all arguments above, the following attributes are exported:

* `arn` - The Amazon Resource Name (ARN) of the instance.
* `instance_state` - The state of the instance.
    * _Valid values:_ `pending`, `running`, `shutting-down`, `stopped`, `stopping`, `terminated`
* `primary_network_interface_id` - The ID of the instance's primary network interface.
* `private_dns` - The private DNS name assigned to the instance. For EC2-VPC, this is only available if you've enabled DNS hostnames for your VPC.
* `public_dns` - The public DNS name assigned to the instance. For EC2-VPC, this is only available if you've enabled DNS hostnames for your VPC.
* `public_ip` - The public IP address assigned to the instance, if applicable.

    ~> **Note** If you are using [`aws_eip`](eip.md) with your instance, you should refer to the EIP's address directly and not use `public_ip` as this field will change after the EIP is attached.

* `security_groups` - List of security group names associated with the instance.
* `tags_all` - Map of tags assigned to the instance, including those inherited from the provider [`default_tags` configuration block][default-tags].

For `ebs_block_device`, in addition to the arguments above, the following attribute is exported:

* `volume_id` - The ID of the volume. For example, the ID can be accessed like this, `aws_instance.web.ebs_block_device.2.volume_id`.

For `root_block_device`, in addition to the arguments above, the following attributes are exported:

* `volume_id` - The ID of the volume. For example, the ID can be accessed like this, `aws_instance.web.root_block_device.0.volume_id`.
* `device_name` - Device name, for example, `disk1`.

### Unsupported attributes

~> **Note** These attributes may be present in the `terraform.tfstate` file, but they have preset values and cannot be specified in configuration files.

The following attributes are not currently supported:

`capacity_reservation_specification`, `cpu_core_count`, `cpu_threads_per_core`, `credit_specification`, `ebs_block_device.kms_key_id`, `ebs_block_device.throughput`, `ebs_optimized`, `enclave_options`, `get_password_data`, `hibernation`, `iam_instance_profile`, `ipv6_address_count`, `ipv6_addresses`, `maintenance_options`, `metadata_options`, `network_interface.network_card_index`, `outpost_arn`, `password_data`, `placement_partition_number`, `root_block_device.kms_key_id`, `root_block_device.throughput`.

## Timeouts

The `timeouts` block allows you to specify [timeouts] for certain actions:

* `create` - (Default `10 minutes`) Used when launching the instance (until it reaches the initial `running` state).
* `update` - (Default `10 minutes`) Used when stopping and starting the instance when necessary during update, for example, when changing instance type.
* `delete` - (Default `20 minutes`) Used when terminating the instance.

## Import

Instances can be imported using `id`, e.g.,

```
$ terraform import aws_instance.web i-12345678
```
