---
subcategory: "VPC (Virtual Private Cloud)"
layout: "aws"
page_title: "aws_network_interface"
description: |-
  Manages an elastic network interface (ENI).
---

[default-tags]: https://www.terraform.io/docs/providers/aws/index.html#default_tags-configuration-block
[network-interfaces]: https://docs.k2.cloud/en/services/networking/interfaces/operations.html

# Resource: aws_network_interface

Manages an elastic network interface (ENI) resource.

For more information about network interfaces, see the documentation on [Network interfaces][network-interfaces].

## Example usage

### Basic example

```terraform
resource "aws_network_interface" "example" {
  subnet_id   = "subnet-12345678"
  private_ips = ["10.0.31.50"]

  attachment {
    instance     = "i-12345678"
    device_index = 1
  }
}
```

## Argument reference

The following arguments are required:

* `subnet_id` - (Required, Forces new resource, String) The ID of the subnet where the ENI has to be created.

The following arguments are optional:

* `attachment` - (Optional, Editable, [Block](#attachment)) One or more ENI attachments.
* `description` - (Optional, Editable, String) The description for the network interface.
* `private_ip_list` - (Optional, Editable, List of strings) The list of private IPs to assign to the ENI **in sequential order**.
    * _Constraints:_ Only a single value can be specified, otherwise, an error will be returned
* `private_ips` - (Optional, Editable, List of strings) The list of private IPs to assign to the ENI **without regard to order**.
    * _Constraints:_ Only a single value can be specified, otherwise, an error will be returned
* `security_groups` - (Optional, Editable, List of strings) The list of security group IDs to assign to the ENI.
* `source_dest_check` - (Optional, Editable, Boolean) Indicates whether the network interface must perform source/destination checking.
    * _Default value:_ `true`
* `tags` - (Optional, Editable, Map of strings) Key-value pairs to assign to the resource.
    If the [`default_tags` configuration block][default-tags] block is used within a provider configuration, the tags with matching keys will overwrite those defined at the provider level.

~> **Note** The `private_ip_list` and `private_ips` cannot be used together in one resource. Choose the preferred one before using.

### attachment

The `attachment` block has the following structure:

* `instance` - (Required, Editable, String) The ID of the instance to attach to.
* `device_index` - (Required, Editable, Integer) The index of the network interface.

## Attribute reference

### Supported attributes

In addition to all arguments above, the following attributes are exported:

* `arn` - (String) The Amazon Resource Name (ARN) of the network interface.
* `id` - (String) The ID of the network interface.
* `mac_address` - (String) The MAC address of the network interface.
* `owner_id` - (String) The ID of the project that owns the network interface.
* `private_dns_name` - (String) The private DNS name of the network interface (IPv4).
* `tags_all` - (Map of strings) Key-value pairs assigned to the resource, including any tags inherited from the [`default_tags` configuration block][default-tags] if used within a provider configuration.

### Unsupported attributes

~> **Note** These attributes may be present in the `terraform.tfstate` file, but they have preset values and cannot be specified in configuration files.

The following attributes are not currently supported:

`interface_type`, `ipv4_prefix_count`, `ipv4_prefixes`, `ipv6_address_count`, `ipv6_address_list_enable`, `ipv6_address_list`, `ipv6_addresses`, `ipv6_prefix_count`, `ipv6_prefixes`, `outpost_arn`, `private_ip_list_enable`, `private_ips_count`.

## Timeouts

Timeouts usage for network interfaces is not currently supported.

## Import

Network interfaces can be imported using `id`, for example:

```
$ terraform import aws_network_interface.test eni-12345678
```

