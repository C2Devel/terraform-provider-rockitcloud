---
subcategory: "PaaS"
page_title: "aws_paas_pgsql_database"
description: |-
  Manages a database of a PaaS PostgreSQL service.
---

# Resource: aws_paas_pgsql_database

Manages a database of a PaaS PostgreSQL service as a standalone resource.

~> **Important** Set `databases_managed_externally = true` in the parent `aws_paas_service` and remove the `pgsql.database` blocks.
Otherwise, the service manages the full list of databases and deletes databases created by this resource on the next apply.
See [PostgreSQL users and databases as standalone resources](paas_service.md#postgresql-users-and-databases-as-standalone-resources).

-> **Note** Each operation switches the service to the `UPDATING` status until the change is applied.
Operations on users and databases of the same service are performed one at a time.
See [Timeouts](#timeouts) for the number of operations per apply.

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

resource "aws_paas_pgsql_user" "app_owner" {
  service_id = aws_paas_service.pgsql.id
  name       = "app_owner"
  password   = var.app_owner_password
}

resource "aws_paas_pgsql_user" "app_reader" {
  service_id = aws_paas_service.pgsql.id
  name       = "app_reader"
  password   = var.app_reader_password
}

resource "aws_paas_pgsql_database" "app" {
  service_id = aws_paas_service.pgsql.id
  name       = "app"
  owner      = aws_paas_pgsql_user.app_owner.name
  encoding   = "UTF8"
  locale     = "ru_RU.UTF-8"
  extensions = ["pgcrypto", "pg_trgm"]
  users      = [aws_paas_pgsql_user.app_reader.name]
}
```

Refer to the users with resource attributes, as in the example above, so that Terraform creates the users
before the database and deletes the database before the users.

## Argument Reference

The following arguments are required:

* `name` - (Required, Forces new resource) The database name.
* `owner` - (Required, Editable) The name of the user who is the database owner. The user must exist in the service.
  The owner cannot be deleted as long as it owns the database.
* `service_id` - (Required, Forces new resource) The ID of the PaaS PostgreSQL service.

The following arguments are optional:

* `backup_enabled` - (Optional, Editable) Indicates whether backup is enabled for the database.
  Backup must be configured with the `backup_settings` block of the parent `aws_paas_service`.
    * _Default value:_ `false`
* `backup_id` - (Optional, Forces new resource) The ID of the backup to create the database from.
* `backup_db_name` - (Optional, Forces new resource) The name of a database from the backup specified in `backup_id`.
* `encoding` - (Optional, Forces new resource) The database encoding. The value is compared case-insensitively.
  If the argument is not specified, the encoding chosen by the cloud is used.
* `extensions` - (Optional, Editable) Set of extensions for the database.
  Remove the argument or set it to `[]` to remove all previously installed extensions.
    * _Valid values:_ `address_standardizer`, `address_standardizer_data_us`, `amcheck`, `autoinc`, `bloom`, `btree_gin`, `btree_gist`, `citext`,`cube`, `dblink`, `dict_int`, `dict_xsyn`, `earthdistance`, `fuzzystrmatch`, `hstore`, `intarray`, `isn`, `lo`, `ltree`, `moddatetime`, `pg_buffercache`, `pg_trgm`, `pg_visibility `, `pgcrypto`, `pgrowlocks`, `pgstattuple`, `postgis`, `postgis_tiger_geocoder`, `postgis_topology`, `postgres_fdw`, `seg`, `tablefunc`, `tcn`, `timescaledb`,  `tsm_system_rows`, `tsm_system_time`, `unaccent`, `uuid-ossp`, `xml2`
* `locale` - (Optional, Forces new resource) The database locale.
  If the argument is not specified, the locale chosen by the cloud is used.
* `users` - (Optional, Editable) Set of names of the users that have access to the database. The maximum number of users is 1000.
  Remove the argument or set it to `[]` to revoke access of all users.
    * _Constraints:_ Must not contain the `owner`, because the owner already has access to the database.
      When the owner is changed to one of the `users`, remove this user from `users` in the same change

~> **Note** Changing `encoding`, `locale`, `backup_id` or `backup_db_name` deletes the database with all its data
and creates a new one.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `database_id` - The ID of the database in the service.
* `id` - The ID of the resource in the format `<service_id>/<database_id>`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/blocks/resources/syntax.html#operation-timeouts) for certain actions:

* `create` - (Default `30 minutes`) How long to wait for the service to be updated after the database is created.
* `update` - (Default `30 minutes`) How long to wait for the service to be updated after the database is modified.
* `delete` - (Default `30 minutes`) How long to wait for the service to be updated after the database is deleted.

~> **Note** The timeout of each database includes the time spent waiting for operations on other users and databases
of the same service. One operation takes about 45 seconds, so the default timeouts allow approximately 35–40 operations
on users and databases of one service in a single apply.
To manage more objects, apply the changes in smaller batches (for example, with `-target` or by adding databases gradually),
or increase the timeouts for your environment.

## Import

PaaS PostgreSQL database can be imported using `id`, e.g.,

```
$ terraform import aws_paas_pgsql_database.example fm-cluster-12345678/paas-db-12345678
```

To move databases created with the `pgsql.database` blocks of `aws_paas_service` to this resource without recreation, see
[Moving existing users and databases](paas_service.md#moving-existing-users-and-databases).
