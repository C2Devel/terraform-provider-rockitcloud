---
subcategory: "EC2 (Elastic Compute Cloud)"
layout: "aws"
page_title: "aws_key_pair"
description: |-
  Manages a key pair.
---

[default-tags]: https://registry.terraform.io/providers/hashicorp/aws/latest/docs#default_tags-configuration-block

# Resource: aws_key_pair

Manages a key pair.
Currently, this resource requires an existing user-supplied key pair.
This key pair's public key will be registered to allow logging in to EC2 instances.

When importing an existing key pair, the public key material may be in any format supported by AWS.
Supported public key material formats are:

* OpenSSH public key format (the format in ~/.ssh/authorized_keys)
* Base64 encoded DER format
* SSH public key file format as specified in RFC4716

## Example Usage

```terraform
resource "aws_key_pair" "deployer" {
  key_name   = "deployer-key"
  public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQD3F6tyPEFEzV0LX3X8BsXdMsQz1x2cEikKDEY0aIj41qgxMCP/iteneqXSIFZBp5vizPvaoIR3Um9xK7PGoW8giupGn+EPuxIA4cDM4vzOqOkiMPhz5XK0whEjkVzTo4+S0puvDZuwIsdiW9mxhJc7tgBNL0cYlWSYVkz4G/fslNfRPW5mYAM49f4fhtxPb5ok4Q2Lg9dPKVHO/Bgeu5woMc7RY0p1ej6D4CKFE6lymSDJpW0YHX/wqE9+cfEauh7xZcG0q9t2ta6F6fmX0agvpFyZo8aFbXeUBr7osSCJNgvavWbM/06niWrOvYX2xwWdhXmXSrbX8ZbabVohBK41 email@example.com"
}
```

## Argument Reference

The following arguments are supported:

* `public_key` - (Required) The public key material.
* `key_name` - (Optional) The name for the key pair.
    _Constraints:_ If neither `key_name` nor `key_name_prefix` is provided, Terraform will create a unique key name using the prefix `terraform-`
* `key_name_prefix` - (Optional) Creates a unique name beginning with the specified prefix.
    _Constraints:_ Conflicts with `key_name`.
    If neither `key_name` nor `key_name_prefix` is provided, Terraform will create a unique key name using the prefix `terraform-`
* `tags` - (Optional) Map of tags to assign to the key pair. If a provider [`default_tags` configuration block][default-tags] is used, tags with matching keys will overwrite those defined at the provider level.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `arn` - The Amazon Resource Name (ARN) of the key pair.
* `id` - The ID of the key pair.
* `key_name` - The name of the key pair.
* `key_pair_id` - The ID of the key pair.
* `fingerprint` - The MD5 public key fingerprint as specified in section 4 of RFC 4716.
* `tags_all` - Map of tags assigned to the key pair, including those inherited from the provider [`default_tags` configuration block][default-tags].

## Import

Key pairs can be imported using `key_name`, e.g.,

```
$ terraform import aws_key_pair.deployer deployer-key
```
