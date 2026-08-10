---
subcategory: "EC2 (Elastic Compute Cloud)"
layout: "aws"
page_title: "aws_ami_launch_permission"
description: |-
  Adds a launch permission to an Amazon Machine Image (AMI).
---

# Resource: aws_ami_launch_permission

Adds a launch permission to an Amazon Machine Image (AMI).

## Example usage

### AWS account ID

```terraform
resource "aws_ami_launch_permission" "example" {
  image_id   = "cmi-12345678"
  account_id = "123456789012"
}
```

### Public access

```terraform
# The cloud currently restricts adding public access permissions to images.
# Applying the resource must throw an error.
resource "aws_ami_launch_permission" "example" {
  image_id = "cmi-12345678"
  group    = "all"
}
```

## Argument reference

The following arguments are supported:

* `image_id` - (Required, Forces new resource, String) The ID of the image.
* `account_id` - (Optional, Forces new resource, String) The ID of the project for the launch permission.
    * _Example:_ `project@customer`
* `group` - (Optional, Forces new resource, String) The name of the group for the launch permission.
    * _Valid values:_ `all`

## Attribute reference

### Supported attributes

In addition to all arguments above, the following attributes are exported:

* `id` - (String) The ID of the launch permission.

### Unsupported attributes

~> **Note** These attributes may be present in the `terraform.tfstate` file, but they have preset values and cannot be specified in configuration files.

The following attributes are not currently supported:

`organization_arn`, `organizational_unit_arn`.

## Timeouts

Timeouts usage for launch permission is not currently supported.

## Import

Launch permissions can be imported using the permission ID and image ID separated by a slash (`/`), for example:

```
$ terraform import aws_ami_launch_permission.example 123456789012/cmi-12345678
```

~> **Note** The import format is `[ACCOUNT-ID|GROUP-NAME|ORGANIZATION-ARN|ORGANIZATIONAL-UNIT-ARN]/IMAGE-ID`.
