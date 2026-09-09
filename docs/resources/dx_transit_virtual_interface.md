---
subcategory: "Direct Connect"
layout: "aws"
page_title: "aws_dx_transit_virtual_interface"
description: |-
  Manages a Direct Connect transit virtual interface.
---

[default-tags]: https://registry.terraform.io/providers/hashicorp/aws/latest/docs#default_tags-configuration-block
[timeouts]: https://developer.hashicorp.com/terraform/plugin/framework/resources/timeouts

# Resource: aws_dx_transit_virtual_interface

Manages a Direct Connect transit virtual interface.
A transit virtual interface is a VLAN that transports traffic from a [Direct Connect gateway](dx_gateway.md) to one or more [transit gateways](ec2_transit_gateway.md).

## Example Usage

```terraform
data "aws_dx_connection" "selected" {
  name = "tf-dxconn-example"
}

resource "aws_dx_gateway" "example" {
  name            = "tf-dxgw-example"
  amazon_side_asn = "64512"
}

resource "aws_dx_transit_virtual_interface" "example" {
  name          = "tf-dxvif-example"
  connection_id = data.aws_dx_connection.selected.id
  dx_gateway_id = aws_dx_gateway.example.id
  vlan          = "4094"
  bgp_asn       = "65352"
}
```

## Argument Reference

The following arguments are supported:

* `address_family` - (Optional) The address family for the BGP peer.
    * _Valid values:_ `ipv4 `
    * _Default value:_ `ipv4`
* `amazon_address` - (Optional) The IPv4 CIDR address of the connection endpoint on the cloud side.
* `bgp_asn` - (Required) The BGP ASN on the client side.
* `bgp_auth_key` - (Optional, Sensitive) The authentication key for BGP configuration.
* `connection_id` - (Required) The ID of the Direct Connect connection (or LAG) on which the virtual interface has to be created.
* `customer_address` - (Optional) The IPv4 CIDR address of the connection endpoint on the client side.
* `dx_gateway_id` - (Required) The ID of the Direct Connect gateway that the virtual interface must be connected to.
* `name` - (Required) The name for the virtual interface.
* `tags` - (Optional, Editable) Map of tags to assign to the virtual interface. If a provider [`default_tags` configuration block][default-tags] is used, tags with matching keys will overwrite those defined at the provider level.
* `vlan` - (Required) The VLAN ID.

## Attribute Reference

### Supported attributes

In addition to all arguments above, the following attributes are exported:

* `amazon_side_asn` - The ASN for the cloud side of the connection.
* `arn` - The Amazon Resource Name (ARN) of the virtual interface.
* `aws_device` - The ID of the device to which the connection (or LAG) is attached.
* `id` - The ID of the virtual interface.
* `tags_all` - Map of tags assigned to the virtual interface, including those inherited from the provider [`default_tags` configuration block][default-tags].

### Unsupported attributes

~> **Note** These attributes may be present in the `terraform.tfstate` file, but they have preset values and cannot be specified in configuration files.

The following attributes are not currently supported:

`jumbo_frame_capable`, `mtu`, `sitelink_enabled`.

## Timeouts

The `timeouts` block allows you to specify [timeouts] for certain actions:

- `create` - (Default `10 minutes`) Timeout for creating virtual interface.
- `update` - (Default `10 minutes`) Timeout for virtual interface modifications.
- `delete` - (Default `10 minutes`) Timeout for destroying virtual interface.

## Import

Direct Connect transit virtual interfaces can be imported using `id`, e.g.,

```
$ terraform import aws_dx_transit_virtual_interface.example dxvif-12345678
```
