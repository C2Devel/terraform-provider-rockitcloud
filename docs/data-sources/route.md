---
subcategory: "VPC (Virtual Private Cloud)"
layout: "aws"
page_title: "aws_route"
description: |-
  Provides information about a route.
---

# Data Source: aws_route

Provides information about a route.

## Example usage

### Basic example

The following example shows how a CIDR value may be used to find the ID of a network interface and read the network interface data using the data source.

```terraform
variable "subnet_id" {}

data "aws_route_table" "selected" {
  subnet_id = var.subnet_id
}

data "aws_route" "selected" {
  route_table_id         = data.aws_route_table.selected.id
  destination_cidr_block = "10.0.1.0/24"
}

data "aws_network_interface" "selected" {
  id = data.aws_route.selected.network_interface_id
}
```

## Argument reference

The arguments of this data source act as filters for querying specific routes in the current region.

~> **Note** The given filters must exactly match the resource whose data will be exported as attributes.

* `route_table_id` - (Required, String) The ID of the specific route table containing the route entry.
* `destination_cidr_block` - (Optional, String) The CIDR block of the route.
* `gateway_id` - (Optional, String) The ID of the internet gateway.
* `instance_id` - (Optional, String) The ID of the instance.
* `nat_gateway_id` - (Optional, String) The ID of the NAT gateway.
* `network_interface_id` - (Optional, String) The ID of the network interface.
* `transit_gateway_id` - (Optional, String) The ID of the transit gateway.

## Attribute reference

### Supported attributes

In addition to all arguments above, the following attributes are exported:

* `id` - (String) The route identifier computed from the routing table identifier and route destination.

### Unsupported attributes

~> **Note** These attributes may be present in the `terraform.tfstate` file, but they have preset values and cannot be specified in configuration files.

The following attributes are not currently supported:

`carrier_gateway_id`, `core_network_arn`, `destination_ipv6_cidr_block`, `destination_prefix_list_id`, `egress_only_gateway_id`, `local_gateway_id`, `vpc_endpoint_id`, `vpc_peering_connection_id`.
