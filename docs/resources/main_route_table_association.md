---
subcategory: "VPC (Virtual Private Cloud)"
layout: "aws"
page_title: "aws_main_route_table_association"
description: |-
  Manages the main routing table of a VPC.
---

[route-tables]: https://docs.k2.cloud/en/services/networking/routetables.html

# Resource: aws_main_route_table_association

Manages the main routing table of a VPC.

~> **Note** **Do not** use both `aws_default_route_table` to manage a default route table **and** `aws_main_route_table_association` with the same VPC due to possible route conflicts. See [aws_default_route_table](default_route_table.md) documentation for more details.
For more information, see the documentation on [route tables][route-tables].
For information about managing normal route tables in Terraform, see [`aws_route_table`](route_table.md).

## Example usage

### Basic example

```terraform
resource "aws_vpc" "example" {
  cidr_block = "10.1.0.0/16"
}

resource "aws_route_table" "example" {
  vpc_id = aws_vpc.example.id
}

resource "aws_main_route_table_association" "example" {
  vpc_id         = aws_vpc.example.id
  route_table_id = aws_route_table.example.id
}
```

## Argument reference

The following arguments are supported:

* `route_table_id` - (Required, Editable, String) The ID of the route table to set as the new main route table for the target VPC.
* `vpc_id` - (Required, Editable, String) The ID of the VPC the main route table should be associated with.

## Attribute reference

In addition to all arguments above, the following attributes are exported:

* `id` - (String) The ID of the route table association.
* `original_route_table_id` - (String) Used internally, see [Notes](#notes) below.

## Notes

On VPC creation, the cloud always creates an initial main route table.
This resource records the ID of that route table under `original_route_table_id`.
The "Delete" action for a `main_route_table_association` consists of resetting this original table as the main route table for the VPC.
You'll see this additional route table in the cloud console.
It must remain intact in order for the `main_route_table_association` delete to work properly.

## Timeouts

Timeouts usage for main route table associations is not currently supported.

## Import

Import of the main route table associations is not currently supported.
