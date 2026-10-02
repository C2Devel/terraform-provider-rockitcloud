---
subcategory: "Transit Gateway"
layout: "aws"
page_title: "aws_ec2_transit_gateway_vpc_attachment"
description: |-
  Manages a transit gateway VPC attachment.
---

[default-tags]: https://registry.terraform.io/providers/hashicorp/aws/latest/docs#default_tags-configuration-block

# Resource: aws_ec2_transit_gateway_vpc_attachment

Manages a transit gateway VPC attachment.

## Example Usage

```terraform
resource "aws_vpc" "example" {
  cidr_block = "172.16.0.0/16"

  tags = {
    Name = "tf-vpc"
  }
}

resource "aws_subnet" "example" {
  vpc_id            = aws_vpc.example.id
  cidr_block        = cidrsubnet(aws_vpc.example.cidr_block, 4, 1)
  availability_zone = "ru-msk-vol52"

  tags = {
    Name = "tf-subnet"
  }
}

resource "aws_ec2_transit_gateway" "example" {
  description = "tf example"

  tags = {
    Name = "tf-tgw"
  }
}

resource "aws_ec2_transit_gateway_vpc_attachment" "example" {
  subnet_ids         = [aws_subnet.example.id]
  transit_gateway_id = aws_ec2_transit_gateway.example.id
  vpc_id             = aws_vpc.example.id

  tags = {
    Name = "tf-tgw-vpc-attachment"
  }
}
```

## Argument Reference

The following arguments are supported:

* `subnet_ids` - (Required) List of subnet IDs.
* `transit_gateway_id` - (Required) The ID of the transit gateway.
* `vpc_id` - (Required) The ID of the VPC.
* `tags` - (Optional)  Map of tags to assign to the transit gateway VPC attachment.
  If a provider [`default_tags` configuration block][default-tags] is used,
  tags with matching keys will overwrite those defined at the provider level.
* `transit_gateway_default_route_table_association` - (Optional) Indicates whether the transit gateway VPC attachment should be associated with the transit gateway default association route table.
    * _Default value:_ `true`
* `transit_gateway_default_route_table_propagation` - (Optional) Indicates whether the transit gateway VPC attachment should propagate routes to the transit gateway default propagation route table.
    * _Default value:_ `true`

~> **Note** `transit_gateway_default_route_table_association` and `transit_gateway_default_route_table_propagation`
cannot be configured for shared transit gateways.

## Attribute Reference

### Supported attributes

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the transit gateway attachment.
* `tags_all` - Map of tags assigned to the transit gateway attachment, including those inherited from the provider [`default_tags` configuration block][default-tags].
* `vpc_owner_id` - The ID of the project that owns the VPC.

### Unsupported attributes

~> **Note** These attributes may be present in the `terraform.tfstate` file, but they have preset values and cannot be specified in configuration files.

The following attributes are not currently supported:

`appliance_mode_support`, `dns_support`, `ipv6_support`.

## Import

The transit gateway VPC attachment can be imported using `id`, e.g.,

```
$ terraform import aws_ec2_transit_gateway_vpc_attachment.example tgw-attach-12345678
```
