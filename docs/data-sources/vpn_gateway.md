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

-> **Note** For convenience, the ID of the VPN gateway is the same as the ID of the VPC, to which it belongs (`vpc-ABCD1234`/`vgw-ABCD1234`).

## Example Usage

```terraform
data "aws_vpn_gateway" "selected" {
  filter {
    name   = "tag:Name"
    values = ["vpn-gw"]
  }
}

output "vpn_gateway_id" {
  value = data.aws_vpn_gateway.selected.id
}
```

## Argument Reference

The arguments of this data source act as filters for querying the available VPN gateways.
The given filters must match exactly one VPN gateway whose data will be exported as attributes.

* `attached_vpc_id` - (Optional) ID of a VPC attached to the specific VPN gateway to retrieve.
* `availability_zone` - (Optional) The availability zone of the specific VPN gateway to retrieve.
* `filter` - (Optional) One or more name/value pairs to use as filters.
    * _Valid values:_ See supported names and values in [EC2 API documentation][describe-vpn-gateways]
* `id` - (Optional) ID of the specific VPN gateway to retrieve.
* `state` - (Optional) The state of the specific VPN gateway to retrieve.
* `tags` - (Optional) Map of tags, each pair of which must exactly match
  a pair on the desired VPN gateway.

## Attribute Reference

All the argument attributes are also exported as result attributes.
