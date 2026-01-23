---
subcategory: "Transit Gateway"
layout: "aws"
page_title: "aws_ec2_transit_gateway_route_table"
description: |-
  Provides information about a transit gateway route table.
---

[describe-tgw-rtb]: https://docs.k2.cloud/en/api/ec2/actions/transit_gateways/DescribeTransitGatewayRouteTables.html

# Data Source: aws_ec2_transit_gateway_route_table

Provides information about a transit gateway route table.

## Example Usage

### By Filter

```terraform
data "aws_ec2_transit_gateway_route_table" "selected" {
  filter {
    name   = "state"
    values = ["available"]
  }

  filter {
    name   = "transit-gateway-id"
    values = ["tgw-12345678"]
  }
}
```

### By Identifier

```terraform
data "aws_ec2_transit_gateway_route_table" "selected" {
  id = "tgw-rtb-12345678"
}
```

## Argument Reference

The following arguments are supported:

* `filter` - (Optional) One or more name/value pairs to use as filters.
    * _Valid values:_ See supported names and values in [EC2 API documentation][describe-tgw-rtb]
* `id` - (Optional) The ID of the transit gateway route table.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `arn` - The Amazon Resource Name (ARN) of the transit gateway route table.
* `default_association_route_table` - Indicates whether this is the default association route table for the transit gateway.
* `default_propagation_route_table` - Indicates whether this is the default propagation route table for the transit gateway.
* `id` - The ID of the transit gateway route table.
* `tags` - Map of tags assigned to the transit gateway route table.
* `transit_gateway_id` - The ID of the transit gateway.
