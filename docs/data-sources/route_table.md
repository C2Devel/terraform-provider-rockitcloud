---
subcategory: "VPC (Virtual Private Cloud)"
layout: "aws"
page_title: "aws_route_table"
description: |-
  Provides information about a route table.
---

[describe-route-tables]: https://docs.k2.cloud/en/api/ec2/actions/routes/DescribeRouteTables.html

# Data Source: aws_route_table

Provides information about a route table.

This data source can be used when a module accepts a subnet ID as an input variable and needs to, for example, add a route to a route table.

## Example usage

### Specific example

The following example shows how to accept a route table ID as an input variable and use this data source to obtain the data necessary to create a route.

```terraform
variable "vpc_id" {}
variable "network_interface_id" {}

data "aws_route_table" "selected" {
  vpc_id = var.vpc_id
}

resource "aws_route" "example" {
  route_table_id         = data.aws_route_table.selected.id
  destination_cidr_block = "10.0.0.0/22"
  network_interface_id   = var.network_interface_id
}
```

## Argument reference

The arguments of this data source act as filters for querying route tables in the current region.

~> **Note** The given filters must exactly match the resource whose data will be exported as attributes.

The following arguments are optional:

* `filter` - (Optional, [Block](#filter)) One or more name/value pairs to use as filters.
    * _Valid values:_ See supported names and values in [EC2 API documentation][describe-route-tables]
* `route_table_id` - (Optional, String) The ID of the route table.
* `subnet_id` - (Optional, String) The ID of the subnet which is associated with the route table (not exported if not passed as a parameter).
* `tags` - (Optional, Map of strings) Key-value pairs. Must exactly match pairs on the required resource.
* `vpc_id` - (Optional, String) The ID of the VPC that owns the desired route table.

### filter

* `name` - (Required, String) The name of the filter.
    * _Constraints:_ Filter names are case-sensitive
* `values` - (Required, List of strings) One or more filter values.
    * _Constraints:_ Filter values are case-sensitive

## Attribute reference

### Supported attributes

In addition to all arguments above, the following attributes are exported:

* `arn` - (String) The Amazon Resource Name (ARN) of the route table.
* `id` - (String) The ID of the route table.
* `associations` - ([Block](#associations)) The list of associations.
* `routes` - ([Block](#routes)) The list of routes.

#### associations

* `main` - (Boolean) Indicates whether the route table is the main route table for the VPC.
* `route_table_association_id` - (String) The ID of the route table association.
* `route_table_id` - (String) The ID of the route table.
* `subnet_id` - (String) The ID of the subnet the route table is associated with. Only set when associated with a subnet.

#### routes

* `cidr_block` - (String) The CIDR block of the route.
* `gateway_id` - (String) The ID of the gateway.
* `instance_id` - (String) The ID of the instance.
* `nat_gateway_id` - (String) The ID of the NAT gateway.
* `network_interface_id` - (String) The ID of the network interface.
* `transit_gateway_id` - (String) The ID of the transit gateway.

### Unsupported attributes

~> **Note** These attributes may be present in the `terraform.tfstate` file, but they have preset values and cannot be specified in configuration files.

The following attributes are not currently supported:

`associations.gateway_id`, `gateway_id`, `owner_id`, `routes.carrier_gateway_id`, `routes.core_network_arn`, `routes.destination_prefix_list_id`, `routes.egress_only_gateway_id`, `routes.ipv6_cidr_block`, `routes.local_gateway_id`, `routes.vpc_endpoint_id`, `routes.vpc_peering_connection_id`.
