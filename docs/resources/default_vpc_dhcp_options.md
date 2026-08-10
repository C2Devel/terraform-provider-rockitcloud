---
subcategory: "VPC (Virtual Private Cloud)"
layout: "aws"
page_title: "aws_default_vpc_dhcp_options"
description: |-
  Manages the default DHCP options set.
---

[rfc-2132]: http://www.ietf.org/rfc/rfc2132.txt
[default-tags]: https://www.terraform.io/docs/providers/aws/index.html#default_tags-configuration-block

# Resource: aws_default_vpc_dhcp_options

Manages the default DHCP options set.

Each cloud region comes with a default set of DHCP options.

~> **Note** **This is an advanced resource**, and has special caveats to be aware of when using it.
Please read this document in its entirety before using this resource.
The `aws_default_vpc_dhcp_options` behaves differently from normal resources, in that Terraform does not _create_ this resource, but instead "adopts" it into management.

## Example usage

### Basic example

```terraform
resource "aws_default_vpc_dhcp_options" "default" {
  tags = {
    Name = "Default DHCP Option Set"
  }
}
```

## Argument reference

The arguments of an `aws_default_vpc_dhcp_options` differ slightly from [`aws_vpc_dhcp_options`](vpc_dhcp_options.md) resources.
Namely, the `domain_name`, `domain_name_servers` and `ntp_servers` arguments are computed.

The following optional arguments are supported:

* `netbios_name_servers` - (Optional, Forces new resource, List of strings) The list of NetBIOS name servers.
* `netbios_node_type` - (Optional, Forces new resource, List of strings) The NetBIOS node type. For more information about these node types, see [RFC 2132][rfc-2132].
    * _Valid values:_ `1`, `2`, `4`, `8`
* `tags` - (Optional, Editable, Map of strings) Key-value pairs to assign to the resource.
    If the [`default_tags` configuration block][default-tags] block is used within a provider configuration, the tags with matching keys will overwrite those defined at the provider level.

## Removing `aws_default_vpc_dhcp_options` from your configuration

The `aws_default_vpc_dhcp_options` resource allows you to manage a region's default DHCP options set, but Terraform cannot destroy it.
Removing this resource from your configuration will remove it from your statefile and management, but will not destroy the DHCP options set.
You can resume managing the DHCP options set via the cloud console.

## Attribute reference

### Supported attributes

In addition to all arguments above, the following attributes are exported:

* `arn` - (String) The Amazon Resource Name (ARN) of the DHCP options set.
* `domain_name` - (String) The IP address of the domain name server.
* `domain_name_servers` - (String) The list of domain name servers.
* `id` - (String) The ID of the DHCP options set.
* `ntp_servers` - (String) The list of NTP servers.
* `tags_all` - (Map of strings) Key-value pairs assigned to the resource, including any tags inherited from the [`default_tags` configuration block][default-tags] if used within a provider configuration.

### Unsupported attributes

~> **Note** These attributes may be present in the `terraform.tfstate` file, but they have preset values and cannot be specified in configuration files.

The following attribute is not currently supported: `owner_id`.

## Timeouts

Timeouts usage for the VPC DHCP options is not currently supported.

## Import

Default VPC DHCP options can be imported using `id`, for example:

```
$ terraform import aws_default_vpc_dhcp_options.default_options dopt-12345678
```
