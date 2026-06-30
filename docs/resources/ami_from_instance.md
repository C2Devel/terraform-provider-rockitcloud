---
subcategory: "EC2 (Elastic Compute Cloud)"
layout: "aws"
page_title: "aws_ami_from_instance"
description: |-
  Creates an Amazon Machine Image (AMI) from an instance.
---

[default-tags]: https://registry.terraform.io/providers/hashicorp/aws/latest/docs#default_tags-configuration-block
[timeouts]: https://developer.hashicorp.com/terraform/plugin/framework/resources/timeouts

# Resource: aws_ami_from_instance

Creates an Amazon Machine Image (AMI) from an instance.

The created image will refer to implicitly created snapshots of the instance's
EBS volumes and mimick its assigned block device configuration at the time
the resource is created.

This resource is best applied to an instance that is stopped when this instance
is created, so that the contents of the created image are predictable. When
applied to an instance that is running, *the instance will be stopped before taking
the snapshots and then started back up again*, resulting in a period of
downtime.

Note that the source instance is inspected only at the initial creation of this
resource. Ongoing updates to the referenced instance will not be propagated into
the generated image. Users may taint or otherwise recreate the resource in order
to produce a fresh snapshot.

## Example usage

```terraform
resource "aws_ami_from_instance" "example" {
  name               = "terraform-example"
  source_instance_id = "i-12345678"
}
```

## Argument reference

The following arguments are supported:

* `name` - (Required, Forces new resource, String) A region-unique name for the image.
* `source_instance_id` - (Required, Forces new resource, String) The ID of the instance to use as the basis of the image.
* `tags` - (Optional, Editable, Map of strings) Key-value pairs to assign to the image. If the [`default_tags` configuration block][default-tags] is used within a provider configuration, the tags with matching keys will overwrite those defined at the provider level.

## Attribute reference

### Supported attributes

In addition to all arguments above, the following attributes are exported:

* `arn` - (String) The Amazon Resource Name (ARN) of the image.
* `id` - (String) The ID of the created image.

This resource also exports a full set of attributes corresponding to the arguments of the
[`aws_ami`](ami.md) resource, allowing the properties of the created image to be used elsewhere in the configuration.

### Unsupported attributes

~> **Note** This attribute may be present in the `terraform.tfstate` file, but it has a preset value and cannot be specified in configuration files.

The following attribute is not currently supported: `snapshot_without_reboot`.

## Timeouts

The `timeouts` block allows you to specify [timeouts] for certain actions:

* `create` - (Default `40 minutes`) Used when creating the image.
* `update` - (Default `40 minutes`) Used when updating the image.
* `delete` - (Default `90 minutes`) Used when deregistering the image.

## Import

Importing of the AMI from instance is not currently supported.
