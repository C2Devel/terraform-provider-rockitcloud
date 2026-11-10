---
subcategory: "Secrets Manager"
layout: "aws"
page_title: "aws_secretsmanager_secret_version"
description: |-
  Provides information about a Secrets Manager secret version.
---

# Data Source: aws_secretsmanager_secret_version

Provides information about a Secrets Manager secret version.

## Example Usage

### Current version

```terraform
data "aws_secretsmanager_secret_version" "example" {
  secret_id = aws_secretsmanager_secret.example.id
}
```

### Specific version

```terraform
data "aws_secretsmanager_secret_version" "example" {
  secret_id  = aws_secretsmanager_secret.example.id
  version_id = "version-id"
}
```

## Argument Reference

The following arguments are supported:

* `secret_id` - (Required) The ID (UUID) or name of the secret.
* `version_id` - (Optional) The unique identifier (UUID) of the version to retrieve. If not specified, the version with the `version_stage` label is retrieved.
* `version_stage` - (Optional) The staging label of the version to retrieve. Defaults to `CLOUDCURRENT`.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The identifier of the version, in the format `secret_id|version`.
* `secret_string` - The string value of the secret.
* `secret_binary` - The binary value of the secret.
* `version_stages` - A list of staging labels attached to the version.
