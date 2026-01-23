---
subcategory: "EC2 (Elastic Compute Cloud)"
layout: "aws"
page_title: "aws_launch_template"
description: |-
  Provides information about a launch template.
---

[describe-lts]: https://docs.k2.cloud/en/api/ec2/actions/launch_templates/DescribeLaunchTemplates.html

# Data Source: aws_launch_template

Provides information about a launch template.

## Example Usage

```terraform
data "aws_launch_template" "example" {
  name = "tf-lt"
}
```

### Search by Filter

```terraform
data "aws_launch_template" "example" {
  filter {
    name   = "launch-template-name"
    values = ["some-template"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `filter` - (Optional) One or more name/value pairs to use as filters.
    * _Valid values:_ See supported names and values in [EC2 API documentation][describe-lts]
* `id` - (Optional) The ID of the specific launch template to retrieve.
* `name` - (Optional) The name of the launch template.
* `tags` - (Optional) Map of tags, each pair of which must exactly match a pair on the desired launch template.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the launch template.

This resource also exports a full set of attributes corresponding to the arguments of the [`aws_launch_template`](../resources/launch_template.md) resource.
