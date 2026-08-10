---
subcategory: "EC2 (Elastic Compute Cloud)"
layout: "aws"
page_title: "aws_ec2_tag"
description: |-
  Manages an individual EC2 resource tag
---

[ignore-tags]: https://registry.terraform.io/providers/hashicorp/aws/latest/docs#ignore_tags-configuration-block

# Resource: aws_ec2_tag

Manages an individual EC2 resource tag. This resource should only be used when Terraform was not used to create EC2 resources (e.g., images).

~> **Note** This tagging resource should not be combined with the Terraform resource for managing the parent resource. For example, using `aws_vpc` and `aws_ec2_tag` to manage tags of the same VPC will cause a perpetual difference where the `aws_vpc` resource will try to remove the tag being added by the `aws_ec2_tag` resource.

~> **Note** This tagging resource does not use the [provider `ignore_tags` configuration][ignore-tags].

## Example usage

```terraform
resource "aws_ec2_tag" "example" {
  resource_id = "vol-12345678"
  key         = "tag-from-tf"
  value       = "tf-tag"
}
```

## Argument reference

The following arguments are supported:

* `resource_id` - (Required, Editable, String) The ID of the EC2 resource to manage the tag for.
* `key` - (Required, Editable, String) The tag name.
* `value` - (Required, Editable, String) The value of the tag.

## Attribute reference

In addition to all arguments above, the following attributes are exported:

* `id` - (String) `resource_id` and `key` separated by a comma (`,`).

## Timeouts

Timeouts usage for the EC2 tag is not currently supported.

## Import

`aws_ec2_tag` can be imported using `id`, for example:

```
$ terraform import aws_ec2_tag.example tgw-attach-12345678,Name
```
