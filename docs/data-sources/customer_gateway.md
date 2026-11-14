---
subcategory: "VPN (Site-to-Site)"
layout: "aws"
page_title: "aws_customer_gateway"
description: |-
  Provides information about a customer gateway.
---

[describe-customer-gateways]: https://docs.k2.cloud/en/api/ec2/actions/customer_gateways/DescribeCustomerGateways.html

# Data Source: aws_customer_gateway

Provides information about a customer gateway.

-> **Note** For convenience, the ID of the VPN gateway is the same as the ID of the VPC to which it belongs (`vpc-ABCD1234`/`vgw-ABCD1234`).

## Example usage

### Basic example

```terraform
data "aws_customer_gateway" "selected" {
  filter {
    name   = "tag:Name"
    values = ["foo-prod"]
  }
}

resource "aws_vpc" "example" {
  cidr_block = "10.1.0.0/16"
}

resource "aws_vpn_connection" "example" {
  vpn_gateway_id      = aws_vpc.example.id
  customer_gateway_id = data.aws_customer_gateway.selected.id
  type                = data.aws_customer_gateway.selected.type
}
```

## Argument reference

The arguments of this data source act as filters for querying the available customer gateways.
The given filters must match exactly one customer gateway whose data will be exported as attributes.

* `filter` - (Optional, [Block](#filter)) One or more name/value pairs to use as filters.
    * _Valid values:_ See supported names and values in [EC2 API documentation][describe-customer-gateways]
* `id` - (Optional, String) The ID of the customer gateway.

### filter

* `name` - (Required, String) The name of the filter.
    * _Constraints:_ Filter names are case-sensitive
* `values` - (Required, Set of strings) One or more filter values.
    * _Constraints:_ Filter values are case-sensitive

## Attribute reference

In addition to all arguments above, the following attributes are exported:

* `arn` - (String) The ARN of the customer gateway.
* `bgp_asn` - (Integer) The Border Gateway Protocol (BGP) Autonomous System Number (ASN) of the gateway.
* `ip_address` - (String) The IPv4 address of the gateway's Internet-routable external interface.
* `tags` - (Map of strings) Key-value pairs assigned to the customer gateway.
* `type` - (String) The type of the customer gateway.
    * _Valid values:_ `ipsec.1` or `ipsec.legacy`
