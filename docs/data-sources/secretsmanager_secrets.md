---
subcategory: "Secrets Manager"
layout: "aws"
page_title: "aws_secretsmanager_secrets"
description: |-
  Provides information about Secrets Manager secrets.
---

# Data Source: aws_secretsmanager_secrets

Provides information about Secrets Manager secrets.

## Example Usage

```terraform
data "aws_secretsmanager_secrets" "example" {}
```

## Argument Reference

The following arguments are supported:

* `filter` - (Optional) Configuration block(s) for filtering. Detailed below.
* `include_planned_deletion` - (Optional) Whether to include secrets scheduled for deletion. Defaults to `false`.
* `sort_by` - (Optional) The field to sort the results by. Valid values are `created-date`, `last-accessed-date`, `last-changed-date`, `name`.

### filter Configuration Block

The `filter` configuration block supports the following arguments:

* `key` - (Required) The filter key. Valid values are `name`, `description`, `all`.
* `values` - (Required) Set of values that are accepted for the given filter key.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The region of the secrets.
* `names` - Set of names of the secrets.
