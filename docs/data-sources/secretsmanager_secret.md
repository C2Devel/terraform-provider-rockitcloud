---
subcategory: "Secrets Manager"
layout: "aws"
page_title: "aws_secretsmanager_secret"
description: |-
  Provides information about a Secrets Manager secret.
---

# Data Source: aws_secretsmanager_secret

Provides information about a Secrets Manager secret.

## Example Usage

```terraform
data "aws_secretsmanager_secret" "example" {
  name = "example"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the secret.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique identifier (UUID) of the secret.
* `secret_id` - The unique identifier (UUID) of the secret.
* `description` - The description of the secret.
* `kms_key_id` - The ARN or alias of the KMS key used to encrypt the secret value.
