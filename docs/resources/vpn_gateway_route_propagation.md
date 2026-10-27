---
subcategory: "VPN (Site-to-Site)"
layout: "aws"
page_title: "aws_vpn_gateway_route_propagation"
description: |-
  Manages the route propagation from a VPN gateway to a route table.
---

[timeouts]: https://developer.hashicorp.com/terraform/plugin/framework/resources/timeouts

# Resource: aws_vpn_gateway_route_propagation

Manages the route propagation from a VPN gateway to a route table.
The routes received from the customer gateway over BGP are installed into the route table.

~> **Important** The route propagation can be enabled in one route table of a VPC only.
An attempt to enable it in another route table fails with the `Propagation.AlreadyEnabled` error.

~> **Note** This resource should not be used with a route table that has the `propagating_vgws` argument set.
If that argument is set, any route propagation not explicitly listed in its value will be removed.

-> **Note** For convenience, the ID of the VPN gateway is the same as the ID of the VPC to which it belongs (`vpc-ABCD1234`/`vgw-ABCD1234`).
This resource requires the ID of the VPN gateway in the `vgw-ABCD1234` form.

## Example usage

```terraform
resource "aws_vpc" "example" {
  cidr_block = "10.1.0.0/16"
}

resource "aws_route_table" "example" {
  vpc_id = aws_vpc.example.id
}

resource "aws_vpn_gateway_route_propagation" "example" {
  vpn_gateway_id = replace(aws_vpc.example.id, "/vpc/", "vgw")
  route_table_id = aws_route_table.example.id
}
```

## Argument reference

The following arguments are required:

* `route_table_id` - (Required, Forces new resource, String) The ID of the route table to propagate the routes into.
* `vpn_gateway_id` - (Required, Forces new resource, String) The ID of the VPN gateway to propagate the routes from.

## Attribute reference

In addition to all arguments above, the following attributes are exported:

* `id` - (String) The ID of the route propagation, in the `<vpn_gateway_id>_<route_table_id>` form.

## Timeouts

The `timeouts` block allows you to specify [timeouts] for certain actions:

* `create` - (Default `2 minutes`) Used for enabling the route propagation.
* `delete` - (Default `2 minutes`) Used for disabling the route propagation.

## Import

Importing is not currently supported.
