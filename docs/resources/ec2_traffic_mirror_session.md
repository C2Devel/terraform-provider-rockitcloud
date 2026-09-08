---
subcategory: "VPC (Virtual Private Cloud)"
layout: "aws"
page_title: "aws_ec2_traffic_mirror_session"
description: |-
  Manages a traffic mirror session.
---

[default-tags]: https://registry.terraform.io/providers/hashicorp/aws/latest/docs#default_tags-configuration-block
[traffic-mirroring]: https://docs.k2.cloud/en/services/interconnect/traffic_mirroring.html

# Resource: aws_ec2_traffic_mirror_session

Manages a traffic mirror session.
For details about traffic mirroring, see the [user documentation][traffic-mirroring].

## Example usage

To create a basic traffic mirror session, use:

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

resource "aws_instance" "src" {
  ami           = var.ami
  instance_type = var.instance_type
  subnet_id     = aws_subnet.sub1.id
}

resource "aws_ec2_traffic_mirror_filter" "filter" {
  description = "traffic mirror filter - terraform example"
}

resource "aws_ec2_traffic_mirror_target" "target" {
  description          = "ENI target"
  network_interface_id = aws_instance.dst.primary_network_interface_id
}

resource "aws_ec2_traffic_mirror_session" "session" {
  description              = "traffic mirror session - terraform example"
  network_interface_id     = aws_instance.src.primary_network_interface_id
  session_number           = 1
  traffic_mirror_filter_id = aws_ec2_traffic_mirror_filter.filter.id
  traffic_mirror_target_id = aws_ec2_traffic_mirror_target.target.id
}
```

## Argument reference

The following arguments are required:

* `network_interface_id` - (Required, Forces new resource, String) The ID of the source network interface.
* `session_number` - (Required, Forces new resource, Integer) The session number determines the order in which sessions are evaluated when the interface is used by multiple sessions.
  The first session with a matching filter is the one that mirrors the packets.
* `traffic_mirror_filter_id` - (Required, Forces new resource, String) The ID of the traffic mirror filter to be used.
* `traffic_mirror_target_id` - (Required, Forces new resource, String) The ID of the traffic mirror target to be used.

The following arguments are optional:

* `description` - (Optional, Forces new resource, String) The description of the traffic mirror session.
* `tags` - (Optional, Editable, Map of strings) Key-value pairs to assign to the traffic mirror session.
  If the [`default_tags` configuration block][default-tags] is used within a provider configuration, the tags with matching keys will overwrite those defined at the provider level.

## Attribute reference

### Supported attributes

In addition to all arguments above, the following attributes are exported:

* `arn` - (String) The Amazon Resource Name (ARN) of the traffic mirror session.
* `id` - (String) The ID of the traffic mirror session.
* `owner_id` - (String) The ID of the project the traffic mirror session belongs to.
* `tags_all` - (Map of strings) Key-value pairs assigned to the traffic mirror session, including any tags inherited from the [`default_tags` configuration block][default-tags] if used within a provider configuration.

### Unsupported attributes

~> **Note** These attributes may be present in the `terraform.tfstate` file, but they have preset values and cannot be specified in configuration files.

The following attributes are not currently supported:

`packet_length`, `virtual_network_id`.

## Import

In Terraform v1.5.0 or later, traffic mirror session can be imported by `id` using the `import` block.

```terraform
import {
  to = aws_ec2_traffic_mirror_session.session
  id = "tms-12345678"
}
```

In older Terraform versions, the traffic mirror session can be imported by its `id` using `terraform import`, for example:

```console
% terraform import aws_ec2_traffic_mirror_session.session tms-12345678
```
