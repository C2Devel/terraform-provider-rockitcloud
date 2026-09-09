---
subcategory: "Transit Gateway"
layout: "aws"
page_title: "aws_ec2_transit_gateway_route_table"
description: |-
  Manages a transit gateway route table.
---

[default-tags]: https://registry.terraform.io/providers/hashicorp/aws/latest/docs#default_tags-configuration-block

# Resource: aws_ec2_transit_gateway_route_table

Manages a transit gateway route table.

## Example Usage

```terraform
resource "aws_ec2_transit_gateway" "example" {
  description = "tf example"

  tags = {
    Name = "tf-tgw"
  }
}

resource "aws_ec2_transit_gateway_route_table" "example" {
  transit_gateway_id = aws_ec2_transit_gateway.example.id

  tags = {
    Name = "tf-rtb"
  }
}
```

## Argument Reference

The following arguments are supported:

* `transit_gateway_id` - (Required) The ID of the transit gateway.
* `tags` - (Optional) Map of tags to assign to the transit gateway route table.
  If a provider [`default_tags` configuration block][default-tags] is used,
  tags with matching keys will overwrite those defined at the provider level.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `arn` - The Amazon Resource Name (ARN) of the transit gateway route table.
* `default_association_route_table` - Indicates whether this is the default association route table for the transit gateway.
* `default_propagation_route_table` - Indicates whether this is the default propagation route table for the transit gateway.
* `id` - The ID of the transit gateway route table.
* `tags_all` - Map of tags assigned to the transit gateway route table, including those inherited from the provider [`default_tags` configuration block][default-tags].

## Import

The transit gateway route table can be imported using `id`, e.g.,

```
$ terraform import aws_ec2_transit_gateway_route_table.example tgw-rtb-12345678
```
