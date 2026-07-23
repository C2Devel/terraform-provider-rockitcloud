---
subcategory: "VPC (Virtual Private Cloud)"
layout: "aws"
page_title: "aws_default_security_group"
description: |-
  Manages the default security group of a VPC.
---

[attribute-as-blocks]: https://www.terraform.io/docs/configuration/attr-as-blocks.html
[default-security-groups]: https://docs.k2.cloud/en/services/security/securitygroups.html#id3
[default-tags]: https://www.terraform.io/docs/providers/aws/index.html#default_tags-configuration-block

# Resource: aws_default_security_group

Manages the default security group of a VPC.
This resource can manage the default security group of the default or a non-default VPC.

~> **Note** This is an advanced resource with special caveats.
Please read this document in its entirety before using this resource.
The `aws_default_security_group` resource behaves differently from normal resources.
Terraform does not _create_ this resource but instead attempts to "adopt" it into management.

When Terraform first adopts the default security group, it **immediately removes all ingress and egress rules in the security group**.
It then creates any rules specified in the configuration.
This way only the rules specified in the configuration are created.

This resource treats its inline rules as absolute; only the rules defined inline are created, and any additions/removals external to this resource will result in diff shown.
For these reasons, this resource is incompatible with the [`aws_security_group_rule`](security_group_rule.md) resource.

For more information about default security groups, see the documentation on [default security groups][default-security-groups].
To manage normal security groups, see the [`aws_security_group`](security_group.md) resource.

## Example usage

### Basic example

The following config gives the default security group the same rules that the cloud provides by default but under Terraform management.
This means that any ingress or egress rules added or changed will be detected as drift.

```terraform
resource "aws_vpc" "mainvpc" {
  cidr_block = "10.1.0.0/16"
}

resource "aws_default_security_group" "example" {
  vpc_id = aws_vpc.mainvpc.id

  ingress {
    protocol  = -1
    self      = true
    from_port = 0
    to_port   = 0
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}
```

### Specific example: deny all egress traffic and allow ingress traffic

The following example denies all egress traffic by omitting any `egress` rules, while including the default `ingress` rule to allow all traffic.

```terraform
resource "aws_vpc" "mainvpc" {
  cidr_block = "10.1.0.0/16"
}

resource "aws_default_security_group" "example" {
  vpc_id = aws_vpc.mainvpc.id

  ingress {
    protocol  = -1
    self      = true
    from_port = 0
    to_port   = 0
  }
}
```

### Removing `aws_default_security_group` from your configuration

Removing this resource from your configuration will remove it from your statefile and management, but will not destroy the security group.
All ingress or egress rules will be left as they are at the time of removal.
You can resume managing them via the cloud console.

## Argument reference

The following arguments are supported:

* `egress` - (Optional, Editable, [Block](#egress-and-ingress)) One or more egress rules (for outgoing traffic).
* `ingress` - (Optional, Editable, [Block](#egress-and-ingress)) One or more ingress rules (for incoming traffic).
* `tags` - (Optional, Editable, Map of strings) Key-value pairs to assign to the resource.
    If the [`default_tags` configuration block][default-tags] block is used within a provider configuration, the tags with matching keys will overwrite those defined at the provider level.
* `vpc_id` - (Optional, Forces new resource, String) The ID of the VPC.

~> **Note** Changing the `vpc_id` argument value will _not_ restore any default security group rules that were modified, added, or removed.
            They will be left in its current state.

### egress and ingress

Both arguments are processed in [attribute-as-blocks mode][attribute-as-blocks].

Both `egress` and `ingress` objects have the same arguments.

The following arguments are required:

* `from_port` - (Required, Editable, Integer) The start of the port range (or ICMP message type number if the `protocol` value is `icmp`).
* `protocol` - (Required, Editable, String) The protocol to match.
    * _Constraints:_
        * If using the `-1` value (semantically equivalent to `all`, which is not a valid value here), you must specify the `from_port` and `to_port` arguments values equal to `0`
        * If the `protocol` value is not `icmp`, `tcp`, `udp`, or `-1`, then refer to the [protocol number](https://www.iana.org/assignments/protocol-numbers/protocol-numbers.xhtml) for detailed information
* `to_port` - (Required, Editable, Integer) The end of the port range (or ICMP message code if the `protocol` value is `icmp`).

The following arguments are optional:

* `cidr_blocks` - (Optional, Editable, List of strings) The list of CIDR blocks.
* `description` - (Optional, Editable, String) The description of this rule.
* `ipv6_cidr_blocks` - (Optional, Editable, List of strings) The list of IPv6 CIDR blocks.
* `security_groups` - (Optional, Editable, List of strings) The list of security group names or IDs.
* `self` - (Optional, Editable, Boolean) Indicates whether the security group itself will be added as a source to this egress rule.
    * _Default value:_ `false`

## Attribute reference

### Supported attributes

In addition to all arguments above, the following attributes are exported:

* `arn` - (String) The Amazon Resource Name (ARN) of the security group.
* `description` - (String) The description of the security group.
* `id` - (String)  The ID of the security group.
* `name` - (String) The name of the security group.
* `owner_id` - (String) The ID of the project that owns the security group.
* `tags_all` - (Map of strings) Key-value pairs assigned to the resource, including any tags inherited from the [`default_tags` configuration block][default-tags] if used within a provider configuration.

### Unsupported attributes

~> **Note** This attribute may be present in the `terraform.tfstate` file, but it has a preset value and cannot be specified in configuration files.

The following attributes are not currently supported:

`egress.prefix_list_ids`, `ingress.prefix_list_ids`, `revoke_rules_on_delete`.

## Timeouts

Timeouts usage for the default security groups is not currently supported.

## Import

Security groups can be imported using the `id`, for example:

```
$ terraform import aws_default_security_group.default_sg sg-12345678
```
