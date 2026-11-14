---
subcategory: "VPN (Site-to-Site)"
layout: "aws"
page_title: "aws_customer_gateway"
description: |-
  Manages a customer gateway.
---

[default-tags]: https://registry.terraform.io/providers/hashicorp/aws/latest/docs#default_tags-configuration-block

# Resource: aws_customer_gateway

Manages a customer gateway.
A customer gateway is the point of connection to the remote infrastructure.
It stores the data about the remote device, and the VPN connections are established between it and the VPC.
For more information, see the [`aws_vpn_connection`](vpn_connection.md) resource.

~> **Note** A customer gateway cannot be deleted while it has a VPN connection.

-> **Note** The BGP ASN of the customer gateway must differ from the BGP ASN of the VPC to which it is connected.

## Example usage

### Basic example

```terraform
resource "aws_customer_gateway" "example" {
  bgp_asn    = 65000
  ip_address = "172.83.124.10"
  type       = "ipsec.1"

  tags = {
    Name = "tf-customer-gateway"
  }
}
```

## Argument reference

The following arguments are required:

* `ip_address` - (Required, Forces new resource, String) The IPv4 address of the gateway's Internet-routable external interface.
* `type` - (Required, Forces new resource, String) The type of the customer gateway.
    * _Valid values:_ `ipsec.1` or `ipsec.legacy`

The following arguments are optional:

* `bgp_asn` - (Optional, Forces new resource, String) The Border Gateway Protocol (BGP) Autonomous System Number (ASN) of the gateway.
    * _Default value:_ `65000`
    * _Valid values:_ From `1` to `2147483647` or from `4200000000` to `4294967294`
* `tags` - (Optional, Editable, Map of strings) Key-value pairs to assign to the customer gateway.
  If the [`default_tags` configuration block][default-tags] is used within a provider configuration, the tags with matching keys will overwrite those defined at the provider level.

## Attribute reference

In addition to all arguments above, the following attributes are exported:

* `arn` - (String) The ARN of the customer gateway.
* `id` - (String) The ID of the customer gateway.
* `tags_all` - (Map of strings) Key-value pairs assigned to the customer gateway, including any tags inherited from the [`default_tags` configuration block][default-tags] if used within a provider configuration.

## Timeouts

Timeouts usage for customer gateways is not currently supported.

## Import

In Terraform v1.5.0 or later, the customer gateway can be imported using `id` in the `import` block, for example:

```terraform
import {
  to = aws_customer_gateway.example
  id = "cgw-12345678"
}
```

In older Terraform versions, the customer gateway can be imported using `id`, for example:

```
$ terraform import aws_customer_gateway.example cgw-12345678
```
