---
subcategory: "Transit Gateway"
layout: "aws"
page_title: "aws_ec2_transit_gateway_dx_gateway_attachment"
description: |-
  Provides information about the attachment of an EC2 transit gateway to a Direct Connect gateway.
---

[describe-transit-gateway-attachments]: https://docs.k2.cloud/en/api/ec2/actions/transit_gateways/DescribeTransitGatewayAttachments.html

# Data Source: aws_ec2_transit_gateway_dx_gateway_attachment

Provides information about the attachment of an EC2 transit gateway to a Direct Connect gateway.

## Example Usage

Using EC2 transit gateway and Direct Connect gateway identifiers to get information on the attachment

```terraform
resource "aws_dx_gateway" "example" {
  name            = "tf-dxgw-example"
  amazon_side_asn = "64512"
}

resource "aws_ec2_transit_gateway" "example" {
}

data "aws_ec2_transit_gateway_dx_gateway_attachment" "selected" {
  transit_gateway_id = aws_ec2_transit_gateway.example.id
  dx_gateway_id      = aws_dx_gateway.example.id
}
```

## Argument Reference

The following arguments are supported:

* `dx_gateway_id` - (Optional) The ID of the Direct Connect gateway.
* `filter` - (Optional) One or more name/value pairs to use as filters.
    * _Valid values:_ See supported names and values in [EC2 API documentation][describe-transit-gateway-attachments]
* `transit_gateway_id` - (Optional) The ID of the EC2 transit gateway.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the EC2 transit gateway attachment.
* `tags` - Map of tags assigned to the EC2 transit gateway attachment.
