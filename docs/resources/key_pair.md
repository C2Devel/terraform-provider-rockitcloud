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

## Example usage

```terraform
resource "aws_key_pair" "deployer" {
  key_name   = "deployer-key"
  public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQD3F6tyPEFEzV0LX3X8BsXdMsQz1x2cEikKDEY0aIj41qgxMCP/iteneqXSIFZBp5vizPvaoIR3Um9xK7PGoW8giupGn+EPuxIA4cDM4vzOqOkiMPhz5XK0whEjkVzTo4+S0puvDZuwIsdiW9mxhJc7tgBNL0cYlWSYVkz4G/fslNfRPW5mYAM49f4fhtxPb5ok4Q2Lg9dPKVHO/Bgeu5woMc7RY0p1ej6D4CKFE6lymSDJpW0YHX/wqE9+cfEauh7xZcG0q9t2ta6F6fmX0agvpFyZo8aFbXeUBr7osSCJNgvavWbM/06niWrOvYX2xwWdhXmXSrbX8ZbabVohBK41 email@example.com"
}
```

## Argument reference

The following arguments are supported:

* `public_key` - (Required, Forces new resource, String) The public key material.
* `key_name` - (Optional, Forces new resource, String) The name for the key pair.
    * _Constraints:_ If neither `key_name` nor `key_name_prefix` is provided, Terraform will create a unique key name using the prefix `terraform-`
* `key_name_prefix` - (Optional, Forces new resource, String) Creates a unique name beginning with the specified prefix.
    * _Constraints:_ Conflicts with `key_name`
    * _Constraints:_ If neither `key_name` nor `key_name_prefix` is provided, Terraform will create a unique key name using the prefix `terraform-`
* `tags` - (Optional, Editable, Map of strings) Key-value pairs to assign to the resource. If the [`default_tags` configuration block][default-tags] is used within a provider configuration, the tags with matching keys will overwrite those defined at the provider level.

## Attribute reference

In addition to all arguments above, the following attributes are exported:

* `arn` - (String) The Amazon Resource Name (ARN) of the key pair.
* `id` - (String) The ID of the key pair.
* `key_name` - (String) The name of the key pair.
* `key_pair_id` - (String) The ID of the key pair.
* `fingerprint` - (String) The MD5 public key fingerprint as specified in section 4 of RFC 4716.
* `tags_all` - (Map of strings) Key-value pairs assigned to the resource, including any tags inherited from the [`default_tags` configuration block][default-tags] if used within a provider configuration.

## Timeouts

Timeouts usage for key pair is not currently supported.

## Import

Key pairs can be imported using `key_name`, for example:

```
$ terraform import aws_key_pair.deployer deployer-key
```
