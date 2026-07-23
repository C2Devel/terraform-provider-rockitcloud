---
subcategory: "VPC (Virtual Private Cloud)"
layout: "aws"
page_title: "aws_internet_gateway"
description: |-
  Manages an internet gateway.
---

[default-tags]: https://www.terraform.io/docs/providers/aws/index.html#default_tags-configuration-block
[internet-gateways]: https://docs.k2.cloud/en/services/networking/igw.html

# Resource: aws_internet_gateway

Manages an internet gateway.

For more information about internet gateways, see the documentation on [Internet gateways][internet-gateways].

## Example usage

### Basic example

```terraform
resource "aws_vpc" "example" {
  cidr_block = "15.0.0.0/16"

  tags = {
    Name = "tf-vpc"
  }
}

resource "aws_internet_gateway" "example" {
  vpc_id = aws_vpc.example.id

  tags = {
    Name = "tf-igw"
  }
}
```

## Argument reference

The following arguments are supported:

* `tags` - (Optional, Editable, Map of strings) Key-value pairs to assign to the resource.
    If the [`default_tags` configuration block][default-tags] block is used within a provider configuration, the tags with matching keys will overwrite those defined at the provider level.
* `vpc_id` - (Optional, Editable, String) The ID of the VPC to which the internet gateway will be attached.
  See the [aws_internet_gateway_attachment](internet_gateway_attachment.md) resource for another way to attach an internet gateway to a VPC.

-> **Note** It's recommended to explicitly specify an internet gateway as a dependency for the resources that require Internet access.
For example:

```terraform
resource "aws_internet_gateway" "example" {
  vpc_id = aws_vpc.example.id

  tags = {
    Name = "tf-igw"
  }
}

resource "aws_instance" "example" {
  # ... other arguments ...

  depends_on = [aws_internet_gateway.example]
}
```

## Attribute reference

In addition to all arguments above, the following attributes are exported:

* `arn` - (String) The Amazon Resource Name (ARN) of the internet gateway.
* `id` - (String) The ID of the internet gateway.
* `owner_id` - (String) The ID of the project that owns the internet gateway.
* `tags_all` - (Map of strings) Key-value pairs assigned to the resource, including any tags inherited from the [`default_tags` configuration block][default-tags] if used within a provider configuration.

## Timeouts

Timeouts usage for internet gateways is not currently supported.

## Import

Internet gateways can be imported using `id`, for example:

```
$ terraform import aws_internet_gateway.example igw-12345678
```
