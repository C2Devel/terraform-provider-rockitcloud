---
subcategory: "PaaS"
layout: "aws"
page_title: "aws_paas_pgsql"
description: |-
  Manages a PaaS PostgreSQL service.
---

[paas]: https://docs.k2.cloud/en/services/paas/index.html
[pgsql-version]: https://docs.k2.cloud/en/api/paas/parameters/pgsql.html#version
[timeouts]: https://developer.hashicorp.com/terraform/plugin/framework/resources/timeouts

# Resource: aws_paas_pgsql

Manages a PaaS PostgreSQL service. For details about PaaS, see the [user documentation][paas].

The resource is a PostgreSQL-only alternative to [`aws_paas_service`](paas_service.md) with the `pgsql` block.
It creates the same service in the cloud and supports the same arguments, but users and databases are always managed by
the standalone [`aws_paas_pgsql_user`](paas_pgsql_user.md) and [`aws_paas_pgsql_database`](paas_pgsql_database.md) resources.
Both resources are supported, so existing configurations with `aws_paas_service` keep working.

~> **Important** A PaaS service requires the VPC to have an internet gateway, a NAT gateway and a default route (`0.0.0.0/0`) to the NAT gateway; otherwise creating the service fails.
It's recommended to specify the route as an explicit dependency via `depends_on`.

## Example Usage

```terraform
resource "aws_paas_pgsql" "example" {
  depends_on = [aws_route.default_route]

  name          = "pgsql"
  instance_type = "c5.large"

  root_volume {
    type = "gp2"
    size = 32
  }

  data_volume {
    type = "gp2"
    size = 32
  }

  security_group_ids           = [aws_vpc.example.default_security_group_id]
  subnet_ids                   = [aws_subnet.example.id]
  delete_interfaces_on_destroy = true

  pgsql {
    version         = "16"
    max_connections = 200
  }
}

resource "aws_paas_pgsql_user" "app_owner" {
  service_id = aws_paas_pgsql.example.id
  name       = "app_owner"
  password   = var.app_owner_password
}

resource "aws_paas_pgsql_database" "app" {
  service_id = aws_paas_pgsql.example.id
  name       = "app"
  owner      = aws_paas_pgsql_user.app_owner.name
}
```

## Argument Reference

The resource supports the [common arguments](paas_service.md#argument-reference) of `aws_paas_service` with the following differences:

* `pgsql` - (Required) PostgreSQL parameters. The block supports the [PostgreSQL arguments](paas_service.md#postgresql-argument-reference)
  of `aws_paas_service`, except for the `user` and `database` blocks. The `version` argument is required,
  the list of supported versions is available in the [user documentation][pgsql-version].
  If `high_availability` is `true`, `pgsql.replication_mode` must be set; this is checked during `terraform plan`.
* Blocks with parameters of other service types (`clickhouse`, `kafka`, `redis` and others) are not supported.
* `additional_roles` and `coordinator` are not supported, because they apply only to ClickHouse and Kafka services.
* `users_managed_externally` and `databases_managed_externally` are not supported, because users and databases are always managed
  by the standalone resources.

As in `aws_paas_service`, the arguments are not editable (changes lead to a new resource) except for the `pgsql` block,
`backup_settings`, `auto_disk_expansion`, `logging` and `monitoring`.

## Attribute Reference

The resource exports the same [attributes](paas_service.md#attribute-reference) as `aws_paas_service`, except for the attributes
of the `pgsql.user` and `pgsql.database` blocks. In addition:

* `additional_roles` and `coordinator` - Always empty for PostgreSQL services.
* `databases_managed_externally` and `users_managed_externally` - Always `true`.

## Timeouts

The `timeouts` block allows you to specify [timeouts] for certain actions:

* `create` - (Default `30 minutes`) How long to wait for the service to be created.
* `update` - (Default `60 minutes`) How long to wait for the service to be updated.
* `delete` - (Default `15 minutes`) How long to wait for the service to be deleted.

## Import

PaaS PostgreSQL service can be imported using `id`, e.g.,

```
$ terraform import aws_paas_pgsql.example fm-cluster-12345678
```

Only PostgreSQL services can be imported. For services of other types, use [`aws_paas_service`](paas_service.md).

Users and databases of the service are not imported with it. Import them into
[`aws_paas_pgsql_user`](paas_pgsql_user.md#import) and [`aws_paas_pgsql_database`](paas_pgsql_database.md#import) resources.

-> **Note** The configuration must match the imported service, otherwise the plan shows a replacement.
Check the arguments that have defaults, for example, `root_volume.type` and `data_volume.type` (the default value is `st2`).

-> **Note** `delete_interfaces_on_destroy` is not stored in the cloud. If it is set in the configuration, the first plan after
import shows an in-place update of this argument. The update does not change the service.

-> **Note** The cloud API does not return `arbitrator_required`, so on import the value is determined from the service instances:
it is `true` if the service has an instance with the `arbitrator` role, and `false` otherwise.

### Moving a service from aws_paas_service

A PostgreSQL service managed by `aws_paas_service` can be moved to `aws_paas_pgsql` without recreation.
Terraform does not support `moved` blocks between different resource types, so use the `removed` and `import` blocks
(Terraform 1.7 or later):

1. Replace the `aws_paas_service` resource with `aws_paas_pgsql` and remove the `pgsql.user`, `pgsql.database`,
   `users_managed_externally` and `databases_managed_externally` arguments.
2. Add the `removed` and `import` blocks:

    ```terraform
    removed {
      from = aws_paas_service.example

      lifecycle {
        destroy = false
      }
    }

    import {
      to = aws_paas_pgsql.example
      id = "fm-cluster-12345678"
    }
    ```

3. Replace `aws_paas_service.example.id` with `aws_paas_pgsql.example.id` in the `aws_paas_pgsql_user` and `aws_paas_pgsql_database`
   resources. The service ID does not change, so these resources are not recreated.
4. If users and databases were managed with the `pgsql.user` and `pgsql.database` blocks, they remain in the cloud.
   Add `aws_paas_pgsql_user` and `aws_paas_pgsql_database` resources and import them, see
   [Moving existing users and databases](paas_service.md#moving-existing-users-and-databases).
5. Run `terraform plan` and make sure that the service is imported and not replaced, and `aws_paas_service` is only removed
   from the state. Then run `terraform apply` and remove the `removed` and `import` blocks.
