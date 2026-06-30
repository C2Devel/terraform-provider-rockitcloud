---
subcategory: "VPC (Virtual Private Cloud)"
layout: "aws"
page_title: "aws_ec2_traffic_mirror_filter_rule"
description: |-
  Manages a traffic mirror filter rule.
---

[protocol-numbers]: https://www.iana.org/assignments/protocol-numbers/protocol-numbers.xhtml
[traffic-mirroring]: https://docs.k2.cloud/en/services/interconnect/traffic_mirroring.html

# Resource: aws_ec2_traffic_mirror_filter_rule

Manages a traffic mirror filter rule. For details about traffic mirroring, see the [user documentation][traffic-mirroring].

## Example usage

To create a basic traffic mirror filter rule, use:

```terraform
resource "aws_ec2_traffic_mirror_filter" "filter" {
  description = "traffic mirror filter - terraform example"
}

resource "aws_ec2_traffic_mirror_filter_rule" "ruleout" {
  description              = "test rule"
  traffic_mirror_filter_id = aws_ec2_traffic_mirror_filter.filter.id
  destination_cidr_block   = "10.0.0.0/8"
  source_cidr_block        = "10.0.0.0/8"
  rule_number              = 1
  rule_action              = "accept"
  traffic_direction        = "egress"
}

resource "aws_ec2_traffic_mirror_filter_rule" "rulein" {
  description              = "test rule"
  traffic_mirror_filter_id = aws_ec2_traffic_mirror_filter.filter.id
  destination_cidr_block   = "10.0.0.0/8"
  source_cidr_block        = "10.0.0.0/8"
  rule_number              = 1
  rule_action              = "accept"
  traffic_direction        = "ingress"
  protocol                 = 6

  destination_port_range {
    from_port = 22
    to_port   = 53
  }

  source_port_range {
    from_port = 0
    to_port   = 10
  }
}
```

## Argument reference

The following arguments are supported:

* `destination_cidr_block` - (Required, Editable, String) The destination CIDR block to assign to the traffic mirror rule.
* `rule_action` - (Required, Editable, String) The action to take on the filtered traffic.
    * _Valid values:_ `accept`, `reject`
* `rule_number` - (Required, Editable, Integer) The number of the traffic mirror rule. This number must be unique for each traffic mirror rule in a given direction. The rules are processed in ascending order by rule number.
    * _Valid values:_ From 1 to 128
* `source_cidr_block` - (Required, Editable, String) The source CIDR block to assign to the traffic mirror rule.
* `traffic_direction` - (Required, Editable, String) The direction of traffic to be captured.
    * _Valid values:_ `egress`, `ingress`
* `traffic_mirror_filter_id` - (Required, Forces new resource, String) The ID of the traffic mirror filter to which this rule should be added.
* `description` - (Optional, Editable, String) The description of the traffic mirror filter rule.
* `destination_port_range` - (Optional, Editable, [Block](#destination_port_range)) The destination port range.
    * _Constraints:_ Supported only when the `protocol` is set to TCP(6) or UDP(17).
* `protocol` - (Optional, Editable, Integer) The protocol number to assign to the traffic mirror rule. For information about the protocol value, see [Protocol Numbers][protocol-numbers] on the Internet Assigned Numbers Authority (IANA) website.
    * _Example:_ `17` (UDP)
* `source_port_range` - (Optional, Editable, [Block](#source_port_range)) The source port range.
    * _Constraints:_ Supported only when the `protocol` is set to TCP(6) or UDP(17).

### destination_port_range

The `destination_port_range` block has the following structure:

* `from_port` - (Optional, Editable, Integer) The starting port of the range.
    * _Valid values:_ From 0 to 65535
* `to_port` - (Optional, Editable, Integer) The ending port of the range.
    * _Valid values:_ From 0 to 65535

### source_port_range

The `source_port_range` block has the following structure:

* `from_port` - (Optional, Editable, Integer) The starting port of the range.
    * _Valid values:_ From 0 to 65535
* `to_port` - (Optional, Editable, Integer) The ending port of the range.
    * _Valid values:_ From 0 to 65535

## Attribute reference

In addition to all arguments above, the following attributes are exported:

* `arn` - (String) The Amazon Resource Name (ARN) of the traffic mirror filter rule.
* `id` - (String) The ID of the traffic mirror filter rule.

## Import

In Terraform v1.5.0 or later, traffic mirror filter rule can be imported by `id` using the `import` block.

```terraform
import {
  to = aws_ec2_traffic_mirror_filter_rule.rule
  id = "tmf-12345678:tmfr-12345678"
}
```

In older Terraform versions, the traffic mirror filter rule can be imported by `traffic_mirror_filter_id` and its `id` separated by `:` using `terraform import`, e.g.:

```console
% terraform import aws_ec2_traffic_mirror_filter_rule.rule tmf-12345678:tmfr-12345678
```
