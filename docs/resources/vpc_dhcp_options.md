---
subcategory: "VPC (Virtual Private Cloud)"
layout: "aws"
page_title: "aws_vpc_dhcp_options"
description: |-
  Manages a DHCP options set for a VPC.
---

[default-tags]: https://www.terraform.io/docs/providers/aws/index.html#default_tags-configuration-block
[dhcp-options]: https://docs.k2.cloud/en/services/networking/dhcpattrs.html
[rfc-2132]: http://www.ietf.org/rfc/rfc2132.txt

# Resource: aws_vpc_dhcp_options

Manages a DHCP options set for a VPC.

For more information, see the documentation on [DHCP options][dhcp-options].

## Example usage

### Basic example

```terraform
resource "aws_vpc_dhcp_options" "dns_resolver" {
  domain_name_servers = ["8.8.8.8", "8.8.4.4"]
}
```

### Specific example

```terraform
resource "aws_vpc_dhcp_options" "foo" {
  domain_name          = "service.consul"
  domain_name_servers  = ["127.0.0.1", "10.0.0.2"]
  ntp_servers          = ["127.0.0.1"]
  netbios_name_servers = ["127.0.0.1"]
  netbios_node_type    = 2

  tags = {
    Name = "foo-name"
  }
}
```

## Argument reference

The following arguments are supported:

* `domain_name` - (Optional, Forces new resource, String) The suffix domain name to use by default when resolving non-fully qualified domain names.
    In other words, this is what ends up being the `search` value in the `/etc/resolv.conf` file.
* `domain_name_servers` - (Optional, Forces new resource, List of strings) The list of IP addresses of domain name servers or `AmazonProvidedDNS`.
    We recommend using only one of the two parameters.
    * _List size:_ From 0 to 4 elements
* `netbios_name_servers` - (Optional, Forces new resource, List of strings) The list of NetBIOS name servers.
    * _List size:_ From 0 to 4 elements
* `netbios_node_type` - (Optional, Forces new resource, String) The NetBIOS node type. For more information about these node types, see [RFC 2132][rfc-2132].
    * _Valid values:_ `1`, `2`, `4`, `8`
* `ntp_servers` - (Optional, Forces new resource, List of strings) The list of NTP servers to configure.
    * _List size:_ From 0 to 4 elements
* `tags` - (Optional, Editable, Map of strings) Key-value pairs to assign to the resource.
    If the [`default_tags` configuration block][default-tags] block is used within a provider configuration, the tags with matching keys will overwrite those defined at the provider level.

## Notes

* All arguments are optional, but you have to specify at least one argument.
* To actually use the DHCP options set, you need to associate it to a VPC using [`aws_vpc_dhcp_options_association`](vpc_dhcp_options_association.md).
* If you delete a DHCP options set, all VPCs using it will be associated to the `default` DHCP options set.
* In most cases unless you're configuring your own DNS you'll want to set `domain_name_servers` to `AmazonProvidedDNS`.

## Attribute reference

### Supported attributes

In addition to all arguments above, the following attributes are exported:

* `arn` - (String) The Amazon Resource Name (ARN) of the DHCP options set.
* `id` - (String) The ID of the DHCP options set.
* `tags_all` - (Map of strings) Key-value pairs assigned to the resource, including any tags inherited from the [`default_tags` configuration block][default-tags] if used within a provider configuration.

### Unsupported attributes

~> **Note** This attribute may be present in the `terraform.tfstate` file, but it has a preset value and cannot be specified in configuration files.

The following attribute is not currently supported: `owner_id`.

## Timeouts

Timeouts usage for the DHCP options is not currently supported.

## Import

VPC DHCP options can be imported using the `id`, for example:

```
$ terraform import aws_vpc_dhcp_options.my_options dopt-12345678
```

