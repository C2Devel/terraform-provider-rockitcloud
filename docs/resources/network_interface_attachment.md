---
subcategory: "VPC (Virtual Private Cloud)"
layout: "aws"
page_title: "aws_network_interface_attachment"
description: |-
  Attaches an elastic network interface (ENI) with an EC2 instance.
---

# Resource: aws_network_interface_attachment

Attaches an elastic network interface (ENI) with an EC2 instance.

## Example usage

### Basic example

```terraform
resource "aws_network_interface_attachment" "test" {
  instance_id          = "i-12345678"
  network_interface_id = "eni-12345678"
  device_index         = 1
}
```

## Argument reference

The following arguments are supported:

* `instance_id` - (Required, Forces new resource, String) The ID of the instance to attach.
* `network_interface_id` - (Required, Forces new resource, String) The ID of the network interface to attach.
* `device_index` - (Required, Forces new resource, Integer) The network interface index.

## Attribute reference

In addition to all arguments above, the following attributes are exported:

* `attachment_id` - (String) The ID of the network interface attachment.
* `id` - (String) The ID of the network interface attachment.
* `status` - (String) The status of the network interface attachment.

## Timeouts

Timeouts usage for network interface attachments is not currently supported.

## Import

Import for network interface attachments is not currently supported.
