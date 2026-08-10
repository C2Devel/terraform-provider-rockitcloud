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

## Example usage

```terraform
data "aws_launch_template" "example" {
  name = "tf-lt"
}
```

### Search by filter

```terraform
data "aws_launch_template" "example" {
  filter {
    name   = "launch-template-name"
    values = ["some-template"]
  }
}
```

## Argument reference

The following arguments are supported:

* `filter` - (Optional, [Block](#filter)) One or more name/value pairs to use as filters.
    * _Valid values:_ See supported names and values in [EC2 API documentation][describe-lts]
* `id` - (Optional, String) The ID of the specific launch template to retrieve.
* `name` - (Optional, String) The name of the launch template.
* `tags` - (Optional, Map of strings) Key-value pairs. Must exactly match pairs on the desired resource.

### filter

* `name` - (Required, String) The name of the filter.
    * _Constraints:_ Filter names are case-sensitive
* `values` - (Required, List of strings) One or more filter values.
    * _Constraints:_ Filter values are case-sensitive

## Attribute reference

In addition to all arguments above, the following attributes are exported:

* `id` - (String) The ID of the launch template.

This resource also exports a full set of attributes corresponding to the arguments of the [`aws_launch_template`](../resources/launch_template.md) resource.
