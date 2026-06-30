---
subcategory: "EC2 (Elastic Compute Cloud)"
layout: "aws"
page_title: "aws_ami_ids"
description: |-
  Provides a list of image IDs.
---

[describe-images]: https://docs.k2.cloud/en/api/ec2/actions/images/DescribeImages.html

# Data Source: aws_ami_ids

Provides a list of image IDs.

## Example usage

### Basic example

```terraform
data "aws_ami_ids" "example" {
  owners = ["self"]
}
```

## Argument reference

The following arguments are supported:

* `owners` - (Required, List of strings) List of image owners to limit search. At least one value must be specified.
    * _Valid values:_ `project@customer` or `self`
* `executable_users` - (Optional, List of strings) Limit search to project with *explicit* launch permission on the image.
    * _Valid values:_ `all`, `project@customer` or `self`
* `filter` - (Optional, [Block](#filter)) One or more name/value pairs to use as filters.
    * _Valid values:_ See supported names and values in [EC2 API documentation][describe-images]
* `name_regex` - (Optional, String) A regex string to apply to the image list returned by the EC2 API.
  It is recommended to combine this with other options to narrow down the list the EC2 API returns.
* `sort_ascending` - (Optional, Boolean) Indicates whether to sort images by creation time.
    * _Default value:_ `false`

### filter

* `name` - (Required, String) The name of the filter.
    * _Constraints:_ Filter names are case-sensitive
* `values` - (Required, List of strings) One or more filter values.
    * _Constraints:_ Filter values are case-sensitive

## Attribute reference

In addition to all arguments above, the following attribute is exported:

* `ids` - (List of strings) The list of image IDs, sorted by creation time according to `sort_ascending`.
