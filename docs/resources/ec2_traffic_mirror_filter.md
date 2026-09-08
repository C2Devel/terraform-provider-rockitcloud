---
subcategory: "VPC (Virtual Private Cloud)"
layout: "aws"
page_title: "aws_ec2_traffic_mirror_filter"
description: |-
  Manages a traffic mirror filter.
---

[default-tags]: https://registry.terraform.io/providers/hashicorp/aws/latest/docs#default_tags-configuration-block
[traffic-mirroring]: https://docs.k2.cloud/en/services/interconnect/traffic_mirroring.html

# Resource: aws_ec2_traffic_mirror_filter

Manages a traffic mirror filter.
For details about traffic mirroring, see the [user documentation][traffic-mirroring].

## Example usage

To create a basic traffic mirror filter, use:

```terraform
resource "aws_ec2_traffic_mirror_filter" "foo" {
  description = "traffic mirror filter - terraform example"
}
```

## Argument reference

* `description` - (Optional, Forces new resource, String) The description of the filter.
* `tags` - (Optional, Editable, Map of strings) Key-value pairs to assign to the traffic mirror filter.
  If the [`default_tags` configuration block][default-tags] is used within a provider configuration, the tags with matching keys will overwrite those defined at the provider level.

## Attribute reference

### Supported attributes

In addition to all arguments above, the following attributes are exported:

* `arn` - (String) The Amazon Resource Name (ARN) of the traffic mirror filter.
* `id` - (String) The ID of the traffic mirror filter.
* `tags_all` - (Map of strings) Key-value pairs assigned to the traffic mirror filter, including any tags inherited from the [`default_tags` configuration block][default-tags] if used within a provider configuration.

### Unsupported attributes

~> **Note** This attribute may be present in the `terraform.tfstate` file, but it has a preset value and cannot be specified in configuration files.

The following attribute is not currently supported:

`network_services`.

## Timeouts

Timeouts usage for traffic mirror filter is not currently supported.

## Import

In Terraform v1.5.0 or later, traffic mirror filter can be imported by `id` using the `import` block.

```terraform
import {
  to = aws_ec2_traffic_mirror_filter.foo
  id = "tmf-12345678"
}
```

In older Terraform versions, the traffic mirror filter can be imported by its `id` using `terraform import`, for example:

```console
terraform import aws_ec2_traffic_mirror_filter.foo tmf-12345678
```
