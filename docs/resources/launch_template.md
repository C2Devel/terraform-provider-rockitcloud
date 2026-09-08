---
subcategory: "EC2 (Elastic Compute Cloud)"
layout: "aws"
page_title: "aws_launch_template"
description: |-
  Manages an EC2 launch template.
---

[asg-create]: https://docs.k2.cloud/en/services/compute/autoscaling.html#createautoscalinggroup
[default-tags]: https://registry.terraform.io/providers/hashicorp/aws/latest/docs#default_tags-configuration-block
[describe-images]: https://docs.k2.cloud/en/api/ec2/actions/images/DescribeImages.html

# Resource: aws_launch_template

Manages an EC2 launch template. The resource can be used to create instances or Auto Scaling groups.

## Example usage

```terraform
resource "aws_launch_template" "example" {
  name = "tf-lt"

  block_device_mappings {
    device_name = "disk1"

    ebs {
      volume_size = 20
    }
  }

  disable_api_termination = true

  instance_initiated_shutdown_behavior = "terminate"

  image_id      = "cmi-12345678" # add image id, change instance type if needed
  instance_type = "m1.micro"

  monitoring {
    enabled = true
  }

  placement {
    availability_zone = "ru-msk-vol52"
  }

  tag_specifications {
    resource_type = "instance"

    tags = {
      Name = "tf-lt"
    }
  }
}
```

## Argument reference

The following arguments are required:

* `image_id` - (Required, Editable, String) The ID of the image from which to launch the instance.

The following arguments are optional:

* `block_device_mappings` - (Optional, Editable, [Block](#block_device_mappings)) Specify volumes to attach to the instance besides the volumes specified by the image.
* `default_version` - (Optional, Editable, Integer) The default version of the launch template.
    * _Constraints:_ Conflicts with `update_default_version`
* `description` - (Optional, Editable, String) The description of the launch template version.
* `disable_api_termination` - (Optional, Editable, Boolean) Indicates whether to disable the possibility to terminate an instance via the API.
* `instance_initiated_shutdown_behavior` - (Optional, Editable, String) The shutdown behavior for the instance.
    * _Valid values:_ `stop`, `terminate`
* `instance_type` - (Optional, Editable, String) The type of the instance.
* `key_name` - (Optional, Editable, String) The key name to use for the instance.
* `monitoring` - (Optional, Editable, [Block](#monitoring)) The monitoring option for the instance.
* `name` - (Optional, Forces new resource, String) The name of the launch template.
  If you leave this blank, Terraform will auto-generate a unique name.
    * _Constraints:_ Conflicts with `name_prefix`
* `name_prefix` - (Optional, Forces new resource, String) Creates a unique name beginning with the specified prefix.
    * _Constraints:_ Conflicts with `name`
* `network_interfaces` - (Optional, Editable, [Block](#network_interfaces)) Customize network interfaces to be attached at instance boot time.
* `placement` - (Optional, Editable, [Block](#placement)) The placement of the instance.
* `tag_specifications` - (Optional, Editable, [Block](#tag_specifications)) The tags to apply to the resources during launch.
* `tags` - (Optional, Editable, Map of strings) Key-value pairs to assign to the launch template.
  If the [`default_tags` configuration block][default-tags] is used within a provider configuration, the tags with matching keys will overwrite those defined at the provider level.
* `update_default_version` - (Optional, Editable, Boolean) Indicates whether to update the default version on each update.
    * _Constraints:_ Conflicts with `default_version`
* `user_data` - (Optional, Editable, String) The base64-encoded user data to provide when launching the instance.
    * _Constraints:_ The text length must not exceed 16 KB
* `vpc_security_group_ids` - (Optional, Editable, Set of strings) List of security group IDs to associate with.

### block_device_mappings

Configures additional volumes of the instance besides specified by the image.

To find out more information for an existing image to override the configuration, such as `device_name`, use the [EC2 API][describe-images].

The `block_device_mappings` block has the following structure:

* `device_name` - (Optional, Editable, String) The name of the device to mount.
* `ebs` - (Optional, Editable, [Block](#ebs)) Configures EBS volume properties.
* `no_device` - (Optional, Editable, String) Suppresses the specified device included in the block device mapping.

#### ebs

The `ebs` block has the following structure:

* `delete_on_termination` - (Optional, Editable, String) Indicates whether the volume should be destroyed on instance termination.
* `iops` - (Optional, Editable, Integer) The amount of provisioned IOPS.
    * _Constraints:_ This must be set with the `volume_type` of `io2`
* `snapshot_id` - (Optional, Editable, String) The ID of the snapshot to mount.
* `volume_size` - (Optional, Editable, Integer) The size of the volume in GiB.
* `volume_type` - (Optional, Editable, String) The type of the volume.

### monitoring

The `monitoring` block has the following structure:

* `enabled` - (Optional, Editable, Boolean) Indicates whether the launched EC2 instance will have detailed monitoring enabled.

### network_interfaces

Attaches one or more network interfaces to the instance.

For the details about configuring network interfaces when creating an Auto Scaling group, see the [user documentation][asg-create].

The `network_interfaces` block has the following structure:

* `associate_public_ip_address` - (Optional, Editable, String) Indicates whether a public IP address should be associated with the network interface.
    * _Constraints:_ The address will be assigned to the `eth0` interface if there are free allocated external addresses.
      This operation is available only for instances running in a VPC and for new network interfaces.
* `delete_on_termination` - (Optional, Editable, String) Indicates whether the network interface should be destroyed on instance termination.
* `description` - (Optional, Editable, String) Description of the network interface.
* `device_index` - (Optional, Editable, Integer) The integer index of the network interface attachment.
* `network_interface_id` - (Optional, Editable, String) The ID of the network interface to attach.
* `private_ip_address` - (Optional, Editable, String) The primary private IPv4 address.
* `security_groups` - (Optional, Editable, Set of strings) List of security group IDs to associate.
* `subnet_id` - (Optional, Editable, String) The ID of the subnet to associate.

### placement

The placement group of the instance.

The `placement` block has the following structure:

* `affinity` - (Optional, Editable, String) The affinity setting for an instance on a dedicated host.
    * _Default value:_ `default`
    * _Constraints:_ The parameter could be set to `host` only if `tenancy` is `host`
* `availability_zone` - (Optional, Editable, String) The availability zone for the instance.
* `group_name` - (Optional, Editable, String) The name of the placement group for the instance.
* `host_id` - (Optional, Editable, String) The ID of the dedicated host for the instance.
* `tenancy` - (Optional, Editable, String) The tenancy of the instance (if the instance is running in a VPC).
    * _Valid values:_ `default`, `host`
    * _Default value:_ `default`

~> **Note** If you use the `host` value, you may encounter the `NotEnoughResourcesForInstanceType` error when running an instance.
To avoid this, it is recommended to provide either the `subnet_id` argument or the `availability_zone` argument.

### tag_specifications

The tags to apply to the resources during launch. You can tag instances and volumes.

Each `tag_specifications` block has the following structure:

* `resource_type` - (Optional, Editable, String) The type of resource to tag.
    * _Valid values:_ `instance`, `volume`
* `tags` - (Optional, Editable, Map of strings) Map of tags to assign to the resource.

## Attribute reference

### Supported attributes

In addition to all arguments above, the following attributes are exported:

* `arn` - (String) The Amazon Resource Name (ARN) of the launch template.
* `id` - (String) The ID of the launch template.
* `latest_version` - (Integer) The latest version of the launch template.
* `tags_all` - (Map of strings) Key-value pairs assigned to the launch template, including any tags inherited from the [`default_tags` configuration block][default-tags] if used within a provider configuration.

### Unsupported attributes

~> **Note** These attributes may be present in the `terraform.tfstate` file, but they have preset values and cannot be specified in configuration files.

The following attributes are not currently supported:

`block_device_mappings.ebs.encrypted`, `block_device_mappings.ebs.kms_key_id`, `block_device_mappings.ebs.throughput`, `block_device_mappings.virtual_name`, `capacity_reservation_specification`, `cpu_options`, `credit_specification`, `ebs_optimized`, `elastic_gpu_specifications`, `elastic_inference_accelerator`, `enclave_options`, `hibernation_options`, `iam_instance_profile`, `instance_market_options`, `instance_requirements`, `kernel_id`, `license_specification`, `maintenance_options`, `metadata_options`, `network_interfaces.associate_carrier_ip_address`, `network_interfaces.interface_type`, `network_interfaces.ipv4_address_count`, `network_interfaces.ipv4_addresses`, `network_interfaces.ipv4_prefix_count`, `network_interfaces.ipv4_prefixes`, `network_interfaces.ipv6_address_count`, `network_interfaces.ipv6_addresses`, `network_interfaces.ipv6_prefix_count`, `network_interfaces.ipv6_prefixes`, `network_interfaces.network_card_index`, `placement.host_resource_group_arn`, `placement.spread_domain`, `placement.partition_number`, `private_dns_name_options`, `ram_disk_id`, `security_group_names`.

## Timeouts

Timeouts usage for the launch template is not currently supported.

## Import

Launch templates can be imported using `id`, for example:

```
$ terraform import aws_launch_template.web lt-12345678
```