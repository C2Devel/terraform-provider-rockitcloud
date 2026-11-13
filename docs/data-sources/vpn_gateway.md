---
subcategory: "VPN (Site-to-Site)"
layout: "aws"
page_title: "aws_vpn_gateway"
description: |-
  Provides information about a VPN gateway.
---

[describe-vpn-gateways]: https://docs.k2.cloud/en/api/ec2/actions/vpn_gateways/DescribeVpnGateways.html

# Data Source: aws_vpn_gateway

Provides information about a VPN gateway.
A VPN gateway is created automatically for each VPC.

-> **Note** For convenience, the ID of the VPN gateway is the same as the ID of the VPC to which it belongs (`vpc-ABCD1234`/`vgw-ABCD1234`).
Either of them can be specified in the `id` argument, the exported `id` is always in the `vgw-ABCD1234` form.

## Example usage

### Basic example

```terraform
resource "aws_vpc" "example" {
  cidr_block = "10.1.0.0/16"
}

data "aws_vpn_gateway" "example" {
  attached_vpc_id = aws_vpc.example.id
}

output "vpn_gateway_id" {
  value = data.aws_vpn_gateway.example.id
}
```

## Argument reference

The arguments of this data source act as filters for querying the available VPN gateways.
The given filters must match exactly one VPN gateway whose data will be exported as attributes.

* `attached_vpc_id` - (Optional, String) The ID of the VPC attached to the VPN gateway.
* `availability_zone` - (Optional, String) The name of the availability zone in which the VPN gateway is available.
* `filter` - (Optional, [Block](#filter)) One or more name/value pairs to use as filters.
    * _Valid values:_ See supported names and values in [EC2 API documentation][describe-vpn-gateways]
* `id` - (Optional, String) The ID of the VPN gateway or of the VPC to which it belongs.
* `state` - (Optional, String) The state of the VPN gateway.
    * _Valid values:_ `available`

### filter

* `name` - (Required, String) The name of the filter.
    * _Constraints:_ Filter names are case-sensitive
* `values` - (Required, Set of strings) One or more filter values.
    * _Constraints:_ Filter values are case-sensitive

## Attribute reference

In addition to all arguments above, the following attributes are exported:

* `arn` - (String) The ARN of the VPN gateway.
