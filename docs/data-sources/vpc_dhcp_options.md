---
subcategory: "VPC (Virtual Private Cloud)"
layout: "aws"
page_title: "aws_vpc_dhcp_options"
description: |-
  Provides information about the DHCP options configuration.
---

[describe-dhcp-options]: https://docs.k2.cloud/en/api/ec2/actions/dhcp_options/DescribeDhcpOptions.html
[rfc-2132]: http://www.ietf.org/rfc/rfc2132.txt

# Data Source: aws_vpc_dhcp_options

Provides information about the DHCP options configuration.

## Example usage

### Specific example: performing a lookup by DHCP options ID

```terraform
variable dopts_id {}

data "aws_vpc_dhcp_options" "example" {
  dhcp_options_id = var.dopts_id
}
```

### Specific example: performing a lookup by filter

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

## Argument reference

The following arguments are supported:

* `dhcp_options_id` - (Optional, String) The ID of the DHCP options set.
* `filter` - (Optional, [Block](#filter)) One or more name/value pairs to use as filters.
    * _Valid values:_ See supported names and values in [EC2 API documentation][describe-dhcp-options]

### filter

* `name` - (Required, String) The name of the filter.
    * _Constraints:_ Filter names are case-sensitive
* `values` - (Required, List of strings) One or more filter values.
    * _Constraints:_ Filter values are case-sensitive

## Attribute reference

### Supported attributes

In addition to all arguments above, the following attributes are exported:

* `arn` - (String) The Amazon Resource Name (ARN) of the DHCP options set.
* `domain_name` - (String) The suffix domain name to be used when resolving non-fully qualified domain names (for example, the `search` value in the `/etc/resolv.conf` file).
* `domain_name_servers` - (List of strings) The list of name servers set.
* `id` - (String) The ID of the DHCP options set.
* `netbios_name_servers` - (List of strings) The List of NetBIOS name servers.
* `netbios_node_type` - (String) The NetBIOS node type (1, 2, 4, or 8).
    For more information about these node types, see [RFC 2132][rfc-2132].
* `ntp_servers` - (List of strings) The list of NTP servers.
* `tags` - (Map of strings) Key-value pairs assigned to the resource.

### Unsupported attributes

~> **Note** This attribute may be present in the `terraform.tfstate` file, but it has a preset value and cannot be specified in configuration files.

The following attribute is not currently supported: `owner_id`.
