---
subcategory: "VPC (Virtual Private Cloud)"
layout: "aws"
page_title: "aws_security_group"
description: |-
  Manages a security group.
---

[attribute-as-blocks]: https://www.terraform.io/docs/configuration/attr-as-blocks.html
[default-tags]: https://www.terraform.io/docs/providers/aws/index.html#default_tags-configuration-block
[protocol-number]: https://www.iana.org/assignments/protocol-numbers/protocol-numbers.xhtml
[create_before_destroy]: https://developer.hashicorp.com/terraform/language/meta-arguments/lifecycle#create_before_destroy
[security-groups]: https://docs.k2.cloud/en/services/security/securitygroups.html
[timeouts]: https://developer.hashicorp.com/terraform/plugin/framework/resources/timeouts

# Resource: aws_security_group

Manages a security group.

~> **Note** Terraform currently provides both a standalone resource [`aws_security_group_rule`](security_group_rule.md) (a single `ingress` or `egress` rule), and a security group resource with `ingress` and `egress` rules defined inline.
            At this time you cannot use a security group with inline rules in conjunction with any security group rule resources.
            Doing so will cause a conflict of rule settings and will overwrite rules.

For more information about security groups, see the documentation on [Security groups][security-groups].

## Example usage

### Specific example

```terraform
resource "aws_vpc" "main" {
  cidr_block = "10.1.0.0/16"
}

resource "aws_security_group" "allow_tls" {
  name        = "allow_tls"
  description = "Allow TLS inbound traffic"
  vpc_id      = aws_vpc.main.id

  ingress {
    description = "TLS from VPC"
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = [aws_vpc.main.cidr_block]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = "allow_tls"
  }
}
```

~> **Note on egress rules** By default, the cloud creates an `ALLOW ALL` egress rule when creating a new security group inside a VPC.
When creating a new security group inside a VPC, **Terraform will remove this default rule**, and require you specifically re-create it if you desire that rule.
We feel this leads to fewer surprises in terms of controlling your egress rules.
If you desire this rule to be in place, you can use this `egress` block:

```terraform
resource "aws_security_group" "example" {
  # ... other configuration ...

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}
```

### Change of name or name-prefix value

Security group's name cannot be edited after the resource is created.
In fact, the `name` and `name-prefix` arguments force the creation of a new security group resource when they change value.
In that case, Terraform first deletes the existing security group resource and then it creates a new one.
If the existing security group is associated to a network interface resource, the deletion cannot complete.
The reason is that network interface resources cannot be left with no security group attached and the new one is not yet available at that point.

It is required to invert the default behavior of Terraform.
That is, first the new security group resource must be created, then associated to possible network interface resources and finally the old security group can be detached and deleted.
To force this behavior, you must set the [create_before_destroy] property:

```terraform
resource "aws_security_group" "sg_with_changeable_name" {
  name = "changeable-name"
  # ... other configuration ...

  lifecycle {
    # Necessary if changing 'name' or 'name_prefix' properties.
    create_before_destroy = true
  }
}
```

## Argument reference

The following arguments are supported:

* `description` - (Optional, Editable, String) The description of the security group.
    * _Default value:_ `Managed by Terraform`
* `egress` - (Optional, Editable, [Block](#egress)) One or more egress rules (for outgoing traffic).
* `ingress` - (Optional, Editable, [Block](#ingress)) One or more ingress rules (for incoming traffic).
* `name` - (Optional, Forces new resource, String) The name of the security group. If omitted, Terraform will assign a random unique name.
    * _Constraints:_ Conflicts with `name`.
* `name_prefix` - (Optional, Forces new resource, String) This argument allows to create a unique name beginning with the specified prefix.
    * _Constraints:_ Conflicts with `name`.
* `revoke_rules_on_delete` - (Optional, Editable, Boolean) The argument that instructs Terraform to revoke all the security groups attached ingress and egress rules before deleting the rule itself.
    * _Default value:_ `false`
* `tags` - (Optional, Editable, Map of strings) Key-value pairs to assign to the resource.
    If the [`default_tags` configuration block][default-tags] block is used within a provider configuration, the tags with matching keys will overwrite those defined at the provider level.
* `vpc_id` - (Optional, Forces new resource, String) The ID of the VPC.

~> **Note** The `name` and `name_prefix` arguments cannot be specified within one configuration due to incompatibility.

### egress

This argument is processed in [attribute-as-blocks mode][attribute-as-blocks].

The following arguments are required:

* `from_port` - (Required, Editable, Integer) The start of the port range (or ICMP message type number if the `protocol` value is `icmp`).
* `protocol` - (Required, Editable, String) The protocol to match.
    * _Constraints:_
        * If using the `-1` value (semantically equivalent to `all`, which is not a valid value here), you must specify the `from_port` and `to_port` arguments values equal to `0`
        * If the `protocol` value is not `icmp`, `tcp`, `udp`, or `-1`, then refer to the [protocol number][protocol-number] for detailed information
* `to_port` - (Required, Editable, Integer) The end of the port range (or ICMP message code if the `protocol` value is `icmp`).

The following arguments are optional:

* `cidr_blocks` - (Optional, Editable, List of strings) The list of CIDR blocks.
* `description` - (Optional, Editable, String) The description of this egress rule.
* `ipv6_cidr_blocks` - (Optional, Editable, List of strings) The list of IPv6 CIDR blocks.
* `security_groups` - (Optional, Editable, List of strings) The list of security group IDs.
* `self` - (Optional, Editable, Boolean) Indicates whether the security group itself will be added as a source to this egress rule.
    * _Default value:_ `false`

### ingress

This argument is processed in [attribute-as-blocks mode][attribute-as-blocks].

The following arguments are required:

* `from_port` - (Required, Editable, Integer) The start of the port range (or ICMP message type number if the `protocol` value is `icmp`).
* `protocol` - (Required, Editable, String) The protocol to match.
    * _Constraints:_
        * If using the `-1` value (semantically equivalent to `all`, which is not a valid value here), you must specify the `from_port` and `to_port` arguments values equal to `0`
        * If the `protocol` value is not `icmp`, `tcp`, `udp`, or `-1`, then refer to the [protocol number][protocol-number] for detailed information
* `to_port` - (Required, Editable, Integer) The end of the port range (or ICMP message code if the `protocol` value is `icmp`).

The following arguments are optional:

* `cidr_blocks` - (Optional, Editable, List of strings) The list of CIDR blocks.
* `description` - (Optional, Editable, String) The description of this ingress rule.
* `ipv6_cidr_blocks` - (Optional, Editable, List of strings) The list of IPv6 CIDR blocks.
* `security_groups` - (Optional, Editable, List of strings) The list of security group IDs.
* `self` - (Optional, Editable, Boolean) Indicates whether the security group itself will be added as a source to this ingress rule.
    * _Default value:_ `false`

## Attribute reference

### Supported attributes

In addition to all arguments above, the following attributes are exported:

* `arn` - (String) The Amazon Resource Name (ARN) of the security group.
* `id` - (String) The ID of the security group.
* `owner_id` - (String) The ID of the project that owns this security group.
* `tags_all` - (Map of strings) Key-value pairs assigned to the resource, including any tags inherited from the [`default_tags` configuration block][default-tags] if used within a provider configuration.

### Unsupported attributes

~> **Note** This attribute may be present in the `terraform.tfstate` file, but it has a preset value and cannot be specified in configuration files.

The following attributes are not currently supported:

`egress.prefix_list_ids`, `ingress.prefix_list_ids`.

## Timeouts

The `timeouts` block allows you to specify [timeouts] for certain actions:

- `create` - (Default `10m`) How long to wait for a security group to be created.
- `delete` - (Default `15m`) How long to wait for a security group to be deleted.

## Import

Security groups can be imported using the `id`, for example:

```
$ terraform import aws_security_group.elb_sg sg-12345678
```


