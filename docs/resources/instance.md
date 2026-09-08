---
subcategory: "EC2 (Elastic Compute Cloud)"
layout: "aws"
page_title: "aws_instance"
description: |-
  Manages an EC2 instance.
---

[default-tags]: https://registry.terraform.io/providers/hashicorp/aws/latest/docs#default_tags-configuration-block
[provisioning]: https://developer.hashicorp.com/terraform/language/provisioners
[timeouts]: https://developer.hashicorp.com/terraform/plugin/framework/resources/timeouts

# Resource: aws_instance

Manages an EC2 instance.
This allows instances to be created, updated, and deleted.
Instances also support [provisioning].

## Example usage

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

### Network example

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

## Argument reference

* `affinity` - (Optional, Forces new resource, String) The affinity setting for an instance on a dedicated host.
    * _Valid values:_ `default`, `host`
    * _Constraints:_ The parameter could be set to `host` only if `tenancy` is `host`
* `ami` - (Optional, Forces new resource, String) An image to use for the instance.
  If an image is specified in the launch template, the `ami` setting will override it.
    * _Constraints:_ Required unless `launch_template` is specified
* `associate_public_ip_address` - (Optional, Forces new resource, Boolean) Indicates whether to associate a public IP address with an instance in a VPC.
    * _Constraints:_ Conflicts with the `network_interface` argument

    ~> **Note** The address will be assigned to the `eth0` interface only if there are free allocated external addresses.
    This operation is available only for instances running in the VPC and for new network interfaces.

* `availability_zone` - (Optional, Forces new resource, String) An availability zone to start the instance in.
* `disable_api_termination` - (Optional, Editable, Boolean) Indicates whether the possibility to terminate an instance via the API is disabled.
* `ebs_block_device` - (Optional, Editable, [Block](#ebs_block_device)) One or more configuration blocks with additional EBS block devices to attach to the instance.
  When accessing this as an attribute reference, it is a set of objects.
    * _Constraints:_ Block device configurations are applied only when the resource is created
* `get_password_data` - (Optional, Editable, Boolean) Indicates whether to retrieve the password data of a Windows instance.
    * _Default value:_ `false`
* `hibernation` - (Optional, Forces new resource, Boolean) Indicates whether the instance is optimized for hibernation.
* `host_id` - (Optional, Forces new resource, String) The ID of the dedicated host that the instance will be assigned to.
* `instance_initiated_shutdown_behavior` - (Optional, Editable, String) The shutdown behavior for the instance.
    * _Valid values:_ `stop`, `terminate`
* `instance_type` - (Optional, Editable, String) The instance type to use for the instance.
  Updates to this field will trigger a stop/start of the EC2 instance.
* `key_name` - (Optional, Forces new resource, String) The key name of the key pair to use for the instance; which can be managed using [the `aws_key_pair` resource](key_pair.md).
* `launch_template` - (Optional, Forces new resource, [Block](#launch_template)) Specifies a launch template to configure the instance.
  Parameters configured on this resource will override the corresponding parameters in the launch template.
* `metadata_options` - (Optional, Editable, [Block](#metadata_options)) Customize the metadata options of the instance.
* `monitoring` - (Optional, Editable, Boolean) Indicates whether detailed monitoring is enabled for the launched EC2 instance.
* `network_interface` - (Optional, Editable, [Block](#network_interface)) Customize network interfaces to be attached at instance boot time.
    * _Constraints:_ Conflicts with the `associate_public_ip_address`, `private_ip`, `secondary_private_ips`, `subnet_id`, `vpc_security_group_ids` arguments
* `placement_group` - (Optional, Forces new resource, String) The placement group to start the instance in.
* `private_ip` - (Optional, Forces new resource, String) A private IP address to associate with the instance in a VPC.
    * _Constraints:_ Conflicts with the `network_interface` argument
* `root_block_device` - (Optional, Editable, [Block](#root_block_device)) The root block device of the instance.
  When accessing this as an attribute reference, it is a list containing one object.
* `secondary_private_ips` - (Optional, Editable, Set of strings) A list of secondary private IPv4 addresses to assign to the instance's primary network interface in a VPC.
    * _Constraints:_
        * Conflicts with the `network_interface` argument
        * Only the primary private IP address can be specified
* `source_dest_check` - (Optional, Editable, Boolean) Indicates whether the traffic is routed to the instance when the destination address does not match the instance.
    * _Default value:_ `true`
* `subnet_id` - (Optional, Forces new resource, String) The ID of a subnet to launch in.
* `tags` - (Optional, Editable, Map of strings) Key-value pairs to assign to the instance.
  If the [`default_tags` configuration block][default-tags] is used within a provider configuration, the tags with matching keys will overwrite those defined at the provider level.
    * _Constraints:_ These tags apply to the instance and not block storage devices
* `tenancy` - (Optional, Forces new resource, String) The placement type.
    * _Valid values:_ `default`, `host`

    ~> **Note** If you use the `host` value, you may encounter the `NotEnoughResourcesForInstanceType` error when running an instance.
    To avoid this, it is recommended to provide either the `subnet_id` argument or the `availability_zone` argument.

* `user_data` - (Optional, Editable, String) User data to provide when launching the instance.
  Do not pass gzip-compressed data via this argument; see `user_data_base64` instead.
  Updates to this field will trigger a stop/start of the EC2 instance by default.
  If the `user_data_replace_on_change` is set then updates to this field will trigger a destroy and recreate.
    * _Constraints:_ Conflicts with the `user_data_base64` argument
* `user_data_base64` - (Optional, Editable, String) Can be used instead of `user_data` to pass base64-encoded binary data directly.
  Use this instead of `user_data` whenever the value is not a valid UTF-8 string.
  For example, gzip-encoded user data must be base64-encoded and passed via this argument to avoid corruption.
  Updates to this field will trigger a stop/start of the EC2 instance by default.
  If the `user_data_replace_on_change` is set then updates to this field will trigger a destroy and recreate.
    * _Constraints:_ Conflicts with the `user_data` argument
* `user_data_replace_on_change` - (Optional, Editable, Boolean) Indicates whether the resource will be destroyed and recreated when the `user_data` or `user_data_base64` argument changes.
    * _Default value:_ `false`
* `volume_tags` - (Optional, Editable, Map of strings) A map of tags to assign to root and EBS volumes when the instance is created.

    ~> **Note** Do not use `volume_tags` if you plan to manage block device tags outside the `aws_instance` configuration, such as using `tags` in an [`aws_ebs_volume`](ebs_volume.md) resource attached via [`aws_volume_attachment`](volume_attachment.md).
    Doing so will result in resource cycling and inconsistent behavior.

* `vpc_security_group_ids` - (Optional, Editable, Set of strings) A list of security group IDs to associate with.
    * _Constraints:_ Conflicts with the `network_interface` argument

### ebs_block_device

The following arguments are required:

* `device_name` - (Required, Forces new resource, String) The name of the device to mount.

The following arguments are optional:

* `delete_on_termination` - (Optional, Forces new resource, Boolean) Indicates whether the volume should be destroyed on instance termination.
    * _Default value:_ `true`
* `iops` - (Optional, Forces new resource, Integer) The amount of provisioned IOPS.
    * _Constraints:_ Only valid for the volume type `io2`
* `snapshot_id` - (Optional, Forces new resource, String) The ID of the snapshot to mount.
* `tags` - (Optional, Editable, Map of strings) A map of tags to assign to the device.
* `volume_size` - (Optional, Forces new resource, Integer) The size of the volume in GiB.
* `volume_type` - (Optional, Forces new resource, String) The type of the volume.

~> **Note** Currently, changes to the `ebs_block_device` configuration of _existing_ resources cannot be automatically detected by Terraform.
To manage changes and attachments of an EBS block to an instance, use the [`aws_ebs_volume`](ebs_volume.md) and [`aws_volume_attachment`](volume_attachment.md) resources instead.
If you use `ebs_block_device` on an `aws_instance`, Terraform will assume management over the full set of non-root EBS block devices for the instance, treating additional block devices as drift.
For this reason, `ebs_block_device` cannot be mixed with external `aws_ebs_volume` and `aws_volume_attachment` resources for a given instance.

### network_interface

Each of the `network_interface` blocks attaches a network interface to an EC2 instance during boot time.
However, because the network interface is attached at boot time, replacing/modifying the network interface **will** trigger a recreation of the EC2 instance.
If you should need at any point to detach/modify/re-attach a network interface to the instance, use the [`aws_network_interface`](network_interface.md) or [`aws_network_interface_attachment`](network_interface_attachment.md) resources instead.

The `network_interface` configuration block _does_, however, allow users to supply their own network interface to be used as the default network interface on an EC2 instance, attached at `eth0`.

The following arguments are required:

* `device_index` - (Required, Forces new resource, Integer) The integer index of the network interface attachment.
* `network_interface_id` - (Required, Forces new resource, String) The ID of the network interface to attach.

The following arguments are optional:

* `delete_on_termination` - (Optional, Forces new resource, Boolean) Indicates whether to delete the network interface on instance termination.
    * _Default value:_ `false`
    * _Constraints:_ Currently, the only valid value is `false`, as this option is only supported when creating new network interfaces during instance launching
* `network_card_index` - (Optional, Forces new resource, Integer) The index of the network card.
    * _Default value:_ `0`

### launch_template

~> **Note** Launch template parameters will be used only once when the instance is created.
If you want to update existing instance you need to change parameters directly.
Updating the launch template specification will force a new instance.

Any other instance parameters that you specify will override the same parameters in the launch template.

The `launch_template` block has the following structure:

* `id` - (Optional, Forces new resource, String) The ID of the launch template.
* `name` - (Optional, Forces new resource, String) The name of the launch template.
* `version` - (Optional, Editable, String) The template version.
    * _Valid values:_ A version number, `$Default`, `$Latest`
    * _Default value:_ `$Default`

### metadata_options

The `metadata_options` block has the following structure:

* `http_endpoint` - (Optional, Editable, String) Indicates whether the metadata service is available.
* `http_put_response_hop_limit` - (Optional, Editable, Integer) The desired HTTP PUT response hop limit for instance metadata requests.
    * _Valid values:_ From 1 to 64
* `http_tokens` - (Optional, Editable, String) Indicates whether the metadata service requires a session token.
* `instance_metadata_tags` - (Optional, Editable, String) Indicates whether access to instance metadata tags is enabled.
    * _Default value:_ `disabled`

### root_block_device

The `root_block_device` block has the following structure:

* `delete_on_termination` - (Optional, Editable, Boolean) Indicates whether the volume should be destroyed on instance termination.
    * _Default value:_ `true`
* `iops` - (Optional, Editable, Integer) The amount of provisioned IOPS.
    * _Constraints:_ Only valid for `volume_type` of `io2`
* `tags` - (Optional, Editable, Map of strings) A map of tags to assign to the device.
* `volume_size` - (Optional, Editable, Integer) The size of the volume in GiB.
* `volume_type` - (Optional, Editable, String) The type of the volume.

## Attribute reference

### Supported attributes

In addition to all arguments above, the following attributes are exported:

* `arn` - (String) The Amazon Resource Name (ARN) of the instance.
* `instance_state` - (String) The state of the instance.
    * _Valid values:_ `pending`, `running`, `shutting-down`, `stopped`, `stopping`, `terminated`
* `primary_network_interface_id` - (String) The ID of the instance's primary network interface.
* `private_dns` - (String) The private DNS name assigned to the instance.
  For EC2-VPC, this is only available if you've enabled DNS hostnames for your VPC.
* `public_dns` - (String) The public DNS name assigned to the instance.
  For EC2-VPC, this is only available if you've enabled DNS hostnames for your VPC.
* `public_ip` - (String) The public IP address assigned to the instance, if applicable.

    ~> **Note** If you are using [`aws_eip`](eip.md) with your instance, you should refer to the EIP's address directly and not use `public_ip` as this field will change after the EIP is attached.

* `security_groups` - (Set of strings) A list of security group names associated with the instance.
* `tags_all` - (Map of strings) Key-value pairs assigned to the instance, including any tags inherited from the [`default_tags` configuration block][default-tags] if used within a provider configuration.

For `ebs_block_device`, in addition to the arguments above, the following attribute is exported:

* `volume_id` - (String) The ID of the volume.
    * _Example:_ `aws_instance.web.ebs_block_device.2.volume_id`

For `root_block_device`, in addition to the arguments above, the following attributes are exported:

* `volume_id` - (String) The ID of the volume.
    * _Example:_ `aws_instance.web.root_block_device.0.volume_id`
* `device_name` - (String) The device name.
    * _Example:_ `disk1`

### Unsupported attributes

~> **Note** These attributes may be present in the `terraform.tfstate` file, but they have preset values and cannot be specified in configuration files.

The following attributes are not currently supported:

`capacity_reservation_specification`, `cpu_core_count`, `cpu_threads_per_core`, `credit_specification`, `ebs_block_device.encrypted`, `ebs_block_device.kms_key_id`, `ebs_block_device.throughput`, `ebs_optimized`, `enclave_options`, `ephemeral_block_device`, `hibernation`, `iam_instance_profile`, `ipv6_address_count`, `ipv6_addresses`, `maintenance_options`, `outpost_arn`, `password_data`, `placement_partition_number`, `root_block_device.encrypted`, `root_block_device.kms_key_id`, `root_block_device.throughput`.

## Timeouts

The `timeouts` block allows you to specify [timeouts] for certain actions:

* `create` - (Default `10 minutes`) Used when launching the instance (until it reaches the initial `running` state).
* `update` - (Default `10 minutes`) Used when stopping and starting the instance when necessary during update, for example, when changing instance type.
* `delete` - (Default `20 minutes`) Used when terminating the instance.

## Import

Instance can be imported using `id`, for example:

```
$ terraform import aws_instance.web i-12345678
```
