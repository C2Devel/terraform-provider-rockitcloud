---
subcategory: "PaaS"
page_title: "aws_paas_pgsql_user"
description: |-
  Manages a user of a PaaS PostgreSQL service.
---

# Resource: aws_paas_pgsql_user

Manages a user of a PaaS PostgreSQL service as a standalone resource.

~> **Important** Set `users_managed_externally = true` in the parent `aws_paas_service` and remove the `pgsql.user` blocks.
Otherwise, the service manages the full list of users and deletes users created by this resource on the next apply.
See [PostgreSQL users and databases as standalone resources](paas_service.md#postgresql-users-and-databases-as-standalone-resources).
For a service managed by [`aws_paas_pgsql`](paas_pgsql.md), no additional settings are needed.

-> **Note** Each operation switches the service to the `UPDATING` status until the change is applied.
Operations on users and databases of the same service are performed one at a time, so they are processed sequentially
even if Terraform starts them in parallel. See [Timeouts](#timeouts) for the number of users per apply.

-> **Note** A user that owns a database cannot be deleted: the cloud returns the `CannotDeleteUser` error.
Manage databases with [`aws_paas_pgsql_database`](paas_pgsql_database.md) and refer to the owner with
`owner = aws_paas_pgsql_user.<name>.name`, so that Terraform deletes the databases before their owners.

## Example Usage

```terraform
resource "aws_paas_service" "pgsql" {
  # ...
  users_managed_externally     = true
  databases_managed_externally = true

  pgsql {
    version = "16"
  }
}

resource "aws_paas_pgsql_user" "app_reader" {
  service_id = aws_paas_service.pgsql.id
  name       = "app_reader"
  password   = var.app_reader_password
  roles      = ["pg_monitor"]
}
```

To grant the user access to a database, add it to the `users` argument of [`aws_paas_pgsql_database`](paas_pgsql_database.md).

## Argument Reference

The following arguments are supported:

* `name` - (Required, Forces new resource) The PostgreSQL user name.
* `password` - (Required, Editable) The PostgreSQL user password. The value must not contain `'`, `"`, `` ` `` and `\`.
* `roles` - (Optional, Editable) Set of predefined PostgreSQL roles assigned to the user.
  Requires PaaS environment `paas_v4_3` or later. Roles apply across the service, rather than to an individual database.
  Remove the argument or set it to `[]` to revoke all previously assigned roles.
    * _Valid values:_ `pg_monitor`, `pg_read_all_settings`, `pg_read_all_stats`, `pg_stat_scan_tables`, `pg_signal_backend`,
      `pg_create_subscription` (PostgreSQL 16 or later), `pg_maintain` (PostgreSQL 17 or later).
* `service_id` - (Required, Forces new resource) The ID of the PaaS PostgreSQL service.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the resource in the format `<service_id>/<user_id>`.
* `user_id` - The ID of the user in the service.

~> **Note** The cloud API returns the user password in plain text, so changes made outside of Terraform are detected,
and the password is restored on import. The password is stored in the Terraform state.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/blocks/resources/syntax.html#operation-timeouts) for certain actions:

* `create` - (Default `30 minutes`) How long to wait for the service to be updated after the user is created.
* `update` - (Default `30 minutes`) How long to wait for the service to be updated after the user is modified.
* `delete` - (Default `30 minutes`) How long to wait for the service to be updated after the user is deleted.

~> **Note** The timeout of each user includes the time spent waiting for operations on other users and databases
of the same service. One operation takes about 45 seconds, so the default timeouts allow creating, modifying or deleting
approximately 35–40 users and databases of one service in a single apply.
To manage more users, apply the changes in smaller batches (for example, with `-target` or by adding users gradually),
or increase the timeouts for your environment:

```terraform
resource "aws_paas_pgsql_user" "example" {
  # ...

  timeouts {
    create = "60m"
    update = "60m"
    delete = "60m"
  }
}
```

## Import

PaaS PostgreSQL user can be imported using `id`, e.g.,

```
$ terraform import aws_paas_pgsql_user.example fm-cluster-12345678/paas-user-12345678
```

To move users created with the `pgsql.user` blocks of `aws_paas_service` to this resource without recreation, see
[Moving existing users and databases](paas_service.md#moving-existing-users-and-databases).
