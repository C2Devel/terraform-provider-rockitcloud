---
subcategory: "Secrets Manager"
layout: "aws"
page_title: "aws_secretsmanager_secret_version"
description: |-
  Manages a Secrets Manager secret version.
---

# Resource: aws_secretsmanager_secret_version

Manages a Secrets Manager secret version.

## Example Usage

### String value

```terraform
resource "aws_secretsmanager_secret_version" "example" {
  secret_id     = aws_secretsmanager_secret.example.id
  secret_string = "example-value"
}
```

### Binary value

```terraform
resource "aws_secretsmanager_secret_version" "example" {
  secret_id     = aws_secretsmanager_secret.example.id
  secret_binary = base64encode("example-binary-value")
}
```

## Argument Reference

The following arguments are supported:

* `secret_id` - (Required, ForceNew) The ID (UUID) or name of the secret to which the version belongs.
* `secret_string` - (Optional, ForceNew, Sensitive) The string value of the secret. Conflicts with `secret_binary`.
* `secret_binary` - (Optional, ForceNew, Sensitive) The binary value of the secret, base64-encoded. Conflicts with `secret_string`.
* `version_stages` - (Optional) A list of staging labels to attach to the version. If not specified, the `CLOUDCURRENT` label is assigned.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The identifier of the version, in the format `secret_id|version_id`.
* `version_id` - The unique identifier (UUID) of the version.

## Import

Secrets Manager secret version can be imported using the `secret_id|version_id` identifier, e.g.,

```
$ terraform import aws_secretsmanager_secret_version.example secret-id|version-id
