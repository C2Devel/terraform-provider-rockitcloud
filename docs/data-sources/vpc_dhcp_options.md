---
subcategory: "VPC (Virtual Private Cloud)"
layout: "aws"
page_title: "aws_vpc_dhcp_options"
description: |-
  Provides information about an DHCP options configuration.
---

[describe-dhcp-options]: https://docs.k2.cloud/en/api/ec2/actions/dhcp_options/DescribeDhcpOptions.html
[rfc-2132]: http://www.ietf.org/rfc/rfc2132.txt

# Data Source: aws_vpc_dhcp_options

Provides information about an DHCP options configuration.

## Example Usage

### Lookup by DHCP options ID

```terraform
variable dopts_id {}

data "aws_vpc_dhcp_options" "example" {
  dhcp_options_id = var.dopts_id
}
```

### Lookup by Filter

```terraform
data "aws_vpc_dhcp_options" "example" {
  filter {
    name   = "key"
    values = ["domain-name"]
  }

  filter {
    name   = "value"
    values = ["example.com"]
  }
}
```

## Argument Reference

* `dhcp_options_id` - (Optional) The DHCP options ID.
* `filter` - (Optional) One or more name/value pairs to use as filters.
    * _Valid values:_ See supported names and values in [EC2 API documentation][describe-dhcp-options]

## Attribute Reference

### Supported attributes

In addition to all arguments above, the following attributes are exported:

* `arn` - The Amazon Resource Name (ARN) of the DHCP options set.
* `dhcp_options_id` - DHCP options ID.
* `domain_name` - The suffix domain name to used when resolving non fully qualified domain names e.g., the `search` value in the `/etc/resolv.conf` file.
* `domain_name_servers` - List of name servers.
* `id` - DHCP options ID.
* `netbios_name_servers` - List of NETBIOS name servers.
* `netbios_node_type` - The NetBIOS node type (1, 2, 4, or 8). For more information about these node types, see [RFC 2132][rfc-2132].
* `ntp_servers` - List of NTP servers.
* `tags` - Map of tags assigned to the DHCP options set.

### Unsupported attributes

~> **Note** This attribute may be present in the `terraform.tfstate` file, but it has a preset value and cannot be specified in configuration files.

The following attribute is not currently supported: `owner_id`.
