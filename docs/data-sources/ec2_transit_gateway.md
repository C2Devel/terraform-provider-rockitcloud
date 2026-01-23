---
subcategory: "Transit Gateway"
layout: "aws"
page_title: "aws_ec2_transit_gateway"
description: |-
  Provides information about a transit gateway.
---

[describe-tgw]: https://docs.k2.cloud/en/api/ec2/actions/transit_gateways/DescribeTransitGateways.html

# Data Source: aws_ec2_transit_gateway

Provides information about a transit gateway.

## Example Usage

### By Filter

```terraform
data "aws_ec2_transit_gateway" "selected" {
  filter {
    name   = "owner-id"
    values = ["project@customer"]
  }
}
```

### By Identifier

```terraform
data "aws_ec2_transit_gateway" "selected" {
  id = "tgw-12345678"
}
```

## Argument Reference

The following arguments are supported:

* `filter` - (Optional) One or more name/value pairs to use as filters.
    * _Valid values:_ See supported names and values in [EC2 API documentation][describe-tgw]
* `id` - (Optional) The ID of the transit gateway.

## Attribute Reference

### Supported attributes

In addition to all arguments above, the following attributes are exported:

* `association_default_route_table_id` - The ID of the default association route table.
* `default_route_table_association` - Indicates whether the association with default association route table is created automatically.
* `default_route_table_propagation` - Indicates whether the routes are automatically propagated to the default propagation route table.
* `description` - The description of the transit gateway
* `owner_id` - The ID of the project that owns the transit gateway.
* `propagation_default_route_table_id` - The ID of the default propagation route table.
* `shared_owners` - List of project IDs that are granted access to the transit gateway.
* `tags` - Map of tags assigned to the transit gateway.

### Unsupported attributes

~> **Note** These attributes may be present in the `terraform.tfstate` file, but they have preset values and cannot be specified in configuration files.

The following attributes are not currently supported:

`amazon_side_asn`, `arn`, `auto_accept_shared_attachments`, `dns_support`, `multicast_support`, `transit_gateway_cidr_blocks`, `vpn_ecmp_support`.
