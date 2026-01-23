---
subcategory: "VPC (Virtual Private Cloud)"
layout: "aws"
page_title: "aws_internet_gateway"
description: |-
  Provides information about an internet gateway.
---

[describe-igws]: https://docs.k2.cloud/en/api/ec2/actions/vpcs/DescribeInternetGateways.html

# Data Source: aws_internet_gateway

Provides information about an internet gateway.

## Example Usage

```terraform
data "aws_internet_gateway" "selected" {
  filter {
    name   = "attachment.vpc-id"
    values = ["vpc-12345678"]
  }
}
```

## Argument Reference

The arguments of this data source act as filters for querying the available
internet gateway. The given filters must match exactly one
internet gateway whose data will be exported as attributes.

* `filter` - (Optional) One or more name/value pairs to use as filters.
    * _Valid values:_ See supported names and values in [EC2 API documentation][describe-igws]
* `internet_gateway_id` - (Optional) The ID of the internet gateway.
* `tags` - (Optional) Map of tags. Each tag must exactly match a tag on the desired internet gateway.

## Attribute Reference

All arguments except `filter` block are also exported as result attributes.
If any fields are missing from the configuration,
then this data source will populate them with data for the selected internet gateway.

* `arn` - The Amazon Resource Name (ARN) of the internet gateway.
* `attachments` - List of VPC attachments to the internet gateway. It can contain 0 or 1 element.
  The structure of this block is [described below](#attachments).
* `id` - The ID of the internet gateway.
* `owner_id` - The ID of the project that the internet gateway belongs to.

### attachments

* `state` - The current state of the VPC attachment to the internet gateway.
* `vpc_id` - The ID of the attached VPC.
