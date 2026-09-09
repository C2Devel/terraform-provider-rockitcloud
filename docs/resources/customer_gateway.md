---
subcategory: "VPN (Site-to-Site)"
layout: "aws"
page_title: "aws_customer_gateway"
description: |-
  Manages a customer gateway inside a VPC.
---

# Resource: aws_customer_gateway

Manages a customer gateway inside a VPC. These objects can be connected to VPN gateways via VPN connections, and allow you to establish tunnels between your network and the VPC.

## Example Usage

```terraform
resource "aws_customer_gateway" "main" {
  bgp_asn    = 65000
  ip_address = "172.83.124.10"
  type       = "ipsec.1"

  tags = {
    Name = "main-customer-gateway"
  }
}
```

## Argument Reference

The following arguments are supported:

* `bgp_asn` - (Required) The gateway's Border Gateway Protocol (BGP) Autonomous System Number (ASN).
* `ip_address` - (Required) The IP address of the gateway's Internet-routable external interface.
* `type` - (Required) The type of customer gateway.
    * _Valid values:_ `ipsec.1`, `ipsec.legacy`
* `tags` - (Optional) Map of tags to assign to the gateway. If a provider [`default_tags` configuration block][default-tags] is used, tags with matching keys will overwrite those defined at the provider level.

## Attribute Reference

### Supported attributes

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the gateway.
* `arn` - The Amazon Resource Name (ARN) of the customer gateway.
* `tags_all` - Map of tags assigned to the customer gateway, including those inherited from the provider [`default_tags` configuration block][default-tags].

### Unsupported attributes

~> **Note** These attributes may be present in the `terraform.tfstate` file, but they have preset values and cannot be specified in configuration files.

The following attributes are not currently supported:

`certificate_arn`, `device_name`.

## Import

Customer gateways can be imported using `id`, e.g.,

```
$ terraform import aws_customer_gateway.main cgw-12345678
```

[default-tags]: https://registry.terraform.io/providers/hashicorp/aws/latest/docs#default_tags-configuration-block
