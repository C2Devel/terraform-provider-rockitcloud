---
subcategory: "VPC (Virtual Private Cloud)"
layout: "aws"
page_title: "aws_ec2_traffic_mirror_target"
description: |-
  Manages a traffic mirror target.
---

[default-tags]: https://registry.terraform.io/providers/hashicorp/aws/latest/docs#default_tags-configuration-block
[traffic-mirroring]: https://docs.k2.cloud/en/services/interconnect/traffic_mirroring.html

# Resource: aws_ec2_traffic_mirror_target

Manages a traffic mirror target. For details about traffic mirroring, see the [user documentation][traffic-mirroring].

## Example usage

To create a basic traffic mirror target, use:

```terraform
variable ami {}
variable instance_type {}

data "aws_availability_zones" "azs" {
  state = "available"
}

resource "aws_vpc" "vpc" {
  cidr_block = "10.0.0.0/16"
}

resource "aws_subnet" "sub1" {
  vpc_id            = aws_vpc.vpc.id
  cidr_block        = "10.0.0.0/24"
  availability_zone = data.aws_availability_zones.azs.names[0]
}

resource "aws_instance" "dst" {
  ami           = var.ami
  instance_type = var.instance_type
  subnet_id     = aws_subnet.sub1.id
}

resource "aws_ec2_traffic_mirror_target" "eni" {
  description          = "ENI target"
  network_interface_id = aws_instance.dst.primary_network_interface_id
}
```

## Argument reference

The following arguments are required:

* `network_interface_id` - (Required, Forces new resource, String) The ID of the network interface that is associated with the target.

The following arguments are optional:

* `description` - (Optional, Forces new resource, String) The description of the traffic mirror target.
* `tags` - (Optional, Editable, Map of strings) Key-value pairs to assign to the traffic mirror target.
  If the [`default_tags` configuration block][default-tags] is used within a provider configuration, the tags with matching keys will overwrite those defined at the provider level.

## Attribute reference

### Supported attributes

In addition to all arguments above, the following attributes are exported:

* `arn` - (String) The Amazon Resource Name (ARN) of the traffic mirror target.
* `id` - (String) The ID of the traffic mirror target.
* `owner_id` - (String) The ID of the project the traffic mirror target belongs to.
* `tags_all` - (Map of strings) Key-value pairs assigned to the traffic mirror target, including any tags inherited from the [`default_tags` configuration block][default-tags] if used within a provider configuration.

### Unsupported attributes

~> **Note** This attribute may be present in the `terraform.tfstate` file, but it has a preset value and cannot be specified in configuration files.

The following attribute is not currently supported:

`network_load_balancer_arn`.

## Import

In Terraform v1.5.0 or later, traffic mirror target can be imported by `id` using the `import` block.

```terraform
import {
  to = aws_ec2_traffic_mirror_target.target
  id = "tmt-12345678"
}
```

In older Terraform versions, the traffic mirror target can be imported by its `id` using `terraform import`, for example:

```console
% terraform import aws_ec2_traffic_mirror_target.target tmt-12345678
```
