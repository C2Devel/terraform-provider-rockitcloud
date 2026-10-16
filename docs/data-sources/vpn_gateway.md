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

~> **Important** The platform currently ignores the filters of the VPN gateway search.
The `attached_vpc_id`, `filter` and `state` arguments do not narrow the result: if the project has one VPN gateway, it is returned whatever the filters are, otherwise the data source fails with the `multiple EC2 VPN Gateways matched` error.
Use the `id` argument.

## Example usage

```terraform
resource "aws_vpc" "example" {
  cidr_block = "10.1.0.0/16"
}

data "aws_vpn_gateway" "example" {
  id = aws_vpc.example.id
}

output "vpn_gateway_id" {
  value = data.aws_vpn_gateway.example.id
}
```

## Argument reference

The arguments of this data source act as filters for querying the available VPN gateways.
The given filters must match exactly one VPN gateway whose data will be exported as attributes.

* `attached_vpc_id` - (Optional, String) The ID of the VPC attached to the VPN gateway.
* `filter` - (Optional, [Block](#filter)) The filter of the search.
* `id` - (Optional, String) The ID of the VPN gateway or of the VPC to which it belongs.
* `state` - (Optional, String) The state of the VPN gateway.
    * _Valid values:_ `pending`, `available`, `deleting`

### filter

Each block has the following structure:

* `name` - (Required, String) The name of the filter.
  See the supported names in the [EC2 API documentation][describe-vpn-gateways].
* `values` - (Required, Set of strings) The values of the filter.
  A VPN gateway is selected if any of the values matches.

### Unsupported arguments

~> **Note** The platform does not support these arguments, so they have no effect on the search.
They are also exported with preset values, which are `0` for `amazon_side_asn` and empty values for the others.

The following arguments are not currently supported:

`amazon_side_asn`, `availability_zone`, `tags`.

## Attribute reference

### Supported attributes

In addition to all arguments above, the following attributes are exported:

* `arn` - (String) The ARN of the VPN gateway.
* `attached_vpc_id` - (String) The ID of the VPC attached to the VPN gateway.
* `id` - (String) The ID of the VPN gateway.
* `state` - (String) The state of the VPN gateway.
