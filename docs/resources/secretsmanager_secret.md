---
subcategory: "Secrets Manager"
layout: "aws"
page_title: "aws_secretsmanager_secret"
description: |-
  Manages a Secrets Manager secret.
---

# Resource: aws_secretsmanager_secret

Manages a Secrets Manager secret.

## Example Usage

### Basic

```terraform
resource "aws_secretsmanager_secret" "example" {
  name = "example"
}
```

### With a value

```terraform
resource "aws_secretsmanager_secret" "example" {
  name = "example"
}

resource "aws_secretsmanager_secret_version" "example" {
  secret_id     = aws_secretsmanager_secret.example.id
  secret_string = "example-value"
}
```

### With a generated name

```terraform
resource "aws_secretsmanager_secret" "example" {
  name_prefix = "example-"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional, ForceNew) The name of the secret. Conflicts with `name_prefix`.
* `name_prefix` - (Optional, ForceNew) Creates a unique name beginning with the specified prefix. Conflicts with `name`.
* `description` - (Optional) The description of the secret.
* `kms_key_id` - (Optional) The ARN or alias of the KMS key used to encrypt the secret value. If not specified, the system key is used.
* `recovery_window_in_days` - (Optional) The number of days that Secrets Manager waits before deleting the secret. Can be `0` for immediate deletion, or `7` to `30`. Defaults to `30`.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique identifier (UUID) of the secret.
* `secret_id` - The unique identifier (UUID) of the secret.

## Import

Secrets Manager secret can be imported using `secret_id` (UUID) or `name`, e.g.,

```
$ terraform import aws_secretsmanager_secret.example example
