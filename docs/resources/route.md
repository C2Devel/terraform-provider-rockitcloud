---
subcategory: "VPC (Virtual Private Cloud)"
layout: "aws"
page_title: "aws_route"
description: |-
  Creates a routing entry in a VPC routing table.
---

[timeouts]: https://developer.hashicorp.com/terraform/plugin/framework/resources/timeouts

# Resource: aws_route

Creates a routing table entry (a route) in a VPC routing table.

~> **Note on route tables and routes** Terraform currently provides both a standalone route resource and a [`aws_route_table`](route_table.md) resource with routes defined inline. At this time you cannot use a route table with inline routes in conjunction with any route resources. Doing so will cause a conflict of rule settings and will overwrite rules.

## Example usage

### Basic example

```terraform
variable "route_table_id" {}
variable "network_interface_id" {}

resource "aws_route" "example" {
  route_table_id         = var.route_table_id
  destination_cidr_block = "10.0.0.0/22"
  network_interface_id   = var.network_interface_id
}
```

## Argument reference

The following argument is supported:

* `route_table_id` - (Required, Forces new resource, String) The ID of the routing table.

The following destination argument must be supplied:

* `destination_cidr_block` - (Required, Forces new resource, String) The destination CIDR block.

One of the following target arguments must be supplied:

* `gateway_id` - (Optional, Editable, String) The ID of the internet gateway.
* `nat_gateway_id` - (Optional, Editable, String) The ID of the NAT gateway.
* `network_interface_id` - (Optional, Editable, String) The ID of the network interface.
* `transit_gateway_id` - (Optional, Editable, String) The ID of the transit gateway.

~> **Note** The `instance_id` argument was removed and cannot be specified in configuration files anymore.
To route traffic to an instance, use `network_interface_id`.

## Attribute reference

### Supported attributes

~> **Note** Only the arguments that are configured (one of the above) will be exported as an attribute once the resource is created.

In addition to all arguments above, the following attributes are exported:

* `id` - (String) The route identifier computed from the routing table identifier and route destination.
* `instance_id` - (String) The ID of the instance the target network interface is attached to.
* `instance_owner_id` - (String) The ID of the project that owns the instance.
* `origin` - (String) Describes how the route was created - by `CreateRouteTable`, `CreateRoute` or `EnableVgwRoutePropagation`.
* `state` - (String) The state of the route - `active` or `blackhole`.

### Unsupported attributes

~> **Note** These attributes may be present in the `terraform.tfstate` file, but they have preset values and cannot be specified in configuration files.

The following attributes are not currently supported:

`carrier_gateway_id`, `core_network_arn`, `destination_ipv6_cidr_block`, `destination_prefix_list_id`, `egress_only_gateway_id`, `local_gateway_id`, `vpc_endpoint_id`, `vpc_peering_connection_id`.

## Timeouts

The `timeouts` block allows you to specify [timeouts] for certain actions:

- `create` - (Default `5 minutes`) Used for route creation.
- `update` - (Default `2 minutes`) Used for updating the route.
- `delete` - (Default `5 minutes`) Used for route deletion.

## Import

Individual routes can be imported using `ROUTETABLEID_DESTINATION`.
For example, import a route in route table `rtb-12345678` with an IPv4 destination CIDR of `10.1.0.0/16` like this:

```console
$ terraform import aws_route.my_route rtb-12345678_10.1.0.0/16
```
