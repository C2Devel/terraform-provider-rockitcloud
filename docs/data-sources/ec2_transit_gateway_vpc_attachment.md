---
subcategory: "Transit Gateway"
layout: "aws"
page_title: "aws_ec2_transit_gateway_vpc_attachment"
description: |-
  Provides information about a transit gateway VPC attachment.
---

[describe-tgw-vpc-attachments]: https://docs.k2.cloud/en/api/ec2/actions/transit_gateways/DescribeTransitGatewayVpcAttachments.html

# Data Source: aws_ec2_transit_gateway_vpc_attachment

Provides information about a transit gateway VPC attachment.

## Example Usage

### By Filter

```terraform
data "aws_ec2_transit_gateway_vpc_attachment" "selected" {
  filter {
    name   = "vpc-id"
    values = ["vpc-12345678"]
  }
}
```

### By Identifier

```terraform
data "aws_ec2_transit_gateway_vpc_attachment" "selected" {
  id = "tgw-attach-12345678"
}
```

## Argument Reference

The following arguments are supported:

* `filter` - (Optional) One or more name/value pairs to use as filters.
    * _Valid values:_ See supported names and values in [EC2 API documentation][describe-tgw-vpc-attachments]
* `id` - (Optional) The ID of the transit gateway VPC attachment.

## Attribute Reference

### Supported attributes

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the transit gateway attachment.
* `subnet_ids` - List of subnet IDs.
* `tags` - Map of tags assigned to the transit gateway VPC attachment.
* `transit_gateway_id` - The ID of the transit gateway.
* `vpc_id` - The ID of the VPC.
* `vpc_owner_id` - The ID of the project that owns the VPC.

### Unsupported attributes

~> **Note** These attributes may be present in the `terraform.tfstate` file, but they have preset values and cannot be specified in configuration files.

The following attributes are not currently supported:

`appliance_mode_support`, `dns_support`, `ipv6_support`.
